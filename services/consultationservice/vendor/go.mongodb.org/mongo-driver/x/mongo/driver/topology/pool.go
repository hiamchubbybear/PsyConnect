





package topology

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)


const (
	poolPaused int = iota
	poolReady
	poolClosed
)



var ErrPoolNotPaused = PoolError("only a paused pool can be marked ready")


var ErrPoolClosed = PoolError("attempted to check out a connection from closed connection pool")


var ErrConnectionClosed = ConnectionError{ConnectionID: "<closed>", message: "connection is closed"}


var ErrWrongPool = PoolError("connection does not belong to this pool")


type PoolError string

func (pe PoolError) Error() string { return string(pe) }



type poolClearedError struct {
	err     error
	address address.Address
}

func (pce poolClearedError) Error() string {
	return fmt.Sprintf(
		"connection pool for %v was cleared because another operation failed with: %v",
		pce.address,
		pce.err)
}


func (poolClearedError) Retryable() bool { return true }


var _ driver.RetryablePoolError = poolClearedError{}


type poolConfig struct {
	Address          address.Address
	MinPoolSize      uint64
	MaxPoolSize      uint64
	MaxConnecting    uint64
	MaxIdleTime      time.Duration
	MaintainInterval time.Duration
	LoadBalanced     bool
	PoolMonitor      *event.PoolMonitor
	Logger           *logger.Logger
	handshakeErrFn   func(error, uint64, *primitive.ObjectID)
}

type pool struct {
	
	
	
	

	nextID                       uint64 
	pinnedCursorConnections      uint64
	pinnedTransactionConnections uint64

	address       address.Address
	minSize       uint64
	maxSize       uint64
	maxConnecting uint64
	loadBalanced  bool
	monitor       *event.PoolMonitor
	logger        *logger.Logger

	
	
	handshakeErrFn func(error, uint64, *primitive.ObjectID)

	connOpts   []ConnectionOption
	generation *poolGenerationMap

	maintainInterval time.Duration   
	maintainReady    chan struct{}   
	backgroundDone   *sync.WaitGroup 

	stateMu      sync.RWMutex 
	state        int          
	lastClearErr error        

	
	
	
	
	createConnectionsCond *sync.Cond
	cancelBackgroundCtx   context.CancelFunc     
	conns                 map[uint64]*connection 
	newConnWait           wantConnQueue          

	idleMu       sync.Mutex    
	idleConns    []*connection 
	idleConnWait wantConnQueue 
}


func (p *pool) getState() int {
	p.stateMu.RLock()
	defer p.stateMu.RUnlock()

	return p.state
}

func mustLogPoolMessage(pool *pool) bool {
	return pool.logger != nil && pool.logger.LevelComponentEnabled(
		logger.LevelDebug, logger.ComponentConnection)
}

func logPoolMessage(pool *pool, msg string, keysAndValues ...interface{}) {
	host, port, err := net.SplitHostPort(pool.address.String())
	if err != nil {
		host = pool.address.String()
		port = ""
	}

	pool.logger.Print(logger.LevelDebug,
		logger.ComponentConnection,
		msg,
		logger.SerializeConnection(logger.Connection{
			Message:    msg,
			ServerHost: host,
			ServerPort: port,
		}, keysAndValues...)...)

}

type reason struct {
	loggerConn string
	event      string
}


func connectionPerished(conn *connection) (reason, bool) {
	switch {
	case conn.closed():
		
		
		return reason{
			loggerConn: logger.ReasonConnClosedError,
			event:      event.ReasonError,
		}, true
	case conn.idleTimeoutExpired():
		return reason{
			loggerConn: logger.ReasonConnClosedIdle,
			event:      event.ReasonIdle,
		}, true
	case conn.pool.stale(conn):
		return reason{
			loggerConn: logger.ReasonConnClosedStale,
			event:      event.ReasonStale,
		}, true
	}

	return reason{}, false
}


