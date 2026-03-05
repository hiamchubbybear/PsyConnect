

package grpc

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/experimental/stats"
	"google.golang.org/grpc/internal"
	"google.golang.org/grpc/internal/balancer/gracefulswitch"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/grpcsync"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

var (
	setConnectedAddress = internal.SetConnectedAddress.(func(*balancer.SubConnState, resolver.Address))
	
	
	noOpRegisterHealthListenerFn = func(_ context.Context, listener func(balancer.SubConnState)) func() {
		listener(balancer.SubConnState{ConnectivityState: connectivity.Ready})
		return func() {}
	}
)















type ccBalancerWrapper struct {
	internal.EnforceClientConnEmbedding
	
	
	cc               *ClientConn
	opts             balancer.BuildOptions
	serializer       *grpcsync.CallbackSerializer
	serializerCancel context.CancelFunc

	
	
	curBalancerName string
	balancer        *gracefulswitch.Balancer

	
	
	mu     sync.Mutex
	closed bool
}




func newCCBalancerWrapper(cc *ClientConn) *ccBalancerWrapper {
	ctx, cancel := context.WithCancel(cc.ctx)
	ccb := &ccBalancerWrapper{
		cc: cc,
		opts: balancer.BuildOptions{
			DialCreds:       cc.dopts.copts.TransportCredentials,
			CredsBundle:     cc.dopts.copts.CredsBundle,
			Dialer:          cc.dopts.copts.Dialer,
			Authority:       cc.authority,
			CustomUserAgent: cc.dopts.copts.UserAgent,
			ChannelzParent:  cc.channelz,
			Target:          cc.parsedTarget,
		},
		serializer:       grpcsync.NewCallbackSerializer(ctx),
		serializerCancel: cancel,
	}
	ccb.balancer = gracefulswitch.NewBalancer(ccb, ccb.opts)
	return ccb
}

func (ccb *ccBalancerWrapper) MetricsRecorder() stats.MetricsRecorder {
	return ccb.cc.metricsRecorderList
}




func (ccb *ccBalancerWrapper) updateClientConnState(ccs *balancer.ClientConnState) error {
	errCh := make(chan error)
	uccs := func(ctx context.Context) {
		defer close(errCh)
		if ctx.Err() != nil || ccb.balancer == nil {
			return
		}
		name := gracefulswitch.ChildName(ccs.BalancerConfig)
		if ccb.curBalancerName != name {
			ccb.curBalancerName = name
			channelz.Infof(logger, ccb.cc.channelz, "Channel switches to new LB policy %q", name)
		}
		err := ccb.balancer.UpdateClientConnState(*ccs)
		if logger.V(2) && err != nil {
			logger.Infof("error from balancer.UpdateClientConnState: %v", err)
		}
		errCh <- err
	}
	onFailure := func() { close(errCh) }

	
	
	
	
	
	
	ccb.serializer.ScheduleOr(uccs, onFailure)
	return <-errCh
}



func (ccb *ccBalancerWrapper) resolverError(err error) {
	ccb.serializer.TrySchedule(func(ctx context.Context) {
		if ctx.Err() != nil || ccb.balancer == nil {
			return
		}
		ccb.balancer.ResolverError(err)
	})
}




func (ccb *ccBalancerWrapper) close() {
	ccb.mu.Lock()
	ccb.closed = true
	ccb.mu.Unlock()
	channelz.Info(logger, ccb.cc.channelz, "ccBalancerWrapper: closing")
	ccb.serializer.TrySchedule(func(context.Context) {
		if ccb.balancer == nil {
			return
		}
		ccb.balancer.Close()
		ccb.balancer = nil
	})
	ccb.serializerCancel()
}


func (ccb *ccBalancerWrapper) exitIdle() {
	ccb.serializer.TrySchedule(func(ctx context.Context) {
		if ctx.Err() != nil || ccb.balancer == nil {
			return
		}
		ccb.balancer.ExitIdle()
	})
}

