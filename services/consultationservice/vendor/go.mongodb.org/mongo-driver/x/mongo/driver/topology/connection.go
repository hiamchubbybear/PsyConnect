





package topology

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/internal/csot"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/ocsp"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)


const (
	connDisconnected int64 = iota
	connConnected
	connInitialized
)

var globalConnectionID uint64 = 1

var (
	defaultMaxMessageSize        uint32 = 48000000
	errResponseTooLarge                 = errors.New("length of read message too large")
	errLoadBalancedStateMismatch        = errors.New("driver attempted to initialize in load balancing mode, but the server does not support this mode")
)

func nextConnectionID() uint64 { return atomic.AddUint64(&globalConnectionID, 1) }

type connection struct {
	
	
	
	state int64

	id                   string
	nc                   net.Conn 
	addr                 address.Address
	idleTimeout          time.Duration
	idleStart            atomic.Value 
	readTimeout          time.Duration
	writeTimeout         time.Duration
	desc                 description.Server
	helloRTT             time.Duration
	compressor           wiremessage.CompressorID
	zliblevel            int
	zstdLevel            int
	connectDone          chan struct{}
	config               *connectionConfig
	cancelConnectContext context.CancelFunc
	connectContextMade   chan struct{}
	canStream            bool
	currentlyStreaming   bool
	connectContextMutex  sync.Mutex
	cancellationListener cancellationListener
	serverConnectionID   *int64 

	
	pool *pool

	
	driverConnectionID uint64
	generation         uint64

	
	
	awaitRemainingBytes *int32

	
	
	oidcTokenGenID uint64
}


func newConnection(addr address.Address, opts ...ConnectionOption) *connection {
	cfg := newConnectionConfig(opts...)

	id := fmt.Sprintf("%s[-%d]", addr, nextConnectionID())

	c := &connection{
		id:                   id,
		addr:                 addr,
		idleTimeout:          cfg.idleTimeout,
		readTimeout:          cfg.readTimeout,
		writeTimeout:         cfg.writeTimeout,
		connectDone:          make(chan struct{}),
		config:               cfg,
		connectContextMade:   make(chan struct{}),
		cancellationListener: newCancellListener(),
	}
	
	
	if !c.config.loadBalanced {
		c.setGenerationNumber()
	}
	atomic.StoreInt64(&c.state, connInitialized)

	return c
}



func (c *connection) setGenerationNumber() {
	if c.config.getGenerationFn != nil {
		c.generation = c.config.getGenerationFn(c.desc.ServiceID)
	}
}



func (c *connection) hasGenerationNumber() bool {
	if !c.config.loadBalanced {
		
		return true
	}

	
	
	return c.desc.LoadBalanced()
}

func configureTLS(ctx context.Context,
	tlsConnSource tlsConnectionSource,
	nc net.Conn,
	addr address.Address,
	config *tls.Config,
	ocspOpts *ocsp.VerifyOptions,
) (net.Conn, error) {
	
	if config.ServerName == "" {
		hostname := addr.String()
		colonPos := strings.LastIndex(hostname, ":")
		if colonPos == -1 {
			colonPos = len(hostname)
		}

		hostname = hostname[:colonPos]
		config.ServerName = hostname
	}

	client := tlsConnSource.Client(nc, config)
	if err := clientHandshake(ctx, client); err != nil {
		return nil, err
	}

	
	if !config.InsecureSkipVerify {
		if ocspErr := ocsp.Verify(ctx, client.ConnectionState(), ocspOpts); ocspErr != nil {
			return nil, ocspErr
		}
	}
	return client, nil
}




