

package grpc

import (
	"context"
	"errors"
	"io"
	"math"
	rand "math/rand/v2"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/internal"
	"google.golang.org/grpc/internal/balancerload"
	"google.golang.org/grpc/internal/binarylog"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/grpcutil"
	imetadata "google.golang.org/grpc/internal/metadata"
	iresolver "google.golang.org/grpc/internal/resolver"
	"google.golang.org/grpc/internal/serviceconfig"
	istatus "google.golang.org/grpc/internal/status"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/mem"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

var metadataFromOutgoingContextRaw = internal.FromOutgoingContextRaw.(func(context.Context) (metadata.MD, [][]string, bool))








type StreamHandler func(srv any, stream ServerStream) error




type StreamDesc struct {
	
	
	StreamName string        
	Handler    StreamHandler 

	
	
	
	ServerStreams bool 
	ClientStreams bool 
}




type Stream interface {
	
	Context() context.Context
	
	SendMsg(m any) error
	
	RecvMsg(m any) error
}





type ClientStream interface {
	
	
	
	
	Header() (metadata.MD, error)
	
	
	
	Trailer() metadata.MD
	
	
	
	CloseSend() error
	
	
	
	
	Context() context.Context
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	SendMsg(m any) error
	
	
	
	
	
	
	
	
	RecvMsg(m any) error
}

















func (cc *ClientConn) NewStream(ctx context.Context, desc *StreamDesc, method string, opts ...CallOption) (ClientStream, error) {
	
	
	opts = combine(cc.dopts.callOptions, opts)

	if cc.dopts.streamInt != nil {
		return cc.dopts.streamInt(ctx, desc, cc, method, newClientStream, opts...)
	}
	return newClientStream(ctx, desc, cc, method, opts...)
}


func NewClientStream(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, opts ...CallOption) (ClientStream, error) {
	return cc.NewStream(ctx, desc, method, opts...)
}

func newClientStream(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, opts ...CallOption) (_ ClientStream, err error) {
	
	
	
	if err := cc.idlenessMgr.OnCallBegin(); err != nil {
		return nil, err
	}
	
	
	opts = append([]CallOption{OnFinish(func(error) { cc.idlenessMgr.OnCallEnd() })}, opts...)

	if md, added, ok := metadataFromOutgoingContextRaw(ctx); ok {
		
		if err := imetadata.Validate(md); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		
		for _, kvs := range added {
			for i := 0; i < len(kvs); i += 2 {
				if err := imetadata.ValidatePair(kvs[i], kvs[i+1]); err != nil {
					return nil, status.Error(codes.Internal, err.Error())
				}
			}
		}
	}
	if channelz.IsOn() {
		cc.incrCallsStarted()
		defer func() {
			if err != nil {
				cc.incrCallsFailed()
			}
		}()
	}
	
	
	if err := cc.waitForResolvedAddrs(ctx); err != nil {
		return nil, err
	}

	var mc serviceconfig.MethodConfig
	var onCommit func()
	newStream := func(ctx context.Context, done func()) (iresolver.ClientStream, error) {
		return newClientStreamWithParams(ctx, desc, cc, method, mc, onCommit, done, opts...)
	}

	rpcInfo := iresolver.RPCInfo{Context: ctx, Method: method}
	rpcConfig, err := cc.safeConfigSelector.SelectConfig(rpcInfo)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			
			if istatus.IsRestrictedControlPlaneCode(st) {
				err = status.Errorf(codes.Internal, "config selector returned illegal status: %v", err)
			}
			return nil, err
		}
		return nil, toRPCErr(err)
	}

	if rpcConfig != nil {
		if rpcConfig.Context != nil {
			ctx = rpcConfig.Context
		}
		mc = rpcConfig.MethodConfig
		onCommit = rpcConfig.OnCommitted
		if rpcConfig.Interceptor != nil {
			rpcInfo.Context = nil
			ns := newStream
			newStream = func(ctx context.Context, done func()) (iresolver.ClientStream, error) {
				cs, err := rpcConfig.Interceptor.NewStream(ctx, rpcInfo, done, ns)
				if err != nil {
					return nil, toRPCErr(err)
				}
				return cs, nil
			}
		}
	}

	return newStream(ctx, func() {})
}