func newPool(config poolConfig, connOpts ...ConnectionOption) *pool {
	if config.MaxIdleTime != time.Duration(0) {
		connOpts = append(connOpts, WithIdleTimeout(func(_ time.Duration) time.Duration { return config.MaxIdleTime }))
	}

	var maxConnecting uint64 = 2
	if config.MaxConnecting > 0 {
		maxConnecting = config.MaxConnecting
	}

	maintainInterval := 10 * time.Second
	if config.MaintainInterval != 0 {
		maintainInterval = config.MaintainInterval
	}

	pool := &pool{
		address:               config.Address,
		minSize:               config.MinPoolSize,
		maxSize:               config.MaxPoolSize,
		maxConnecting:         maxConnecting,
		loadBalanced:          config.LoadBalanced,
		monitor:               config.PoolMonitor,
		logger:                config.Logger,
		handshakeErrFn:        config.handshakeErrFn,
		connOpts:              connOpts,
		generation:            newPoolGenerationMap(),
		state:                 poolPaused,
		maintainInterval:      maintainInterval,
		maintainReady:         make(chan struct{}, 1),
		backgroundDone:        &sync.WaitGroup{},
		createConnectionsCond: sync.NewCond(&sync.Mutex{}),
		conns:                 make(map[uint64]*connection, config.MaxPoolSize),
		idleConns:             make([]*connection, 0, config.MaxPoolSize),
	}
	
	if pool.maxSize != 0 && pool.minSize > pool.maxSize {
		pool.minSize = pool.maxSize
	}
	pool.connOpts = append(pool.connOpts, withGenerationNumberFn(func(_ generationNumberFn) generationNumberFn { return pool.getGenerationForNewConnection }))

	pool.generation.connect()

	
	
	
	var ctx context.Context
	ctx, pool.cancelBackgroundCtx = context.WithCancel(context.Background())

	for i := 0; i < int(pool.maxConnecting); i++ {
		pool.backgroundDone.Add(1)
		go pool.createConnections(ctx, pool.backgroundDone)
	}

	
	
	if maintainInterval > 0 {
		pool.backgroundDone.Add(1)
		go pool.maintain(ctx, pool.backgroundDone)
	}

	if mustLogPoolMessage(pool) {
		keysAndValues := logger.KeyValues{
			logger.KeyMaxIdleTimeMS, config.MaxIdleTime.Milliseconds(),
			logger.KeyMinPoolSize, config.MinPoolSize,
			logger.KeyMaxPoolSize, config.MaxPoolSize,
			logger.KeyMaxConnecting, config.MaxConnecting,
		}

		logPoolMessage(pool, logger.ConnectionPoolCreated, keysAndValues...)
	}

	if pool.monitor != nil {
		pool.monitor.Event(&event.PoolEvent{
			Type: event.PoolCreated,
			PoolOptions: &event.MonitorPoolOptions{
				MaxPoolSize: config.MaxPoolSize,
				MinPoolSize: config.MinPoolSize,
			},
			Address: pool.address.String(),
		})
	}

	return pool
}


func (p *pool) stale(conn *connection) bool {
	return conn == nil || p.generation.stale(conn.desc.ServiceID, conn.generation)
}




func (p *pool) ready() error {
	
	p.stateMu.Lock()
	if p.state == poolReady {
		p.stateMu.Unlock()
		return nil
	}
	if p.state != poolPaused {
		p.stateMu.Unlock()
		return ErrPoolNotPaused
	}
	p.lastClearErr = nil
	p.state = poolReady
	p.stateMu.Unlock()

	if mustLogPoolMessage(p) {
		logPoolMessage(p, logger.ConnectionPoolReady)
	}

	
	
	if p.monitor != nil {
		p.monitor.Event(&event.PoolEvent{
			Type:    event.PoolReady,
			Address: p.address.String(),
		})
	}

	
	select {
	case p.maintainReady <- struct{}{}:
	default:
	}

	return nil
}




