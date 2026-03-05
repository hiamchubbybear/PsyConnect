





package http2

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math"
	"math/bits"
	mathrand "math/rand"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/textproto"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/http/httpguts"
	"golang.org/x/net/http2/hpack"
	"golang.org/x/net/idna"
	"golang.org/x/net/internal/httpcommon"
)

const (
	
	
	transportDefaultConnFlow = 1 << 30

	
	
	
	transportDefaultStreamFlow = 4 << 20

	defaultUserAgent = "Go-http-client/2.0"

	
	
	
	initialMaxConcurrentStreams = 100

	
	
	defaultMaxConcurrentStreams = 1000
)





type Transport struct {
	
	
	
	
	
	
	
	DialTLSContext func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error)

	
	
	
	
	
	
	
	
	DialTLS func(network, addr string, cfg *tls.Config) (net.Conn, error)

	
	
	TLSClientConfig *tls.Config

	
	
	ConnPool ClientConnPool

	
	
	
	
	
	
	
	
	DisableCompression bool

	
	
	AllowHTTP bool

	
	
	
	
	
	
	
	MaxHeaderListSize uint32

	
	
	
	
	
	
	
	MaxReadFrameSize uint32

	
	
	
	
	
	MaxDecoderHeaderTableSize uint32

	
	
	
	
	MaxEncoderHeaderTableSize uint32

	
	
	
	
	
	
	
	
	StrictMaxConcurrentStreams bool

	
	
	
	
	IdleConnTimeout time.Duration

	
	
	
	
	
	
	ReadIdleTimeout time.Duration

	
	
	
	PingTimeout time.Duration

	
	
	
	WriteByteTimeout time.Duration

	
	
	
	
	CountError func(errType string)

	
	
	
	t1 *http.Transport

	connPoolOnce  sync.Once
	connPoolOrDef ClientConnPool 

	*transportTestHooks
}





type transportTestHooks struct {
	newclientconn func(*ClientConn)
	group         synctestGroupInterface
}

func (t *Transport) markNewGoroutine() {
	if t != nil && t.transportTestHooks != nil {
		t.transportTestHooks.group.Join()
	}
}

func (t *Transport) now() time.Time {
	if t != nil && t.transportTestHooks != nil {
		return t.transportTestHooks.group.Now()
	}
	return time.Now()
}

func (t *Transport) timeSince(when time.Time) time.Duration {
	if t != nil && t.transportTestHooks != nil {
		return t.now().Sub(when)
	}
	return time.Since(when)
}


func (t *Transport) newTimer(d time.Duration) timer {
	if t.transportTestHooks != nil {
		return t.transportTestHooks.group.NewTimer(d)
	}
	return timeTimer{time.NewTimer(d)}
}


func (t *Transport) afterFunc(d time.Duration, f func()) timer {
	if t.transportTestHooks != nil {
		return t.transportTestHooks.group.AfterFunc(d, f)
	}
	return timeTimer{time.AfterFunc(d, f)}
}

func (t *Transport) contextWithTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	if t.transportTestHooks != nil {
		return t.transportTestHooks.group.ContextWithTimeout(ctx, d)
	}
	return context.WithTimeout(ctx, d)
}

func (t *Transport) maxHeaderListSize() uint32 {
	n := int64(t.MaxHeaderListSize)
	if t.t1 != nil && t.t1.MaxResponseHeaderBytes != 0 {
		n = t.t1.MaxResponseHeaderBytes
		if n > 0 {
			n = adjustHTTP1MaxHeaderSize(n)
		}
	}
	if n <= 0 {
		return 10 << 20
	}
	if n >= 0xffffffff {
		return 0
	}
	return uint32(n)
}

func (t *Transport) disableCompression() bool {
	return t.DisableCompression || (t.t1 != nil && t.t1.DisableCompression)
}





func ConfigureTransport(t1 *http.Transport) error {
	_, err := ConfigureTransports(t1)
	return err
}




func ConfigureTransports(t1 *http.Transport) (*Transport, error) {
	return configureTransports(t1)
}

func configureTransports(t1 *http.Transport) (*Transport, error) {
	connPool := new(clientConnPool)
	t2 := &Transport{
		ConnPool: noDialClientConnPool{connPool},
		t1:       t1,
	}
	connPool.t = t2
	if err := registerHTTPSProtocol(t1, noDialH2RoundTripper{t2}); err != nil {
		return nil, err
	}
	if t1.TLSClientConfig == nil {
		t1.TLSClientConfig = new(tls.Config)
	}
	if !strSliceContains(t1.TLSClientConfig.NextProtos, "h2") {
		t1.TLSClientConfig.NextProtos = append([]string{"h2"}, t1.TLSClientConfig.NextProtos...)
	}
	if !strSliceContains(t1.TLSClientConfig.NextProtos, "http/1.1") {
		t1.TLSClientConfig.NextProtos = append(t1.TLSClientConfig.NextProtos, "http/1.1")
	}
	upgradeFn := func(scheme, authority string, c net.Conn) http.RoundTripper {
		addr := authorityAddr(scheme, authority)
		if used, err := connPool.addConnIfNeeded(addr, t2, c); err != nil {
			go c.Close()
			return erringRoundTripper{err}
		} else if !used {
			
			
			
			
			go c.Close()
		}
		if scheme == "http" {
			return (*unencryptedTransport)(t2)
		}
		return t2
	}
	if t1.TLSNextProto == nil {
		t1.TLSNextProto = make(map[string]func(string, *tls.Conn) http.RoundTripper)
	}
	t1.TLSNextProto[NextProtoTLS] = func(authority string, c *tls.Conn) http.RoundTripper {
		return upgradeFn("https", authority, c)
	}
	
	t1.TLSNextProto[nextProtoUnencryptedHTTP2] = func(authority string, c *tls.Conn) http.RoundTripper {
		nc, err := unencryptedNetConnFromTLSConn(c)
		if err != nil {
			go c.Close()
			return erringRoundTripper{err}
		}
		return upgradeFn("http", authority, nc)
	}
	return t2, nil
}



type unencryptedTransport Transport

func (t *unencryptedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return (*Transport)(t).RoundTripOpt(req, RoundTripOpt{allowHTTP: true})
}

func (t *Transport) connPool() ClientConnPool {
	t.connPoolOnce.Do(t.initConnPool)
	return t.connPoolOrDef
}

func (t *Transport) initConnPool() {
	if t.ConnPool != nil {
		t.connPoolOrDef = t.ConnPool
	} else {
		t.connPoolOrDef = &clientConnPool{t: t}
	}
}



type ClientConn struct {
	t             *Transport
	tconn         net.Conn             
	tlsState      *tls.ConnectionState 
	atomicReused  uint32               
	singleUse     bool                 
	getConnCalled bool                 

	
	readerDone chan struct{} 
	readerErr  error         

	idleTimeout time.Duration 
	idleTimer   timer

	mu               sync.Mutex 
	cond             *sync.Cond 
	flow             outflow    
	inflow           inflow     
	doNotReuse       bool       
	closing          bool
	closed           bool
	closedOnIdle     bool                     
	seenSettings     bool                     
	seenSettingsChan chan struct{}            
	wantSettingsAck  bool                     
	goAway           *GoAwayFrame             
	goAwayDebug      string                   
	streams          map[uint32]*clientStream 
	streamsReserved  int                      
	nextStreamID     uint32
	pendingRequests  int                       
	pings            map[[8]byte]chan struct{} 
	br               *bufio.Reader
	lastActive       time.Time
	lastIdle         time.Time 
	
	maxFrameSize                uint32
	maxConcurrentStreams        uint32
	peerMaxHeaderListSize       uint64
	peerMaxHeaderTableSize      uint32
	initialWindowSize           uint32
	initialStreamRecvWindowSize int32
	readIdleTimeout             time.Duration
	pingTimeout                 time.Duration
	extendedConnectAllowed      bool

	
	
	
	
	
	
	
	
	rstStreamPingsBlocked bool

	
	
	
	
	
	
	pendingResets int

	
	
	
	reqHeaderMu chan struct{}

	
	
	
	wmu  sync.Mutex
	bw   *bufio.Writer
	fr   *Framer
	werr error        
	hbuf bytes.Buffer 
	henc *hpack.Encoder
}



type clientStream struct {
	cc *ClientConn

	
	ctx       context.Context
	reqCancel <-chan struct{}

	trace         *httptrace.ClientTrace 
	ID            uint32
	bufPipe       pipe 
	requestedGzip bool
	isHead        bool

	abortOnce sync.Once
	abort     chan struct{} 
	abortErr  error         

	peerClosed chan struct{} 
	donec      chan struct{} 
	on100      chan struct{} 

	respHeaderRecv chan struct{}  
	res            *http.Response 

	flow        outflow 
	inflow      inflow  
	bytesRemain int64   
	readErr     error   

	reqBody              io.ReadCloser
	reqBodyContentLength int64         
	reqBodyClosed        chan struct{} 

	
	sentEndStream bool 
	sentHeaders   bool

	
	firstByte       bool  
	pastHeaders     bool  
	pastTrailers    bool  
	readClosed      bool  
	readAborted     bool  
	totalHeaderSize int64 

	trailer    http.Header  
	resTrailer *http.Header 
}

var got1xxFuncForTests func(int, textproto.MIMEHeader) error



func (cs *clientStream) get1xxTraceFunc() func(int, textproto.MIMEHeader) error {
	if fn := got1xxFuncForTests; fn != nil {
		return fn
	}
	return traceGot1xxResponseFunc(cs.trace)
}

func (cs *clientStream) abortStream(err error) {
	cs.cc.mu.Lock()
	defer cs.cc.mu.Unlock()
	cs.abortStreamLocked(err)
}

