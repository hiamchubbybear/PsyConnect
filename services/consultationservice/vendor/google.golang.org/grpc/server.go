

package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
	estats "google.golang.org/grpc/experimental/stats"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal"
	"google.golang.org/grpc/internal/binarylog"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/grpcsync"
	"google.golang.org/grpc/internal/grpcutil"
	istats "google.golang.org/grpc/internal/stats"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/mem"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
)

const (
	defaultServerMaxReceiveMessageSize = 1024 * 1024 * 4
	defaultServerMaxSendMessageSize    = math.MaxInt32

	
	
	
	
	
	listenerAddressForServeHTTP = "listenerAddressForServeHTTP"
)

func init() {
	internal.GetServerCredentials = func(srv *Server) credentials.TransportCredentials {
		return srv.opts.creds
	}
	internal.IsRegisteredMethod = func(srv *Server, method string) bool {
		return srv.isRegisteredMethod(method)
	}
	internal.ServerFromContext = serverFromContext
	internal.AddGlobalServerOptions = func(opt ...ServerOption) {
		globalServerOptions = append(globalServerOptions, opt...)
	}
	internal.ClearGlobalServerOptions = func() {
		globalServerOptions = nil
	}
	internal.BinaryLogger = binaryLogger
	internal.JoinServerOptions = newJoinServerOption
	internal.BufferPool = bufferPool
	internal.MetricsRecorderForServer = func(srv *Server) estats.MetricsRecorder {
		return istats.NewMetricsRecorderList(srv.opts.statsHandlers)
	}
}

var statusOK = status.New(codes.OK, "")
var logger = grpclog.Component("core")


type MethodHandler func(srv any, ctx context.Context, dec func(any) error, interceptor UnaryServerInterceptor) (any, error)


type MethodDesc struct {
	MethodName string
	Handler    MethodHandler
}


type ServiceDesc struct {
	ServiceName string
	
	
	HandlerType any
	Methods     []MethodDesc
	Streams     []StreamDesc
	Metadata    any
}



type serviceInfo struct {
	
	serviceImpl any
	methods     map[string]*MethodDesc
	streams     map[string]*StreamDesc
	mdata       any
}


type Server struct {
	opts serverOptions

	mu  sync.Mutex 
	lis map[net.Listener]bool
	
	
	
	conns    map[string]map[transport.ServerTransport]bool
	serve    bool
	drain    bool
	cv       *sync.Cond              
	services map[string]*serviceInfo 
	events   traceEventLog

	quit               *grpcsync.Event
	done               *grpcsync.Event
	channelzRemoveOnce sync.Once
	serveWG            sync.WaitGroup 
	handlersWG         sync.WaitGroup 

	channelz *channelz.Server

	serverWorkerChannel      chan func()
	serverWorkerChannelClose func()
}

type serverOptions struct {
	creds                 credentials.TransportCredentials
	codec                 baseCodec
	cp                    Compressor
	dc                    Decompressor
	unaryInt              UnaryServerInterceptor
	streamInt             StreamServerInterceptor
	chainUnaryInts        []UnaryServerInterceptor
	chainStreamInts       []StreamServerInterceptor
	binaryLogger          binarylog.Logger
	inTapHandle           tap.ServerInHandle
	statsHandlers         []stats.Handler
	maxConcurrentStreams  uint32
	maxReceiveMessageSize int
	maxSendMessageSize    int
	unknownStreamDesc     *StreamDesc
	keepaliveParams       keepalive.ServerParameters
	keepalivePolicy       keepalive.EnforcementPolicy
	initialWindowSize     int32
	initialConnWindowSize int32
	writeBufferSize       int
	readBufferSize        int
	sharedWriteBuffer     bool
	connectionTimeout     time.Duration
	maxHeaderListSize     *uint32
	headerTableSize       *uint32
	numServerWorkers      uint32
	bufferPool            mem.BufferPool
	waitForHandlers       bool
}

var defaultServerOptions = serverOptions{
	maxConcurrentStreams:  math.MaxUint32,
	maxReceiveMessageSize: defaultServerMaxReceiveMessageSize,
	maxSendMessageSize:    defaultServerMaxSendMessageSize,
	connectionTimeout:     120 * time.Second,
	writeBufferSize:       defaultWriteBufSize,
	readBufferSize:        defaultReadBufSize,
	bufferPool:            mem.DefaultBufferPool(),
}
var globalServerOptions []ServerOption


type ServerOption interface {
	apply(*serverOptions)
}








type EmptyServerOption struct{}

func (EmptyServerOption) apply(*serverOptions) {}



type funcServerOption struct {
	f func(*serverOptions)
}

func (fdo *funcServerOption) apply(do *serverOptions) {
	fdo.f(do)
}

func newFuncServerOption(f func(*serverOptions)) *funcServerOption {
	return &funcServerOption{
		f: f,
	}
}



type joinServerOption struct {
	opts []ServerOption
}

func (mdo *joinServerOption) apply(do *serverOptions) {
	for _, opt := range mdo.opts {
		opt.apply(do)
	}
}

func newJoinServerOption(opts ...ServerOption) ServerOption {
	return &joinServerOption{opts: opts}
}









func SharedWriteBuffer(val bool) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.sharedWriteBuffer = val
	})
}





func WriteBufferSize(s int) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.writeBufferSize = s
	})
}





func ReadBufferSize(s int) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.readBufferSize = s
	})
}



func InitialWindowSize(s int32) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.initialWindowSize = s
	})
}



func InitialConnWindowSize(s int32) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.initialConnWindowSize = s
	})
}


func KeepaliveParams(kp keepalive.ServerParameters) ServerOption {
	if kp.Time > 0 && kp.Time < internal.KeepaliveMinServerPingTime {
		logger.Warning("Adjusting keepalive ping interval to minimum period of 1s")
		kp.Time = internal.KeepaliveMinServerPingTime
	}

	return newFuncServerOption(func(o *serverOptions) {
		o.keepaliveParams = kp
	})
}


func KeepaliveEnforcementPolicy(kep keepalive.EnforcementPolicy) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.keepalivePolicy = kep
	})
}










func CustomCodec(codec Codec) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.codec = newCodecV0Bridge(codec)
	})
}
























func ForceServerCodec(codec encoding.Codec) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.codec = newCodecV1Bridge(codec)
	})
}










