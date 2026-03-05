

package transport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	rand "math/rand/v2"
	"net"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	"google.golang.org/grpc/internal/grpclog"
	"google.golang.org/grpc/internal/grpcutil"
	"google.golang.org/grpc/internal/pretty"
	"google.golang.org/grpc/internal/syscall"
	"google.golang.org/grpc/mem"
	"google.golang.org/protobuf/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/grpcsync"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
)

var (
	
	
	ErrIllegalHeaderWrite = status.Error(codes.Internal, "transport: SendHeader called multiple times")
	
	
	ErrHeaderListSizeLimitViolation = status.Error(codes.Internal, "transport: trying to send header list size larger than the limit set by peer")
)



var serverConnectionCounter uint64


type http2Server struct {
	lastRead        int64 
	done            chan struct{}
	conn            net.Conn
	loopy           *loopyWriter
	readerDone      chan struct{} 
	loopyWriterDone chan struct{}
	peer            peer.Peer
	inTapHandle     tap.ServerInHandle
	framer          *framer
	
	maxStreams uint32
	
	
	controlBuf *controlBuffer
	fc         *trInFlow
	stats      []stats.Handler
	
	kp keepalive.ServerParameters
	
	kep keepalive.EnforcementPolicy
	
	lastPingAt time.Time
	
	pingStrikes uint8
	
	
	
	resetPingStrikes      uint32 
	initialWindowSize     int32
	bdpEst                *bdpEstimator
	maxSendHeaderListSize *uint32

	mu sync.Mutex 

	
	
	
	
	
	
	drainEvent    *grpcsync.Event
	state         transportState
	activeStreams map[uint32]*ServerStream
	
	
	
	
	idle time.Time

	
	channelz   *channelz.Socket
	bufferPool mem.BufferPool

	connectionID uint64

	
	
	maxStreamMu sync.Mutex
	maxStreamID uint32 

	logger *grpclog.PrefixLogger
}