func (cs *clientStream) abortStreamLocked(err error) {
	cs.abortOnce.Do(func() {
		cs.abortErr = err
		close(cs.abort)
	})
	if cs.reqBody != nil {
		cs.closeReqBodyLocked()
	}
	
	if cs.cc.cond != nil {
		
		cs.cc.cond.Broadcast()
	}
}

func (cs *clientStream) abortRequestBodyWrite() {
	cc := cs.cc
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if cs.reqBody != nil && cs.reqBodyClosed == nil {
		cs.closeReqBodyLocked()
		cc.cond.Broadcast()
	}
}

func (cs *clientStream) closeReqBodyLocked() {
	if cs.reqBodyClosed != nil {
		return
	}
	cs.reqBodyClosed = make(chan struct{})
	reqBodyClosed := cs.reqBodyClosed
	go func() {
		cs.cc.t.markNewGoroutine()
		cs.reqBody.Close()
		close(reqBodyClosed)
	}()
}

type stickyErrWriter struct {
	group   synctestGroupInterface
	conn    net.Conn
	timeout time.Duration
	err     *error
}

func (sew stickyErrWriter) Write(p []byte) (n int, err error) {
	if *sew.err != nil {
		return 0, *sew.err
	}
	n, err = writeWithByteTimeout(sew.group, sew.conn, sew.timeout, p)
	*sew.err = err
	return n, err
}







type noCachedConnError struct{}

func (noCachedConnError) IsHTTP2NoCachedConnError() {}
func (noCachedConnError) Error() string             { return "http2: no cached connection was available" }




func isNoCachedConnError(err error) bool {
	_, ok := err.(interface{ IsHTTP2NoCachedConnError() })
	return ok
}

var ErrNoCachedConn error = noCachedConnError{}


type RoundTripOpt struct {
	
	
	
	
	OnlyCachedConn bool

	allowHTTP bool 
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.RoundTripOpt(req, RoundTripOpt{})
}



func authorityAddr(scheme string, authority string) (addr string) {
	host, port, err := net.SplitHostPort(authority)
	if err != nil { 
		host = authority
		port = ""
	}
	if port == "" { 
		port = "443"
		if scheme == "http" {
			port = "80"
		}
	}
	if a, err := idna.ToASCII(host); err == nil {
		host = a
	}
	
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		return host + ":" + port
	}
	return net.JoinHostPort(host, port)
}


func (t *Transport) RoundTripOpt(req *http.Request, opt RoundTripOpt) (*http.Response, error) {
	switch req.URL.Scheme {
	case "https":
		
	case "http":
		if !t.AllowHTTP && !opt.allowHTTP {
			return nil, errors.New("http2: unencrypted HTTP/2 not enabled")
		}
	default:
		return nil, errors.New("http2: unsupported scheme")
	}

	addr := authorityAddr(req.URL.Scheme, req.URL.Host)
	for retry := 0; ; retry++ {
		cc, err := t.connPool().GetClientConn(req, addr)
		if err != nil {
			t.vlogf("http2: Transport failed to get client conn for %s: %v", addr, err)
			return nil, err
		}
		reused := !atomic.CompareAndSwapUint32(&cc.atomicReused, 0, 1)
		traceGotConn(req, cc, reused)
		res, err := cc.RoundTrip(req)
		if err != nil && retry <= 6 {
			roundTripErr := err
			if req, err = shouldRetryRequest(req, err); err == nil {
				
				if retry == 0 {
					t.vlogf("RoundTrip retrying after failure: %v", roundTripErr)
					continue
				}
				backoff := float64(uint(1) << (uint(retry) - 1))
				backoff += backoff * (0.1 * mathrand.Float64())
				d := time.Second * time.Duration(backoff)
				tm := t.newTimer(d)
				select {
				case <-tm.C():
					t.vlogf("RoundTrip retrying after failure: %v", roundTripErr)
					continue
				case <-req.Context().Done():
					tm.Stop()
					err = req.Context().Err()
				}
			}
		}
		if err == errClientConnNotEstablished {
			
			
			
			
			
			
			
			
			
			
			if cc.idleTimer != nil {
				cc.idleTimer.Stop()
			}
			t.connPool().MarkDead(cc)
		}
		if err != nil {
			t.vlogf("RoundTrip failure: %v", err)
			return nil, err
		}
		return res, nil
	}
}




func (t *Transport) CloseIdleConnections() {
	if cp, ok := t.connPool().(clientConnPoolIdleCloser); ok {
		cp.closeIdleConnections()
	}
}

var (
	errClientConnClosed         = errors.New("http2: client conn is closed")
	errClientConnUnusable       = errors.New("http2: client conn not usable")
	errClientConnNotEstablished = errors.New("http2: client conn could not be established")
	errClientConnGotGoAway      = errors.New("http2: Transport received Server's graceful shutdown GOAWAY")
)





func shouldRetryRequest(req *http.Request, err error) (*http.Request, error) {
	if !canRetryError(err) {
		return nil, err
	}
	
	
	if req.Body == nil || req.Body == http.NoBody {
		return req, nil
	}

	
	
	if req.GetBody != nil {
		body, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		newReq := *req
		newReq.Body = body
		return &newReq, nil
	}

	
	
	
	if err == errClientConnUnusable {
		return req, nil
	}

	return nil, fmt.Errorf("http2: Transport: cannot retry err [%v] after Request.Body was written; define Request.GetBody to avoid this error", err)
}

func canRetryError(err error) bool {
	if err == errClientConnUnusable || err == errClientConnGotGoAway {
		return true
	}
	if se, ok := err.(StreamError); ok {
		if se.Code == ErrCodeProtocol && se.Cause == errFromPeer {
			
			return true
		}
		return se.Code == ErrCodeRefusedStream
	}
	return false
}

func (t *Transport) dialClientConn(ctx context.Context, addr string, singleUse bool) (*ClientConn, error) {
	if t.transportTestHooks != nil {
		return t.newClientConn(nil, singleUse)
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	tconn, err := t.dialTLS(ctx, "tcp", addr, t.newTLSConfig(host))
	if err != nil {
		return nil, err
	}
	return t.newClientConn(tconn, singleUse)
}

func (t *Transport) newTLSConfig(host string) *tls.Config {
	cfg := new(tls.Config)
	if t.TLSClientConfig != nil {
		*cfg = *t.TLSClientConfig.Clone()
	}
	if !strSliceContains(cfg.NextProtos, NextProtoTLS) {
		cfg.NextProtos = append([]string{NextProtoTLS}, cfg.NextProtos...)
	}
	if cfg.ServerName == "" {
		cfg.ServerName = host
	}
	return cfg
}

func (t *Transport) dialTLS(ctx context.Context, network, addr string, tlsCfg *tls.Config) (net.Conn, error) {
	if t.DialTLSContext != nil {
		return t.DialTLSContext(ctx, network, addr, tlsCfg)
	} else if t.DialTLS != nil {
		return t.DialTLS(network, addr, tlsCfg)
	}

	tlsCn, err := t.dialTLSWithContext(ctx, network, addr, tlsCfg)
	if err != nil {
		return nil, err
	}
	state := tlsCn.ConnectionState()
	if p := state.NegotiatedProtocol; p != NextProtoTLS {
		return nil, fmt.Errorf("http2: unexpected ALPN protocol %q; want %q", p, NextProtoTLS)
	}
	if !state.NegotiatedProtocolIsMutual {
		return nil, errors.New("http2: could not negotiate protocol mutually")
	}
	return tlsCn, nil
}



func (t *Transport) disableKeepAlives() bool {
	return t.t1 != nil && t.t1.DisableKeepAlives
}

func (t *Transport) expectContinueTimeout() time.Duration {
	if t.t1 == nil {
		return 0
	}
	return t.t1.ExpectContinueTimeout
}

func (t *Transport) NewClientConn(c net.Conn) (*ClientConn, error) {
	return t.newClientConn(c, t.disableKeepAlives())
}

func (t *Transport) newClientConn(c net.Conn, singleUse bool) (*ClientConn, error) {
	conf := configFromTransport(t)
	cc := &ClientConn{
		t:                           t,
		tconn:                       c,
		readerDone:                  make(chan struct{}),
		nextStreamID:                1,
		maxFrameSize:                16 << 10, 
		initialWindowSize:           65535,    
		initialStreamRecvWindowSize: conf.MaxUploadBufferPerStream,
		maxConcurrentStreams:        initialMaxConcurrentStreams, 
		peerMaxHeaderListSize:       0xffffffffffffffff,          
		streams:                     make(map[uint32]*clientStream),
		singleUse:                   singleUse,
		seenSettingsChan:            make(chan struct{}),
		wantSettingsAck:             true,
		readIdleTimeout:             conf.SendPingTimeout,
		pingTimeout:                 conf.PingTimeout,
		pings:                       make(map[[8]byte]chan struct{}),
		reqHeaderMu:                 make(chan struct{}, 1),
		lastActive:                  t.now(),
	}
	var group synctestGroupInterface
	if t.transportTestHooks != nil {
		t.markNewGoroutine()
		t.transportTestHooks.newclientconn(cc)
		c = cc.tconn
		group = t.group
	}
	if VerboseLogs {
		t.vlogf("http2: Transport creating client conn %p to %v", cc, c.RemoteAddr())
	}

	cc.cond = sync.NewCond(&cc.mu)
	cc.flow.add(int32(initialWindowSize))

	
	
	cc.bw = bufio.NewWriter(stickyErrWriter{
		group:   group,
		conn:    c,
		timeout: conf.WriteByteTimeout,
		err:     &cc.werr,
	})
	cc.br = bufio.NewReader(c)
	cc.fr = NewFramer(cc.bw, cc.br)
	cc.fr.SetMaxReadFrameSize(conf.MaxReadFrameSize)
	if t.CountError != nil {
		cc.fr.countError = t.CountError
	}
	maxHeaderTableSize := conf.MaxDecoderHeaderTableSize
	cc.fr.ReadMetaHeaders = hpack.NewDecoder(maxHeaderTableSize, nil)
	cc.fr.MaxHeaderListSize = t.maxHeaderListSize()

	cc.henc = hpack.NewEncoder(&cc.hbuf)
	cc.henc.SetMaxDynamicTableSizeLimit(conf.MaxEncoderHeaderTableSize)
	cc.peerMaxHeaderTableSize = initialHeaderTableSize

	if cs, ok := c.(connectionStater); ok {
		state := cs.ConnectionState()
		cc.tlsState = &state
	}

	initialSettings := []Setting{
		{ID: SettingEnablePush, Val: 0},
		{ID: SettingInitialWindowSize, Val: uint32(cc.initialStreamRecvWindowSize)},
	}
	initialSettings = append(initialSettings, Setting{ID: SettingMaxFrameSize, Val: conf.MaxReadFrameSize})
	if max := t.maxHeaderListSize(); max != 0 {
		initialSettings = append(initialSettings, Setting{ID: SettingMaxHeaderListSize, Val: max})
	}
	if maxHeaderTableSize != initialHeaderTableSize {
		initialSettings = append(initialSettings, Setting{ID: SettingHeaderTableSize, Val: maxHeaderTableSize})
	}

	cc.bw.Write(clientPreface)
	cc.fr.WriteSettings(initialSettings...)
	cc.fr.WriteWindowUpdate(0, uint32(conf.MaxUploadBufferPerConnection))
	cc.inflow.init(conf.MaxUploadBufferPerConnection + initialWindowSize)
	cc.bw.Flush()
	if cc.werr != nil {
		cc.Close()
		return nil, cc.werr
	}

	
	if d := t.idleConnTimeout(); d != 0 {
		cc.idleTimeout = d
		cc.idleTimer = t.afterFunc(d, cc.onIdleTimeout)
	}

	go cc.readLoop()
	return cc, nil
}

func (cc *ClientConn) healthCheck() {
	pingTimeout := cc.pingTimeout
	
	
	ctx, cancel := cc.t.contextWithTimeout(context.Background(), pingTimeout)
	defer cancel()
	cc.vlogf("http2: Transport sending health check")
	err := cc.Ping(ctx)
	if err != nil {
		cc.vlogf("http2: Transport health check failure: %v", err)
		cc.closeForLostPing()
	} else {
		cc.vlogf("http2: Transport health check success")
	}
}


func (cc *ClientConn) SetDoNotReuse() {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.doNotReuse = true
}

func (cc *ClientConn) setGoAway(f *GoAwayFrame) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	old := cc.goAway
	cc.goAway = f

	
	if cc.goAwayDebug == "" {
		cc.goAwayDebug = string(f.DebugData())
	}
	if old != nil && old.ErrCode != ErrCodeNo {
		cc.goAway.ErrCode = old.ErrCode
	}
	last := f.LastStreamID
	for streamID, cs := range cc.streams {
		if streamID <= last {
			
			
			
			continue
		}
		if streamID == 1 && cc.goAway.ErrCode != ErrCodeNo {
			
			
			
			cs.abortStreamLocked(fmt.Errorf("http2: Transport received GOAWAY from server ErrCode:%v", cc.goAway.ErrCode))
		} else {
			
			
			cs.abortStreamLocked(errClientConnGotGoAway)
		}
	}
}