func ForceServerCodecV2(codecV2 encoding.CodecV2) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.codec = codecV2
	})
}









func RPCCompressor(cp Compressor) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.cp = cp
	})
}







func RPCDecompressor(dc Decompressor) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.dc = dc
	})
}





func MaxMsgSize(m int) ServerOption {
	return MaxRecvMsgSize(m)
}



func MaxRecvMsgSize(m int) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.maxReceiveMessageSize = m
	})
}



func MaxSendMsgSize(m int) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.maxSendMessageSize = m
	})
}



func MaxConcurrentStreams(n uint32) ServerOption {
	if n == 0 {
		n = math.MaxUint32
	}
	return newFuncServerOption(func(o *serverOptions) {
		o.maxConcurrentStreams = n
	})
}


func Creds(c credentials.TransportCredentials) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.creds = c
	})
}




func UnaryInterceptor(i UnaryServerInterceptor) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		if o.unaryInt != nil {
			panic("The unary server interceptor was already set and may not be reset.")
		}
		o.unaryInt = i
	})
}





func ChainUnaryInterceptor(interceptors ...UnaryServerInterceptor) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.chainUnaryInts = append(o.chainUnaryInts, interceptors...)
	})
}



func StreamInterceptor(i StreamServerInterceptor) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		if o.streamInt != nil {
			panic("The stream server interceptor was already set and may not be reset.")
		}
		o.streamInt = i
	})
}





func ChainStreamInterceptor(interceptors ...StreamServerInterceptor) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.chainStreamInts = append(o.chainStreamInts, interceptors...)
	})
}








func InTapHandle(h tap.ServerInHandle) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		if o.inTapHandle != nil {
			panic("The tap handle was already set and may not be reset.")
		}
		o.inTapHandle = h
	})
}


func StatsHandler(h stats.Handler) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		if h == nil {
			logger.Error("ignoring nil parameter in grpc.StatsHandler ServerOption")
			
			
			return
		}
		o.statsHandlers = append(o.statsHandlers, h)
	})
}



func binaryLogger(bl binarylog.Logger) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.binaryLogger = bl
	})
}







func UnknownServiceHandler(streamHandler StreamHandler) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.unknownStreamDesc = &StreamDesc{
			StreamName: "unknown_service_handler",
			Handler:    streamHandler,
			
			ClientStreams: true,
			ServerStreams: true,
		}
	})
}










func ConnectionTimeout(d time.Duration) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.connectionTimeout = d
	})
}



type MaxHeaderListSizeServerOption struct {
	MaxHeaderListSize uint32
}

func (o MaxHeaderListSizeServerOption) apply(so *serverOptions) {
	so.maxHeaderListSize = &o.MaxHeaderListSize
}



func MaxHeaderListSize(s uint32) ServerOption {
	return MaxHeaderListSizeServerOption{
		MaxHeaderListSize: s,
	}
}








func HeaderTableSize(s uint32) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.headerTableSize = &s
	})
}










func NumStreamWorkers(numServerWorkers uint32) ServerOption {
	
	
	
	
	return newFuncServerOption(func(o *serverOptions) {
		o.numServerWorkers = numServerWorkers
	})
}










func WaitForHandlers(w bool) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.waitForHandlers = w
	})
}

func bufferPool(bufferPool mem.BufferPool) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.bufferPool = bufferPool
	})
}






const serverWorkerResetThreshold = 1 << 16







func (s *Server) serverWorker() {
	for completed := 0; completed < serverWorkerResetThreshold; completed++ {
		f, ok := <-s.serverWorkerChannel
		if !ok {
			return
		}
		f()
	}
	go s.serverWorker()
}



func (s *Server) initServerWorkers() {
	s.serverWorkerChannel = make(chan func())
	s.serverWorkerChannelClose = sync.OnceFunc(func() {
		close(s.serverWorkerChannel)
	})
	for i := uint32(0); i < s.opts.numServerWorkers; i++ {
		go s.serverWorker()
	}
}



func NewServer(opt ...ServerOption) *Server {
	opts := defaultServerOptions
	for _, o := range globalServerOptions {
		o.apply(&opts)
	}
	for _, o := range opt {
		o.apply(&opts)
	}
	s := &Server{
		lis:      make(map[net.Listener]bool),
		opts:     opts,
		conns:    make(map[string]map[transport.ServerTransport]bool),
		services: make(map[string]*serviceInfo),
		quit:     grpcsync.NewEvent(),
		done:     grpcsync.NewEvent(),
		channelz: channelz.RegisterServer(""),
	}
	chainUnaryServerInterceptors(s)
	chainStreamServerInterceptors(s)
	s.cv = sync.NewCond(&s.mu)
	if EnableTracing {
		_, file, line, _ := runtime.Caller(1)
		s.events = newTraceEventLog("grpc.Server", fmt.Sprintf("%s:%d", file, line))
	}

	if s.opts.numServerWorkers > 0 {
		s.initServerWorkers()
	}

	channelz.Info(logger, s.channelz, "Server created")
	return s
}



func (s *Server) printf(format string, a ...any) {
	if s.events != nil {
		s.events.Printf(format, a...)
	}
}



func (s *Server) errorf(format string, a ...any) {
	if s.events != nil {
		s.events.Errorf(format, a...)
	}
}




type ServiceRegistrar interface {
	
	
	
	
	
	RegisterService(desc *ServiceDesc, impl any)
}





func (s *Server) RegisterService(sd *ServiceDesc, ss any) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			logger.Fatalf("grpc: Server.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		}
	}
	s.register(sd, ss)
}

func (s *Server) register(sd *ServiceDesc, ss any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.printf("RegisterService(%q)", sd.ServiceName)
	if s.serve {
		logger.Fatalf("grpc: Server.RegisterService after Server.Serve for %q", sd.ServiceName)
	}
	if _, ok := s.services[sd.ServiceName]; ok {
		logger.Fatalf("grpc: Server.RegisterService found duplicate service registration for %q", sd.ServiceName)
	}
	info := &serviceInfo{
		serviceImpl: ss,
		methods:     make(map[string]*MethodDesc),
		streams:     make(map[string]*StreamDesc),
		mdata:       sd.Metadata,
	}
	for i := range sd.Methods {
		d := &sd.Methods[i]
		info.methods[d.MethodName] = d
	}
	for i := range sd.Streams {
		d := &sd.Streams[i]
		info.streams[d.StreamName] = d
	}
	s.services[sd.ServiceName] = info
}