func (c *connection) connect(ctx context.Context) (err error) {
	if !atomic.CompareAndSwapInt64(&c.state, connInitialized, connConnected) {
		return nil
	}

	defer close(c.connectDone)

	
	
	defer func() {
		if err != nil {
			atomic.StoreInt64(&c.state, connDisconnected)

			if c.nc != nil {
				_ = c.nc.Close()
			}
		}
	}()

	
	
	
	
	
	
	
	
	
	
	c.connectContextMutex.Lock()
	var handshakeCtx context.Context
	handshakeCtx, c.cancelConnectContext = context.WithCancel(ctx)
	c.connectContextMutex.Unlock()

	dialCtx := handshakeCtx
	var dialCancel context.CancelFunc
	if c.config.connectTimeout != 0 {
		dialCtx, dialCancel = context.WithTimeout(handshakeCtx, c.config.connectTimeout)
		defer dialCancel()
	}

	defer func() {
		var cancelFn context.CancelFunc

		c.connectContextMutex.Lock()
		cancelFn = c.cancelConnectContext
		c.cancelConnectContext = nil
		c.connectContextMutex.Unlock()

		if cancelFn != nil {
			cancelFn()
		}
	}()

	close(c.connectContextMade)

	
	tempNc, err := c.config.dialer.DialContext(dialCtx, c.addr.Network(), c.addr.String())
	if err != nil {
		return ConnectionError{Wrapped: err, init: true}
	}
	c.nc = tempNc

	if c.config.tlsConfig != nil {
		tlsConfig := c.config.tlsConfig.Clone()

		
		
		ocspOpts := &ocsp.VerifyOptions{
			Cache:                   c.config.ocspCache,
			DisableEndpointChecking: c.config.disableOCSPEndpointCheck,
			HTTPClient:              c.config.httpClient,
		}
		tlsNc, err := configureTLS(dialCtx, c.config.tlsConnectionSource, c.nc, c.addr, tlsConfig, ocspOpts)
		if err != nil {
			return ConnectionError{Wrapped: err, init: true}
		}
		c.nc = tlsNc
	}

	
	handshaker := c.config.handshaker
	if handshaker == nil {
		return nil
	}

	var handshakeInfo driver.HandshakeInformation
	handshakeStartTime := time.Now()
	handshakeConn := initConnection{c}
	handshakeInfo, err = handshaker.GetHandshakeInformation(handshakeCtx, c.addr, handshakeConn)
	if err == nil {
		
		
		c.desc = handshakeInfo.Description
		c.serverConnectionID = handshakeInfo.ServerConnectionID
		c.helloRTT = time.Since(handshakeStartTime)

		
		
		if c.config.loadBalanced && c.desc.ServiceID == nil {
			err = errLoadBalancedStateMismatch
		}
	}
	if err == nil {
		
		
		
		if c.config.loadBalanced {
			c.setGenerationNumber()
		}

		
		
		err = handshaker.FinishHandshake(handshakeCtx, handshakeConn)
	}

	
	if err != nil {
		return ConnectionError{Wrapped: err, init: true}
	}

	if len(c.desc.Compression) > 0 {
	clientMethodLoop:
		for _, method := range c.config.compressors {
			for _, serverMethod := range c.desc.Compression {
				if method != serverMethod {
					continue
				}

				switch strings.ToLower(method) {
				case "snappy":
					c.compressor = wiremessage.CompressorSnappy
				case "zlib":
					c.compressor = wiremessage.CompressorZLib
					c.zliblevel = wiremessage.DefaultZlibLevel
					if c.config.zlibLevel != nil {
						c.zliblevel = *c.config.zlibLevel
					}
				case "zstd":
					c.compressor = wiremessage.CompressorZstd
					c.zstdLevel = wiremessage.DefaultZstdLevel
					if c.config.zstdLevel != nil {
						c.zstdLevel = *c.config.zstdLevel
					}
				}
				break clientMethodLoop
			}
		}
	}
	return nil
}

func (c *connection) wait() {
	if c.connectDone != nil {
		<-c.connectDone
	}
}

func (c *connection) closeConnectContext() {
	<-c.connectContextMade
	var cancelFn context.CancelFunc

	c.connectContextMutex.Lock()
	cancelFn = c.cancelConnectContext
	c.cancelConnectContext = nil
	c.connectContextMutex.Unlock()

	if cancelFn != nil {
		cancelFn()
	}
}

func (c *connection) cancellationListenerCallback() {
	_ = c.close()
}

func transformNetworkError(ctx context.Context, originalError error, contextDeadlineUsed bool) error {
	if originalError == nil {
		return nil
	}

	
	if errors.Is(ctx.Err(), context.Canceled) {
		return ctx.Err()
	}

	
	
	if !contextDeadlineUsed {
		return originalError
	}
	if netErr, ok := originalError.(net.Error); ok && netErr.Timeout() {
		return fmt.Errorf("%w: %s", context.DeadlineExceeded, originalError.Error())
	}

	return originalError
}