func NewServerTransport(conn net.Conn, config *ServerConfig) (_ ServerTransport, err error) {
	var authInfo credentials.AuthInfo
	rawConn := conn
	if config.Credentials != nil {
		var err error
		conn, authInfo, err = config.Credentials.ServerHandshake(rawConn)
		if err != nil {
			
			
			
			
			if err == credentials.ErrConnDispatched || err == io.EOF {
				return nil, err
			}
			return nil, connectionErrorf(false, err, "ServerHandshake(%q) failed: %v", rawConn.RemoteAddr(), err)
		}
	}
	writeBufSize := config.WriteBufferSize
	readBufSize := config.ReadBufferSize
	maxHeaderListSize := defaultServerMaxHeaderListSize
	if config.MaxHeaderListSize != nil {
		maxHeaderListSize = *config.MaxHeaderListSize
	}
	framer := newFramer(conn, writeBufSize, readBufSize, config.SharedWriteBuffer, maxHeaderListSize)
	
	isettings := []http2.Setting{{
		ID:  http2.SettingMaxFrameSize,
		Val: http2MaxFrameLen,
	}}
	if config.MaxStreams != math.MaxUint32 {
		isettings = append(isettings, http2.Setting{
			ID:  http2.SettingMaxConcurrentStreams,
			Val: config.MaxStreams,
		})
	}
	dynamicWindow := true
	iwz := int32(initialWindowSize)
	if config.InitialWindowSize >= defaultWindowSize {
		iwz = config.InitialWindowSize
		dynamicWindow = false
	}
	icwz := int32(initialWindowSize)
	if config.InitialConnWindowSize >= defaultWindowSize {
		icwz = config.InitialConnWindowSize
		dynamicWindow = false
	}
	if iwz != defaultWindowSize {
		isettings = append(isettings, http2.Setting{
			ID:  http2.SettingInitialWindowSize,
			Val: uint32(iwz)})
	}
	if config.MaxHeaderListSize != nil {
		isettings = append(isettings, http2.Setting{
			ID:  http2.SettingMaxHeaderListSize,
			Val: *config.MaxHeaderListSize,
		})
	}
	if config.HeaderTableSize != nil {
		isettings = append(isettings, http2.Setting{
			ID:  http2.SettingHeaderTableSize,
			Val: *config.HeaderTableSize,
		})
	}
	if err := framer.fr.WriteSettings(isettings...); err != nil {
		return nil, connectionErrorf(false, err, "transport: %v", err)
	}
	
	if delta := uint32(icwz - defaultWindowSize); delta > 0 {
		if err := framer.fr.WriteWindowUpdate(0, delta); err != nil {
			return nil, connectionErrorf(false, err, "transport: %v", err)
		}
	}
	kp := config.KeepaliveParams
	if kp.MaxConnectionIdle == 0 {
		kp.MaxConnectionIdle = defaultMaxConnectionIdle
	}
	if kp.MaxConnectionAge == 0 {
		kp.MaxConnectionAge = defaultMaxConnectionAge
	}
	
	kp.MaxConnectionAge += getJitter(kp.MaxConnectionAge)
	if kp.MaxConnectionAgeGrace == 0 {
		kp.MaxConnectionAgeGrace = defaultMaxConnectionAgeGrace
	}
	if kp.Time == 0 {
		kp.Time = defaultServerKeepaliveTime
	}
	if kp.Timeout == 0 {
		kp.Timeout = defaultServerKeepaliveTimeout
	}
	if kp.Time != infinity {
		if err = syscall.SetTCPUserTimeout(rawConn, kp.Timeout); err != nil {
			return nil, connectionErrorf(false, err, "transport: failed to set TCP_USER_TIMEOUT: %v", err)
		}
	}
	kep := config.KeepalivePolicy
	if kep.MinTime == 0 {
		kep.MinTime = defaultKeepalivePolicyMinTime
	}

	done := make(chan struct{})
	peer := peer.Peer{
		Addr:      conn.RemoteAddr(),
		LocalAddr: conn.LocalAddr(),
		AuthInfo:  authInfo,
	}
	t := &http2Server{
		done:              done,
		conn:              conn,
		peer:              peer,
		framer:            framer,
		readerDone:        make(chan struct{}),
		loopyWriterDone:   make(chan struct{}),
		maxStreams:        config.MaxStreams,
		inTapHandle:       config.InTapHandle,
		fc:                &trInFlow{limit: uint32(icwz)},
		state:             reachable,
		activeStreams:     make(map[uint32]*ServerStream),
		stats:             config.StatsHandlers,
		kp:                kp,
		idle:              time.Now(),
		kep:               kep,
		initialWindowSize: iwz,
		bufferPool:        config.BufferPool,
	}
	var czSecurity credentials.ChannelzSecurityValue
	if au, ok := authInfo.(credentials.ChannelzSecurityInfo); ok {
		czSecurity = au.GetSecurityValue()
	}
	t.channelz = channelz.RegisterSocket(
		&channelz.Socket{
			SocketType:       channelz.SocketTypeNormal,
			Parent:           config.ChannelzParent,
			SocketMetrics:    channelz.SocketMetrics{},
			EphemeralMetrics: t.socketMetrics,
			LocalAddr:        t.peer.LocalAddr,
			RemoteAddr:       t.peer.Addr,
			SocketOptions:    channelz.GetSocketOption(t.conn),
			Security:         czSecurity,
		},
	)
	t.logger = prefixLoggerForServerTransport(t)

	t.controlBuf = newControlBuffer(t.done)
	if dynamicWindow {
		t.bdpEst = &bdpEstimator{
			bdp:               initialWindowSize,
			updateFlowControl: t.updateFlowControl,
		}
	}

	t.connectionID = atomic.AddUint64(&serverConnectionCounter, 1)
	t.framer.writer.Flush()

	defer func() {
		if err != nil {
			t.Close(err)
		}
	}()

	
	preface := make([]byte, len(clientPreface))
	if _, err := io.ReadFull(t.conn, preface); err != nil {
		
		
		
		
		
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, connectionErrorf(false, err, "transport: http2Server.HandleStreams failed to receive the preface from client: %v", err)
	}
	if !bytes.Equal(preface, clientPreface) {
		return nil, connectionErrorf(false, nil, "transport: http2Server.HandleStreams received bogus greeting from client: %q", preface)
	}

	frame, err := t.framer.fr.ReadFrame()
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return nil, err
	}
	if err != nil {
		return nil, connectionErrorf(false, err, "transport: http2Server.HandleStreams failed to read initial settings frame: %v", err)
	}
	atomic.StoreInt64(&t.lastRead, time.Now().UnixNano())
	sf, ok := frame.(*http2.SettingsFrame)
	if !ok {
		return nil, connectionErrorf(false, nil, "transport: http2Server.HandleStreams saw invalid preface type %T from client", frame)
	}
	t.handleSettings(sf)

	go func() {
		t.loopy = newLoopyWriter(serverSide, t.framer, t.controlBuf, t.bdpEst, t.conn, t.logger, t.outgoingGoAwayHandler, t.bufferPool)
		err := t.loopy.run()
		close(t.loopyWriterDone)
		if !isIOError(err) {
			
			
			
			
			
			
			
			
			
			timer := time.NewTimer(time.Second)
			defer timer.Stop()
			select {
			case <-t.readerDone:
			case <-timer.C:
			}
			t.conn.Close()
		}
	}()
	go t.keepalive()
	return t, nil
}