type MethodInfo struct {
	
	Name string
	
	IsClientStream bool
	
	IsServerStream bool
}


type ServiceInfo struct {
	Methods []MethodInfo
	
	Metadata any
}



func (s *Server) GetServiceInfo() map[string]ServiceInfo {
	ret := make(map[string]ServiceInfo)
	for n, srv := range s.services {
		methods := make([]MethodInfo, 0, len(srv.methods)+len(srv.streams))
		for m := range srv.methods {
			methods = append(methods, MethodInfo{
				Name:           m,
				IsClientStream: false,
				IsServerStream: false,
			})
		}
		for m, d := range srv.streams {
			methods = append(methods, MethodInfo{
				Name:           m,
				IsClientStream: d.ClientStreams,
				IsServerStream: d.ServerStreams,
			})
		}

		ret[n] = ServiceInfo{
			Methods:  methods,
			Metadata: srv.mdata,
		}
	}
	return ret
}



var ErrServerStopped = errors.New("grpc: the server has been stopped")

type listenSocket struct {
	net.Listener
	channelz *channelz.Socket
}

func (l *listenSocket) Close() error {
	err := l.Listener.Close()
	channelz.RemoveEntry(l.channelz.ID)
	channelz.Info(logger, l.channelz, "ListenSocket deleted")
	return err
}



















func (s *Server) Serve(lis net.Listener) error {
	s.mu.Lock()
	s.printf("serving")
	s.serve = true
	if s.lis == nil {
		
		s.mu.Unlock()
		lis.Close()
		return ErrServerStopped
	}

	s.serveWG.Add(1)
	defer func() {
		s.serveWG.Done()
		if s.quit.HasFired() {
			
			<-s.done.Done()
		}
	}()

	ls := &listenSocket{
		Listener: lis,
		channelz: channelz.RegisterSocket(&channelz.Socket{
			SocketType:    channelz.SocketTypeListen,
			Parent:        s.channelz,
			RefName:       lis.Addr().String(),
			LocalAddr:     lis.Addr(),
			SocketOptions: channelz.GetSocketOption(lis)},
		),
	}
	s.lis[ls] = true

	defer func() {
		s.mu.Lock()
		if s.lis != nil && s.lis[ls] {
			ls.Close()
			delete(s.lis, ls)
		}
		s.mu.Unlock()
	}()

	s.mu.Unlock()
	channelz.Info(logger, ls.channelz, "ListenSocket created")

	var tempDelay time.Duration 
	for {
		rawConn, err := lis.Accept()
		if err != nil {
			if ne, ok := err.(interface {
				Temporary() bool
			}); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				s.mu.Lock()
				s.printf("Accept error: %v; retrying in %v", err, tempDelay)
				s.mu.Unlock()
				timer := time.NewTimer(tempDelay)
				select {
				case <-timer.C:
				case <-s.quit.Done():
					timer.Stop()
					return nil
				}
				continue
			}
			s.mu.Lock()
			s.printf("done serving; Accept = %v", err)
			s.mu.Unlock()

			if s.quit.HasFired() {
				return nil
			}
			return err
		}
		tempDelay = 0
		
		
		
		
		
		s.serveWG.Add(1)
		go func() {
			s.handleRawConn(lis.Addr().String(), rawConn)
			s.serveWG.Done()
		}()
	}
}



func (s *Server) handleRawConn(lisAddr string, rawConn net.Conn) {
	if s.quit.HasFired() {
		rawConn.Close()
		return
	}
	rawConn.SetDeadline(time.Now().Add(s.opts.connectionTimeout))

	
	st := s.newHTTP2Transport(rawConn)
	rawConn.SetDeadline(time.Time{})
	if st == nil {
		return
	}

	if cc, ok := rawConn.(interface {
		PassServerTransport(transport.ServerTransport)
	}); ok {
		cc.PassServerTransport(st)
	}

	if !s.addConn(lisAddr, st) {
		return
	}
	go func() {
		s.serveStreams(context.Background(), st, rawConn)
		s.removeConn(lisAddr, st)
	}()
}



func (s *Server) newHTTP2Transport(c net.Conn) transport.ServerTransport {
	config := &transport.ServerConfig{
		MaxStreams:            s.opts.maxConcurrentStreams,
		ConnectionTimeout:     s.opts.connectionTimeout,
		Credentials:           s.opts.creds,
		InTapHandle:           s.opts.inTapHandle,
		StatsHandlers:         s.opts.statsHandlers,
		KeepaliveParams:       s.opts.keepaliveParams,
		KeepalivePolicy:       s.opts.keepalivePolicy,
		InitialWindowSize:     s.opts.initialWindowSize,
		InitialConnWindowSize: s.opts.initialConnWindowSize,
		WriteBufferSize:       s.opts.writeBufferSize,
		ReadBufferSize:        s.opts.readBufferSize,
		SharedWriteBuffer:     s.opts.sharedWriteBuffer,
		ChannelzParent:        s.channelz,
		MaxHeaderListSize:     s.opts.maxHeaderListSize,
		HeaderTableSize:       s.opts.headerTableSize,
		BufferPool:            s.opts.bufferPool,
	}
	st, err := transport.NewServerTransport(c, config)
	if err != nil {
		s.mu.Lock()
		s.errorf("NewServerTransport(%q) failed: %v", c.RemoteAddr(), err)
		s.mu.Unlock()
		
		
		if err != credentials.ErrConnDispatched {
			
			if err != io.EOF {
				channelz.Info(logger, s.channelz, "grpc: Server.Serve failed to create ServerTransport: ", err)
			}
			c.Close()
		}
		return nil
	}

	return st
}