func (c *connection) writeWireMessage(ctx context.Context, wm []byte) error {
	var err error
	if atomic.LoadInt64(&c.state) != connConnected {
		return ConnectionError{
			ConnectionID: c.id,
			message:      "connection is closed",
		}
	}

	var deadline time.Time
	if c.writeTimeout != 0 {
		deadline = time.Now().Add(c.writeTimeout)
	}

	var contextDeadlineUsed bool
	if dl, ok := ctx.Deadline(); ok && (deadline.IsZero() || dl.Before(deadline)) {
		contextDeadlineUsed = true
		deadline = dl
	}

	if err := c.nc.SetWriteDeadline(deadline); err != nil {
		return ConnectionError{ConnectionID: c.id, Wrapped: err, message: "failed to set write deadline"}
	}

	err = c.write(ctx, wm)
	if err != nil {
		c.close()
		return ConnectionError{
			ConnectionID: c.id,
			Wrapped:      transformNetworkError(ctx, err, contextDeadlineUsed),
			message:      "unable to write wire message to network",
		}
	}

	return nil
}

func (c *connection) write(ctx context.Context, wm []byte) (err error) {
	go c.cancellationListener.Listen(ctx, c.cancellationListenerCallback)
	defer func() {
		
		
		
		

		if aborted := c.cancellationListener.StopListening(); aborted && err == nil {
			err = context.Canceled
		}
	}()

	_, err = c.nc.Write(wm)
	return err
}


func (c *connection) readWireMessage(ctx context.Context) ([]byte, error) {
	if atomic.LoadInt64(&c.state) != connConnected {
		return nil, ConnectionError{
			ConnectionID: c.id,
			message:      "connection is closed",
		}
	}

	var deadline time.Time
	if c.readTimeout != 0 {
		deadline = time.Now().Add(c.readTimeout)
	}

	var contextDeadlineUsed bool
	if dl, ok := ctx.Deadline(); ok && (deadline.IsZero() || dl.Before(deadline)) {
		contextDeadlineUsed = true
		deadline = dl
	}

	if err := c.nc.SetReadDeadline(deadline); err != nil {
		return nil, ConnectionError{ConnectionID: c.id, Wrapped: err, message: "failed to set read deadline"}
	}

	dst, errMsg, err := c.read(ctx)
	if err != nil {
		if c.awaitRemainingBytes == nil {
			
			
			
			c.close()
		}
		message := errMsg
		if errors.Is(err, io.EOF) {
			message = "socket was unexpectedly closed"
		}
		return nil, ConnectionError{
			ConnectionID: c.id,
			Wrapped:      transformNetworkError(ctx, err, contextDeadlineUsed),
			message:      message,
		}
	}

	return dst, nil
}

func (c *connection) parseWmSizeBytes(wmSizeBytes [4]byte) (int32, error) {
	
	size := int32(binary.LittleEndian.Uint32(wmSizeBytes[:]))

	if size < 4 {
		return 0, fmt.Errorf("malformed message length: %d", size)
	}
	
	
	maxMessageSize := c.desc.MaxMessageSize
	if maxMessageSize == 0 {
		maxMessageSize = defaultMaxMessageSize
	}
	if uint32(size) > maxMessageSize {
		return 0, errResponseTooLarge
	}

	return size, nil
}

func (c *connection) read(ctx context.Context) (bytesRead []byte, errMsg string, err error) {
	go c.cancellationListener.Listen(ctx, c.cancellationListenerCallback)
	defer func() {
		
		
		

		if aborted := c.cancellationListener.StopListening(); aborted && err == nil {
			errMsg = "unable to read server response"
			err = context.Canceled
		}
	}()

	isCSOTTimeout := func(err error) bool {
		
		
		
		
		nerr := net.Error(nil)
		return errors.As(err, &nerr) && nerr.Timeout() && csot.IsTimeoutContext(ctx)
	}

	
	
	var sizeBuf [4]byte

	
	
	
	n, err := io.ReadFull(c.nc, sizeBuf[:])
	if err != nil {
		if l := int32(n); l == 0 && isCSOTTimeout(err) {
			c.awaitRemainingBytes = &l
		}
		return nil, "incomplete read of message header", err
	}
	size, err := c.parseWmSizeBytes(sizeBuf)
	if err != nil {
		return nil, err.Error(), err
	}

	dst := make([]byte, size)
	copy(dst, sizeBuf[:])

	n, err = io.ReadFull(c.nc, dst[4:])
	if err != nil {
		remainingBytes := size - 4 - int32(n)
		if remainingBytes > 0 && isCSOTTimeout(err) {
			c.awaitRemainingBytes = &remainingBytes
		}
		return dst, "incomplete read of full message", err
	}

	return dst, "", nil
}