func (t *http2Server) operateHeaders(ctx context.Context, frame *http2.MetaHeadersFrame, handle func(*ServerStream)) error {
	
	t.maxStreamMu.Lock()
	defer t.maxStreamMu.Unlock()

	streamID := frame.Header().StreamID

	
	
	if frame.Truncated {
		t.controlBuf.put(&cleanupStream{
			streamID: streamID,
			rst:      true,
			rstCode:  http2.ErrCodeFrameSize,
			onWrite:  func() {},
		})
		return nil
	}

	if streamID%2 != 1 || streamID <= t.maxStreamID {
		
		return fmt.Errorf("received an illegal stream id: %v. headers frame: %+v", streamID, frame)
	}
	t.maxStreamID = streamID

	buf := newRecvBuffer()
	s := &ServerStream{
		Stream: &Stream{
			id:  streamID,
			buf: buf,
			fc:  &inFlow{limit: uint32(t.initialWindowSize)},
		},
		st:               t,
		headerWireLength: int(frame.Header().Length),
	}
	var (
		
		isGRPC      = false
		contentType = ""
		mdata       = make(metadata.MD, len(frame.Fields))
		httpMethod  string
		
		protocolError bool
		headerError   *status.Status

		timeoutSet bool
		timeout    time.Duration
	)

	for _, hf := range frame.Fields {
		switch hf.Name {
		case "content-type":
			contentSubtype, validContentType := grpcutil.ContentSubtype(hf.Value)
			if !validContentType {
				contentType = hf.Value
				break
			}
			mdata[hf.Name] = append(mdata[hf.Name], hf.Value)
			s.contentSubtype = contentSubtype
			isGRPC = true

		case "grpc-accept-encoding":
			mdata[hf.Name] = append(mdata[hf.Name], hf.Value)
			if hf.Value == "" {
				continue
			}
			compressors := hf.Value
			if s.clientAdvertisedCompressors != "" {
				compressors = s.clientAdvertisedCompressors + "," + compressors
			}
			s.clientAdvertisedCompressors = compressors
		case "grpc-encoding":
			s.recvCompress = hf.Value
		case ":method":
			httpMethod = hf.Value
		case ":path":
			s.method = hf.Value
		case "grpc-timeout":
			timeoutSet = true
			var err error
			if timeout, err = decodeTimeout(hf.Value); err != nil {
				headerError = status.Newf(codes.Internal, "malformed grpc-timeout: %v", err)
			}
		
		
		case "connection":
			if t.logger.V(logLevel) {
				t.logger.Infof("Received a HEADERS frame with a :connection header which makes the request malformed, as per the HTTP/2 spec")
			}
			protocolError = true
		default:
			if isReservedHeader(hf.Name) && !isWhitelistedHeader(hf.Name) {
				break
			}
			v, err := decodeMetadataHeader(hf.Name, hf.Value)
			if err != nil {
				headerError = status.Newf(codes.Internal, "malformed binary metadata %q in header %q: %v", hf.Value, hf.Name, err)
				t.logger.Warningf("Failed to decode metadata header (%q, %q): %v", hf.Name, hf.Value, err)
				break
			}
			mdata[hf.Name] = append(mdata[hf.Name], v)
		}
	}

	
	
	
	
	
	if len(mdata[":authority"]) > 1 || len(mdata["host"]) > 1 {
		errMsg := fmt.Sprintf("num values of :authority: %v, num values of host: %v, both must only have 1 value as per HTTP/2 spec", len(mdata[":authority"]), len(mdata["host"]))
		if t.logger.V(logLevel) {
			t.logger.Infof("Aborting the stream early: %v", errMsg)
		}
		t.controlBuf.put(&earlyAbortStream{
			httpStatus:     http.StatusBadRequest,
			streamID:       streamID,
			contentSubtype: s.contentSubtype,
			status:         status.New(codes.Internal, errMsg),
			rst:            !frame.StreamEnded(),
		})
		return nil
	}

	if protocolError {
		t.controlBuf.put(&cleanupStream{
			streamID: streamID,
			rst:      true,
			rstCode:  http2.ErrCodeProtocol,
			onWrite:  func() {},
		})
		return nil
	}
	if !isGRPC {
		t.controlBuf.put(&earlyAbortStream{
			httpStatus:     http.StatusUnsupportedMediaType,
			streamID:       streamID,
			contentSubtype: s.contentSubtype,
			status:         status.Newf(codes.InvalidArgument, "invalid gRPC request content-type %q", contentType),
			rst:            !frame.StreamEnded(),
		})
		return nil
	}
	if headerError != nil {
		t.controlBuf.put(&earlyAbortStream{
			httpStatus:     http.StatusBadRequest,
			streamID:       streamID,
			contentSubtype: s.contentSubtype,
			status:         headerError,
			rst:            !frame.StreamEnded(),
		})
		return nil
	}

	
	if len(mdata[":authority"]) == 0 {
		
		
		if host, ok := mdata["host"]; ok {
			mdata[":authority"] = host
			delete(mdata, "host")
		}
	} else {
		
		delete(mdata, "host")
	}

	if frame.StreamEnded() {
		
		s.state = streamReadDone
	}
	if timeoutSet {
		s.ctx, s.cancel = context.WithTimeout(ctx, timeout)
	} else {
		s.ctx, s.cancel = context.WithCancel(ctx)
	}

	
	if len(mdata) > 0 {
		s.ctx = metadata.NewIncomingContext(s.ctx, mdata)
	}
	t.mu.Lock()
	if t.state != reachable {
		t.mu.Unlock()
		s.cancel()
		return nil
	}
	if uint32(len(t.activeStreams)) >= t.maxStreams {
		t.mu.Unlock()
		t.controlBuf.put(&cleanupStream{
			streamID: streamID,
			rst:      true,
			rstCode:  http2.ErrCodeRefusedStream,
			onWrite:  func() {},
		})
		s.cancel()
		return nil
	}
	if httpMethod != http.MethodPost {
		t.mu.Unlock()
		errMsg := fmt.Sprintf("Received a HEADERS frame with :method %q which should be POST", httpMethod)
		if t.logger.V(logLevel) {
			t.logger.Infof("Aborting the stream early: %v", errMsg)
		}
		t.controlBuf.put(&earlyAbortStream{
			httpStatus:     http.StatusMethodNotAllowed,
			streamID:       streamID,
			contentSubtype: s.contentSubtype,
			status:         status.New(codes.Internal, errMsg),
			rst:            !frame.StreamEnded(),
		})
		s.cancel()
		return nil
	}
	if t.inTapHandle != nil {
		var err error
		if s.ctx, err = t.inTapHandle(s.ctx, &tap.Info{FullMethodName: s.method, Header: mdata}); err != nil {
			t.mu.Unlock()
			if t.logger.V(logLevel) {
				t.logger.Infof("Aborting the stream early due to InTapHandle failure: %v", err)
			}
			stat, ok := status.FromError(err)
			if !ok {
				stat = status.New(codes.PermissionDenied, err.Error())
			}
			t.controlBuf.put(&earlyAbortStream{
				httpStatus:     http.StatusOK,
				streamID:       s.id,
				contentSubtype: s.contentSubtype,
				status:         stat,
				rst:            !frame.StreamEnded(),
			})
			return nil
		}
	}
	t.activeStreams[streamID] = s
	if len(t.activeStreams) == 1 {
		t.idle = time.Time{}
	}
	t.mu.Unlock()
	if channelz.IsOn() {
		t.channelz.SocketMetrics.StreamsStarted.Add(1)
		t.channelz.SocketMetrics.LastRemoteStreamCreatedTimestamp.Store(time.Now().UnixNano())
	}
	s.requestRead = func(n int) {
		t.adjustWindow(s, uint32(n))
	}
	s.ctxDone = s.ctx.Done()
	s.wq = newWriteQuota(defaultWriteQuota, s.ctxDone)
	s.trReader = &transportReader{
		reader: &recvBufferReader{
			ctx:     s.ctx,
			ctxDone: s.ctxDone,
			recv:    s.buf,
		},
		windowHandler: func(n int) {
			t.updateWindow(s, uint32(n))
		},
	}
	
	t.controlBuf.put(&registerStream{
		streamID: s.id,
		wq:       s.wq,
	})
	handle(s)
	return nil
}