func (s *Server) serveStreams(ctx context.Context, st transport.ServerTransport, rawConn net.Conn) {
	ctx = transport.SetConnection(ctx, rawConn)
	ctx = peer.NewContext(ctx, st.Peer())
	for _, sh := range s.opts.statsHandlers {
		ctx = sh.TagConn(ctx, &stats.ConnTagInfo{
			RemoteAddr: st.Peer().Addr,
			LocalAddr:  st.Peer().LocalAddr,
		})
		sh.HandleConn(ctx, &stats.ConnBegin{})
	}

	defer func() {
		st.Close(errors.New("finished serving streams for the server transport"))
		for _, sh := range s.opts.statsHandlers {
			sh.HandleConn(ctx, &stats.ConnEnd{})
		}
	}()

	streamQuota := newHandlerQuota(s.opts.maxConcurrentStreams)
	st.HandleStreams(ctx, func(stream *transport.ServerStream) {
		s.handlersWG.Add(1)
		streamQuota.acquire()
		f := func() {
			defer streamQuota.release()
			defer s.handlersWG.Done()
			s.handleStream(st, stream)
		}

		if s.opts.numServerWorkers > 0 {
			select {
			case s.serverWorkerChannel <- f:
				return
			default:
				
			}
		}
		go f()
	})
}

var _ http.Handler = (*Server)(nil)





























func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	st, err := transport.NewServerHandlerTransport(w, r, s.opts.statsHandlers, s.opts.bufferPool)
	if err != nil {
		
		
		return
	}
	if !s.addConn(listenerAddressForServeHTTP, st) {
		return
	}
	defer s.removeConn(listenerAddressForServeHTTP, st)
	s.serveStreams(r.Context(), st, nil)
}

func (s *Server) addConn(addr string, st transport.ServerTransport) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conns == nil {
		st.Close(errors.New("Server.addConn called when server has already been stopped"))
		return false
	}
	if s.drain {
		
		
		st.Drain("")
	}

	if s.conns[addr] == nil {
		
		s.conns[addr] = make(map[transport.ServerTransport]bool)
	}
	s.conns[addr][st] = true
	return true
}

func (s *Server) removeConn(addr string, st transport.ServerTransport) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conns := s.conns[addr]
	if conns != nil {
		delete(conns, st)
		if len(conns) == 0 {
			
			
			
			delete(s.conns, addr)
		}
		s.cv.Broadcast()
	}
}

func (s *Server) incrCallsStarted() {
	s.channelz.ServerMetrics.CallsStarted.Add(1)
	s.channelz.ServerMetrics.LastCallStartedTimestamp.Store(time.Now().UnixNano())
}

func (s *Server) incrCallsSucceeded() {
	s.channelz.ServerMetrics.CallsSucceeded.Add(1)
}

func (s *Server) incrCallsFailed() {
	s.channelz.ServerMetrics.CallsFailed.Add(1)
}

func (s *Server) sendResponse(ctx context.Context, stream *transport.ServerStream, msg any, cp Compressor, opts *transport.WriteOptions, comp encoding.Compressor) error {
	data, err := encode(s.getCodec(stream.ContentSubtype()), msg)
	if err != nil {
		channelz.Error(logger, s.channelz, "grpc: server failed to encode response: ", err)
		return err
	}

	compData, pf, err := compress(data, cp, comp, s.opts.bufferPool)
	if err != nil {
		data.Free()
		channelz.Error(logger, s.channelz, "grpc: server failed to compress response: ", err)
		return err
	}

	hdr, payload := msgHeader(data, compData, pf)

	defer func() {
		compData.Free()
		data.Free()
		
		
	}()

	dataLen := data.Len()
	payloadLen := payload.Len()
	
	if payloadLen > s.opts.maxSendMessageSize {
		return status.Errorf(codes.ResourceExhausted, "grpc: trying to send message larger than max (%d vs. %d)", payloadLen, s.opts.maxSendMessageSize)
	}
	err = stream.Write(hdr, payload, opts)
	if err == nil {
		if len(s.opts.statsHandlers) != 0 {
			for _, sh := range s.opts.statsHandlers {
				sh.HandleRPC(ctx, outPayload(false, msg, dataLen, payloadLen, time.Now()))
			}
		}
	}
	return err
}


func chainUnaryServerInterceptors(s *Server) {
	
	
	interceptors := s.opts.chainUnaryInts
	if s.opts.unaryInt != nil {
		interceptors = append([]UnaryServerInterceptor{s.opts.unaryInt}, s.opts.chainUnaryInts...)
	}

	var chainedInt UnaryServerInterceptor
	if len(interceptors) == 0 {
		chainedInt = nil
	} else if len(interceptors) == 1 {
		chainedInt = interceptors[0]
	} else {
		chainedInt = chainUnaryInterceptors(interceptors)
	}

	s.opts.unaryInt = chainedInt
}

func chainUnaryInterceptors(interceptors []UnaryServerInterceptor) UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (any, error) {
		return interceptors[0](ctx, req, info, getChainUnaryHandler(interceptors, 0, info, handler))
	}
}

func getChainUnaryHandler(interceptors []UnaryServerInterceptor, curr int, info *UnaryServerInfo, finalHandler UnaryHandler) UnaryHandler {
	if curr == len(interceptors)-1 {
		return finalHandler
	}
	return func(ctx context.Context, req any) (any, error) {
		return interceptors[curr+1](ctx, req, info, getChainUnaryHandler(interceptors, curr+1, info, finalHandler))
	}
}