func (p *pool) close(ctx context.Context) {
	p.stateMu.Lock()
	if p.state == poolClosed {
		p.stateMu.Unlock()
		return
	}
	p.state = poolClosed
	p.stateMu.Unlock()

	
	
	
	
	
	
	p.createConnectionsCond.L.Lock()
	p.cancelBackgroundCtx()
	p.createConnectionsCond.Broadcast()
	p.createConnectionsCond.L.Unlock()

	
	p.backgroundDone.Wait()

	p.generation.disconnect()

	if ctx == nil {
		ctx = context.Background()
	}

	
	
	
	if _, ok := ctx.Deadline(); ok {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

	graceful:
		for {
			if p.totalConnectionCount() == p.availableConnectionCount() {
				break graceful
			}

			select {
			case <-ticker.C:
			case <-ctx.Done():
				break graceful
			default:
			}
		}
	}

	
	
	p.idleMu.Lock()
	for _, conn := range p.idleConns {
		_ = p.removeConnection(conn, reason{
			loggerConn: logger.ReasonConnClosedPoolClosed,
			event:      event.ReasonPoolClosed,
		}, nil)
		_ = p.closeConnection(conn) 
	}
	p.idleConns = p.idleConns[:0]
	for {
		w := p.idleConnWait.popFront()
		if w == nil {
			break
		}
		w.tryDeliver(nil, ErrPoolClosed)
	}
	p.idleMu.Unlock()

	
	
	
	p.createConnectionsCond.L.Lock()
	conns := make([]*connection, 0, len(p.conns))
	for _, conn := range p.conns {
		conns = append(conns, conn)
	}
	for {
		w := p.newConnWait.popFront()
		if w == nil {
			break
		}
		w.tryDeliver(nil, ErrPoolClosed)
	}
	p.createConnectionsCond.L.Unlock()

	if mustLogPoolMessage(p) {
		logPoolMessage(p, logger.ConnectionPoolClosed)
	}

	if p.monitor != nil {
		p.monitor.Event(&event.PoolEvent{
			Type:    event.PoolClosedEvent,
			Address: p.address.String(),
		})
	}

	
	
	for _, conn := range conns {
		_ = p.removeConnection(conn, reason{
			loggerConn: logger.ReasonConnClosedPoolClosed,
			event:      event.ReasonPoolClosed,
		}, nil)
		_ = p.closeConnection(conn) 
	}
}

func (p *pool) pinConnectionToCursor() {
	atomic.AddUint64(&p.pinnedCursorConnections, 1)
}

func (p *pool) unpinConnectionFromCursor() {
	
	atomic.AddUint64(&p.pinnedCursorConnections, ^uint64(0))
}

func (p *pool) pinConnectionToTransaction() {
	atomic.AddUint64(&p.pinnedTransactionConnections, 1)
}

func (p *pool) unpinConnectionFromTransaction() {
	
	atomic.AddUint64(&p.pinnedTransactionConnections, ^uint64(0))
}