func (t *http2Server) HandleStreams(ctx context.Context, handle func(*ServerStream)) {
	defer func() {
		close(t.readerDone)
		<-t.loopyWriterDone
	}()
	for {
		t.controlBuf.throttle()
		frame, err := t.framer.fr.ReadFrame()
		atomic.StoreInt64(&t.lastRead, time.Now().UnixNano())
		if err != nil {
			if se, ok := err.(http2.StreamError); ok {
				if t.logger.V(logLevel) {
					t.logger.Warningf("Encountered http2.StreamError: %v", se)
				}
				t.mu.Lock()
				s := t.activeStreams[se.StreamID]
				t.mu.Unlock()
				if s != nil {
					t.closeStream(s, true, se.Code, false)
				} else {
					t.controlBuf.put(&cleanupStream{
						streamID: se.StreamID,
						rst:      true,
						rstCode:  se.Code,
						onWrite:  func() {},
					})
				}
				continue
			}
			t.Close(err)
			return
		}
		switch frame := frame.(type) {
		case *http2.MetaHeadersFrame:
			if err := t.operateHeaders(ctx, frame, handle); err != nil {
				
				
				t.controlBuf.put(&goAway{
					code:      http2.ErrCodeProtocol,
					debugData: []byte(err.Error()),
					closeConn: err,
				})
				continue
			}
		case *http2.DataFrame:
			t.handleData(frame)
		case *http2.RSTStreamFrame:
			t.handleRSTStream(frame)
		case *http2.SettingsFrame:
			t.handleSettings(frame)
		case *http2.PingFrame:
			t.handlePing(frame)
		case *http2.WindowUpdateFrame:
			t.handleWindowUpdate(frame)
		case *http2.GoAwayFrame:
			
		default:
			if t.logger.V(logLevel) {
				t.logger.Infof("Received unsupported frame type %T", frame)
			}
		}
	}
}

func (t *http2Server) getStream(f http2.Frame) (*ServerStream, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.activeStreams == nil {
		
		return nil, false
	}
	s, ok := t.activeStreams[f.Header().StreamID]
	if !ok {
		
		return nil, false
	}
	return s, true
}




func (t *http2Server) adjustWindow(s *ServerStream, n uint32) {
	if w := s.fc.maybeAdjust(n); w > 0 {
		t.controlBuf.put(&outgoingWindowUpdate{streamID: s.id, increment: w})
	}

}




func (t *http2Server) updateWindow(s *ServerStream, n uint32) {
	if w := s.fc.onRead(n); w > 0 {
		t.controlBuf.put(&outgoingWindowUpdate{streamID: s.id,
			increment: w,
		})
	}
}