func (cc *ClientConn) CanTakeNewRequest() bool {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc.canTakeNewRequestLocked()
}




func (cc *ClientConn) ReserveNewRequest() bool {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if st := cc.idleStateLocked(); !st.canTakeNewRequest {
		return false
	}
	cc.streamsReserved++
	return true
}


type ClientConnState struct {
	
	Closed bool

	
	
	
	
	Closing bool

	
	StreamsActive int

	
	
	StreamsReserved int

	
	
	
	StreamsPending int

	
	
	
	MaxConcurrentStreams uint32

	
	
	LastIdle time.Time
}


func (cc *ClientConn) State() ClientConnState {
	cc.wmu.Lock()
	maxConcurrent := cc.maxConcurrentStreams
	if !cc.seenSettings {
		maxConcurrent = 0
	}
	cc.wmu.Unlock()

	cc.mu.Lock()
	defer cc.mu.Unlock()
	return ClientConnState{
		Closed:               cc.closed,
		Closing:              cc.closing || cc.singleUse || cc.doNotReuse || cc.goAway != nil,
		StreamsActive:        len(cc.streams) + cc.pendingResets,
		StreamsReserved:      cc.streamsReserved,
		StreamsPending:       cc.pendingRequests,
		LastIdle:             cc.lastIdle,
		MaxConcurrentStreams: maxConcurrent,
	}
}



type clientConnIdleState struct {
	canTakeNewRequest bool
}

func (cc *ClientConn) idleState() clientConnIdleState {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc.idleStateLocked()
}

func (cc *ClientConn) idleStateLocked() (st clientConnIdleState) {
	if cc.singleUse && cc.nextStreamID > 1 {
		return
	}
	var maxConcurrentOkay bool
	if cc.t.StrictMaxConcurrentStreams {
		
		
		
		
		maxConcurrentOkay = true
	} else {
		
		
		
		
		
		
		maxConcurrentOkay = cc.currentRequestCountLocked() < int(cc.maxConcurrentStreams)
	}

	st.canTakeNewRequest = cc.goAway == nil && !cc.closed && !cc.closing && maxConcurrentOkay &&
		!cc.doNotReuse &&
		int64(cc.nextStreamID)+2*int64(cc.pendingRequests) < math.MaxInt32 &&
		!cc.tooIdleLocked()

	
	
	
	
	
	
	
	if cc.nextStreamID == 1 && cc.streamsReserved == 0 && cc.closed && !cc.closedOnIdle {
		st.canTakeNewRequest = true
	}

	return
}



func (cc *ClientConn) currentRequestCountLocked() int {
	return len(cc.streams) + cc.streamsReserved + cc.pendingResets
}

func (cc *ClientConn) canTakeNewRequestLocked() bool {
	st := cc.idleStateLocked()
	return st.canTakeNewRequest
}



func (cc *ClientConn) tooIdleLocked() bool {
	
	
	
	
	return cc.idleTimeout != 0 && !cc.lastIdle.IsZero() && cc.t.timeSince(cc.lastIdle.Round(0)) > cc.idleTimeout
}







func (cc *ClientConn) onIdleTimeout() {
	cc.closeIfIdle()
}

func (cc *ClientConn) closeConn() {
	t := time.AfterFunc(250*time.Millisecond, cc.forceCloseConn)
	defer t.Stop()
	cc.tconn.Close()
}



func (cc *ClientConn) forceCloseConn() {
	tc, ok := cc.tconn.(*tls.Conn)
	if !ok {
		return
	}
	if nc := tc.NetConn(); nc != nil {
		nc.Close()
	}
}

func (cc *ClientConn) closeIfIdle() {
	cc.mu.Lock()
	if len(cc.streams) > 0 || cc.streamsReserved > 0 {
		cc.mu.Unlock()
		return
	}
	cc.closed = true
	cc.closedOnIdle = true
	nextID := cc.nextStreamID
	
	cc.mu.Unlock()

	if VerboseLogs {
		cc.vlogf("http2: Transport closing idle conn %p (forSingleUse=%v, maxStream=%v)", cc, cc.singleUse, nextID-2)
	}
	cc.closeConn()
}

func (cc *ClientConn) isDoNotReuseAndIdle() bool {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc.doNotReuse && len(cc.streams) == 0
}

var shutdownEnterWaitStateHook = func() {}


func (cc *ClientConn) Shutdown(ctx context.Context) error {
	if err := cc.sendGoAway(); err != nil {
		return err
	}
	
	done := make(chan struct{})
	cancelled := false 
	go func() {
		cc.t.markNewGoroutine()
		cc.mu.Lock()
		defer cc.mu.Unlock()
		for {
			if len(cc.streams) == 0 || cc.closed {
				cc.closed = true
				close(done)
				break
			}
			if cancelled {
				break
			}
			cc.cond.Wait()
		}
	}()
	shutdownEnterWaitStateHook()
	select {
	case <-done:
		cc.closeConn()
		return nil
	case <-ctx.Done():
		cc.mu.Lock()
		
		cancelled = true
		cc.cond.Broadcast()
		cc.mu.Unlock()
		return ctx.Err()
	}
}

func (cc *ClientConn) sendGoAway() error {
	cc.mu.Lock()
	closing := cc.closing
	cc.closing = true
	maxStreamID := cc.nextStreamID
	cc.mu.Unlock()
	if closing {
		
		return nil
	}

	cc.wmu.Lock()
	defer cc.wmu.Unlock()
	
	if err := cc.fr.WriteGoAway(maxStreamID, ErrCodeNo, nil); err != nil {
		return err
	}
	if err := cc.bw.Flush(); err != nil {
		return err
	}
	
	return nil
}



func (cc *ClientConn) closeForError(err error) {
	cc.mu.Lock()
	cc.closed = true
	for _, cs := range cc.streams {
		cs.abortStreamLocked(err)
	}
	cc.cond.Broadcast()
	cc.mu.Unlock()
	cc.closeConn()
}




func (cc *ClientConn) Close() error {
	err := errors.New("http2: client connection force closed via ClientConn.Close")
	cc.closeForError(err)
	return nil
}


func (cc *ClientConn) closeForLostPing() {
	err := errors.New("http2: client connection lost")
	if f := cc.t.CountError; f != nil {
		f("conn_close_lost_ping")
	}
	cc.closeForError(err)
}



var errRequestCanceled = errors.New("net/http: request canceled")

func (cc *ClientConn) responseHeaderTimeout() time.Duration {
	if cc.t.t1 != nil {
		return cc.t.t1.ResponseHeaderTimeout
	}
	
	
	
	
	return 0
}