func (c *connection) close() error {
	
	if !atomic.CompareAndSwapInt64(&c.state, connConnected, connDisconnected) {
		return nil
	}

	var err error
	if c.nc != nil {
		err = c.nc.Close()
	}

	return err
}


func (c *connection) closed() bool {
	return atomic.LoadInt64(&c.state) == connDisconnected
}

func (c *connection) idleTimeoutExpired() bool {
	if c.idleTimeout == 0 {
		return false
	}

	idleStart, ok := c.idleStart.Load().(time.Time)
	return ok && idleStart.Add(c.idleTimeout).Before(time.Now())
}

func (c *connection) bumpIdleStart() {
	if c.idleTimeout > 0 {
		c.idleStart.Store(time.Now())
	}
}

func (c *connection) setCanStream(canStream bool) {
	c.canStream = canStream
}

func (c *connection) setStreaming(streaming bool) {
	c.currentlyStreaming = streaming
}

func (c *connection) getCurrentlyStreaming() bool {
	return c.currentlyStreaming
}

func (c *connection) setSocketTimeout(timeout time.Duration) {
	c.readTimeout = timeout
	c.writeTimeout = timeout
}



func (c *connection) DriverConnectionID() uint64 {
	return c.driverConnectionID
}

func (c *connection) ID() string {
	return c.id
}

func (c *connection) ServerConnectionID() *int64 {
	return c.serverConnectionID
}

func (c *connection) OIDCTokenGenID() uint64 {
	return c.oidcTokenGenID
}

func (c *connection) SetOIDCTokenGenID(genID uint64) {
	c.oidcTokenGenID = genID
}




type initConnection struct{ *connection }

var _ driver.Connection = initConnection{}
var _ driver.StreamerConnection = initConnection{}

func (c initConnection) Description() description.Server {
	if c.connection == nil {
		return description.Server{}
	}
	return c.connection.desc
}
func (c initConnection) Close() error             { return nil }
func (c initConnection) ID() string               { return c.id }
func (c initConnection) Address() address.Address { return c.addr }
func (c initConnection) Stale() bool              { return false }
func (c initConnection) LocalAddress() address.Address {
	if c.connection == nil || c.nc == nil {
		return address.Address("0.0.0.0")
	}
	return address.Address(c.nc.LocalAddr().String())
}
func (c initConnection) WriteWireMessage(ctx context.Context, wm []byte) error {
	return c.writeWireMessage(ctx, wm)
}
func (c initConnection) ReadWireMessage(ctx context.Context) ([]byte, error) {
	return c.readWireMessage(ctx)
}
func (c initConnection) SetStreaming(streaming bool) {
	c.setStreaming(streaming)
}
func (c initConnection) CurrentlyStreaming() bool {
	return c.getCurrentlyStreaming()
}
func (c initConnection) SupportsStreaming() bool {
	return c.canStream
}




type Connection struct {
	connection    *connection
	refCount      int
	cleanupPoolFn func()

	oidcTokenGenID uint64

	
	
	cleanupServerFn func()

	mu sync.RWMutex
}

var _ driver.Connection = (*Connection)(nil)
var _ driver.Expirable = (*Connection)(nil)
var _ driver.PinnedConnection = (*Connection)(nil)


func (c *Connection) WriteWireMessage(ctx context.Context, wm []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil {
		return ErrConnectionClosed
	}
	return c.connection.writeWireMessage(ctx, wm)
}



func (c *Connection) ReadWireMessage(ctx context.Context) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil {
		return nil, ErrConnectionClosed
	}
	return c.connection.readWireMessage(ctx)
}