func (t *http2Server) updateFlowControl(n uint32) {
	t.mu.Lock()
	for _, s := range t.activeStreams {
		s.fc.newLimit(n)
	}
	t.initialWindowSize = int32(n)
	t.mu.Unlock()
	t.controlBuf.put(&outgoingWindowUpdate{
		streamID:  0,
		increment: t.fc.newLimit(n),
	})
	t.controlBuf.put(&outgoingSettings{
		ss: []http2.Setting{
			{
				ID:  http2.SettingInitialWindowSize,
				Val: n,
			},
		},
	})

}

func (t *http2Server) handleData(f *http2.DataFrame) {
	size := f.Header().Length
	var sendBDPPing bool
	if t.bdpEst != nil {
		sendBDPPing = t.bdpEst.add(size)
	}
	
	
	
	
	
	
	
	
	if w := t.fc.onData(size); w > 0 {
		t.controlBuf.put(&outgoingWindowUpdate{
			streamID:  0,
			increment: w,
		})
	}
	if sendBDPPing {
		
		
		if w := t.fc.reset(); w > 0 {
			t.controlBuf.put(&outgoingWindowUpdate{
				streamID:  0,
				increment: w,
			})
		}
		t.controlBuf.put(bdpPing)
	}
	
	s, ok := t.getStream(f)
	if !ok {
		return
	}
	if s.getState() == streamReadDone {
		t.closeStream(s, true, http2.ErrCodeStreamClosed, false)
		return
	}
	if size > 0 {
		if err := s.fc.onData(size); err != nil {
			t.closeStream(s, true, http2.ErrCodeFlowControl, false)
			return
		}
		if f.Header().Flags.Has(http2.FlagDataPadded) {
			if w := s.fc.onRead(size - uint32(len(f.Data()))); w > 0 {
				t.controlBuf.put(&outgoingWindowUpdate{s.id, w})
			}
		}
		
		
		
		if len(f.Data()) > 0 {
			pool := t.bufferPool
			if pool == nil {
				
				
				pool = mem.DefaultBufferPool()
			}
			s.write(recvMsg{buffer: mem.Copy(f.Data(), pool)})
		}
	}
	if f.StreamEnded() {
		
		s.compareAndSwapState(streamActive, streamReadDone)
		s.write(recvMsg{err: io.EOF})
	}
}

func (t *http2Server) handleRSTStream(f *http2.RSTStreamFrame) {
	
	if s, ok := t.getStream(f); ok {
		t.closeStream(s, false, 0, false)
		return
	}
	
	t.controlBuf.put(&cleanupStream{
		streamID: f.Header().StreamID,
		rst:      false,
		rstCode:  0,
		onWrite:  func() {},
	})
}

func (t *http2Server) handleSettings(f *http2.SettingsFrame) {
	if f.IsAck() {
		return
	}
	var ss []http2.Setting
	var updateFuncs []func()
	f.ForeachSetting(func(s http2.Setting) error {
		switch s.ID {
		case http2.SettingMaxHeaderListSize:
			updateFuncs = append(updateFuncs, func() {
				t.maxSendHeaderListSize = new(uint32)
				*t.maxSendHeaderListSize = s.Val
			})
		default:
			ss = append(ss, s)
		}
		return nil
	})
	t.controlBuf.executeAndPut(func() bool {
		for _, f := range updateFuncs {
			f()
		}
		return true
	}, &incomingSettings{
		ss: ss,
	})
}

const (
	maxPingStrikes     = 2
	defaultPingTimeout = 2 * time.Hour
)

func (t *http2Server) handlePing(f *http2.PingFrame) {
	if f.IsAck() {
		if f.Data == goAwayPing.data && t.drainEvent != nil {
			t.drainEvent.Fire()
			return
		}
		
		if t.bdpEst != nil {
			t.bdpEst.calculate(f.Data)
		}
		return
	}
	pingAck := &ping{ack: true}
	copy(pingAck.data[:], f.Data[:])
	t.controlBuf.put(pingAck)

	now := time.Now()
	defer func() {
		t.lastPingAt = now
	}()
	
	
	
	if atomic.CompareAndSwapUint32(&t.resetPingStrikes, 1, 0) {
		t.pingStrikes = 0
		return
	}
	t.mu.Lock()
	ns := len(t.activeStreams)
	t.mu.Unlock()
	if ns < 1 && !t.kep.PermitWithoutStream {
		
		
		if t.lastPingAt.Add(defaultPingTimeout).After(now) {
			t.pingStrikes++
		}
	} else {
		
		if t.lastPingAt.Add(t.kep.MinTime).After(now) {
			t.pingStrikes++
		}
	}

	if t.pingStrikes > maxPingStrikes {
		
		t.controlBuf.put(&goAway{code: http2.ErrCodeEnhanceYourCalm, debugData: []byte("too_many_pings"), closeConn: errors.New("got too many pings from the client")})
	}
}

func (t *http2Server) handleWindowUpdate(f *http2.WindowUpdateFrame) {
	t.controlBuf.put(&incomingWindowUpdate{
		streamID:  f.Header().StreamID,
		increment: f.Increment,
	})
}