func actualContentLength(req *http.Request) int64 {
	if req.Body == nil || req.Body == http.NoBody {
		return 0
	}
	if req.ContentLength != 0 {
		return req.ContentLength
	}
	return -1
}

func (cc *ClientConn) decrStreamReservations() {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.decrStreamReservationsLocked()
}

func (cc *ClientConn) decrStreamReservationsLocked() {
	if cc.streamsReserved > 0 {
		cc.streamsReserved--
	}
}

func (cc *ClientConn) RoundTrip(req *http.Request) (*http.Response, error) {
	return cc.roundTrip(req, nil)
}

func (cc *ClientConn) roundTrip(req *http.Request, streamf func(*clientStream)) (*http.Response, error) {
	ctx := req.Context()
	cs := &clientStream{
		cc:                   cc,
		ctx:                  ctx,
		reqCancel:            req.Cancel,
		isHead:               req.Method == "HEAD",
		reqBody:              req.Body,
		reqBodyContentLength: actualContentLength(req),
		trace:                httptrace.ContextClientTrace(ctx),
		peerClosed:           make(chan struct{}),
		abort:                make(chan struct{}),
		respHeaderRecv:       make(chan struct{}),
		donec:                make(chan struct{}),
	}

	cs.requestedGzip = httpcommon.IsRequestGzip(req.Method, req.Header, cc.t.disableCompression())

	go cs.doRequest(req, streamf)

	waitDone := func() error {
		select {
		case <-cs.donec:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case <-cs.reqCancel:
			return errRequestCanceled
		}
	}

	handleResponseHeaders := func() (*http.Response, error) {
		res := cs.res
		if res.StatusCode > 299 {
			
			
			
			
			
			
			
			
			
			cs.abortRequestBodyWrite()
		}
		res.Request = req
		res.TLS = cc.tlsState
		if res.Body == noBody && actualContentLength(req) == 0 {
			
			
			
			if err := waitDone(); err != nil {
				return nil, err
			}
		}
		return res, nil
	}

	cancelRequest := func(cs *clientStream, err error) error {
		cs.cc.mu.Lock()
		bodyClosed := cs.reqBodyClosed
		cs.cc.mu.Unlock()
		
		
		
		
		
		
		
		
		
		
		
		
		
		if bodyClosed != nil {
			<-bodyClosed
		}
		return err
	}

	for {
		select {
		case <-cs.respHeaderRecv:
			return handleResponseHeaders()
		case <-cs.abort:
			select {
			case <-cs.respHeaderRecv:
				
				
				
				
				return handleResponseHeaders()
			default:
				waitDone()
				return nil, cs.abortErr
			}
		case <-ctx.Done():
			err := ctx.Err()
			cs.abortStream(err)
			return nil, cancelRequest(cs, err)
		case <-cs.reqCancel:
			cs.abortStream(errRequestCanceled)
			return nil, cancelRequest(cs, errRequestCanceled)
		}
	}
}




func (cs *clientStream) doRequest(req *http.Request, streamf func(*clientStream)) {
	cs.cc.t.markNewGoroutine()
	err := cs.writeRequest(req, streamf)
	cs.cleanupWriteRequest(err)
}

var errExtendedConnectNotSupported = errors.New("net/http: extended connect not supported by peer")








func (cs *clientStream) writeRequest(req *http.Request, streamf func(*clientStream)) (err error) {
	cc := cs.cc
	ctx := cs.ctx

	
	
	var isExtendedConnect bool
	if req.Method == "CONNECT" && req.Header.Get(":protocol") != "" {
		isExtendedConnect = true
	}

	
	
	
	if cc.reqHeaderMu == nil {
		panic("RoundTrip on uninitialized ClientConn") 
	}
	if isExtendedConnect {
		select {
		case <-cs.reqCancel:
			return errRequestCanceled
		case <-ctx.Done():
			return ctx.Err()
		case <-cc.seenSettingsChan:
			if !cc.extendedConnectAllowed {
				return errExtendedConnectNotSupported
			}
		}
	}
	select {
	case cc.reqHeaderMu <- struct{}{}:
	case <-cs.reqCancel:
		return errRequestCanceled
	case <-ctx.Done():
		return ctx.Err()
	}

	cc.mu.Lock()
	if cc.idleTimer != nil {
		cc.idleTimer.Stop()
	}
	cc.decrStreamReservationsLocked()
	if err := cc.awaitOpenSlotForStreamLocked(cs); err != nil {
		cc.mu.Unlock()
		<-cc.reqHeaderMu
		return err
	}
	cc.addStreamLocked(cs) 
	if isConnectionCloseRequest(req) {
		cc.doNotReuse = true
	}
	cc.mu.Unlock()

	if streamf != nil {
		streamf(cs)
	}

	continueTimeout := cc.t.expectContinueTimeout()
	if continueTimeout != 0 {
		if !httpguts.HeaderValuesContainsToken(req.Header["Expect"], "100-continue") {
			continueTimeout = 0
		} else {
			cs.on100 = make(chan struct{}, 1)
		}
	}

	
	
	
	
	err = cs.encodeAndWriteHeaders(req)
	<-cc.reqHeaderMu
	if err != nil {
		return err
	}

	hasBody := cs.reqBodyContentLength != 0
	if !hasBody {
		cs.sentEndStream = true
	} else {
		if continueTimeout != 0 {
			traceWait100Continue(cs.trace)
			timer := time.NewTimer(continueTimeout)
			select {
			case <-timer.C:
				err = nil
			case <-cs.on100:
				err = nil
			case <-cs.abort:
				err = cs.abortErr
			case <-ctx.Done():
				err = ctx.Err()
			case <-cs.reqCancel:
				err = errRequestCanceled
			}
			timer.Stop()
			if err != nil {
				traceWroteRequest(cs.trace, err)
				return err
			}
		}

		if err = cs.writeRequestBody(req); err != nil {
			if err != errStopReqBodyWrite {
				traceWroteRequest(cs.trace, err)
				return err
			}
		} else {
			cs.sentEndStream = true
		}
	}

	traceWroteRequest(cs.trace, err)

	var respHeaderTimer <-chan time.Time
	var respHeaderRecv chan struct{}
	if d := cc.responseHeaderTimeout(); d != 0 {
		timer := cc.t.newTimer(d)
		defer timer.Stop()
		respHeaderTimer = timer.C()
		respHeaderRecv = cs.respHeaderRecv
	}
	
	
	
	for {
		select {
		case <-cs.peerClosed:
			return nil
		case <-respHeaderTimer:
			return errTimeout
		case <-respHeaderRecv:
			respHeaderRecv = nil
			respHeaderTimer = nil 
		case <-cs.abort:
			return cs.abortErr
		case <-ctx.Done():
			return ctx.Err()
		case <-cs.reqCancel:
			return errRequestCanceled
		}
	}
}

func (cs *clientStream) encodeAndWriteHeaders(req *http.Request) error {
	cc := cs.cc
	ctx := cs.ctx

	cc.wmu.Lock()
	defer cc.wmu.Unlock()

	
	select {
	case <-cs.abort:
		return cs.abortErr
	case <-ctx.Done():
		return ctx.Err()
	case <-cs.reqCancel:
		return errRequestCanceled
	default:
	}

	
	
	
	
	
	cc.hbuf.Reset()
	res, err := encodeRequestHeaders(req, cs.requestedGzip, cc.peerMaxHeaderListSize, func(name, value string) {
		cc.writeHeader(name, value)
	})
	if err != nil {
		return fmt.Errorf("http2: %w", err)
	}
	hdrs := cc.hbuf.Bytes()

	
	endStream := !res.HasBody && !res.HasTrailers
	cs.sentHeaders = true
	err = cc.writeHeaders(cs.ID, endStream, int(cc.maxFrameSize), hdrs)
	traceWroteHeaders(cs.trace)
	return err
}

func encodeRequestHeaders(req *http.Request, addGzipHeader bool, peerMaxHeaderListSize uint64, headerf func(name, value string)) (httpcommon.EncodeHeadersResult, error) {
	return httpcommon.EncodeHeaders(req.Context(), httpcommon.EncodeHeadersParam{
		Request: httpcommon.Request{
			Header:              req.Header,
			Trailer:             req.Trailer,
			URL:                 req.URL,
			Host:                req.Host,
			Method:              req.Method,
			ActualContentLength: actualContentLength(req),
		},
		AddGzipHeader:         addGzipHeader,
		PeerMaxHeaderListSize: peerMaxHeaderListSize,
		DefaultUserAgent:      defaultUserAgent,
	}, headerf)
}





func (cs *clientStream) cleanupWriteRequest(err error) {
	cc := cs.cc

	if cs.ID == 0 {
		
		cc.decrStreamReservations()
	}

	
	
	
	
	cc.mu.Lock()
	mustCloseBody := false
	if cs.reqBody != nil && cs.reqBodyClosed == nil {
		mustCloseBody = true
		cs.reqBodyClosed = make(chan struct{})
	}
	bodyClosed := cs.reqBodyClosed
	closeOnIdle := cc.singleUse || cc.doNotReuse || cc.t.disableKeepAlives() || cc.goAway != nil
	cc.mu.Unlock()
	if mustCloseBody {
		cs.reqBody.Close()
		close(bodyClosed)
	}
	if bodyClosed != nil {
		<-bodyClosed
	}

	if err != nil && cs.sentEndStream {
		
		
		
		select {
		case <-cs.peerClosed:
			err = nil
		default:
		}
	}
	if err != nil {
		cs.abortStream(err) 
		if cs.sentHeaders {
			if se, ok := err.(StreamError); ok {
				if se.Cause != errFromPeer {
					cc.writeStreamReset(cs.ID, se.Code, false, err)
				}
			} else {
				
				
				
				
				
				
				
				
				
				
				
				
				
				
				
				ping := false
				if !closeOnIdle {
					cc.mu.Lock()
					
					
					if !cc.rstStreamPingsBlocked {
						if cc.pendingResets == 0 {
							ping = true
						}
						cc.pendingResets++
					}
					cc.mu.Unlock()
				}
				cc.writeStreamReset(cs.ID, ErrCodeCancel, ping, err)
			}
		}
		cs.bufPipe.CloseWithError(err) 
	} else {
		if cs.sentHeaders && !cs.sentEndStream {
			cc.writeStreamReset(cs.ID, ErrCodeNo, false, nil)
		}
		cs.bufPipe.CloseWithError(errRequestCanceled)
	}
	if cs.ID != 0 {
		cc.forgetStreamID(cs.ID)
	}

	cc.wmu.Lock()
	werr := cc.werr
	cc.wmu.Unlock()
	if werr != nil {
		cc.Close()
	}

	close(cs.donec)
}