func (p *pool) checkOut(ctx context.Context) (conn *connection, err error) {
	if mustLogPoolMessage(p) {
		logPoolMessage(p, logger.ConnectionCheckoutStarted)
	}

	
	
	if p.monitor != nil {
		p.monitor.Event(&event.PoolEvent{
			Type:    event.GetStarted,
			Address: p.address.String(),
		})
	}

	start := time.Now()
	
	
	
	
	
	p.stateMu.RLock()
	switch p.state {
	case poolClosed:
		p.stateMu.RUnlock()

		duration := time.Since(start)
		if mustLogPoolMessage(p) {
			keysAndValues := logger.KeyValues{
				logger.KeyDurationMS, duration.Milliseconds(),
				logger.KeyReason, logger.ReasonConnCheckoutFailedPoolClosed,
			}

			logPoolMessage(p, logger.ConnectionCheckoutFailed, keysAndValues...)
		}

		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:     event.GetFailed,
				Address:  p.address.String(),
				Duration: duration,
				Reason:   event.ReasonPoolClosed,
			})
		}
		return nil, ErrPoolClosed
	case poolPaused:
		err := poolClearedError{err: p.lastClearErr, address: p.address}
		p.stateMu.RUnlock()

		duration := time.Since(start)
		if mustLogPoolMessage(p) {
			keysAndValues := logger.KeyValues{
				logger.KeyDurationMS, duration.Milliseconds(),
				logger.KeyReason, logger.ReasonConnCheckoutFailedError,
			}

			logPoolMessage(p, logger.ConnectionCheckoutFailed, keysAndValues...)
		}

		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:     event.GetFailed,
				Address:  p.address.String(),
				Duration: duration,
				Reason:   event.ReasonConnectionErrored,
				Error:    err,
			})
		}
		return nil, err
	}

	if ctx == nil {
		ctx = context.Background()
	}

	
	
	
	
	w := newWantConn()
	defer func() {
		if err != nil {
			w.cancel(p, err)
		}
	}()

	
	
	
	if delivered := p.getOrQueueForIdleConn(w); delivered {
		
		
		p.stateMu.RUnlock()

		duration := time.Since(start)
		if w.err != nil {
			if mustLogPoolMessage(p) {
				keysAndValues := logger.KeyValues{
					logger.KeyDurationMS, duration.Milliseconds(),
					logger.KeyReason, logger.ReasonConnCheckoutFailedError,
				}

				logPoolMessage(p, logger.ConnectionCheckoutFailed, keysAndValues...)
			}

			if p.monitor != nil {
				p.monitor.Event(&event.PoolEvent{
					Type:     event.GetFailed,
					Address:  p.address.String(),
					Duration: duration,
					Reason:   event.ReasonConnectionErrored,
					Error:    w.err,
				})
			}
			return nil, w.err
		}

		duration = time.Since(start)
		if mustLogPoolMessage(p) {
			keysAndValues := logger.KeyValues{
				logger.KeyDriverConnectionID, w.conn.driverConnectionID,
				logger.KeyDurationMS, duration.Milliseconds(),
			}

			logPoolMessage(p, logger.ConnectionCheckedOut, keysAndValues...)
		}

		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:         event.GetSucceeded,
				Address:      p.address.String(),
				ConnectionID: w.conn.driverConnectionID,
				Duration:     duration,
			})
		}

		return w.conn, nil
	}

	
	
	p.queueForNewConn(w)
	p.stateMu.RUnlock()

	
	waitQueueStart := time.Now()
	select {
	case <-w.ready:
		if w.err != nil {
			duration := time.Since(start)
			if mustLogPoolMessage(p) {
				keysAndValues := logger.KeyValues{
					logger.KeyDurationMS, duration.Milliseconds(),
					logger.KeyReason, logger.ReasonConnCheckoutFailedError,
					logger.KeyError, w.err.Error(),
				}

				logPoolMessage(p, logger.ConnectionCheckoutFailed, keysAndValues...)
			}

			if p.monitor != nil {
				p.monitor.Event(&event.PoolEvent{
					Type:     event.GetFailed,
					Address:  p.address.String(),
					Duration: duration,
					Reason:   event.ReasonConnectionErrored,
					Error:    w.err,
				})
			}

			return nil, w.err
		}

		duration := time.Since(start)
		if mustLogPoolMessage(p) {
			keysAndValues := logger.KeyValues{
				logger.KeyDriverConnectionID, w.conn.driverConnectionID,
				logger.KeyDurationMS, duration.Milliseconds(),
			}

			logPoolMessage(p, logger.ConnectionCheckedOut, keysAndValues...)
		}

		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:         event.GetSucceeded,
				Address:      p.address.String(),
				ConnectionID: w.conn.driverConnectionID,
				Duration:     duration,
			})
		}
		return w.conn, nil
	case <-ctx.Done():
		waitQueueDuration := time.Since(waitQueueStart)

		duration := time.Since(start)
		if mustLogPoolMessage(p) {
			keysAndValues := logger.KeyValues{
				logger.KeyDurationMS, duration.Milliseconds(),
				logger.KeyReason, logger.ReasonConnCheckoutFailedTimout,
			}

			logPoolMessage(p, logger.ConnectionCheckoutFailed, keysAndValues...)
		}

		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:     event.GetFailed,
				Address:  p.address.String(),
				Duration: duration,
				Reason:   event.ReasonTimedOut,
				Error:    ctx.Err(),
			})
		}

		err := WaitQueueTimeoutError{
			Wrapped:              ctx.Err(),
			maxPoolSize:          p.maxSize,
			totalConnections:     p.totalConnectionCount(),
			availableConnections: p.availableConnectionCount(),
			waitDuration:         waitQueueDuration,
		}
		if p.loadBalanced {
			err.pinnedConnections = &pinnedConnections{
				cursorConnections:      atomic.LoadUint64(&p.pinnedCursorConnections),
				transactionConnections: atomic.LoadUint64(&p.pinnedTransactionConnections),
			}
		}
		return nil, err
	}
}