func (c *Connection) CompressWireMessage(src, dst []byte) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil {
		return dst, ErrConnectionClosed
	}
	if c.connection.compressor == wiremessage.CompressorNoOp {
		return append(dst, src...), nil
	}
	_, reqid, respto, origcode, rem, ok := wiremessage.ReadHeader(src)
	if !ok {
		return dst, errors.New("wiremessage is too short to compress, less than 16 bytes")
	}
	idx, dst := wiremessage.AppendHeaderStart(dst, reqid, respto, wiremessage.OpCompressed)
	dst = wiremessage.AppendCompressedOriginalOpCode(dst, origcode)
	dst = wiremessage.AppendCompressedUncompressedSize(dst, int32(len(rem)))
	dst = wiremessage.AppendCompressedCompressorID(dst, c.connection.compressor)
	opts := driver.CompressionOpts{
		Compressor: c.connection.compressor,
		ZlibLevel:  c.connection.zliblevel,
		ZstdLevel:  c.connection.zstdLevel,
	}
	compressed, err := driver.CompressPayload(rem, opts)
	if err != nil {
		return nil, err
	}
	dst = wiremessage.AppendCompressedCompressedMessage(dst, compressed)
	return bsoncore.UpdateLength(dst, idx, int32(len(dst[idx:]))), nil
}


func (c *Connection) Description() description.Server {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil {
		return description.Server{}
	}
	return c.connection.desc
}



func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connection == nil || c.refCount > 0 {
		return nil
	}

	return c.cleanupReferences()
}


func (c *Connection) Expire() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connection == nil {
		return nil
	}

	_ = c.connection.close()
	return c.cleanupReferences()
}

func (c *Connection) cleanupReferences() error {
	err := c.connection.pool.checkIn(c.connection)
	if c.cleanupPoolFn != nil {
		c.cleanupPoolFn()
		c.cleanupPoolFn = nil
	}
	if c.cleanupServerFn != nil {
		c.cleanupServerFn()
		c.cleanupServerFn = nil
	}
	c.connection = nil
	return err
}


func (c *Connection) Alive() bool {
	return c.connection != nil
}


func (c *Connection) ID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil {
		return "<closed>"
	}
	return c.connection.id
}


func (c *Connection) ServerConnectionID() *int64 {
	if c.connection == nil {
		return nil
	}
	return c.connection.serverConnectionID
}


func (c *Connection) Stale() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connection.pool.stale(c.connection)
}


func (c *Connection) Address() address.Address {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil {
		return address.Address("0.0.0.0")
	}
	return c.connection.addr
}


func (c *Connection) LocalAddress() address.Address {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil || c.connection.nc == nil {
		return address.Address("0.0.0.0")
	}
	return address.Address(c.connection.nc.LocalAddr().String())
}


func (c *Connection) PinToCursor() error {
	return c.pin("cursor", c.connection.pool.pinConnectionToCursor, c.connection.pool.unpinConnectionFromCursor)
}


func (c *Connection) PinToTransaction() error {
	return c.pin("transaction", c.connection.pool.pinConnectionToTransaction, c.connection.pool.unpinConnectionFromTransaction)
}

func (c *Connection) pin(reason string, updatePoolFn, cleanupPoolFn func()) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connection == nil {
		return fmt.Errorf("attempted to pin a connection for a %s, but the connection has already been returned to the pool", reason)
	}

	
	
	if c.refCount == 0 {
		updatePoolFn()
		c.cleanupPoolFn = cleanupPoolFn
	}
	c.refCount++
	return nil
}


func (c *Connection) UnpinFromCursor() error {
	return c.unpin("cursor")
}


func (c *Connection) UnpinFromTransaction() error {
	return c.unpin("transaction")
}

func (c *Connection) unpin(reason string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connection == nil {
		
		return nil
	}
	if c.refCount == 0 {
		return fmt.Errorf("attempted to unpin a connection from a %s, but the connection is not pinned by any resources", reason)
	}

	c.refCount--
	return nil
}



func (c *Connection) DriverConnectionID() uint64 {
	return c.connection.DriverConnectionID()
}


func (c *Connection) OIDCTokenGenID() uint64 {
	return c.oidcTokenGenID
}


func (c *Connection) SetOIDCTokenGenID(genID uint64) {
	c.oidcTokenGenID = genID
}





type cancellListener struct {
	aborted bool
	done    chan struct{}
}


func newCancellListener() *cancellListener {
	return &cancellListener{
		done: make(chan struct{}),
	}
}






func (c *cancellListener) Listen(ctx context.Context, abortFn func()) {
	c.aborted = false

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.Canceled) {
			c.aborted = true
			abortFn()
		}

		<-c.done
	case <-c.done:
	}
}




func (c *cancellListener) StopListening() bool {
	c.done <- struct{}{}
	return c.aborted
}