func newClientStreamWithParams(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, mc serviceconfig.MethodConfig, onCommit, doneFunc func(), opts ...CallOption) (_ iresolver.ClientStream, err error) {
	callInfo := defaultCallInfo()
	if mc.WaitForReady != nil {
		callInfo.failFast = !*mc.WaitForReady
	}

	
	
	
	
	
	var cancel context.CancelFunc
	if mc.Timeout != nil && *mc.Timeout >= 0 {
		ctx, cancel = context.WithTimeout(ctx, *mc.Timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	for _, o := range opts {
		if err := o.before(callInfo); err != nil {
			return nil, toRPCErr(err)
		}
	}
	callInfo.maxSendMessageSize = getMaxSize(mc.MaxReqSize, callInfo.maxSendMessageSize, defaultClientMaxSendMessageSize)
	callInfo.maxReceiveMessageSize = getMaxSize(mc.MaxRespSize, callInfo.maxReceiveMessageSize, defaultClientMaxReceiveMessageSize)
	if err := setCallInfoCodec(callInfo); err != nil {
		return nil, err
	}

	callHdr := &transport.CallHdr{
		Host:           cc.authority,
		Method:         method,
		ContentSubtype: callInfo.contentSubtype,
		DoneFunc:       doneFunc,
	}

	
	
	
	
	var compressorV0 Compressor
	var compressorV1 encoding.Compressor
	if ct := callInfo.compressorName; ct != "" {
		callHdr.SendCompress = ct
		if ct != encoding.Identity {
			compressorV1 = encoding.GetCompressor(ct)
			if compressorV1 == nil {
				return nil, status.Errorf(codes.Internal, "grpc: Compressor is not installed for requested grpc-encoding %q", ct)
			}
		}
	} else if cc.dopts.compressorV0 != nil {
		callHdr.SendCompress = cc.dopts.compressorV0.Type()
		compressorV0 = cc.dopts.compressorV0
	}
	if callInfo.creds != nil {
		callHdr.Creds = callInfo.creds
	}

	cs := &clientStream{
		callHdr:      callHdr,
		ctx:          ctx,
		methodConfig: &mc,
		opts:         opts,
		callInfo:     callInfo,
		cc:           cc,
		desc:         desc,
		codec:        callInfo.codec,
		compressorV0: compressorV0,
		compressorV1: compressorV1,
		cancel:       cancel,
		firstAttempt: true,
		onCommit:     onCommit,
	}
	if !cc.dopts.disableRetry {
		cs.retryThrottler = cc.retryThrottler.Load().(*retryThrottler)
	}
	if ml := binarylog.GetMethodLogger(method); ml != nil {
		cs.binlogs = append(cs.binlogs, ml)
	}
	if cc.dopts.binaryLogger != nil {
		if ml := cc.dopts.binaryLogger.GetMethodLogger(method); ml != nil {
			cs.binlogs = append(cs.binlogs, ml)
		}
	}

	
	
	op := func(a *csAttempt) error {
		if err := a.getTransport(); err != nil {
			return err
		}
		if err := a.newStream(); err != nil {
			return err
		}
		
		
		
		cs.attempt = a
		return nil
	}
	if err := cs.withRetry(op, func() { cs.bufferForRetryLocked(0, op, nil) }); err != nil {
		return nil, err
	}

	if len(cs.binlogs) != 0 {
		md, _ := metadata.FromOutgoingContext(ctx)
		logEntry := &binarylog.ClientHeader{
			OnClientSide: true,
			Header:       md,
			MethodName:   method,
			Authority:    cs.cc.authority,
		}
		if deadline, ok := ctx.Deadline(); ok {
			logEntry.Timeout = time.Until(deadline)
			if logEntry.Timeout < 0 {
				logEntry.Timeout = 0
			}
		}
		for _, binlog := range cs.binlogs {
			binlog.Log(cs.ctx, logEntry)
		}
	}

	if desc != unaryStreamDesc {
		
		
		
		
		
		go func() {
			select {
			case <-cc.ctx.Done():
				cs.finish(ErrClientConnClosing)
			case <-ctx.Done():
				cs.finish(toRPCErr(ctx.Err()))
			}
		}()
	}
	return cs, nil
}


func (cs *clientStream) newAttemptLocked(isTransparent bool) (*csAttempt, error) {
	if err := cs.ctx.Err(); err != nil {
		return nil, toRPCErr(err)
	}
	if err := cs.cc.ctx.Err(); err != nil {
		return nil, ErrClientConnClosing
	}

	ctx := newContextWithRPCInfo(cs.ctx, cs.callInfo.failFast, cs.callInfo.codec, cs.compressorV0, cs.compressorV1)
	method := cs.callHdr.Method
	var beginTime time.Time
	shs := cs.cc.dopts.copts.StatsHandlers
	for _, sh := range shs {
		ctx = sh.TagRPC(ctx, &stats.RPCTagInfo{FullMethodName: method, FailFast: cs.callInfo.failFast})
		beginTime = time.Now()
		begin := &stats.Begin{
			Client:                    true,
			BeginTime:                 beginTime,
			FailFast:                  cs.callInfo.failFast,
			IsClientStream:            cs.desc.ClientStreams,
			IsServerStream:            cs.desc.ServerStreams,
			IsTransparentRetryAttempt: isTransparent,
		}
		sh.HandleRPC(ctx, begin)
	}

	var trInfo *traceInfo
	if EnableTracing {
		trInfo = &traceInfo{
			tr: newTrace("grpc.Sent."+methodFamily(method), method),
			firstLine: firstLine{
				client: true,
			},
		}
		if deadline, ok := ctx.Deadline(); ok {
			trInfo.firstLine.deadline = time.Until(deadline)
		}
		trInfo.tr.LazyLog(&trInfo.firstLine, false)
		ctx = newTraceContext(ctx, trInfo.tr)
	}

	if cs.cc.parsedTarget.URL.Scheme == internal.GRPCResolverSchemeExtraMetadata {
		
		
		ctx = grpcutil.WithExtraMetadata(ctx, metadata.Pairs(
			"content-type", grpcutil.ContentType(cs.callHdr.ContentSubtype),
		))
	}

	return &csAttempt{
		ctx:            ctx,
		beginTime:      beginTime,
		cs:             cs,
		decompressorV0: cs.cc.dopts.dc,
		statsHandlers:  shs,
		trInfo:         trInfo,
	}, nil
}

func (a *csAttempt) getTransport() error {
	cs := a.cs

	var err error
	a.transport, a.pickResult, err = cs.cc.getTransport(a.ctx, cs.callInfo.failFast, cs.callHdr.Method)
	if err != nil {
		if de, ok := err.(dropError); ok {
			err = de.error
			a.drop = true
		}
		return err
	}
	if a.trInfo != nil {
		a.trInfo.firstLine.SetRemoteAddr(a.transport.RemoteAddr())
	}
	return nil
}

func (a *csAttempt) newStream() error {
	cs := a.cs
	cs.callHdr.PreviousAttempts = cs.numRetries

	
	
	
	
	if a.pickResult.Metadata != nil {
		
		
		
		
		
		
		
		
		md, _ := metadata.FromOutgoingContext(a.ctx)
		md = metadata.Join(md, a.pickResult.Metadata)
		a.ctx = metadata.NewOutgoingContext(a.ctx, md)
	}

	s, err := a.transport.NewStream(a.ctx, cs.callHdr)
	if err != nil {
		nse, ok := err.(*transport.NewStreamError)
		if !ok {
			
			return err
		}

		if nse.AllowTransparentRetry {
			a.allowTransparentRetry = true
		}

		
		return toRPCErr(nse.Err)
	}
	a.transportStream = s
	a.ctx = s.Context()
	a.parser = &parser{r: s, bufferPool: a.cs.cc.dopts.copts.BufferPool}
	return nil
}


type clientStream struct {
	callHdr  *transport.CallHdr
	opts     []CallOption
	callInfo *callInfo
	cc       *ClientConn
	desc     *StreamDesc

	codec        baseCodec
	compressorV0 Compressor
	compressorV1 encoding.Compressor

	cancel context.CancelFunc 

	sentLast bool 

	methodConfig *MethodConfig

	ctx context.Context 

	retryThrottler *retryThrottler 

	binlogs []binarylog.MethodLogger
	
	
	
	
	
	
	serverHeaderBinlogged bool

	mu                      sync.Mutex
	firstAttempt            bool 
	numRetries              int  
	numRetriesSincePushback int  
	finished                bool 
	
	
	
	
	
	
	
	attempt *csAttempt
	
	committed        bool 
	onCommit         func()
	replayBuffer     []replayOp 
	replayBufferSize int        
}

type replayOp struct {
	op      func(a *csAttempt) error
	cleanup func()
}



type csAttempt struct {
	ctx             context.Context
	cs              *clientStream
	transport       transport.ClientTransport
	transportStream *transport.ClientStream
	parser          *parser
	pickResult      balancer.PickResult

	finished        bool
	decompressorV0  Decompressor
	decompressorV1  encoding.Compressor
	decompressorSet bool

	mu sync.Mutex 
	
	
	
	trInfo *traceInfo

	statsHandlers []stats.Handler
	beginTime     time.Time

	
	allowTransparentRetry bool
	
	drop bool
}

func (cs *clientStream) commitAttemptLocked() {
	if !cs.committed && cs.onCommit != nil {
		cs.onCommit()
	}
	cs.committed = true
	for _, op := range cs.replayBuffer {
		if op.cleanup != nil {
			op.cleanup()
		}
	}
	cs.replayBuffer = nil
}

func (cs *clientStream) commitAttempt() {
	cs.mu.Lock()
	cs.commitAttemptLocked()
	cs.mu.Unlock()
}




func (a *csAttempt) shouldRetry(err error) (bool, error) {
	cs := a.cs

	if cs.finished || cs.committed || a.drop {
		
		return false, err
	}
	if a.transportStream == nil && a.allowTransparentRetry {
		return true, nil
	}
	
	unprocessed := false
	if a.transportStream != nil {
		<-a.transportStream.Done()
		unprocessed = a.transportStream.Unprocessed()
	}
	if cs.firstAttempt && unprocessed {
		
		return true, nil
	}
	if cs.cc.dopts.disableRetry {
		return false, err
	}

	pushback := 0
	hasPushback := false
	if a.transportStream != nil {
		if !a.transportStream.TrailersOnly() {
			return false, err
		}

		
		
		sps := a.transportStream.Trailer()["grpc-retry-pushback-ms"]
		if len(sps) == 1 {
			var e error
			if pushback, e = strconv.Atoi(sps[0]); e != nil || pushback < 0 {
				channelz.Infof(logger, cs.cc.channelz, "Server retry pushback specified to abort (%q).", sps[0])
				cs.retryThrottler.throttle() 
				return false, err
			}
			hasPushback = true
		} else if len(sps) > 1 {
			channelz.Warningf(logger, cs.cc.channelz, "Server retry pushback specified multiple values (%q); not retrying.", sps)
			cs.retryThrottler.throttle() 
			return false, err
		}
	}

	var code codes.Code
	if a.transportStream != nil {
		code = a.transportStream.Status().Code()
	} else {
		code = status.Code(err)
	}

	rp := cs.methodConfig.RetryPolicy
	if rp == nil || !rp.RetryableStatusCodes[code] {
		return false, err
	}

	
	
	if cs.retryThrottler.throttle() {
		return false, err
	}
	if cs.numRetries+1 >= rp.MaxAttempts {
		return false, err
	}

	var dur time.Duration
	if hasPushback {
		dur = time.Millisecond * time.Duration(pushback)
		cs.numRetriesSincePushback = 0
	} else {
		fact := math.Pow(rp.BackoffMultiplier, float64(cs.numRetriesSincePushback))
		cur := min(float64(rp.InitialBackoff)*fact, float64(rp.MaxBackoff))
		
		cur *= 0.8 + 0.4*rand.Float64()
		dur = time.Duration(int64(cur))
		cs.numRetriesSincePushback++
	}

	
	
	t := time.NewTimer(dur)
	select {
	case <-t.C:
		cs.numRetries++
		return false, nil
	case <-cs.ctx.Done():
		t.Stop()
		return false, status.FromContextError(cs.ctx.Err()).Err()
	}
}


func (cs *clientStream) retryLocked(attempt *csAttempt, lastErr error) error {
	for {
		attempt.finish(toRPCErr(lastErr))
		isTransparent, err := attempt.shouldRetry(lastErr)
		if err != nil {
			cs.commitAttemptLocked()
			return err
		}
		cs.firstAttempt = false
		attempt, err = cs.newAttemptLocked(isTransparent)
		if err != nil {
			
			
			return err
		}
		
		
		if lastErr = cs.replayBufferLocked(attempt); lastErr == nil {
			return nil
		}
	}
}

func (cs *clientStream) Context() context.Context {
	cs.commitAttempt()
	
	
	if cs.attempt.transportStream != nil {
		return cs.attempt.transportStream.Context()
	}
	return cs.ctx
}

func (cs *clientStream) withRetry(op func(a *csAttempt) error, onSuccess func()) error {
	cs.mu.Lock()
	for {
		if cs.committed {
			cs.mu.Unlock()
			
			
			
			
			return toRPCErr(op(cs.attempt))
		}
		if len(cs.replayBuffer) == 0 {
			
			
			
			
			var err error
			if cs.attempt, err = cs.newAttemptLocked(false ); err != nil {
				cs.mu.Unlock()
				cs.finish(err)
				return err
			}
		}
		a := cs.attempt
		cs.mu.Unlock()
		err := op(a)
		cs.mu.Lock()
		if a != cs.attempt {
			
			continue
		}
		if err == io.EOF {
			<-a.transportStream.Done()
		}
		if err == nil || (err == io.EOF && a.transportStream.Status().Code() == codes.OK) {
			onSuccess()
			cs.mu.Unlock()
			return err
		}
		if err := cs.retryLocked(a, err); err != nil {
			cs.mu.Unlock()
			return err
		}
	}
}

func (cs *clientStream) Header() (metadata.MD, error) {
	var m metadata.MD
	err := cs.withRetry(func(a *csAttempt) error {
		var err error
		m, err = a.transportStream.Header()
		return toRPCErr(err)
	}, cs.commitAttemptLocked)

	if m == nil && err == nil {
		
		err = io.EOF
	}

	if err != nil {
		cs.finish(err)
		
		return nil, nil
	}

	if len(cs.binlogs) != 0 && !cs.serverHeaderBinlogged && m != nil {
		
		
		logEntry := &binarylog.ServerHeader{
			OnClientSide: true,
			Header:       m,
			PeerAddr:     nil,
		}
		if peer, ok := peer.FromContext(cs.Context()); ok {
			logEntry.PeerAddr = peer.Addr
		}
		cs.serverHeaderBinlogged = true
		for _, binlog := range cs.binlogs {
			binlog.Log(cs.ctx, logEntry)
		}
	}

	return m, nil
}

func (cs *clientStream) Trailer() metadata.MD {
	
	
	
	
	
	
	
	cs.commitAttempt()
	if cs.attempt.transportStream == nil {
		return nil
	}
	return cs.attempt.transportStream.Trailer()
}

func (cs *clientStream) replayBufferLocked(attempt *csAttempt) error {
	for _, f := range cs.replayBuffer {
		if err := f.op(attempt); err != nil {
			return err
		}
	}
	return nil
}

func (cs *clientStream) bufferForRetryLocked(sz int, op func(a *csAttempt) error, cleanup func()) {
	
	if cs.committed {
		return
	}
	cs.replayBufferSize += sz
	if cs.replayBufferSize > cs.callInfo.maxRetryRPCBufferSize {
		cs.commitAttemptLocked()
		cleanup()
		return
	}
	cs.replayBuffer = append(cs.replayBuffer, replayOp{op: op, cleanup: cleanup})
}

func (cs *clientStream) SendMsg(m any) (err error) {
	defer func() {
		if err != nil && err != io.EOF {
			
			
			
			
			
			cs.finish(err)
		}
	}()
	if cs.sentLast {
		return status.Errorf(codes.Internal, "SendMsg called after CloseSend")
	}
	if !cs.desc.ClientStreams {
		cs.sentLast = true
	}

	
	hdr, data, payload, pf, err := prepareMsg(m, cs.codec, cs.compressorV0, cs.compressorV1, cs.cc.dopts.copts.BufferPool)
	if err != nil {
		return err
	}

	defer func() {
		data.Free()
		
		
		if pf.isCompressed() {
			payload.Free()
		}
	}()

	dataLen := data.Len()
	payloadLen := payload.Len()
	
	if payloadLen > *cs.callInfo.maxSendMessageSize {
		return status.Errorf(codes.ResourceExhausted, "trying to send message larger than max (%d vs. %d)", payloadLen, *cs.callInfo.maxSendMessageSize)
	}

	
	
	payload.Ref()
	op := func(a *csAttempt) error {
		return a.sendMsg(m, hdr, payload, dataLen, payloadLen)
	}

	
	
	
	
	onSuccessCalled := false
	err = cs.withRetry(op, func() {
		cs.bufferForRetryLocked(len(hdr)+payloadLen, op, payload.Free)
		onSuccessCalled = true
	})
	if !onSuccessCalled {
		payload.Free()
	}
	if len(cs.binlogs) != 0 && err == nil {
		cm := &binarylog.ClientMessage{
			OnClientSide: true,
			Message:      data.Materialize(),
		}
		for _, binlog := range cs.binlogs {
			binlog.Log(cs.ctx, cm)
		}
	}
	return err
}

func (cs *clientStream) RecvMsg(m any) error {
	if len(cs.binlogs) != 0 && !cs.serverHeaderBinlogged {
		
		cs.Header()
	}
	var recvInfo *payloadInfo
	if len(cs.binlogs) != 0 {
		recvInfo = &payloadInfo{}
		defer recvInfo.free()
	}
	err := cs.withRetry(func(a *csAttempt) error {
		return a.recvMsg(m, recvInfo)
	}, cs.commitAttemptLocked)
	if len(cs.binlogs) != 0 && err == nil {
		sm := &binarylog.ServerMessage{
			OnClientSide: true,
			Message:      recvInfo.uncompressedBytes.Materialize(),
		}
		for _, binlog := range cs.binlogs {
			binlog.Log(cs.ctx, sm)
		}
	}
	if err != nil || !cs.desc.ServerStreams {
		
		cs.finish(err)
	}
	return err
}

func (cs *clientStream) CloseSend() error {
	if cs.sentLast {
		
		return nil
	}
	cs.sentLast = true
	op := func(a *csAttempt) error {
		a.transportStream.Write(nil, nil, &transport.WriteOptions{Last: true})
		
		
		
		
		return nil
	}
	cs.withRetry(op, func() { cs.bufferForRetryLocked(0, op, nil) })
	if len(cs.binlogs) != 0 {
		chc := &binarylog.ClientHalfClose{
			OnClientSide: true,
		}
		for _, binlog := range cs.binlogs {
			binlog.Log(cs.ctx, chc)
		}
	}
	
	return nil
}

func (cs *clientStream) finish(err error) {
	if err == io.EOF {
		
		err = nil
	}
	cs.mu.Lock()
	if cs.finished {
		cs.mu.Unlock()
		return
	}
	cs.finished = true
	for _, onFinish := range cs.callInfo.onFinish {
		onFinish(err)
	}
	cs.commitAttemptLocked()
	if cs.attempt != nil {
		cs.attempt.finish(err)
		
		if cs.attempt.transportStream != nil {
			for _, o := range cs.opts {
				o.after(cs.callInfo, cs.attempt)
			}
		}
	}

	cs.mu.Unlock()
	
	if len(cs.binlogs) != 0 {
		switch err {
		case errContextCanceled, errContextDeadline, ErrClientConnClosing:
			c := &binarylog.Cancel{
				OnClientSide: true,
			}
			for _, binlog := range cs.binlogs {
				binlog.Log(cs.ctx, c)
			}
		default:
			logEntry := &binarylog.ServerTrailer{
				OnClientSide: true,
				Trailer:      cs.Trailer(),
				Err:          err,
			}
			if peer, ok := peer.FromContext(cs.Context()); ok {
				logEntry.PeerAddr = peer.Addr
			}
			for _, binlog := range cs.binlogs {
				binlog.Log(cs.ctx, logEntry)
			}
		}
	}
	if err == nil {
		cs.retryThrottler.successfulRPC()
	}
	if channelz.IsOn() {
		if err != nil {
			cs.cc.incrCallsFailed()
		} else {
			cs.cc.incrCallsSucceeded()
		}
	}
	cs.cancel()
}

func (a *csAttempt) sendMsg(m any, hdr []byte, payld mem.BufferSlice, dataLength, payloadLength int) error {
	cs := a.cs
	if a.trInfo != nil {
		a.mu.Lock()
		if a.trInfo.tr != nil {
			a.trInfo.tr.LazyLog(&payload{sent: true, msg: m}, true)
		}
		a.mu.Unlock()
	}
	if err := a.transportStream.Write(hdr, payld, &transport.WriteOptions{Last: !cs.desc.ClientStreams}); err != nil {
		if !cs.desc.ClientStreams {
			
			
			
			return nil
		}
		return io.EOF
	}
	if len(a.statsHandlers) != 0 {
		for _, sh := range a.statsHandlers {
			sh.HandleRPC(a.ctx, outPayload(true, m, dataLength, payloadLength, time.Now()))
		}
	}
	return nil
}

func (a *csAttempt) recvMsg(m any, payInfo *payloadInfo) (err error) {
	cs := a.cs
	if len(a.statsHandlers) != 0 && payInfo == nil {
		payInfo = &payloadInfo{}
		defer payInfo.free()
	}

	if !a.decompressorSet {
		
		if ct := a.transportStream.RecvCompress(); ct != "" && ct != encoding.Identity {
			if a.decompressorV0 == nil || a.decompressorV0.Type() != ct {
				
				
				a.decompressorV0 = nil
				a.decompressorV1 = encoding.GetCompressor(ct)
			}
		} else {
			
			a.decompressorV0 = nil
		}
		
		a.decompressorSet = true
	}
	if err := recv(a.parser, cs.codec, a.transportStream, a.decompressorV0, m, *cs.callInfo.maxReceiveMessageSize, payInfo, a.decompressorV1, false); err != nil {
		if err == io.EOF {
			if statusErr := a.transportStream.Status().Err(); statusErr != nil {
				return statusErr
			}
			return io.EOF 
		}

		return toRPCErr(err)
	}
	if a.trInfo != nil {
		a.mu.Lock()
		if a.trInfo.tr != nil {
			a.trInfo.tr.LazyLog(&payload{sent: false, msg: m}, true)
		}
		a.mu.Unlock()
	}
	for _, sh := range a.statsHandlers {
		sh.HandleRPC(a.ctx, &stats.InPayload{
			Client:           true,
			RecvTime:         time.Now(),
			Payload:          m,
			WireLength:       payInfo.compressedLength + headerLen,
			CompressedLength: payInfo.compressedLength,
			Length:           payInfo.uncompressedBytes.Len(),
		})
	}
	if cs.desc.ServerStreams {
		
		return nil
	}
	
	
	if err := recv(a.parser, cs.codec, a.transportStream, a.decompressorV0, m, *cs.callInfo.maxReceiveMessageSize, nil, a.decompressorV1, false); err == io.EOF {
		return a.transportStream.Status().Err() 
	} else if err != nil {
		return toRPCErr(err)
	}
	return toRPCErr(errors.New("grpc: client streaming protocol violation: get <nil>, want <EOF>"))
}

func (a *csAttempt) finish(err error) {
	a.mu.Lock()
	if a.finished {
		a.mu.Unlock()
		return
	}
	a.finished = true
	if err == io.EOF {
		
		err = nil
	}
	var tr metadata.MD
	if a.transportStream != nil {
		a.transportStream.Close(err)
		tr = a.transportStream.Trailer()
	}

	if a.pickResult.Done != nil {
		br := false
		if a.transportStream != nil {
			br = a.transportStream.BytesReceived()
		}
		a.pickResult.Done(balancer.DoneInfo{
			Err:           err,
			Trailer:       tr,
			BytesSent:     a.transportStream != nil,
			BytesReceived: br,
			ServerLoad:    balancerload.Parse(tr),
		})
	}
	for _, sh := range a.statsHandlers {
		end := &stats.End{
			Client:    true,
			BeginTime: a.beginTime,
			EndTime:   time.Now(),
			Trailer:   tr,
			Error:     err,
		}
		sh.HandleRPC(a.ctx, end)
	}
	if a.trInfo != nil && a.trInfo.tr != nil {
		if err == nil {
			a.trInfo.tr.LazyPrintf("RPC: [OK]")
		} else {
			a.trInfo.tr.LazyPrintf("RPC: [%v]", err)
			a.trInfo.tr.SetError()
		}
		a.trInfo.tr.Finish()
		a.trInfo.tr = nil
	}
	a.mu.Unlock()
}












func newNonRetryClientStream(ctx context.Context, desc *StreamDesc, method string, t transport.ClientTransport, ac *addrConn, opts ...CallOption) (_ ClientStream, err error) {
	if t == nil {
		
		return nil, errors.New("transport provided is nil")
	}
	
	c := &callInfo{}

	
	
	
	
	
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	for _, o := range opts {
		if err := o.before(c); err != nil {
			return nil, toRPCErr(err)
		}
	}
	c.maxReceiveMessageSize = getMaxSize(nil, c.maxReceiveMessageSize, defaultClientMaxReceiveMessageSize)
	c.maxSendMessageSize = getMaxSize(nil, c.maxSendMessageSize, defaultServerMaxSendMessageSize)
	if err := setCallInfoCodec(c); err != nil {
		return nil, err
	}

	callHdr := &transport.CallHdr{
		Host:           ac.cc.authority,
		Method:         method,
		ContentSubtype: c.contentSubtype,
	}

	
	
	
	
	var cp Compressor
	var comp encoding.Compressor
	if ct := c.compressorName; ct != "" {
		callHdr.SendCompress = ct
		if ct != encoding.Identity {
			comp = encoding.GetCompressor(ct)
			if comp == nil {
				return nil, status.Errorf(codes.Internal, "grpc: Compressor is not installed for requested grpc-encoding %q", ct)
			}
		}
	} else if ac.cc.dopts.compressorV0 != nil {
		callHdr.SendCompress = ac.cc.dopts.compressorV0.Type()
		cp = ac.cc.dopts.compressorV0
	}
	if c.creds != nil {
		callHdr.Creds = c.creds
	}

	
	as := &addrConnStream{
		callHdr:          callHdr,
		ac:               ac,
		ctx:              ctx,
		cancel:           cancel,
		opts:             opts,
		callInfo:         c,
		desc:             desc,
		codec:            c.codec,
		sendCompressorV0: cp,
		sendCompressorV1: comp,
		transport:        t,
	}

	s, err := as.transport.NewStream(as.ctx, as.callHdr)
	if err != nil {
		err = toRPCErr(err)
		return nil, err
	}
	as.transportStream = s
	as.parser = &parser{r: s, bufferPool: ac.dopts.copts.BufferPool}
	ac.incrCallsStarted()
	if desc != unaryStreamDesc {
		
		
		
		
		
		
		
		go func() {
			ac.mu.Lock()
			acCtx := ac.ctx
			ac.mu.Unlock()
			select {
			case <-acCtx.Done():
				as.finish(status.Error(codes.Canceled, "grpc: the SubConn is closing"))
			case <-ctx.Done():
				as.finish(toRPCErr(ctx.Err()))
			}
		}()
	}
	return as, nil
}

type addrConnStream struct {
	transportStream  *transport.ClientStream
	ac               *addrConn
	callHdr          *transport.CallHdr
	cancel           context.CancelFunc
	opts             []CallOption
	callInfo         *callInfo
	transport        transport.ClientTransport
	ctx              context.Context
	sentLast         bool
	desc             *StreamDesc
	codec            baseCodec
	sendCompressorV0 Compressor
	sendCompressorV1 encoding.Compressor
	decompressorSet  bool
	decompressorV0   Decompressor
	decompressorV1   encoding.Compressor
	parser           *parser

	
	mu       sync.Mutex
	finished bool
}

func (as *addrConnStream) Header() (metadata.MD, error) {
	m, err := as.transportStream.Header()
	if err != nil {
		as.finish(toRPCErr(err))
	}
	return m, err
}

func (as *addrConnStream) Trailer() metadata.MD {
	return as.transportStream.Trailer()
}

func (as *addrConnStream) CloseSend() error {
	if as.sentLast {
		
		return nil
	}
	as.sentLast = true

	as.transportStream.Write(nil, nil, &transport.WriteOptions{Last: true})
	
	
	
	
	return nil
}

func (as *addrConnStream) Context() context.Context {
	return as.transportStream.Context()
}

func (as *addrConnStream) SendMsg(m any) (err error) {
	defer func() {
		if err != nil && err != io.EOF {
			
			
			
			
			
			as.finish(err)
		}
	}()
	if as.sentLast {
		return status.Errorf(codes.Internal, "SendMsg called after CloseSend")
	}
	if !as.desc.ClientStreams {
		as.sentLast = true
	}

	
	hdr, data, payload, pf, err := prepareMsg(m, as.codec, as.sendCompressorV0, as.sendCompressorV1, as.ac.dopts.copts.BufferPool)
	if err != nil {
		return err
	}

	defer func() {
		data.Free()
		
		
		if pf.isCompressed() {
			payload.Free()
		}
	}()

	
	if payload.Len() > *as.callInfo.maxSendMessageSize {
		return status.Errorf(codes.ResourceExhausted, "trying to send message larger than max (%d vs. %d)", payload.Len(), *as.callInfo.maxSendMessageSize)
	}

	if err := as.transportStream.Write(hdr, payload, &transport.WriteOptions{Last: !as.desc.ClientStreams}); err != nil {
		if !as.desc.ClientStreams {
			
			
			
			return nil
		}
		return io.EOF
	}

	return nil
}

func (as *addrConnStream) RecvMsg(m any) (err error) {
	defer func() {
		if err != nil || !as.desc.ServerStreams {
			
			as.finish(err)
		}
	}()

	if !as.decompressorSet {
		
		if ct := as.transportStream.RecvCompress(); ct != "" && ct != encoding.Identity {
			if as.decompressorV0 == nil || as.decompressorV0.Type() != ct {
				
				
				as.decompressorV0 = nil
				as.decompressorV1 = encoding.GetCompressor(ct)
			}
		} else {
			
			as.decompressorV0 = nil
		}
		
		as.decompressorSet = true
	}
	if err := recv(as.parser, as.codec, as.transportStream, as.decompressorV0, m, *as.callInfo.maxReceiveMessageSize, nil, as.decompressorV1, false); err != nil {
		if err == io.EOF {
			if statusErr := as.transportStream.Status().Err(); statusErr != nil {
				return statusErr
			}
			return io.EOF 
		}
		return toRPCErr(err)
	}

	if as.desc.ServerStreams {
		
		return nil
	}

	
	
	if err := recv(as.parser, as.codec, as.transportStream, as.decompressorV0, m, *as.callInfo.maxReceiveMessageSize, nil, as.decompressorV1, false); err == io.EOF {
		return as.transportStream.Status().Err() 
	} else if err != nil {
		return toRPCErr(err)
	}
	return toRPCErr(errors.New("grpc: client streaming protocol violation: get <nil>, want <EOF>"))
}

func (as *addrConnStream) finish(err error) {
	as.mu.Lock()
	if as.finished {
		as.mu.Unlock()
		return
	}
	as.finished = true
	if err == io.EOF {
		
		err = nil
	}
	if as.transportStream != nil {
		as.transportStream.Close(err)
	}

	if err != nil {
		as.ac.incrCallsFailed()
	} else {
		as.ac.incrCallsSucceeded()
	}
	as.cancel()
	as.mu.Unlock()
}







type ServerStream interface {
	
	
	
	
	
	
	SetHeader(metadata.MD) error
	
	
	
	SendHeader(metadata.MD) error
	
	
	SetTrailer(metadata.MD)
	
	Context() context.Context
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	SendMsg(m any) error
	
	
	
	
	
	
	
	
	RecvMsg(m any) error
}


type serverStream struct {
	ctx   context.Context
	s     *transport.ServerStream
	p     *parser
	codec baseCodec

	compressorV0   Compressor
	compressorV1   encoding.Compressor
	decompressorV0 Decompressor
	decompressorV1 encoding.Compressor

	sendCompressorName string

	maxReceiveMessageSize int
	maxSendMessageSize    int
	trInfo                *traceInfo

	statsHandler []stats.Handler

	binlogs []binarylog.MethodLogger
	
	
	
	
	
	
	serverHeaderBinlogged bool

	mu sync.Mutex 
}

func (ss *serverStream) Context() context.Context {
	return ss.ctx
}

func (ss *serverStream) SetHeader(md metadata.MD) error {
	if md.Len() == 0 {
		return nil
	}
	err := imetadata.Validate(md)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return ss.s.SetHeader(md)
}

func (ss *serverStream) SendHeader(md metadata.MD) error {
	err := imetadata.Validate(md)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	err = ss.s.SendHeader(md)
	if len(ss.binlogs) != 0 && !ss.serverHeaderBinlogged {
		h, _ := ss.s.Header()
		sh := &binarylog.ServerHeader{
			Header: h,
		}
		ss.serverHeaderBinlogged = true
		for _, binlog := range ss.binlogs {
			binlog.Log(ss.ctx, sh)
		}
	}
	return err
}

func (ss *serverStream) SetTrailer(md metadata.MD) {
	if md.Len() == 0 {
		return
	}
	if err := imetadata.Validate(md); err != nil {
		logger.Errorf("stream: failed to validate md when setting trailer, err: %v", err)
	}
	ss.s.SetTrailer(md)
}

func (ss *serverStream) SendMsg(m any) (err error) {
	defer func() {
		if ss.trInfo != nil {
			ss.mu.Lock()
			if ss.trInfo.tr != nil {
				if err == nil {
					ss.trInfo.tr.LazyLog(&payload{sent: true, msg: m}, true)
				} else {
					ss.trInfo.tr.LazyLog(&fmtStringer{"%v", []any{err}}, true)
					ss.trInfo.tr.SetError()
				}
			}
			ss.mu.Unlock()
		}
		if err != nil && err != io.EOF {
			st, _ := status.FromError(toRPCErr(err))
			ss.s.WriteStatus(st)
			
			
			
			
			
			
		}
	}()

	
	
	if sendCompressorsName := ss.s.SendCompress(); sendCompressorsName != ss.sendCompressorName {
		ss.compressorV1 = encoding.GetCompressor(sendCompressorsName)
		ss.sendCompressorName = sendCompressorsName
	}

	
	hdr, data, payload, pf, err := prepareMsg(m, ss.codec, ss.compressorV0, ss.compressorV1, ss.p.bufferPool)
	if err != nil {
		return err
	}

	defer func() {
		data.Free()
		
		
		if pf.isCompressed() {
			payload.Free()
		}
	}()

	dataLen := data.Len()
	payloadLen := payload.Len()

	
	if payloadLen > ss.maxSendMessageSize {
		return status.Errorf(codes.ResourceExhausted, "trying to send message larger than max (%d vs. %d)", payloadLen, ss.maxSendMessageSize)
	}
	if err := ss.s.Write(hdr, payload, &transport.WriteOptions{Last: false}); err != nil {
		return toRPCErr(err)
	}

	if len(ss.binlogs) != 0 {
		if !ss.serverHeaderBinlogged {
			h, _ := ss.s.Header()
			sh := &binarylog.ServerHeader{
				Header: h,
			}
			ss.serverHeaderBinlogged = true
			for _, binlog := range ss.binlogs {
				binlog.Log(ss.ctx, sh)
			}
		}
		sm := &binarylog.ServerMessage{
			Message: data.Materialize(),
		}
		for _, binlog := range ss.binlogs {
			binlog.Log(ss.ctx, sm)
		}
	}
	if len(ss.statsHandler) != 0 {
		for _, sh := range ss.statsHandler {
			sh.HandleRPC(ss.s.Context(), outPayload(false, m, dataLen, payloadLen, time.Now()))
		}
	}
	return nil
}

func (ss *serverStream) RecvMsg(m any) (err error) {
	defer func() {
		if ss.trInfo != nil {
			ss.mu.Lock()
			if ss.trInfo.tr != nil {
				if err == nil {
					ss.trInfo.tr.LazyLog(&payload{sent: false, msg: m}, true)
				} else if err != io.EOF {
					ss.trInfo.tr.LazyLog(&fmtStringer{"%v", []any{err}}, true)
					ss.trInfo.tr.SetError()
				}
			}
			ss.mu.Unlock()
		}
		if err != nil && err != io.EOF {
			st, _ := status.FromError(toRPCErr(err))
			ss.s.WriteStatus(st)
			
			
			
			
			
			
		}
	}()
	var payInfo *payloadInfo
	if len(ss.statsHandler) != 0 || len(ss.binlogs) != 0 {
		payInfo = &payloadInfo{}
		defer payInfo.free()
	}
	if err := recv(ss.p, ss.codec, ss.s, ss.decompressorV0, m, ss.maxReceiveMessageSize, payInfo, ss.decompressorV1, true); err != nil {
		if err == io.EOF {
			if len(ss.binlogs) != 0 {
				chc := &binarylog.ClientHalfClose{}
				for _, binlog := range ss.binlogs {
					binlog.Log(ss.ctx, chc)
				}
			}
			return err
		}
		if err == io.ErrUnexpectedEOF {
			err = status.Error(codes.Internal, io.ErrUnexpectedEOF.Error())
		}
		return toRPCErr(err)
	}
	if len(ss.statsHandler) != 0 {
		for _, sh := range ss.statsHandler {
			sh.HandleRPC(ss.s.Context(), &stats.InPayload{
				RecvTime:         time.Now(),
				Payload:          m,
				Length:           payInfo.uncompressedBytes.Len(),
				WireLength:       payInfo.compressedLength + headerLen,
				CompressedLength: payInfo.compressedLength,
			})
		}
	}
	if len(ss.binlogs) != 0 {
		cm := &binarylog.ClientMessage{
			Message: payInfo.uncompressedBytes.Materialize(),
		}
		for _, binlog := range ss.binlogs {
			binlog.Log(ss.ctx, cm)
		}
	}
	return nil
}



func MethodFromServerStream(stream ServerStream) (string, bool) {
	return Method(stream.Context())
}






func prepareMsg(m any, codec baseCodec, cp Compressor, comp encoding.Compressor, pool mem.BufferPool) (hdr []byte, data, payload mem.BufferSlice, pf payloadFormat, err error) {
	if preparedMsg, ok := m.(*PreparedMsg); ok {
		return preparedMsg.hdr, preparedMsg.encodedData, preparedMsg.payload, preparedMsg.pf, nil
	}
	
	
	data, err = encode(codec, m)
	if err != nil {
		return nil, nil, nil, 0, err
	}
	compData, pf, err := compress(data, cp, comp, pool)
	if err != nil {
		data.Free()
		return nil, nil, nil, 0, err
	}
	hdr, payload = msgHeader(data, compData, pf)
	return hdr, data, payload, pf, nil
}