func (p *pool) closeConnection(conn *connection) error {
	if conn.pool != p {
		return ErrWrongPool
	}

	if atomic.LoadInt64(&conn.state) == connConnected {
		conn.closeConnectContext()
		conn.wait() 
	}

	err := conn.close()
	if err != nil {
		return ConnectionError{ConnectionID: conn.id, Wrapped: err, message: "failed to close net.Conn"}
	}

	return nil
}

func (p *pool) getGenerationForNewConnection(serviceID *primitive.ObjectID) uint64 {
	return p.generation.addConnection(serviceID)
}


func (p *pool) removeConnection(conn *connection, reason reason, err error) error {
	if conn == nil {
		return nil
	}

	if conn.pool != p {
		return ErrWrongPool
	}

	p.createConnectionsCond.L.Lock()
	_, ok := p.conns[conn.driverConnectionID]
	if !ok {
		
		
		p.createConnectionsCond.L.Unlock()
		return nil
	}
	delete(p.conns, conn.driverConnectionID)
	
	
	p.createConnectionsCond.Signal()
	p.createConnectionsCond.L.Unlock()

	
	
	
	if conn.hasGenerationNumber() {
		p.generation.removeConnection(conn.desc.ServiceID)
	}

	if mustLogPoolMessage(p) {
		keysAndValues := logger.KeyValues{
			logger.KeyDriverConnectionID, conn.driverConnectionID,
			logger.KeyReason, reason.loggerConn,
		}

		if err != nil {
			keysAndValues.Add(logger.KeyError, err.Error())
		}

		logPoolMessage(p, logger.ConnectionClosed, keysAndValues...)
	}

	if p.monitor != nil {
		p.monitor.Event(&event.PoolEvent{
			Type:         event.ConnectionClosed,
			Address:      p.address.String(),
			ConnectionID: conn.driverConnectionID,
			Reason:       reason.event,
			Error:        err,
		})
	}

	return nil
}

var (
	
	
	
	
	
	
	BGReadTimeout = 1 * time.Second

	
	
	
	
	
	BGReadCallback func(addr string, start, read time.Time, errs []error, connClosed bool)
)








func bgRead(pool *pool, conn *connection, size int32) {
	var err error
	start := time.Now()

	defer func() {
		read := time.Now()
		errs := make([]error, 0)
		connClosed := false
		if err != nil {
			errs = append(errs, err)
			connClosed = true
			err = conn.close()
			if err != nil {
				errs = append(errs, fmt.Errorf("error closing conn after reading: %w", err))
			}
		}

		
		
		
		err = pool.checkInNoEvent(conn)
		if err != nil {
			errs = append(errs, fmt.Errorf("error checking in: %w", err))
		}

		if BGReadCallback != nil {
			BGReadCallback(conn.addr.String(), start, read, errs, connClosed)
		}
	}()

	err = conn.nc.SetReadDeadline(time.Now().Add(BGReadTimeout))
	if err != nil {
		err = fmt.Errorf("error setting a read deadline: %w", err)
		return
	}

	if size == 0 {
		var sizeBuf [4]byte
		_, err = io.ReadFull(conn.nc, sizeBuf[:])
		if err != nil {
			err = fmt.Errorf("error reading the message size: %w", err)
			return
		}
		size, err = conn.parseWmSizeBytes(sizeBuf)
		if err != nil {
			return
		}
		size -= 4
	}
	_, err = io.CopyN(io.Discard, conn.nc, int64(size))
	if err != nil {
		err = fmt.Errorf("error discarding %d byte message: %w", size, err)
	}
}



func (p *pool) checkIn(conn *connection) error {
	if conn == nil {
		return nil
	}
	if conn.pool != p {
		return ErrWrongPool
	}

	if mustLogPoolMessage(p) {
		keysAndValues := logger.KeyValues{
			logger.KeyDriverConnectionID, conn.driverConnectionID,
		}

		logPoolMessage(p, logger.ConnectionCheckedIn, keysAndValues...)
	}

	if p.monitor != nil {
		p.monitor.Event(&event.PoolEvent{
			Type:         event.ConnectionReturned,
			ConnectionID: conn.driverConnectionID,
			Address:      conn.addr.String(),
		})
	}

	return p.checkInNoEvent(conn)
}