func (cc *ClientConn) awaitOpenSlotForStreamLocked(cs *clientStream) error {
	for {
		if cc.closed && cc.nextStreamID == 1 && cc.streamsReserved == 0 {
			
			
			return errClientConnNotEstablished
		}
		cc.lastActive = cc.t.now()
		if cc.closed || !cc.canTakeNewRequestLocked() {
			return errClientConnUnusable
		}
		cc.lastIdle = time.Time{}
		if cc.currentRequestCountLocked() < int(cc.maxConcurrentStreams) {
			return nil
		}
		cc.pendingRequests++
		cc.cond.Wait()
		cc.pendingRequests--
		select {
		case <-cs.abort:
			return cs.abortErr
		default:
		}
	}
}


func (cc *ClientConn) writeHeaders(streamID uint32, endStream bool, maxFrameSize int, hdrs []byte) error {
	first := true 
	for len(hdrs) > 0 && cc.werr == nil {
		chunk := hdrs
		if len(chunk) > maxFrameSize {
			chunk = chunk[:maxFrameSize]
		}
		hdrs = hdrs[len(chunk):]
		endHeaders := len(hdrs) == 0
		if first {
			cc.fr.WriteHeaders(HeadersFrameParam{
				StreamID:      streamID,
				BlockFragment: chunk,
				EndStream:     endStream,
				EndHeaders:    endHeaders,
			})
			first = false
		} else {
			cc.fr.WriteContinuation(streamID, endHeaders, chunk)
		}
	}
	cc.bw.Flush()
	return cc.werr
}


var (
	
	errStopReqBodyWrite = errors.New("http2: aborting request body write")

	
	errStopReqBodyWriteAndCancel = errors.New("http2: canceling request")

	errReqBodyTooLong = errors.New("http2: request body larger than specified content length")
)






func (cs *clientStream) frameScratchBufferLen(maxFrameSize int) int {
	const max = 512 << 10
	n := int64(maxFrameSize)
	if n > max {
		n = max
	}
	if cl := cs.reqBodyContentLength; cl != -1 && cl+1 < n {
		
		
		
		
		n = cl + 1
	}
	if n < 1 {
		return 1
	}
	return int(n) 
}









var bufPools [7]sync.Pool 
func bufPoolIndex(size int) int {
	if size <= 16384 {
		return 0
	}
	size -= 1
	bits := bits.Len(uint(size))
	index := bits - 14
	if index >= len(bufPools) {
		return len(bufPools) - 1
	}
	return index
}

func (cs *clientStream) writeRequestBody(req *http.Request) (err error) {
	cc := cs.cc
	body := cs.reqBody
	sentEnd := false 

	hasTrailers := req.Trailer != nil
	remainLen := cs.reqBodyContentLength
	hasContentLen := remainLen != -1

	cc.mu.Lock()
	maxFrameSize := int(cc.maxFrameSize)
	cc.mu.Unlock()

	
	scratchLen := cs.frameScratchBufferLen(maxFrameSize)
	var buf []byte
	index := bufPoolIndex(scratchLen)
	if bp, ok := bufPools[index].Get().(*[]byte); ok && len(*bp) >= scratchLen {
		defer bufPools[index].Put(bp)
		buf = *bp
	} else {
		buf = make([]byte, scratchLen)
		defer bufPools[index].Put(&buf)
	}

	var sawEOF bool
	for !sawEOF {
		n, err := body.Read(buf)
		if hasContentLen {
			remainLen -= int64(n)
			if remainLen == 0 && err == nil {
				
				
				
				
				
				
				
				var scratch [1]byte
				var n1 int
				n1, err = body.Read(scratch[:])
				remainLen -= int64(n1)
			}
			if remainLen < 0 {
				err = errReqBodyTooLong
				return err
			}
		}
		if err != nil {
			cc.mu.Lock()
			bodyClosed := cs.reqBodyClosed != nil
			cc.mu.Unlock()
			switch {
			case bodyClosed:
				return errStopReqBodyWrite
			case err == io.EOF:
				sawEOF = true
				err = nil
			default:
				return err
			}
		}

		remain := buf[:n]
		for len(remain) > 0 && err == nil {
			var allowed int32
			allowed, err = cs.awaitFlowControl(len(remain))
			if err != nil {
				return err
			}
			cc.wmu.Lock()
			data := remain[:allowed]
			remain = remain[allowed:]
			sentEnd = sawEOF && len(remain) == 0 && !hasTrailers
			err = cc.fr.WriteData(cs.ID, sentEnd, data)
			if err == nil {
				
				
				
				
				
				
				err = cc.bw.Flush()
			}
			cc.wmu.Unlock()
		}
		if err != nil {
			return err
		}
	}

	if sentEnd {
		
		
		
		return nil
	}

	
	
	
	cc.mu.Lock()
	trailer := req.Trailer
	err = cs.abortErr
	cc.mu.Unlock()
	if err != nil {
		return err
	}

	cc.wmu.Lock()
	defer cc.wmu.Unlock()
	var trls []byte
	if len(trailer) > 0 {
		trls, err = cc.encodeTrailers(trailer)
		if err != nil {
			return err
		}
	}

	
	
	if len(trls) > 0 {
		err = cc.writeHeaders(cs.ID, true, maxFrameSize, trls)
	} else {
		err = cc.fr.WriteData(cs.ID, true, nil)
	}
	if ferr := cc.bw.Flush(); ferr != nil && err == nil {
		err = ferr
	}
	return err
}





func (cs *clientStream) awaitFlowControl(maxBytes int) (taken int32, err error) {
	cc := cs.cc
	ctx := cs.ctx
	cc.mu.Lock()
	defer cc.mu.Unlock()
	for {
		if cc.closed {
			return 0, errClientConnClosed
		}
		if cs.reqBodyClosed != nil {
			return 0, errStopReqBodyWrite
		}
		select {
		case <-cs.abort:
			return 0, cs.abortErr
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-cs.reqCancel:
			return 0, errRequestCanceled
		default:
		}
		if a := cs.flow.available(); a > 0 {
			take := a
			if int(take) > maxBytes {

				take = int32(maxBytes) 
			}
			if take > int32(cc.maxFrameSize) {
				take = int32(cc.maxFrameSize)
			}
			cs.flow.take(take)
			return take, nil
		}
		cc.cond.Wait()
	}
}


func (cc *ClientConn) encodeTrailers(trailer http.Header) ([]byte, error) {
	cc.hbuf.Reset()

	hlSize := uint64(0)
	for k, vv := range trailer {
		for _, v := range vv {
			hf := hpack.HeaderField{Name: k, Value: v}
			hlSize += uint64(hf.Size())
		}
	}
	if hlSize > cc.peerMaxHeaderListSize {
		return nil, errRequestHeaderListSize
	}

	for k, vv := range trailer {
		lowKey, ascii := httpcommon.LowerHeader(k)
		if !ascii {
			
			
			continue
		}
		
		
		for _, v := range vv {
			cc.writeHeader(lowKey, v)
		}
	}
	return cc.hbuf.Bytes(), nil
}

func (cc *ClientConn) writeHeader(name, value string) {
	if VerboseLogs {
		log.Printf("http2: Transport encoding header %q = %q", name, value)
	}
	cc.henc.WriteField(hpack.HeaderField{Name: name, Value: value})
}

type resAndError struct {
	_   incomparable
	res *http.Response
	err error
}


func (cc *ClientConn) addStreamLocked(cs *clientStream) {
	cs.flow.add(int32(cc.initialWindowSize))
	cs.flow.setConnFlow(&cc.flow)
	cs.inflow.init(cc.initialStreamRecvWindowSize)
	cs.ID = cc.nextStreamID
	cc.nextStreamID += 2
	cc.streams[cs.ID] = cs
	if cs.ID == 0 {
		panic("assigned stream ID 0")
	}
}

func (cc *ClientConn) forgetStreamID(id uint32) {
	cc.mu.Lock()
	slen := len(cc.streams)
	delete(cc.streams, id)
	if len(cc.streams) != slen-1 {
		panic("forgetting unknown stream id")
	}
	cc.lastActive = cc.t.now()
	if len(cc.streams) == 0 && cc.idleTimer != nil {
		cc.idleTimer.Reset(cc.idleTimeout)
		cc.lastIdle = cc.t.now()
	}
	
	
	cc.cond.Broadcast()

	closeOnIdle := cc.singleUse || cc.doNotReuse || cc.t.disableKeepAlives() || cc.goAway != nil
	if closeOnIdle && cc.streamsReserved == 0 && len(cc.streams) == 0 {
		if VerboseLogs {
			cc.vlogf("http2: Transport closing idle conn %p (forSingleUse=%v, maxStream=%v)", cc, cc.singleUse, cc.nextStreamID-2)
		}
		cc.closed = true
		defer cc.closeConn()
	}

	cc.mu.Unlock()
}


type clientConnReadLoop struct {
	_  incomparable
	cc *ClientConn
}