func appendHeaderFieldsFromMD(headerFields []hpack.HeaderField, md metadata.MD) []hpack.HeaderField {
	for k, vv := range md {
		if isReservedHeader(k) {
			
			continue
		}
		for _, v := range vv {
			headerFields = append(headerFields, hpack.HeaderField{Name: k, Value: encodeMetadataHeader(k, v)})
		}
	}
	return headerFields
}

func (t *http2Server) checkForHeaderListSize(it any) bool {
	if t.maxSendHeaderListSize == nil {
		return true
	}
	hdrFrame := it.(*headerFrame)
	var sz int64
	for _, f := range hdrFrame.hf {
		if sz += int64(f.Size()); sz > int64(*t.maxSendHeaderListSize) {
			if t.logger.V(logLevel) {
				t.logger.Infof("Header list size to send violates the maximum size (%d bytes) set by client", *t.maxSendHeaderListSize)
			}
			return false
		}
	}
	return true
}

func (t *http2Server) streamContextErr(s *ServerStream) error {
	select {
	case <-t.done:
		return ErrConnClosing
	default:
	}
	return ContextErr(s.ctx.Err())
}


func (t *http2Server) writeHeader(s *ServerStream, md metadata.MD) error {
	s.hdrMu.Lock()
	defer s.hdrMu.Unlock()
	if s.getState() == streamDone {
		return t.streamContextErr(s)
	}

	if s.updateHeaderSent() {
		return ErrIllegalHeaderWrite
	}

	if md.Len() > 0 {
		if s.header.Len() > 0 {
			s.header = metadata.Join(s.header, md)
		} else {
			s.header = md
		}
	}
	if err := t.writeHeaderLocked(s); err != nil {
		switch e := err.(type) {
		case ConnectionError:
			return status.Error(codes.Unavailable, e.Desc)
		default:
			return status.Convert(err).Err()
		}
	}
	return nil
}

func (t *http2Server) setResetPingStrikes() {
	atomic.StoreUint32(&t.resetPingStrikes, 1)
}

func (t *http2Server) writeHeaderLocked(s *ServerStream) error {
	
	
	headerFields := make([]hpack.HeaderField, 0, 2) 
	headerFields = append(headerFields, hpack.HeaderField{Name: ":status", Value: "200"})
	headerFields = append(headerFields, hpack.HeaderField{Name: "content-type", Value: grpcutil.ContentType(s.contentSubtype)})
	if s.sendCompress != "" {
		headerFields = append(headerFields, hpack.HeaderField{Name: "grpc-encoding", Value: s.sendCompress})
	}
	headerFields = appendHeaderFieldsFromMD(headerFields, s.header)
	hf := &headerFrame{
		streamID:  s.id,
		hf:        headerFields,
		endStream: false,
		onWrite:   t.setResetPingStrikes,
	}
	success, err := t.controlBuf.executeAndPut(func() bool { return t.checkForHeaderListSize(hf) }, hf)
	if !success {
		if err != nil {
			return err
		}
		t.closeStream(s, true, http2.ErrCodeInternal, false)
		return ErrHeaderListSizeLimitViolation
	}
	for _, sh := range t.stats {
		
		
		outHeader := &stats.OutHeader{
			Header:      s.header.Copy(),
			Compression: s.sendCompress,
		}
		sh.HandleRPC(s.Context(), outHeader)
	}
	return nil
}





func (t *http2Server) writeStatus(s *ServerStream, st *status.Status) error {
	s.hdrMu.Lock()
	defer s.hdrMu.Unlock()

	if s.getState() == streamDone {
		return nil
	}

	
	
	headerFields := make([]hpack.HeaderField, 0, 2) 
	if !s.updateHeaderSent() {                      
		if len(s.header) > 0 { 
			if err := t.writeHeaderLocked(s); err != nil {
				return err
			}
		} else { 
			headerFields = append(headerFields, hpack.HeaderField{Name: ":status", Value: "200"})
			headerFields = append(headerFields, hpack.HeaderField{Name: "content-type", Value: grpcutil.ContentType(s.contentSubtype)})
		}
	}
	headerFields = append(headerFields, hpack.HeaderField{Name: "grpc-status", Value: strconv.Itoa(int(st.Code()))})
	headerFields = append(headerFields, hpack.HeaderField{Name: "grpc-message", Value: encodeGrpcMessage(st.Message())})

	if p := st.Proto(); p != nil && len(p.Details) > 0 {
		
		
		delete(s.trailer, grpcStatusDetailsBinHeader)
		stBytes, err := proto.Marshal(p)
		if err != nil {
			
			t.logger.Errorf("Failed to marshal rpc status: %s, error: %v", pretty.ToJSON(p), err)
		} else {
			headerFields = append(headerFields, hpack.HeaderField{Name: grpcStatusDetailsBinHeader, Value: encodeBinHeader(stBytes)})
		}
	}

	
	headerFields = appendHeaderFieldsFromMD(headerFields, s.trailer)
	trailingHeader := &headerFrame{
		streamID:  s.id,
		hf:        headerFields,
		endStream: true,
		onWrite:   t.setResetPingStrikes,
	}

	success, err := t.controlBuf.executeAndPut(func() bool {
		return t.checkForHeaderListSize(trailingHeader)
	}, nil)
	if !success {
		if err != nil {
			return err
		}
		t.closeStream(s, true, http2.ErrCodeInternal, false)
		return ErrHeaderListSizeLimitViolation
	}
	
	rst := s.getState() == streamActive
	t.finishStream(s, rst, http2.ErrCodeNo, trailingHeader, true)
	for _, sh := range t.stats {
		
		
		sh.HandleRPC(s.Context(), &stats.OutTrailer{
			Trailer: s.trailer.Copy(),
		})
	}
	return nil
}