func (ccb *ccBalancerWrapper) NewSubConn(addrs []resolver.Address, opts balancer.NewSubConnOptions) (balancer.SubConn, error) {
	ccb.cc.mu.Lock()
	defer ccb.cc.mu.Unlock()

	ccb.mu.Lock()
	if ccb.closed {
		ccb.mu.Unlock()
		return nil, fmt.Errorf("balancer is being closed; no new SubConns allowed")
	}
	ccb.mu.Unlock()

	if len(addrs) == 0 {
		return nil, fmt.Errorf("grpc: cannot create SubConn with empty address list")
	}
	ac, err := ccb.cc.newAddrConnLocked(addrs, opts)
	if err != nil {
		channelz.Warningf(logger, ccb.cc.channelz, "acBalancerWrapper: NewSubConn: failed to newAddrConn: %v", err)
		return nil, err
	}
	acbw := &acBalancerWrapper{
		ccb:           ccb,
		ac:            ac,
		producers:     make(map[balancer.ProducerBuilder]*refCountedProducer),
		stateListener: opts.StateListener,
		healthData:    newHealthData(connectivity.Idle),
	}
	ac.acbw = acbw
	return acbw, nil
}

func (ccb *ccBalancerWrapper) RemoveSubConn(balancer.SubConn) {
	
	logger.Errorf("ccb RemoveSubConn(%v) called unexpectedly, sc")
}

func (ccb *ccBalancerWrapper) UpdateAddresses(sc balancer.SubConn, addrs []resolver.Address) {
	acbw, ok := sc.(*acBalancerWrapper)
	if !ok {
		return
	}
	acbw.UpdateAddresses(addrs)
}

func (ccb *ccBalancerWrapper) UpdateState(s balancer.State) {
	ccb.cc.mu.Lock()
	defer ccb.cc.mu.Unlock()
	if ccb.cc.conns == nil {
		
		return
	}

	ccb.mu.Lock()
	if ccb.closed {
		ccb.mu.Unlock()
		return
	}
	ccb.mu.Unlock()
	
	
	
	
	

	
	
	
	ccb.cc.pickerWrapper.updatePicker(s.Picker)
	ccb.cc.csMgr.updateState(s.ConnectivityState)
}

func (ccb *ccBalancerWrapper) ResolveNow(o resolver.ResolveNowOptions) {
	ccb.cc.mu.RLock()
	defer ccb.cc.mu.RUnlock()

	ccb.mu.Lock()
	if ccb.closed {
		ccb.mu.Unlock()
		return
	}
	ccb.mu.Unlock()
	ccb.cc.resolveNowLocked(o)
}

func (ccb *ccBalancerWrapper) Target() string {
	return ccb.cc.target
}



type acBalancerWrapper struct {
	internal.EnforceSubConnEmbedding
	ac            *addrConn          
	ccb           *ccBalancerWrapper 
	stateListener func(balancer.SubConnState)

	producersMu sync.Mutex
	producers   map[balancer.ProducerBuilder]*refCountedProducer

	
	healthMu sync.Mutex
	
	
	
	healthData *healthData
}


type healthData struct {
	
	
	
	connectivityState connectivity.State
	
	
	
	closeHealthProducer func()
}

func newHealthData(s connectivity.State) *healthData {
	return &healthData{
		connectivityState:   s,
		closeHealthProducer: func() {},
	}
}



func (acbw *acBalancerWrapper) updateState(s connectivity.State, curAddr resolver.Address, err error) {
	acbw.ccb.serializer.TrySchedule(func(ctx context.Context) {
		if ctx.Err() != nil || acbw.ccb.balancer == nil {
			return
		}
		
		acbw.closeProducers()

		
		
		
		scs := balancer.SubConnState{ConnectivityState: s, ConnectionError: err}
		if s == connectivity.Ready {
			setConnectedAddress(&scs, curAddr)
		}
		
		acbw.healthMu.Lock()
		
		
		
		
		
		
		
		
		
		
		
		
		
		acbw.healthData = newHealthData(scs.ConnectivityState)
		acbw.healthMu.Unlock()

		acbw.stateListener(scs)
	})
}

func (acbw *acBalancerWrapper) String() string {
	return fmt.Sprintf("SubConn(id:%d)", acbw.ac.channelz.ID)
}

func (acbw *acBalancerWrapper) UpdateAddresses(addrs []resolver.Address) {
	acbw.ac.updateAddrs(addrs)
}

func (acbw *acBalancerWrapper) Connect() {
	go acbw.ac.connect()
}