func (s *Server) processUnaryRPC(ctx context.Context, stream *transport.ServerStream, info *serviceInfo, md *MethodDesc, trInfo *traceInfo) (err error) {
	shs := s.opts.statsHandlers
	if len(shs) != 0 || trInfo != nil || channelz.IsOn() {
		if channelz.IsOn() {
			s.incrCallsStarted()
		}
		var statsBegin *stats.Begin
		for _, sh := range shs {
			beginTime := time.Now()
			statsBegin = &stats.Begin{
				BeginTime:      beginTime,
				IsClientStream: false,
				IsServerStream: false,
			}
			sh.HandleRPC(ctx, statsBegin)
		}
		if trInfo != nil {
			trInfo.tr.LazyLog(&trInfo.firstLine, false)
		}
		
		
		
		
		
		
		
		
		
		
		defer func() {
			if trInfo != nil {
				if err != nil && err != io.EOF {
					trInfo.tr.LazyLog(&fmtStringer{"%v", []any{err}}, true)
					trInfo.tr.SetError()
				}
				trInfo.tr.Finish()
			}

			for _, sh := range shs {
				end := &stats.End{
					BeginTime: statsBegin.BeginTime,
					EndTime:   time.Now(),
				}
				if err != nil && err != io.EOF {
					end.Error = toRPCErr(err)
				}
				sh.HandleRPC(ctx, end)
			}

			if channelz.IsOn() {
				if err != nil && err != io.EOF {
					s.incrCallsFailed()
				} else {
					s.incrCallsSucceeded()
				}
			}
		}()
	}
	var binlogs []binarylog.MethodLogger
	if ml := binarylog.GetMethodLogger(stream.Method()); ml != nil {
		binlogs = append(binlogs, ml)
	}
	if s.opts.binaryLogger != nil {
		if ml := s.opts.binaryLogger.GetMethodLogger(stream.Method()); ml != nil {
			binlogs = append(binlogs, ml)
		}
	}
	if len(binlogs) != 0 {
		md, _ := metadata.FromIncomingContext(ctx)
		logEntry := &binarylog.ClientHeader{
			Header:     md,
			MethodName: stream.Method(),
			PeerAddr:   nil,
		}
		if deadline, ok := ctx.Deadline(); ok {
			logEntry.Timeout = time.Until(deadline)
			if logEntry.Timeout < 0 {
				logEntry.Timeout = 0
			}
		}
		if a := md[":authority"]; len(a) > 0 {
			logEntry.Authority = a[0]
		}
		if peer, ok := peer.FromContext(ctx); ok {
			logEntry.PeerAddr = peer.Addr
		}
		for _, binlog := range binlogs {
			binlog.Log(ctx, logEntry)
		}
	}

	
	
	
	
	var comp, decomp encoding.Compressor
	var cp Compressor
	var dc Decompressor
	var sendCompressorName string

	
	
	if rc := stream.RecvCompress(); s.opts.dc != nil && s.opts.dc.Type() == rc {
		dc = s.opts.dc
	} else if rc != "" && rc != encoding.Identity {
		decomp = encoding.GetCompressor(rc)
		if decomp == nil {
			st := status.Newf(codes.Unimplemented, "grpc: Decompressor is not installed for grpc-encoding %q", rc)
			stream.WriteStatus(st)
			return st.Err()
		}
	}

	
	
	
	
	if s.opts.cp != nil {
		cp = s.opts.cp
		sendCompressorName = cp.Type()
	} else if rc := stream.RecvCompress(); rc != "" && rc != encoding.Identity {
		
		comp = encoding.GetCompressor(rc)
		if comp != nil {
			sendCompressorName = comp.Name()
		}
	}

	if sendCompressorName != "" {
		if err := stream.SetSendCompress(sendCompressorName); err != nil {
			return status.Errorf(codes.Internal, "grpc: failed to set send compressor: %v", err)
		}
	}

	var payInfo *payloadInfo
	if len(shs) != 0 || len(binlogs) != 0 {
		payInfo = &payloadInfo{}
		defer payInfo.free()
	}

	d, err := recvAndDecompress(&parser{r: stream, bufferPool: s.opts.bufferPool}, stream, dc, s.opts.maxReceiveMessageSize, payInfo, decomp, true)
	if err != nil {
		if e := stream.WriteStatus(status.Convert(err)); e != nil {
			channelz.Warningf(logger, s.channelz, "grpc: Server.processUnaryRPC failed to write status: %v", e)
		}
		return err
	}
	freed := false
	dataFree := func() {
		if !freed {
			d.Free()
			freed = true
		}
	}
	defer dataFree()
	df := func(v any) error {
		defer dataFree()
		if err := s.getCodec(stream.ContentSubtype()).Unmarshal(d, v); err != nil {
			return status.Errorf(codes.Internal, "grpc: error unmarshalling request: %v", err)
		}

		for _, sh := range shs {
			sh.HandleRPC(ctx, &stats.InPayload{
				RecvTime:         time.Now(),
				Payload:          v,
				Length:           d.Len(),
				WireLength:       payInfo.compressedLength + headerLen,
				CompressedLength: payInfo.compressedLength,
			})
		}
		if len(binlogs) != 0 {
			cm := &binarylog.ClientMessage{
				Message: d.Materialize(),
			}
			for _, binlog := range binlogs {
				binlog.Log(ctx, cm)
			}
		}
		if trInfo != nil {
			trInfo.tr.LazyLog(&payload{sent: false, msg: v}, true)
		}
		return nil
	}
	ctx = NewContextWithServerTransportStream(ctx, stream)
	reply, appErr := md.Handler(info.serviceImpl, ctx, df, s.opts.unaryInt)
	if appErr != nil {
		appStatus, ok := status.FromError(appErr)
		if !ok {
			
			
			appStatus = status.FromContextError(appErr)
			appErr = appStatus.Err()
		}
		if trInfo != nil {
			trInfo.tr.LazyLog(stringer(appStatus.Message()), true)
			trInfo.tr.SetError()
		}
		if e := stream.WriteStatus(appStatus); e != nil {
			channelz.Warningf(logger, s.channelz, "grpc: Server.processUnaryRPC failed to write status: %v", e)
		}
		if len(binlogs) != 0 {
			if h, _ := stream.Header(); h.Len() > 0 {
				
				
				sh := &binarylog.ServerHeader{
					Header: h,
				}
				for _, binlog := range binlogs {
					binlog.Log(ctx, sh)
				}
			}
			st := &binarylog.ServerTrailer{
				Trailer: stream.Trailer(),
				Err:     appErr,
			}
			for _, binlog := range binlogs {
				binlog.Log(ctx, st)
			}
		}
		return appErr
	}
	if trInfo != nil {
		trInfo.tr.LazyLog(stringer("OK"), false)
	}
	opts := &transport.WriteOptions{Last: true}

	
	
	if stream.SendCompress() != sendCompressorName {
		comp = encoding.GetCompressor(stream.SendCompress())
	}
	if err := s.sendResponse(ctx, stream, reply, cp, opts, comp); err != nil {
		if err == io.EOF {
			
			return err
		}
		if sts, ok := status.FromError(err); ok {
			if e := stream.WriteStatus(sts); e != nil {
				channelz.Warningf(logger, s.channelz, "grpc: Server.processUnaryRPC failed to write status: %v", e)
			}
		} else {
			switch st := err.(type) {
			case transport.ConnectionError:
				
			default:
				panic(fmt.Sprintf("grpc: Unexpected error (%T) from sendResponse: %v", st, st))
			}
		}
		if len(binlogs) != 0 {
			h, _ := stream.Header()
			sh := &binarylog.ServerHeader{
				Header: h,
			}
			st := &binarylog.ServerTrailer{
				Trailer: stream.Trailer(),
				Err:     appErr,
			}
			for _, binlog := range binlogs {
				binlog.Log(ctx, sh)
				binlog.Log(ctx, st)
			}
		}
		return err
	}
	if len(binlogs) != 0 {
		h, _ := stream.Header()
		sh := &binarylog.ServerHeader{
			Header: h,
		}
		sm := &binarylog.ServerMessage{
			Message: reply,
		}
		for _, binlog := range binlogs {
			binlog.Log(ctx, sh)
			binlog.Log(ctx, sm)
		}
	}
	if trInfo != nil {
		trInfo.tr.LazyLog(&payload{sent: true, msg: reply}, true)
	}
	
	
	
	if len(binlogs) != 0 {
		st := &binarylog.ServerTrailer{
			Trailer: stream.Trailer(),
			Err:     appErr,
		}
		for _, binlog := range binlogs {
			binlog.Log(ctx, st)
		}
	}
	return stream.WriteStatus(statusOK)
}