func (t *http2Server) write(s *ServerStream, hdr []byte, data mem.BufferSlice, _ *WriteOptions) error {
	reader := data.Reader()

	if !s.isHeaderSent() { 
		if err := t.writeHeader(s, nil); err != nil {
			_ = reader.Close()
			return err
		}
	} else {
		
		if s.getState() == streamDone {
			_ = reader.Close()
			return t.streamContextErr(s)
		}
	}

	df := &dataFrame{
		streamID:    s.id,
		h:           hdr,
		reader:      reader,
		onEachWrite: t.setResetPingStrikes,
	}
	if err := s.wq.get(int32(len(hdr) + df.reader.Remaining())); err != nil {
		_ = reader.Close()
		return t.streamContextErr(s)
	}
	if err := t.controlBuf.put(df); err != nil {
		_ = reader.Close()
		return err
	}
	t.incrMsgSent()
	return nil
}







func (t *http2Server) keepalive() {
	p := &ping{}
	
	outstandingPing := false
	
	
	kpTimeoutLeft := time.Duration(0)
	
	
	prevNano := time.Now().UnixNano()
	
	idleTimer := time.NewTimer(t.kp.MaxConnectionIdle)
	ageTimer := time.NewTimer(t.kp.MaxConnectionAge)
	kpTimer := time.NewTimer(t.kp.Time)
	defer func() {
		
		
		
		idleTimer.Stop()
		ageTimer.Stop()
		kpTimer.Stop()
	}()

	for {
		select {
		case <-idleTimer.C:
			t.mu.Lock()
			idle := t.idle
			if idle.IsZero() { 
				t.mu.Unlock()
				idleTimer.Reset(t.kp.MaxConnectionIdle)
				continue
			}
			val := t.kp.MaxConnectionIdle - time.Since(idle)
			t.mu.Unlock()
			if val <= 0 {
				
				
				t.Drain("max_idle")
				return
			}
			idleTimer.Reset(val)
		case <-ageTimer.C:
			t.Drain("max_age")
			ageTimer.Reset(t.kp.MaxConnectionAgeGrace)
			select {
			case <-ageTimer.C:
				
				if t.logger.V(logLevel) {
					t.logger.Infof("Closing server transport due to maximum connection age")
				}
				t.controlBuf.put(closeConnection{})
			case <-t.done:
			}
			return
		case <-kpTimer.C:
			lastRead := atomic.LoadInt64(&t.lastRead)
			if lastRead > prevNano {
				
				
				
				outstandingPing = false
				kpTimer.Reset(time.Duration(lastRead) + t.kp.Time - time.Duration(time.Now().UnixNano()))
				prevNano = lastRead
				continue
			}
			if outstandingPing && kpTimeoutLeft <= 0 {
				t.Close(fmt.Errorf("keepalive ping not acked within timeout %s", t.kp.Timeout))
				return
			}
			if !outstandingPing {
				if channelz.IsOn() {
					t.channelz.SocketMetrics.KeepAlivesSent.Add(1)
				}
				t.controlBuf.put(p)
				kpTimeoutLeft = t.kp.Timeout
				outstandingPing = true
			}
			
			
			
			
			sleepDuration := min(t.kp.Time, kpTimeoutLeft)
			kpTimeoutLeft -= sleepDuration
			kpTimer.Reset(sleepDuration)
		case <-t.done:
			return
		}
	}
}




func (t *http2Server) Close(err error) {
	t.mu.Lock()
	if t.state == closing {
		t.mu.Unlock()
		return
	}
	if t.logger.V(logLevel) {
		t.logger.Infof("Closing: %v", err)
	}
	t.state = closing
	streams := t.activeStreams
	t.activeStreams = nil
	t.mu.Unlock()
	t.controlBuf.finish()
	close(t.done)
	if err := t.conn.Close(); err != nil && t.logger.V(logLevel) {
		t.logger.Infof("Error closing underlying net.Conn during Close: %v", err)
	}
	channelz.RemoveEntry(t.channelz.ID)
	
	for _, s := range streams {
		s.cancel()
	}
}