func (cc *ClientConn) readLoop() {
	cc.t.markNewGoroutine()
	rl := &clientConnReadLoop{cc: cc}
	defer rl.cleanup()
	cc.readerErr = rl.run()
	if ce, ok := cc.readerErr.(ConnectionError); ok {
		cc.wmu.Lock()
		cc.fr.WriteGoAway(0, ErrCode(ce), nil)
		cc.wmu.Unlock()
	}
}



type GoAwayError struct {
	LastStreamID uint32
	ErrCode      ErrCode
	DebugData    string
}

func (e GoAwayError) Error() string {
	return fmt.Sprintf("http2: server sent GOAWAY and closed the connection; LastStreamID=%v, ErrCode=%v, debug=%q",
		e.LastStreamID, e.ErrCode, e.DebugData)
}

func isEOFOrNetReadError(err error) bool {
	if err == io.EOF {
		return true
	}
	ne, ok := err.(*net.OpError)
	return ok && ne.Op == "read"
}

func (rl *clientConnReadLoop) cleanup() {
	cc := rl.cc
	defer cc.closeConn()
	defer close(cc.readerDone)

	if cc.idleTimer != nil {
		cc.idleTimer.Stop()
	}

	
	
	
	err := cc.readerErr
	cc.mu.Lock()
	if cc.goAway != nil && isEOFOrNetReadError(err) {
		err = GoAwayError{
			LastStreamID: cc.goAway.LastStreamID,
			ErrCode:      cc.goAway.ErrCode,
			DebugData:    cc.goAwayDebug,
		}
	} else if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	cc.closed = true

	
	
	
	
	
	
	unusedWaitTime := 5 * time.Second
	if cc.idleTimeout > 0 && unusedWaitTime > cc.idleTimeout {
		unusedWaitTime = cc.idleTimeout
	}
	idleTime := cc.t.now().Sub(cc.lastActive)
	if atomic.LoadUint32(&cc.atomicReused) == 0 && idleTime < unusedWaitTime && !cc.closedOnIdle {
		cc.idleTimer = cc.t.afterFunc(unusedWaitTime-idleTime, func() {
			cc.t.connPool().MarkDead(cc)
		})
	} else {
		cc.mu.Unlock() 
		cc.t.connPool().MarkDead(cc)
		cc.mu.Lock()
	}

	for _, cs := range cc.streams {
		select {
		case <-cs.peerClosed:
			
			
		default:
			cs.abortStreamLocked(err)
		}
	}
	cc.cond.Broadcast()
	cc.mu.Unlock()

	if !cc.seenSettings {
		
		
		cc.extendedConnectAllowed = true
		close(cc.seenSettingsChan)
	}
}



func (cc *ClientConn) countReadFrameError(err error) {
	f := cc.t.CountError
	if f == nil || err == nil {
		return
	}
	if ce, ok := err.(ConnectionError); ok {
		errCode := ErrCode(ce)
		f(fmt.Sprintf("read_frame_conn_error_%s", errCode.stringToken()))
		return
	}
	if errors.Is(err, io.EOF) {
		f("read_frame_eof")
		return
	}
	if errors.Is(err, io.ErrUnexpectedEOF) {
		f("read_frame_unexpected_eof")
		return
	}
	if errors.Is(err, ErrFrameTooLarge) {
		f("read_frame_too_large")
		return
	}
	f("read_frame_other")
}

func (rl *clientConnReadLoop) run() error {
	cc := rl.cc
	gotSettings := false
	readIdleTimeout := cc.readIdleTimeout
	var t timer
	if readIdleTimeout != 0 {
		t = cc.t.afterFunc(readIdleTimeout, cc.healthCheck)
	}
	for {
		f, err := cc.fr.ReadFrame()
		if t != nil {
			t.Reset(readIdleTimeout)
		}
		if err != nil {
			cc.vlogf("http2: Transport readFrame error on conn %p: (%T) %v", cc, err, err)
		}
		if se, ok := err.(StreamError); ok {
			if cs := rl.streamByID(se.StreamID, notHeaderOrDataFrame); cs != nil {
				if se.Cause == nil {
					se.Cause = cc.fr.errDetail
				}
				rl.endStreamError(cs, se)
			}
			continue
		} else if err != nil {
			cc.countReadFrameError(err)
			return err
		}
		if VerboseLogs {
			cc.vlogf("http2: Transport received %s", summarizeFrame(f))
		}
		if !gotSettings {
			if _, ok := f.(*SettingsFrame); !ok {
				cc.logf("protocol error: received %T before a SETTINGS frame", f)
				return ConnectionError(ErrCodeProtocol)
			}
			gotSettings = true
		}

		switch f := f.(type) {
		case *MetaHeadersFrame:
			err = rl.processHeaders(f)
		case *DataFrame:
			err = rl.processData(f)
		case *GoAwayFrame:
			err = rl.processGoAway(f)
		case *RSTStreamFrame:
			err = rl.processResetStream(f)
		case *SettingsFrame:
			err = rl.processSettings(f)
		case *PushPromiseFrame:
			err = rl.processPushPromise(f)
		case *WindowUpdateFrame:
			err = rl.processWindowUpdate(f)
		case *PingFrame:
			err = rl.processPing(f)
		default:
			cc.logf("Transport: unhandled response frame type %T", f)
		}
		if err != nil {
			if VerboseLogs {
				cc.vlogf("http2: Transport conn %p received error from processing frame %v: %v", cc, summarizeFrame(f), err)
			}
			return err
		}
	}
}

func (rl *clientConnReadLoop) processHeaders(f *MetaHeadersFrame) error {
	cs := rl.streamByID(f.StreamID, headerOrDataFrame)
	if cs == nil {
		
		
		
		return nil
	}
	if cs.readClosed {
		rl.endStreamError(cs, StreamError{
			StreamID: f.StreamID,
			Code:     ErrCodeProtocol,
			Cause:    errors.New("protocol error: headers after END_STREAM"),
		})
		return nil
	}
	if !cs.firstByte {
		if cs.trace != nil {
			
			
			
			
			traceFirstResponseByte(cs.trace)
		}
		cs.firstByte = true
	}
	if !cs.pastHeaders {
		cs.pastHeaders = true
	} else {
		return rl.processTrailers(cs, f)
	}

	res, err := rl.handleResponse(cs, f)
	if err != nil {
		if _, ok := err.(ConnectionError); ok {
			return err
		}
		
		rl.endStreamError(cs, StreamError{
			StreamID: f.StreamID,
			Code:     ErrCodeProtocol,
			Cause:    err,
		})
		return nil 
	}
	if res == nil {
		
		return nil
	}
	cs.resTrailer = &res.Trailer
	cs.res = res
	close(cs.respHeaderRecv)
	if f.StreamEnded() {
		rl.endStream(cs)
	}
	return nil
}







func (rl *clientConnReadLoop) handleResponse(cs *clientStream, f *MetaHeadersFrame) (*http.Response, error) {
	if f.Truncated {
		return nil, errResponseHeaderListSize
	}

	status := f.PseudoValue("status")
	if status == "" {
		return nil, errors.New("malformed response from server: missing status pseudo header")
	}
	statusCode, err := strconv.Atoi(status)
	if err != nil {
		return nil, errors.New("malformed response from server: malformed non-numeric status pseudo header")
	}

	regularFields := f.RegularFields()
	strs := make([]string, len(regularFields))
	header := make(http.Header, len(regularFields))
	res := &http.Response{
		Proto:      "HTTP/2.0",
		ProtoMajor: 2,
		Header:     header,
		StatusCode: statusCode,
		Status:     status + " " + http.StatusText(statusCode),
	}
	for _, hf := range regularFields {
		key := httpcommon.CanonicalHeader(hf.Name)
		if key == "Trailer" {
			t := res.Trailer
			if t == nil {
				t = make(http.Header)
				res.Trailer = t
			}
			foreachHeaderElement(hf.Value, func(v string) {
				t[httpcommon.CanonicalHeader(v)] = nil
			})
		} else {
			vv := header[key]
			if vv == nil && len(strs) > 0 {
				
				
				
				
				vv, strs = strs[:1:1], strs[1:]
				vv[0] = hf.Value
				header[key] = vv
			} else {
				header[key] = append(vv, hf.Value)
			}
		}
	}

	if statusCode >= 100 && statusCode <= 199 {
		if f.StreamEnded() {
			return nil, errors.New("1xx informational response with END_STREAM flag")
		}
		if fn := cs.get1xxTraceFunc(); fn != nil {
			
			
			
			if err := fn(statusCode, textproto.MIMEHeader(header)); err != nil {
				return nil, err
			}
		} else {
			
			
			
			
			
			
			
			limit := int64(cs.cc.t.maxHeaderListSize())
			if t1 := cs.cc.t.t1; t1 != nil && t1.MaxResponseHeaderBytes > limit {
				limit = t1.MaxResponseHeaderBytes
			}
			for _, h := range f.Fields {
				cs.totalHeaderSize += int64(h.Size())
			}
			if cs.totalHeaderSize > limit {
				if VerboseLogs {
					log.Printf("http2: 1xx informational responses too large")
				}
				return nil, errors.New("header list too large")
			}
		}
		if statusCode == 100 {
			traceGot100Continue(cs.trace)
			select {
			case cs.on100 <- struct{}{}:
			default:
			}
		}
		cs.pastHeaders = false 
		return nil, nil
	}

	res.ContentLength = -1
	if clens := res.Header["Content-Length"]; len(clens) == 1 {
		if cl, err := strconv.ParseUint(clens[0], 10, 63); err == nil {
			res.ContentLength = int64(cl)
		} else {
			
			
		}
	} else if len(clens) > 1 {
		
		
	} else if f.StreamEnded() && !cs.isHead {
		res.ContentLength = 0
	}

	if cs.isHead {
		res.Body = noBody
		return res, nil
	}

	if f.StreamEnded() {
		if res.ContentLength > 0 {
			res.Body = missingBody{}
		} else {
			res.Body = noBody
		}
		return res, nil
	}

	cs.bufPipe.setBuffer(&dataBuffer{expected: res.ContentLength})
	cs.bytesRemain = res.ContentLength
	res.Body = transportResponseBody{cs}

	if cs.requestedGzip && asciiEqualFold(res.Header.Get("Content-Encoding"), "gzip") {
		res.Header.Del("Content-Encoding")
		res.Header.Del("Content-Length")
		res.ContentLength = -1
		res.Body = &gzipReader{body: res.Body}
		res.Uncompressed = true
	}
	return res, nil
}