func chainStreamServerInterceptors(s *Server) {
	
	
	interceptors := s.opts.chainStreamInts
	if s.opts.streamInt != nil {
		interceptors = append([]StreamServerInterceptor{s.opts.streamInt}, s.opts.chainStreamInts...)
	}

	var chainedInt StreamServerInterceptor
	if len(interceptors) == 0 {
		chainedInt = nil
	} else if len(interceptors) == 1 {
		chainedInt = interceptors[0]
	} else {
		chainedInt = chainStreamInterceptors(interceptors)
	}

	s.opts.streamInt = chainedInt
}

func chainStreamInterceptors(interceptors []StreamServerInterceptor) StreamServerInterceptor {
	return func(srv any, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error {
		return interceptors[0](srv, ss, info, getChainStreamHandler(interceptors, 0, info, handler))
	}
}

func getChainStreamHandler(interceptors []StreamServerInterceptor, curr int, info *StreamServerInfo, finalHandler StreamHandler) StreamHandler {
	if curr == len(interceptors)-1 {
		return finalHandler
	}
	return func(srv any, stream ServerStream) error {
		return interceptors[curr+1](srv, stream, info, getChainStreamHandler(interceptors, curr+1, info, finalHandler))
	}
}

func (s *Server) processStreamingRPC(ctx context.Context, stream *transport.ServerStream, info *serviceInfo, sd *StreamDesc, trInfo *traceInfo) (err error) {
	if channelz.IsOn() {
		s.incrCallsStarted()
	}
	shs := s.opts.statsHandlers
	var statsBegin *stats.Begin
	if len(shs) != 0 {
		beginTime := time.Now()
		statsBegin = &stats.Begin{
			BeginTime:      beginTime,
			IsClientStream: sd.ClientStreams,
			IsServerStream: sd.ServerStreams,
		}
		for _, sh := range shs {
			sh.HandleRPC(ctx, statsBegin)
		}
	}
	ctx = NewContextWithServerTransportStream(ctx, stream)
	ss := &serverStream{
		ctx:                   ctx,
		s:                     stream,
		p:                     &parser{r: stream, bufferPool: s.opts.bufferPool},
		codec:                 s.getCodec(stream.ContentSubtype()),
		maxReceiveMessageSize: s.opts.maxReceiveMessageSize,
		maxSendMessageSize:    s.opts.maxSendMessageSize,
		trInfo:                trInfo,
		statsHandler:          shs,
	}

	if len(shs) != 0 || trInfo != nil || channelz.IsOn() {
		
		defer func() {
			if trInfo != nil {
				ss.mu.Lock()
				if err != nil && err != io.EOF {
					ss.trInfo.tr.LazyLog(&fmtStringer{"%v", []any{err}}, true)
					ss.trInfo.tr.SetError()
				}
				ss.trInfo.tr.Finish()
				ss.trInfo.tr = nil
				ss.mu.Unlock()
			}

			if len(shs) != 0 {
				end := &stats.End{
					BeginTime: statsBegin.BeginTime,
					EndTime:   time.Now(),
				}
				if err != nil && err != io.EOF {
					end.Error = toRPCErr(err)
				}
				for _, sh := range shs {
					sh.HandleRPC(ctx, end)
				}
			}

			if channelz.IsOn() {
				if err != nil && err != io.EOF {
					s.incrCallsFailed()
				} else {
					s.incrCallsSucceeded()
				}
			}
		}()
	}

	if ml := binarylog.GetMethodLogger(stream.Method()); ml != nil {
		ss.binlogs = append(ss.binlogs, ml)
	}
	if s.opts.binaryLogger != nil {
		if ml := s.opts.binaryLogger.GetMethodLogger(stream.Method()); ml != nil {
			ss.binlogs = append(ss.binlogs, ml)
		}
	}
	if len(ss.binlogs) != 0 {
		md, _ := metadata.FromIncomingContext(ctx)
		logEntry := &binarylog.ClientHeader{
			Header:     md,
			MethodName: stream.Method(),
			PeerAddr:   nil,
		}
		if deadline, ok := ctx.Deadline(); ok {
			logEntry.Timeout = time.Until(deadline)
			if logEntry.Timeout < 0 {
				logEntry.Timeout = 0
			}
		}
		if a := md[":authority"]; len(a) > 0 {
			logEntry.Authority = a[0]
		}
		if peer, ok := peer.FromContext(ss.Context()); ok {
			logEntry.PeerAddr = peer.Addr
		}
		for _, binlog := range ss.binlogs {
			binlog.Log(ctx, logEntry)
		}
	}

	
	
	if rc := stream.RecvCompress(); s.opts.dc != nil && s.opts.dc.Type() == rc {
		ss.decompressorV0 = s.opts.dc
	} else if rc != "" && rc != encoding.Identity {
		ss.decompressorV1 = encoding.GetCompressor(rc)
		if ss.decompressorV1 == nil {
			st := status.Newf(codes.Unimplemented, "grpc: Decompressor is not installed for grpc-encoding %q", rc)
			ss.s.WriteStatus(st)
			return st.Err()
		}
	}

	
	
	
	
	if s.opts.cp != nil {
		ss.compressorV0 = s.opts.cp
		ss.sendCompressorName = s.opts.cp.Type()
	} else if rc := stream.RecvCompress(); rc != "" && rc != encoding.Identity {
		
		ss.compressorV1 = encoding.GetCompressor(rc)
		if ss.compressorV1 != nil {
			ss.sendCompressorName = rc
		}
	}

	if ss.sendCompressorName != "" {
		if err := stream.SetSendCompress(ss.sendCompressorName); err != nil {
			return status.Errorf(codes.Internal, "grpc: failed to set send compressor: %v", err)
		}
	}

	ss.ctx = newContextWithRPCInfo(ss.ctx, false, ss.codec, ss.compressorV0, ss.compressorV1)

	if trInfo != nil {
		trInfo.tr.LazyLog(&trInfo.firstLine, false)
	}
	var appErr error
	var server any
	if info != nil {
		server = info.serviceImpl
	}
	if s.opts.streamInt == nil {
		appErr = sd.Handler(server, ss)
	} else {
		info := &StreamServerInfo{
			FullMethod:     stream.Method(),
			IsClientStream: sd.ClientStreams,
			IsServerStream: sd.ServerStreams,
		}
		appErr = s.opts.streamInt(server, ss, info, sd.Handler)
	}
	if appErr != nil {
		appStatus, ok := status.FromError(appErr)
		if !ok {
			
			
			appStatus = status.FromContextError(appErr)
			appErr = appStatus.Err()
		}
		if trInfo != nil {
			ss.mu.Lock()
			ss.trInfo.tr.LazyLog(stringer(appStatus.Message()), true)
			ss.trInfo.tr.SetError()
			ss.mu.Unlock()
		}
		if len(ss.binlogs) != 0 {
			st := &binarylog.ServerTrailer{
				Trailer: ss.s.Trailer(),
				Err:     appErr,
			}
			for _, binlog := range ss.binlogs {
				binlog.Log(ctx, st)
			}
		}
		ss.s.WriteStatus(appStatus)
		
		return appErr
	}
	if trInfo != nil {
		ss.mu.Lock()
		ss.trInfo.tr.LazyLog(stringer("OK"), false)
		ss.mu.Unlock()
	}
	if len(ss.binlogs) != 0 {
		st := &binarylog.ServerTrailer{
			Trailer: ss.s.Trailer(),
			Err:     appErr,
		}
		for _, binlog := range ss.binlogs {
			binlog.Log(ctx, st)
		}
	}
	return ss.s.WriteStatus(statusOK)
}

func (s *Server) handleStream(t transport.ServerTransport, stream *transport.ServerStream) {
	ctx := stream.Context()
	ctx = contextWithServer(ctx, s)
	var ti *traceInfo
	if EnableTracing {
		tr := newTrace("grpc.Recv."+methodFamily(stream.Method()), stream.Method())
		ctx = newTraceContext(ctx, tr)
		ti = &traceInfo{
			tr: tr,
			firstLine: firstLine{
				client:     false,
				remoteAddr: t.Peer().Addr,
			},
		}
		if dl, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(dl)
		}
	}

	sm := stream.Method()
	if sm != "" && sm[0] == '/' {
		sm = sm[1:]
	}
	pos := strings.LastIndex(sm, "/")
	if pos == -1 {
		if ti != nil {
			ti.tr.LazyLog(&fmtStringer{"Malformed method name %q", []any{sm}}, true)
			ti.tr.SetError()
		}
		errDesc := fmt.Sprintf("malformed method name: %q", stream.Method())
		if err := stream.WriteStatus(status.New(codes.Unimplemented, errDesc)); err != nil {
			if ti != nil {
				ti.tr.LazyLog(&fmtStringer{"%v", []any{err}}, true)
				ti.tr.SetError()
			}
			channelz.Warningf(logger, s.channelz, "grpc: Server.handleStream failed to write status: %v", err)
		}
		if ti != nil {
			ti.tr.Finish()
		}
		return
	}
	service := sm[:pos]
	method := sm[pos+1:]

	
	if len(s.opts.statsHandlers) > 0 {
		md, _ := metadata.FromIncomingContext(ctx)
		for _, sh := range s.opts.statsHandlers {
			ctx = sh.TagRPC(ctx, &stats.RPCTagInfo{FullMethodName: stream.Method()})
			sh.HandleRPC(ctx, &stats.InHeader{
				FullMethod:  stream.Method(),
				RemoteAddr:  t.Peer().Addr,
				LocalAddr:   t.Peer().LocalAddr,
				Compression: stream.RecvCompress(),
				WireLength:  stream.HeaderWireLength(),
				Header:      md,
			})
		}
	}
	
	
	stream.SetContext(ctx)

	srv, knownService := s.services[service]
	if knownService {
		if md, ok := srv.methods[method]; ok {
			s.processUnaryRPC(ctx, stream, srv, md, ti)
			return
		}
		if sd, ok := srv.streams[method]; ok {
			s.processStreamingRPC(ctx, stream, srv, sd, ti)
			return
		}
	}
	
	if unknownDesc := s.opts.unknownStreamDesc; unknownDesc != nil {
		s.processStreamingRPC(ctx, stream, nil, unknownDesc, ti)
		return
	}
	var errDesc string
	if !knownService {
		errDesc = fmt.Sprintf("unknown service %v", service)
	} else {
		errDesc = fmt.Sprintf("unknown method %v for service %v", method, service)
	}
	if ti != nil {
		ti.tr.LazyPrintf("%s", errDesc)
		ti.tr.SetError()
	}
	if err := stream.WriteStatus(status.New(codes.Unimplemented, errDesc)); err != nil {
		if ti != nil {
			ti.tr.LazyLog(&fmtStringer{"%v", []any{err}}, true)
			ti.tr.SetError()
		}
		channelz.Warningf(logger, s.channelz, "grpc: Server.handleStream failed to write status: %v", err)
	}
	if ti != nil {
		ti.tr.Finish()
	}
}