func (t *http2Server) deleteStream(s *ServerStream, eosReceived bool) {

	t.mu.Lock()
	if _, ok := t.activeStreams[s.id]; ok {
		delete(t.activeStreams, s.id)
		if len(t.activeStreams) == 0 {
			t.idle = time.Now()
		}
	}
	t.mu.Unlock()

	if channelz.IsOn() {
		if eosReceived {
			t.channelz.SocketMetrics.StreamsSucceeded.Add(1)
		} else {
			t.channelz.SocketMetrics.StreamsFailed.Add(1)
		}
	}
}


func (t *http2Server) finishStream(s *ServerStream, rst bool, rstCode http2.ErrCode, hdr *headerFrame, eosReceived bool) {
	
	
	
	s.cancel()

	oldState := s.swapState(streamDone)
	if oldState == streamDone {
		
		return
	}

	hdr.cleanup = &cleanupStream{
		streamID: s.id,
		rst:      rst,
		rstCode:  rstCode,
		onWrite: func() {
			t.deleteStream(s, eosReceived)
		},
	}
	t.controlBuf.put(hdr)
}


func (t *http2Server) closeStream(s *ServerStream, rst bool, rstCode http2.ErrCode, eosReceived bool) {
	
	
	
	s.cancel()

	s.swapState(streamDone)
	t.deleteStream(s, eosReceived)

	t.controlBuf.put(&cleanupStream{
		streamID: s.id,
		rst:      rst,
		rstCode:  rstCode,
		onWrite:  func() {},
	})
}

func (t *http2Server) Drain(debugData string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.drainEvent != nil {
		return
	}
	t.drainEvent = grpcsync.NewEvent()
	t.controlBuf.put(&goAway{code: http2.ErrCodeNo, debugData: []byte(debugData), headsUp: true})
}

var goAwayPing = &ping{data: [8]byte{1, 6, 1, 8, 0, 3, 3, 9}}



func (t *http2Server) outgoingGoAwayHandler(g *goAway) (bool, error) {
	t.maxStreamMu.Lock()
	t.mu.Lock()
	if t.state == closing { 
		t.mu.Unlock()
		t.maxStreamMu.Unlock()
		
		return false, ErrConnClosing
	}
	if !g.headsUp {
		
		t.state = draining
		sid := t.maxStreamID
		retErr := g.closeConn
		if len(t.activeStreams) == 0 {
			retErr = errors.New("second GOAWAY written and no active streams left to process")
		}
		t.mu.Unlock()
		t.maxStreamMu.Unlock()
		if err := t.framer.fr.WriteGoAway(sid, g.code, g.debugData); err != nil {
			return false, err
		}
		t.framer.writer.Flush()
		if retErr != nil {
			return false, retErr
		}
		return true, nil
	}
	t.mu.Unlock()
	t.maxStreamMu.Unlock()
	
	
	
	
	
	
	if err := t.framer.fr.WriteGoAway(math.MaxUint32, http2.ErrCodeNo, g.debugData); err != nil {
		return false, err
	}
	if err := t.framer.fr.WritePing(false, goAwayPing.data); err != nil {
		return false, err
	}
	go func() {
		timer := time.NewTimer(5 * time.Second)
		defer timer.Stop()
		select {
		case <-t.drainEvent.Done():
		case <-timer.C:
		case <-t.done:
			return
		}
		t.controlBuf.put(&goAway{code: g.code, debugData: g.debugData})
	}()
	return false, nil
}

func (t *http2Server) socketMetrics() *channelz.EphemeralSocketMetrics {
	return &channelz.EphemeralSocketMetrics{
		LocalFlowControlWindow:  int64(t.fc.getSize()),
		RemoteFlowControlWindow: t.getOutFlowWindow(),
	}
}

func (t *http2Server) incrMsgSent() {
	if channelz.IsOn() {
		t.channelz.SocketMetrics.MessagesSent.Add(1)
		t.channelz.SocketMetrics.LastMessageSentTimestamp.Add(1)
	}
}

func (t *http2Server) incrMsgRecv() {
	if channelz.IsOn() {
		t.channelz.SocketMetrics.MessagesReceived.Add(1)
		t.channelz.SocketMetrics.LastMessageReceivedTimestamp.Add(1)
	}
}

func (t *http2Server) getOutFlowWindow() int64 {
	resp := make(chan uint32, 1)
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	t.controlBuf.put(&outFlowControlSizeRequest{resp})
	select {
	case sz := <-resp:
		return int64(sz)
	case <-t.done:
		return -1
	case <-timer.C:
		return -2
	}
}


func (t *http2Server) Peer() *peer.Peer {
	return &peer.Peer{
		Addr:      t.peer.Addr,
		LocalAddr: t.peer.LocalAddr,
		AuthInfo:  t.peer.AuthInfo, 
	}
}

func getJitter(v time.Duration) time.Duration {
	if v == infinity {
		return 0
	}
	
	r := int64(v / 10)
	j := rand.Int64N(2*r) - r
	return time.Duration(j)
}

type connectionKey struct{}


func GetConnection(ctx context.Context) net.Conn {
	conn, _ := ctx.Value(connectionKey{}).(net.Conn)
	return conn
}




func SetConnection(ctx context.Context, conn net.Conn) context.Context {
	return context.WithValue(ctx, connectionKey{}, conn)
}