func (p *pool) checkInNoEvent(conn *connection) error {
	if conn == nil {
		return nil
	}
	if conn.pool != p {
		return ErrWrongPool
	}

	
	
	
	
	
	
	
	
	if conn.awaitRemainingBytes != nil {
		size := *conn.awaitRemainingBytes
		conn.awaitRemainingBytes = nil
		go bgRead(p, conn, size)
		return nil
	}

	
	
	
	
	
	
	
	
	conn.bumpIdleStart()

	r, perished := connectionPerished(conn)
	if !perished && conn.pool.getState() == poolClosed {
		perished = true
		r = reason{
			loggerConn: logger.ReasonConnClosedPoolClosed,
			event:      event.ReasonPoolClosed,
		}
	}
	if perished {
		_ = p.removeConnection(conn, r, nil)
		go func() {
			_ = p.closeConnection(conn)
		}()
		return nil
	}

	p.idleMu.Lock()
	defer p.idleMu.Unlock()

	for {
		w := p.idleConnWait.popFront()
		if w == nil {
			break
		}
		if w.tryDeliver(conn, nil) {
			return nil
		}
	}

	for _, idle := range p.idleConns {
		if idle == conn {
			return fmt.Errorf("duplicate idle conn %p in idle connections stack", conn)
		}
	}

	p.idleConns = append(p.idleConns, conn)
	return nil
}


func (p *pool) clear(err error, serviceID *primitive.ObjectID) {
	p.clearImpl(err, serviceID, false)
}


func (p *pool) clearAll(err error, serviceID *primitive.ObjectID) {
	p.clearImpl(err, serviceID, true)
}


func (p *pool) interruptConnections(conns []*connection) {
	for _, conn := range conns {
		_ = p.removeConnection(conn, reason{
			loggerConn: logger.ReasonConnClosedStale,
			event:      event.ReasonStale,
		}, nil)
		go func(c *connection) {
			_ = p.closeConnection(c)
		}(conn)
	}
}








func (p *pool) clearImpl(err error, serviceID *primitive.ObjectID, interruptAllConnections bool) {
	if p.getState() == poolClosed {
		return
	}

	p.generation.clear(serviceID)

	
	
	
	sendEvent := true
	if serviceID == nil {
		
		
		
		p.stateMu.Lock()
		if p.state == poolPaused {
			sendEvent = false
		}
		if p.state == poolReady {
			p.state = poolPaused
		}
		p.lastClearErr = err
		p.stateMu.Unlock()
	}

	if mustLogPoolMessage(p) {
		keysAndValues := logger.KeyValues{
			logger.KeyServiceID, serviceID,
		}

		logPoolMessage(p, logger.ConnectionPoolCleared, keysAndValues...)
	}

	if sendEvent && p.monitor != nil {
		event := &event.PoolEvent{
			Type:         event.PoolCleared,
			Address:      p.address.String(),
			ServiceID:    serviceID,
			Interruption: interruptAllConnections,
			Error:        err,
		}
		p.monitor.Event(event)
	}

	p.removePerishedConns()
	if interruptAllConnections {
		p.createConnectionsCond.L.Lock()
		p.idleMu.Lock()

		idleConns := make(map[*connection]bool, len(p.idleConns))
		for _, idle := range p.idleConns {
			idleConns[idle] = true
		}

		conns := make([]*connection, 0, len(p.conns))
		for _, conn := range p.conns {
			if _, ok := idleConns[conn]; !ok && p.stale(conn) {
				conns = append(conns, conn)
			}
		}

		p.idleMu.Unlock()
		p.createConnectionsCond.L.Unlock()

		p.interruptConnections(conns)
	}

	if serviceID == nil {
		pcErr := poolClearedError{err: err, address: p.address}

		
		p.idleMu.Lock()
		for {
			w := p.idleConnWait.popFront()
			if w == nil {
				break
			}
			w.tryDeliver(nil, pcErr)
		}
		p.idleMu.Unlock()

		
		
		
		p.createConnectionsCond.L.Lock()
		for {
			w := p.newConnWait.popFront()
			if w == nil {
				break
			}
			w.tryDeliver(nil, pcErr)
		}
		p.createConnectionsCond.L.Unlock()
	}
}