type streamKey struct{}








func NewContextWithServerTransportStream(ctx context.Context, stream ServerTransportStream) context.Context {
	return context.WithValue(ctx, streamKey{}, stream)
}












type ServerTransportStream interface {
	Method() string
	SetHeader(md metadata.MD) error
	SendHeader(md metadata.MD) error
	SetTrailer(md metadata.MD) error
}









func ServerTransportStreamFromContext(ctx context.Context) ServerTransportStream {
	s, _ := ctx.Value(streamKey{}).(ServerTransportStream)
	return s
}






func (s *Server) Stop() {
	s.stop(false)
}




func (s *Server) GracefulStop() {
	s.stop(true)
}

func (s *Server) stop(graceful bool) {
	s.quit.Fire()
	defer s.done.Fire()

	s.channelzRemoveOnce.Do(func() { channelz.RemoveEntry(s.channelz.ID) })
	s.mu.Lock()
	s.closeListenersLocked()
	
	
	s.mu.Unlock()
	s.serveWG.Wait()

	s.mu.Lock()
	defer s.mu.Unlock()

	if graceful {
		s.drainAllServerTransportsLocked()
	} else {
		s.closeServerTransportsLocked()
	}

	for len(s.conns) != 0 {
		s.cv.Wait()
	}
	s.conns = nil

	if s.opts.numServerWorkers > 0 {
		
		
		
		
		s.serverWorkerChannelClose()
	}

	if graceful || s.opts.waitForHandlers {
		s.handlersWG.Wait()
	}

	if s.events != nil {
		s.events.Finish()
		s.events = nil
	}
}