func (acbw *acBalancerWrapper) Shutdown() {
	acbw.closeProducers()
	acbw.ccb.cc.removeAddrConn(acbw.ac, errConnDrain)
}




func (acbw *acBalancerWrapper) NewStream(ctx context.Context, desc *StreamDesc, method string, opts ...CallOption) (ClientStream, error) {
	transport := acbw.ac.getReadyTransport()
	if transport == nil {
		return nil, status.Errorf(codes.Unavailable, "SubConn state is not Ready")

	}
	return newNonRetryClientStream(ctx, desc, method, transport, acbw.ac, opts...)
}



func (acbw *acBalancerWrapper) Invoke(ctx context.Context, method string, args any, reply any, opts ...CallOption) error {
	cs, err := acbw.NewStream(ctx, unaryStreamDesc, method, opts...)
	if err != nil {
		return err
	}
	if err := cs.SendMsg(args); err != nil {
		return err
	}
	return cs.RecvMsg(reply)
}

type refCountedProducer struct {
	producer balancer.Producer
	refs     int    
	close    func() 
}

func (acbw *acBalancerWrapper) GetOrBuildProducer(pb balancer.ProducerBuilder) (balancer.Producer, func()) {
	acbw.producersMu.Lock()
	defer acbw.producersMu.Unlock()

	
	pData := acbw.producers[pb]
	if pData == nil {
		
		p, closeFn := pb.Build(acbw)
		pData = &refCountedProducer{producer: p, close: closeFn}
		acbw.producers[pb] = pData
	}
	
	pData.refs++

	
	
	
	unref := func() {
		acbw.producersMu.Lock()
		
		
		
		pData.refs--
		if pData.refs == 0 {
			defer pData.close() 
			delete(acbw.producers, pb)
		}
		acbw.producersMu.Unlock()
	}
	return pData.producer, sync.OnceFunc(unref)
}

func (acbw *acBalancerWrapper) closeProducers() {
	acbw.producersMu.Lock()
	defer acbw.producersMu.Unlock()
	for pb, pData := range acbw.producers {
		pData.refs = 0
		pData.close()
		delete(acbw.producers, pb)
	}
}



type healthProducerRegisterFn = func(context.Context, balancer.SubConn, string, func(balancer.SubConnState)) func()










func (acbw *acBalancerWrapper) healthListenerRegFn() func(context.Context, func(balancer.SubConnState)) func() {
	if acbw.ccb.cc.dopts.disableHealthCheck {
		return noOpRegisterHealthListenerFn
	}
	regHealthLisFn := internal.RegisterClientHealthCheckListener
	if regHealthLisFn == nil {
		
		return noOpRegisterHealthListenerFn
	}
	cfg := acbw.ac.cc.healthCheckConfig()
	if cfg == nil {
		return noOpRegisterHealthListenerFn
	}
	return func(ctx context.Context, listener func(balancer.SubConnState)) func() {
		return regHealthLisFn.(healthProducerRegisterFn)(ctx, acbw, cfg.ServiceName, listener)
	}
}






func (acbw *acBalancerWrapper) RegisterHealthListener(listener func(balancer.SubConnState)) {
	acbw.healthMu.Lock()
	defer acbw.healthMu.Unlock()
	acbw.healthData.closeHealthProducer()
	
	
	
	
	if acbw.healthData.connectivityState != connectivity.Ready {
		return
	}
	
	
	hd := newHealthData(connectivity.Ready)
	acbw.healthData = hd
	if listener == nil {
		return
	}

	registerFn := acbw.healthListenerRegFn()
	acbw.ccb.serializer.TrySchedule(func(ctx context.Context) {
		if ctx.Err() != nil || acbw.ccb.balancer == nil {
			return
		}
		
		acbw.healthMu.Lock()
		defer acbw.healthMu.Unlock()
		if acbw.healthData != hd {
			return
		}
		
		
		listenerWrapper := func(scs balancer.SubConnState) {
			acbw.ccb.serializer.TrySchedule(func(ctx context.Context) {
				if ctx.Err() != nil || acbw.ccb.balancer == nil {
					return
				}
				acbw.healthMu.Lock()
				defer acbw.healthMu.Unlock()
				if acbw.healthData != hd {
					return
				}
				listener(scs)
			})
		}

		hd.closeHealthProducer = registerFn(ctx, listenerWrapper)
	})
}