func (p *pool) getOrQueueForIdleConn(w *wantConn) bool {
	p.idleMu.Lock()
	defer p.idleMu.Unlock()

	
	for len(p.idleConns) > 0 {
		conn := p.idleConns[len(p.idleConns)-1]
		p.idleConns = p.idleConns[:len(p.idleConns)-1]

		if conn == nil {
			continue
		}

		if reason, perished := connectionPerished(conn); perished {
			_ = conn.pool.removeConnection(conn, reason, nil)
			go func() {
				_ = conn.pool.closeConnection(conn)
			}()
			continue
		}

		if !w.tryDeliver(conn, nil) {
			
			p.idleConns = append(p.idleConns, conn)
		}

		
		
		
		return true
	}

	p.idleConnWait.cleanFront()
	p.idleConnWait.pushBack(w)
	return false
}

func (p *pool) queueForNewConn(w *wantConn) {
	p.createConnectionsCond.L.Lock()
	defer p.createConnectionsCond.L.Unlock()

	p.newConnWait.cleanFront()
	p.newConnWait.pushBack(w)
	p.createConnectionsCond.Signal()
}

func (p *pool) totalConnectionCount() int {
	p.createConnectionsCond.L.Lock()
	defer p.createConnectionsCond.L.Unlock()

	return len(p.conns)
}

func (p *pool) availableConnectionCount() int {
	p.idleMu.Lock()
	defer p.idleMu.Unlock()

	return len(p.idleConns)
}


func (p *pool) createConnections(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	
	
	
	condition := func() bool {
		checkOutWaiting := p.newConnWait.len() > 0
		poolHasSpace := p.maxSize == 0 || uint64(len(p.conns)) < p.maxSize
		cancelled := ctx.Err() != nil
		return (checkOutWaiting && poolHasSpace) || cancelled
	}

	
	
	
	
	wait := func() (*wantConn, *connection, bool) {
		p.createConnectionsCond.L.Lock()
		defer p.createConnectionsCond.L.Unlock()

		for !condition() {
			p.createConnectionsCond.Wait()
		}

		if ctx.Err() != nil {
			return nil, nil, false
		}

		p.newConnWait.cleanFront()
		w := p.newConnWait.popFront()
		if w == nil {
			return nil, nil, false
		}

		conn := newConnection(p.address, p.connOpts...)
		conn.pool = p
		conn.driverConnectionID = atomic.AddUint64(&p.nextID, 1)
		p.conns[conn.driverConnectionID] = conn

		return w, conn, true
	}

	for ctx.Err() == nil {
		w, conn, ok := wait()
		if !ok {
			continue
		}

		if mustLogPoolMessage(p) {
			keysAndValues := logger.KeyValues{
				logger.KeyDriverConnectionID, conn.driverConnectionID,
			}

			logPoolMessage(p, logger.ConnectionCreated, keysAndValues...)
		}

		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:         event.ConnectionCreated,
				Address:      p.address.String(),
				ConnectionID: conn.driverConnectionID,
			})
		}

		start := time.Now()
		
		
		err := conn.connect(ctx)
		if err != nil {
			w.tryDeliver(nil, err)

			
			
			
			
			
			
			if p.handshakeErrFn != nil {
				p.handshakeErrFn(err, conn.generation, conn.desc.ServiceID)
			}

			_ = p.removeConnection(conn, reason{
				loggerConn: logger.ReasonConnClosedError,
				event:      event.ReasonError,
			}, err)

			_ = p.closeConnection(conn)

			continue
		}

		duration := time.Since(start)
		if mustLogPoolMessage(p) {
			keysAndValues := logger.KeyValues{
				logger.KeyDriverConnectionID, conn.driverConnectionID,
				logger.KeyDurationMS, duration.Milliseconds(),
			}

			logPoolMessage(p, logger.ConnectionReady, keysAndValues...)
		}

		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:         event.ConnectionReady,
				Address:      p.address.String(),
				ConnectionID: conn.driverConnectionID,
				Duration:     duration,
			})
		}

		if w.tryDeliver(conn, nil) {
			continue
		}

		_ = p.checkInNoEvent(conn)
	}
}