func (s *Server) closeServerTransportsLocked() {
	for _, conns := range s.conns {
		for st := range conns {
			st.Close(errors.New("Server.Stop called"))
		}
	}
}


func (s *Server) drainAllServerTransportsLocked() {
	if !s.drain {
		for _, conns := range s.conns {
			for st := range conns {
				st.Drain("graceful_stop")
			}
		}
		s.drain = true
	}
}


func (s *Server) closeListenersLocked() {
	for lis := range s.lis {
		lis.Close()
	}
	s.lis = nil
}



func (s *Server) getCodec(contentSubtype string) baseCodec {
	if s.opts.codec != nil {
		return s.opts.codec
	}
	if contentSubtype == "" {
		return getCodec(proto.Name)
	}
	codec := getCodec(contentSubtype)
	if codec == nil {
		logger.Warningf("Unsupported codec %q. Defaulting to %q for now. This will start to fail in future releases.", contentSubtype, proto.Name)
		return getCodec(proto.Name)
	}
	return codec
}

type serverKey struct{}


func serverFromContext(ctx context.Context) *Server {
	s, _ := ctx.Value(serverKey{}).(*Server)
	return s
}


func contextWithServer(ctx context.Context, server *Server) context.Context {
	return context.WithValue(ctx, serverKey{}, server)
}




func (s *Server) isRegisteredMethod(serviceMethod string) bool {
	if serviceMethod != "" && serviceMethod[0] == '/' {
		serviceMethod = serviceMethod[1:]
	}
	pos := strings.LastIndex(serviceMethod, "/")
	if pos == -1 { 
		return false
	}
	service := serviceMethod[:pos]
	method := serviceMethod[pos+1:]
	srv, knownService := s.services[service]
	if knownService {
		if _, ok := srv.methods[method]; ok {
			return true
		}
		if _, ok := srv.streams[method]; ok {
			return true
		}
	}
	return false
}





















func SetHeader(ctx context.Context, md metadata.MD) error {
	if md.Len() == 0 {
		return nil
	}
	stream := ServerTransportStreamFromContext(ctx)
	if stream == nil {
		return status.Errorf(codes.Internal, "grpc: failed to fetch the stream from the context %v", ctx)
	}
	return stream.SetHeader(md)
}









func SendHeader(ctx context.Context, md metadata.MD) error {
	stream := ServerTransportStreamFromContext(ctx)
	if stream == nil {
		return status.Errorf(codes.Internal, "grpc: failed to fetch the stream from the context %v", ctx)
	}
	if err := stream.SendHeader(md); err != nil {
		return toRPCErr(err)
	}
	return nil
}
























func SetSendCompressor(ctx context.Context, name string) error {
	stream, ok := ServerTransportStreamFromContext(ctx).(*transport.ServerStream)
	if !ok || stream == nil {
		return fmt.Errorf("failed to fetch the stream from the given context")
	}

	if err := validateSendCompressor(name, stream.ClientAdvertisedCompressors()); err != nil {
		return fmt.Errorf("unable to set send compressor: %w", err)
	}

	return stream.SetSendCompress(name)
}










func ClientSupportedCompressors(ctx context.Context) ([]string, error) {
	stream, ok := ServerTransportStreamFromContext(ctx).(*transport.ServerStream)
	if !ok || stream == nil {
		return nil, fmt.Errorf("failed to fetch the stream from the given context %v", ctx)
	}

	return stream.ClientAdvertisedCompressors(), nil
}







func SetTrailer(ctx context.Context, md metadata.MD) error {
	if md.Len() == 0 {
		return nil
	}
	stream := ServerTransportStreamFromContext(ctx)
	if stream == nil {
		return status.Errorf(codes.Internal, "grpc: failed to fetch the stream from the context %v", ctx)
	}
	return stream.SetTrailer(md)
}



func Method(ctx context.Context) (string, bool) {
	s := ServerTransportStreamFromContext(ctx)
	if s == nil {
		return "", false
	}
	return s.Method(), true
}



func validateSendCompressor(name string, clientCompressors []string) error {
	if name == encoding.Identity {
		return nil
	}

	if !grpcutil.IsCompressorNameRegistered(name) {
		return fmt.Errorf("compressor not registered %q", name)
	}

	for _, c := range clientCompressors {
		if c == name {
			return nil 
		}
	}
	return fmt.Errorf("client does not support compressor %q", name)
}



type atomicSemaphore struct {
	n    atomic.Int64
	wait chan struct{}
}

func (q *atomicSemaphore) acquire() {
	if q.n.Add(-1) < 0 {
		
		<-q.wait
	}
}

func (q *atomicSemaphore) release() {
	
	
	
	
	if q.n.Add(1) <= 0 {
		
		q.wait <- struct{}{}
	}
}

func newHandlerQuota(n uint32) *atomicSemaphore {
	a := &atomicSemaphore{wait: make(chan struct{}, 1)}
	a.n.Store(int64(n))
	return a
}