func (rl *clientConnReadLoop) processTrailers(cs *clientStream, f *MetaHeadersFrame) error {
	if cs.pastTrailers {
		
		return ConnectionError(ErrCodeProtocol)
	}
	cs.pastTrailers = true
	if !f.StreamEnded() {
		
		
		return ConnectionError(ErrCodeProtocol)
	}
	if len(f.PseudoFields()) > 0 {
		
		
		return ConnectionError(ErrCodeProtocol)
	}

	trailer := make(http.Header)
	for _, hf := range f.RegularFields() {
		key := httpcommon.CanonicalHeader(hf.Name)
		trailer[key] = append(trailer[key], hf.Value)
	}
	cs.trailer = trailer

	rl.endStream(cs)
	return nil
}



type transportResponseBody struct {
	cs *clientStream
}

func (b transportResponseBody) Read(p []byte) (n int, err error) {
	cs := b.cs
	cc := cs.cc

	if cs.readErr != nil {
		return 0, cs.readErr
	}
	n, err = b.cs.bufPipe.Read(p)
	if cs.bytesRemain != -1 {
		if int64(n) > cs.bytesRemain {
			n = int(cs.bytesRemain)
			if err == nil {
				err = errors.New("net/http: server replied with more than declared Content-Length; truncated")
				cs.abortStream(err)
			}
			cs.readErr = err
			return int(cs.bytesRemain), err
		}
		cs.bytesRemain -= int64(n)
		if err == io.EOF && cs.bytesRemain > 0 {
			err = io.ErrUnexpectedEOF
			cs.readErr = err
			return n, err
		}
	}
	if n == 0 {
		
		return
	}

	cc.mu.Lock()
	connAdd := cc.inflow.add(n)
	var streamAdd int32
	if err == nil { 
		streamAdd = cs.inflow.add(n)
	}
	cc.mu.Unlock()

	if connAdd != 0 || streamAdd != 0 {
		cc.wmu.Lock()
		defer cc.wmu.Unlock()
		if connAdd != 0 {
			cc.fr.WriteWindowUpdate(0, mustUint31(connAdd))
		}
		if streamAdd != 0 {
			cc.fr.WriteWindowUpdate(cs.ID, mustUint31(streamAdd))
		}
		cc.bw.Flush()
	}
	return
}

var errClosedResponseBody = errors.New("http2: response body closed")

func (b transportResponseBody) Close() error {
	cs := b.cs
	cc := cs.cc

	cs.bufPipe.BreakWithError(errClosedResponseBody)
	cs.abortStream(errClosedResponseBody)

	unread := cs.bufPipe.Len()
	if unread > 0 {
		cc.mu.Lock()
		
		connAdd := cc.inflow.add(unread)
		cc.mu.Unlock()

		
		
		cc.wmu.Lock()
		
		if connAdd > 0 {
			cc.fr.WriteWindowUpdate(0, uint32(connAdd))
		}
		cc.bw.Flush()
		cc.wmu.Unlock()
	}

	select {
	case <-cs.donec:
	case <-cs.ctx.Done():
		
		
		
		return nil
	case <-cs.reqCancel:
		return errRequestCanceled
	}
	return nil
}

func (rl *clientConnReadLoop) processData(f *DataFrame) error {
	cc := rl.cc
	cs := rl.streamByID(f.StreamID, headerOrDataFrame)
	data := f.Data()
	if cs == nil {
		cc.mu.Lock()
		neverSent := cc.nextStreamID
		cc.mu.Unlock()
		if f.StreamID >= neverSent {
			
			cc.logf("http2: Transport received unsolicited DATA frame; closing connection")
			return ConnectionError(ErrCodeProtocol)
		}
		
		
		
		

		
		if f.Length > 0 {
			cc.mu.Lock()
			ok := cc.inflow.take(f.Length)
			connAdd := cc.inflow.add(int(f.Length))
			cc.mu.Unlock()
			if !ok {
				return ConnectionError(ErrCodeFlowControl)
			}
			if connAdd > 0 {
				cc.wmu.Lock()
				cc.fr.WriteWindowUpdate(0, uint32(connAdd))
				cc.bw.Flush()
				cc.wmu.Unlock()
			}
		}
		return nil
	}
	if cs.readClosed {
		cc.logf("protocol error: received DATA after END_STREAM")
		rl.endStreamError(cs, StreamError{
			StreamID: f.StreamID,
			Code:     ErrCodeProtocol,
		})
		return nil
	}
	if !cs.pastHeaders {
		cc.logf("protocol error: received DATA before a HEADERS frame")
		rl.endStreamError(cs, StreamError{
			StreamID: f.StreamID,
			Code:     ErrCodeProtocol,
		})
		return nil
	}
	if f.Length > 0 {
		if cs.isHead && len(data) > 0 {
			cc.logf("protocol error: received DATA on a HEAD request")
			rl.endStreamError(cs, StreamError{
				StreamID: f.StreamID,
				Code:     ErrCodeProtocol,
			})
			return nil
		}
		
		cc.mu.Lock()
		if !takeInflows(&cc.inflow, &cs.inflow, f.Length) {
			cc.mu.Unlock()
			return ConnectionError(ErrCodeFlowControl)
		}
		
		
		var refund int
		if pad := int(f.Length) - len(data); pad > 0 {
			refund += pad
		}

		didReset := false
		var err error
		if len(data) > 0 {
			if _, err = cs.bufPipe.Write(data); err != nil {
				
				
				didReset = true
				refund += len(data)
			}
		}

		sendConn := cc.inflow.add(refund)
		var sendStream int32
		if !didReset {
			sendStream = cs.inflow.add(refund)
		}
		cc.mu.Unlock()

		if sendConn > 0 || sendStream > 0 {
			cc.wmu.Lock()
			if sendConn > 0 {
				cc.fr.WriteWindowUpdate(0, uint32(sendConn))
			}
			if sendStream > 0 {
				cc.fr.WriteWindowUpdate(cs.ID, uint32(sendStream))
			}
			cc.bw.Flush()
			cc.wmu.Unlock()
		}

		if err != nil {
			rl.endStreamError(cs, err)
			return nil
		}
	}

	if f.StreamEnded() {
		rl.endStream(cs)
	}
	return nil
}

func (rl *clientConnReadLoop) endStream(cs *clientStream) {
	
	
	if !cs.readClosed {
		cs.readClosed = true
		
		
		
		
		rl.cc.mu.Lock()
		defer rl.cc.mu.Unlock()
		cs.bufPipe.closeWithErrorAndCode(io.EOF, cs.copyTrailers)
		close(cs.peerClosed)
	}
}

func (rl *clientConnReadLoop) endStreamError(cs *clientStream, err error) {
	cs.readAborted = true
	cs.abortStream(err)
}


const (
	headerOrDataFrame    = true
	notHeaderOrDataFrame = false
)



func (rl *clientConnReadLoop) streamByID(id uint32, headerOrData bool) *clientStream {
	rl.cc.mu.Lock()
	defer rl.cc.mu.Unlock()
	if headerOrData {
		
		
		rl.cc.rstStreamPingsBlocked = false
	}
	cs := rl.cc.streams[id]
	if cs != nil && !cs.readAborted {
		return cs
	}
	return nil
}

func (cs *clientStream) copyTrailers() {
	for k, vv := range cs.trailer {
		t := cs.resTrailer
		if *t == nil {
			*t = make(http.Header)
		}
		(*t)[k] = vv
	}
}

func (rl *clientConnReadLoop) processGoAway(f *GoAwayFrame) error {
	cc := rl.cc
	cc.t.connPool().MarkDead(cc)
	if f.ErrCode != 0 {
		
		cc.vlogf("transport got GOAWAY with error code = %v", f.ErrCode)
		if fn := cc.t.CountError; fn != nil {
			fn("recv_goaway_" + f.ErrCode.stringToken())
		}
	}
	cc.setGoAway(f)
	return nil
}

func (rl *clientConnReadLoop) processSettings(f *SettingsFrame) error {
	cc := rl.cc
	
	
	cc.wmu.Lock()
	defer cc.wmu.Unlock()

	if err := rl.processSettingsNoWrite(f); err != nil {
		return err
	}
	if !f.IsAck() {
		cc.fr.WriteSettingsAck()
		cc.bw.Flush()
	}
	return nil
}