func (p *pool) maintain(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(p.maintainInterval)
	defer ticker.Stop()

	
	
	remove := func(arr []*wantConn, i int) []*wantConn {
		end := len(arr) - 1
		arr[i], arr[end] = arr[end], arr[i]
		return arr[:end]
	}

	
	
	
	
	removeNotWaiting := func(arr []*wantConn) []*wantConn {
		for i := len(arr) - 1; i >= 0; i-- {
			w := arr[i]
			if !w.waiting() {
				arr = remove(arr, i)
			}
		}

		return arr
	}

	wantConns := make([]*wantConn, 0, p.minSize)
	defer func() {
		for _, w := range wantConns {
			w.tryDeliver(nil, ErrPoolClosed)
		}
	}()

	for {
		select {
		case <-ticker.C:
		case <-p.maintainReady:
		case <-ctx.Done():
			return
		}

		
		
		
		
		
		
		p.stateMu.RLock()
		if p.state != poolReady {
			p.stateMu.RUnlock()
			continue
		}

		p.removePerishedConns()

		
		wantConns = removeNotWaiting(wantConns)

		
		
		
		
		
		total := p.totalConnectionCount()
		n := int(p.minSize) - total - len(wantConns)
		if n > 10 {
			n = 10
		}

		for i := 0; i < n; i++ {
			w := newWantConn()
			p.queueForNewConn(w)
			wantConns = append(wantConns, w)

			
			go func() {
				<-w.ready
				if w.conn != nil {
					_ = p.checkInNoEvent(w.conn)
				}
			}()
		}
		p.stateMu.RUnlock()
	}
}

func (p *pool) removePerishedConns() {
	p.idleMu.Lock()
	defer p.idleMu.Unlock()

	for i := range p.idleConns {
		conn := p.idleConns[i]
		if conn == nil {
			continue
		}

		if reason, perished := connectionPerished(conn); perished {
			p.idleConns[i] = nil

			_ = p.removeConnection(conn, reason, nil)
			go func() {
				_ = p.closeConnection(conn)
			}()
		}
	}

	p.idleConns = compact(p.idleConns)
}



func compact(arr []*connection) []*connection {
	offset := 0
	for i := range arr {
		if arr[i] == nil {
			continue
		}
		arr[offset] = arr[i]
		offset++
	}
	return arr[:offset]
}






type wantConn struct {
	ready chan struct{}

	mu   sync.Mutex 
	conn *connection
	err  error
}

func newWantConn() *wantConn {
	return &wantConn{
		ready: make(chan struct{}, 1),
	}
}


func (w *wantConn) waiting() bool {
	select {
	case <-w.ready:
		return false
	default:
		return true
	}
}


func (w *wantConn) tryDeliver(conn *connection, err error) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.conn != nil || w.err != nil {
		return false
	}

	w.conn = conn
	w.err = err
	if w.conn == nil && w.err == nil {
		panic("x/mongo/driver/topology: internal error: misuse of tryDeliver")
	}

	close(w.ready)

	return true
}




func (w *wantConn) cancel(p *pool, err error) {
	if err == nil {
		panic("x/mongo/driver/topology: internal error: misuse of cancel")
	}

	w.mu.Lock()
	if w.conn == nil && w.err == nil {
		close(w.ready) 
	}
	conn := w.conn
	w.conn = nil
	w.err = err
	w.mu.Unlock()

	if conn != nil {
		_ = p.checkInNoEvent(conn)
	}
}



type wantConnQueue struct {
	
	
	
	
	
	
	
	
	
	
	head    []*wantConn
	headPos int
	tail    []*wantConn
}


func (q *wantConnQueue) len() int {
	return len(q.head) - q.headPos + len(q.tail)
}


func (q *wantConnQueue) pushBack(w *wantConn) {
	q.tail = append(q.tail, w)
}


func (q *wantConnQueue) popFront() *wantConn {
	if q.headPos >= len(q.head) {
		if len(q.tail) == 0 {
			return nil
		}
		
		q.head, q.headPos, q.tail = q.tail, 0, q.head[:0]
	}
	w := q.head[q.headPos]
	q.head[q.headPos] = nil
	q.headPos++
	return w
}


func (q *wantConnQueue) peekFront() *wantConn {
	if q.headPos < len(q.head) {
		return q.head[q.headPos]
	}
	if len(q.tail) > 0 {
		return q.tail[0]
	}
	return nil
}


func (q *wantConnQueue) cleanFront() {
	for {
		w := q.peekFront()
		if w == nil || w.waiting() {
			return
		}
		q.popFront()
	}
}