func (rl *clientConnReadLoop) processSettingsNoWrite(f *SettingsFrame) error {
	cc := rl.cc
	cc.mu.Lock()
	defer cc.mu.Unlock()

	if f.IsAck() {
		if cc.wantSettingsAck {
			cc.wantSettingsAck = false
			return nil
		}
		return ConnectionError(ErrCodeProtocol)
	}

	var seenMaxConcurrentStreams bool
	err := f.ForeachSetting(func(s Setting) error {
		switch s.ID {
		case SettingMaxFrameSize:
			cc.maxFrameSize = s.Val
		case SettingMaxConcurrentStreams:
			cc.maxConcurrentStreams = s.Val
			seenMaxConcurrentStreams = true
		case SettingMaxHeaderListSize:
			cc.peerMaxHeaderListSize = uint64(s.Val)
		case SettingInitialWindowSize:
			
			
			
			
			if s.Val > math.MaxInt32 {
				return ConnectionError(ErrCodeFlowControl)
			}

			
			
			
			delta := int32(s.Val) - int32(cc.initialWindowSize)
			for _, cs := range cc.streams {
				cs.flow.add(delta)
			}
			cc.cond.Broadcast()

			cc.initialWindowSize = s.Val
		case SettingHeaderTableSize:
			cc.henc.SetMaxDynamicTableSize(s.Val)
			cc.peerMaxHeaderTableSize = s.Val
		case SettingEnableConnectProtocol:
			if err := s.Valid(); err != nil {
				return err
			}
			
			
			
			
			
			
			
			
			if !cc.seenSettings {
				cc.extendedConnectAllowed = s.Val == 1
			}
		default:
			cc.vlogf("Unhandled Setting: %v", s)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if !cc.seenSettings {
		if !seenMaxConcurrentStreams {
			
			
			
			
			cc.maxConcurrentStreams = defaultMaxConcurrentStreams
		}
		close(cc.seenSettingsChan)
		cc.seenSettings = true
	}

	return nil
}

func (rl *clientConnReadLoop) processWindowUpdate(f *WindowUpdateFrame) error {
	cc := rl.cc
	cs := rl.streamByID(f.StreamID, notHeaderOrDataFrame)
	if f.StreamID != 0 && cs == nil {
		return nil
	}

	cc.mu.Lock()
	defer cc.mu.Unlock()

	fl := &cc.flow
	if cs != nil {
		fl = &cs.flow
	}
	if !fl.add(int32(f.Increment)) {
		
		if cs != nil {
			rl.endStreamError(cs, StreamError{
				StreamID: f.StreamID,
				Code:     ErrCodeFlowControl,
			})
			return nil
		}

		return ConnectionError(ErrCodeFlowControl)
	}
	cc.cond.Broadcast()
	return nil
}

func (rl *clientConnReadLoop) processResetStream(f *RSTStreamFrame) error {
	cs := rl.streamByID(f.StreamID, notHeaderOrDataFrame)
	if cs == nil {
		
		return nil
	}
	serr := streamError(cs.ID, f.ErrCode)
	serr.Cause = errFromPeer
	if f.ErrCode == ErrCodeProtocol {
		rl.cc.SetDoNotReuse()
	}
	if fn := cs.cc.t.CountError; fn != nil {
		fn("recv_rststream_" + f.ErrCode.stringToken())
	}
	cs.abortStream(serr)

	cs.bufPipe.CloseWithError(serr)
	return nil
}


func (cc *ClientConn) Ping(ctx context.Context) error {
	c := make(chan struct{})
	
	var p [8]byte
	for {
		if _, err := rand.Read(p[:]); err != nil {
			return err
		}
		cc.mu.Lock()
		
		if _, found := cc.pings[p]; !found {
			cc.pings[p] = c
			cc.mu.Unlock()
			break
		}
		cc.mu.Unlock()
	}
	var pingError error
	errc := make(chan struct{})
	go func() {
		cc.t.markNewGoroutine()
		cc.wmu.Lock()
		defer cc.wmu.Unlock()
		if pingError = cc.fr.WritePing(false, p); pingError != nil {
			close(errc)
			return
		}
		if pingError = cc.bw.Flush(); pingError != nil {
			close(errc)
			return
		}
	}()
	select {
	case <-c:
		return nil
	case <-errc:
		return pingError
	case <-ctx.Done():
		return ctx.Err()
	case <-cc.readerDone:
		
		return cc.readerErr
	}
}

func (rl *clientConnReadLoop) processPing(f *PingFrame) error {
	if f.IsAck() {
		cc := rl.cc
		cc.mu.Lock()
		defer cc.mu.Unlock()
		
		if c, ok := cc.pings[f.Data]; ok {
			close(c)
			delete(cc.pings, f.Data)
		}
		if cc.pendingResets > 0 {
			
			cc.pendingResets = 0
			cc.rstStreamPingsBlocked = true
			cc.cond.Broadcast()
		}
		return nil
	}
	cc := rl.cc
	cc.wmu.Lock()
	defer cc.wmu.Unlock()
	if err := cc.fr.WritePing(true, f.Data); err != nil {
		return err
	}
	return cc.bw.Flush()
}

func (rl *clientConnReadLoop) processPushPromise(f *PushPromiseFrame) error {
	
	
	
	
	
	
	
	return ConnectionError(ErrCodeProtocol)
}



func (cc *ClientConn) writeStreamReset(streamID uint32, code ErrCode, ping bool, err error) {
	
	
	
	
	cc.wmu.Lock()
	cc.fr.WriteRSTStream(streamID, code)
	if ping {
		var payload [8]byte
		rand.Read(payload[:])
		cc.fr.WritePing(false, payload)
	}
	cc.bw.Flush()
	cc.wmu.Unlock()
}

var (
	errResponseHeaderListSize = errors.New("http2: response header list larger than advertised limit")
	errRequestHeaderListSize  = httpcommon.ErrRequestHeaderListSize
)

func (cc *ClientConn) logf(format string, args ...interface{}) {
	cc.t.logf(format, args...)
}

func (cc *ClientConn) vlogf(format string, args ...interface{}) {
	cc.t.vlogf(format, args...)
}

func (t *Transport) vlogf(format string, args ...interface{}) {
	if VerboseLogs {
		t.logf(format, args...)
	}
}

func (t *Transport) logf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

var noBody io.ReadCloser = noBodyReader{}

type noBodyReader struct{}

func (noBodyReader) Close() error             { return nil }
func (noBodyReader) Read([]byte) (int, error) { return 0, io.EOF }

type missingBody struct{}

func (missingBody) Close() error             { return nil }
func (missingBody) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }

func strSliceContains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

type erringRoundTripper struct{ err error }

func (rt erringRoundTripper) RoundTripErr() error                             { return rt.err }
func (rt erringRoundTripper) RoundTrip(*http.Request) (*http.Response, error) { return nil, rt.err }



type gzipReader struct {
	_    incomparable
	body io.ReadCloser 
	zr   *gzip.Reader  
	zerr error         
}

func (gz *gzipReader) Read(p []byte) (n int, err error) {
	if gz.zerr != nil {
		return 0, gz.zerr
	}
	if gz.zr == nil {
		gz.zr, err = gzip.NewReader(gz.body)
		if err != nil {
			gz.zerr = err
			return 0, err
		}
	}
	return gz.zr.Read(p)
}

func (gz *gzipReader) Close() error {
	if err := gz.body.Close(); err != nil {
		return err
	}
	gz.zerr = fs.ErrClosed
	return nil
}

type errorReader struct{ err error }

func (r errorReader) Read(p []byte) (int, error) { return 0, r.err }



func isConnectionCloseRequest(req *http.Request) bool {
	return req.Close || httpguts.HeaderValuesContainsToken(req.Header["Connection"], "close")
}



func registerHTTPSProtocol(t *http.Transport, rt noDialH2RoundTripper) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	t.RegisterProtocol("https", rt)
	return nil
}





type noDialH2RoundTripper struct{ *Transport }

func (rt noDialH2RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := rt.Transport.RoundTrip(req)
	if isNoCachedConnError(err) {
		return nil, http.ErrSkipAltProtocol
	}
	return res, err
}

func (t *Transport) idleConnTimeout() time.Duration {
	
	
	
	if t.IdleConnTimeout != 0 {
		return t.IdleConnTimeout
	}

	if t.t1 != nil {
		return t.t1.IdleConnTimeout
	}

	return 0
}

func traceGetConn(req *http.Request, hostPort string) {
	trace := httptrace.ContextClientTrace(req.Context())
	if trace == nil || trace.GetConn == nil {
		return
	}
	trace.GetConn(hostPort)
}

func traceGotConn(req *http.Request, cc *ClientConn, reused bool) {
	trace := httptrace.ContextClientTrace(req.Context())
	if trace == nil || trace.GotConn == nil {
		return
	}
	ci := httptrace.GotConnInfo{Conn: cc.tconn}
	ci.Reused = reused
	cc.mu.Lock()
	ci.WasIdle = len(cc.streams) == 0 && reused
	if ci.WasIdle && !cc.lastActive.IsZero() {
		ci.IdleTime = cc.t.timeSince(cc.lastActive)
	}
	cc.mu.Unlock()

	trace.GotConn(ci)
}

func traceWroteHeaders(trace *httptrace.ClientTrace) {
	if trace != nil && trace.WroteHeaders != nil {
		trace.WroteHeaders()
	}
}

func traceGot100Continue(trace *httptrace.ClientTrace) {
	if trace != nil && trace.Got100Continue != nil {
		trace.Got100Continue()
	}
}

func traceWait100Continue(trace *httptrace.ClientTrace) {
	if trace != nil && trace.Wait100Continue != nil {
		trace.Wait100Continue()
	}
}

func traceWroteRequest(trace *httptrace.ClientTrace, err error) {
	if trace != nil && trace.WroteRequest != nil {
		trace.WroteRequest(httptrace.WroteRequestInfo{Err: err})
	}
}

func traceFirstResponseByte(trace *httptrace.ClientTrace) {
	if trace != nil && trace.GotFirstResponseByte != nil {
		trace.GotFirstResponseByte()
	}
}

func traceGot1xxResponseFunc(trace *httptrace.ClientTrace) func(int, textproto.MIMEHeader) error {
	if trace != nil {
		return trace.Got1xxResponse
	}
	return nil
}



func (t *Transport) dialTLSWithContext(ctx context.Context, network, addr string, cfg *tls.Config) (*tls.Conn, error) {
	dialer := &tls.Dialer{
		Config: cfg,
	}
	cn, err := dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}
	tlsCn := cn.(*tls.Conn) 
	return tlsCn, nil
}
