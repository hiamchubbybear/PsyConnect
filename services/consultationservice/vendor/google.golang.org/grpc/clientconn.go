

package grpc

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/url"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/balancer/pickfirst"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/internal"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/grpcsync"
	"google.golang.org/grpc/internal/idle"
	iresolver "google.golang.org/grpc/internal/resolver"
	"google.golang.org/grpc/internal/stats"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
	"google.golang.org/grpc/status"

	_ "google.golang.org/grpc/balancer/roundrobin"           
	_ "google.golang.org/grpc/internal/resolver/passthrough" 
	_ "google.golang.org/grpc/internal/resolver/unix"        
	_ "google.golang.org/grpc/resolver/dns"                  
)

const (
	
	minConnectTimeout = 20 * time.Second
)

var (
	
	
	
	
	
	ErrClientConnClosing = status.Error(codes.Canceled, "grpc: the client connection is closing")
	
	errConnDrain = errors.New("grpc: the connection is drained")
	
	errConnClosing = errors.New("grpc: the connection is closing")
	
	
	errConnIdling = errors.New("grpc: the connection is closing due to channel idleness")
	
	
	invalidDefaultServiceConfigErrPrefix = "grpc: the provided default service config is invalid"
	
	PickFirstBalancerName = pickfirst.Name
)


var (
	
	
	
	errNoTransportSecurity = errors.New("grpc: no transport security set (use grpc.WithTransportCredentials(insecure.NewCredentials()) explicitly or set credentials)")
	
	
	errTransportCredsAndBundle = errors.New("grpc: credentials.Bundle may not be used with individual TransportCredentials")
	
	
	errNoTransportCredsInBundle = errors.New("grpc: credentials.Bundle must return non-nil transport credentials")
	
	
	
	errTransportCredentialsMissing = errors.New("grpc: the credentials require transport level security (use grpc.WithTransportCredentials() to set)")
)

const (
	defaultClientMaxReceiveMessageSize = 1024 * 1024 * 4
	defaultClientMaxSendMessageSize    = math.MaxInt32
	
	defaultWriteBufSize = 32 * 1024
	defaultReadBufSize  = 32 * 1024
)

type defaultConfigSelector struct {
	sc *ServiceConfig
}

func (dcs *defaultConfigSelector) SelectConfig(rpcInfo iresolver.RPCInfo) (*iresolver.RPCConfig, error) {
	return &iresolver.RPCConfig{
		Context:      rpcInfo.Context,
		MethodConfig: getMethodConfig(dcs.sc, rpcInfo.Method),
	}, nil
}



























func NewClient(target string, opts ...DialOption) (conn *ClientConn, err error) {
	cc := &ClientConn{
		target: target,
		conns:  make(map[*addrConn]struct{}),
		dopts:  defaultDialOptions(),
	}

	cc.retryThrottler.Store((*retryThrottler)(nil))
	cc.safeConfigSelector.UpdateConfigSelector(&defaultConfigSelector{nil})
	cc.ctx, cc.cancel = context.WithCancel(context.Background())

	
	disableGlobalOpts := false
	for _, opt := range opts {
		if _, ok := opt.(*disableGlobalDialOptions); ok {
			disableGlobalOpts = true
			break
		}
	}

	if !disableGlobalOpts {
		for _, opt := range globalDialOptions {
			opt.apply(&cc.dopts)
		}
	}

	for _, opt := range opts {
		opt.apply(&cc.dopts)
	}

	
	if err := cc.initParsedTargetAndResolverBuilder(); err != nil {
		return nil, err
	}

	for _, opt := range globalPerTargetDialOptions {
		opt.DialOptionForTarget(cc.parsedTarget.URL).apply(&cc.dopts)
	}

	chainUnaryClientInterceptors(cc)
	chainStreamClientInterceptors(cc)

	if err := cc.validateTransportCredentials(); err != nil {
		return nil, err
	}

	if cc.dopts.defaultServiceConfigRawJSON != nil {
		scpr := parseServiceConfig(*cc.dopts.defaultServiceConfigRawJSON, cc.dopts.maxCallAttempts)
		if scpr.Err != nil {
			return nil, fmt.Errorf("%s: %v", invalidDefaultServiceConfigErrPrefix, scpr.Err)
		}
		cc.dopts.defaultServiceConfig, _ = scpr.Config.(*ServiceConfig)
	}
	cc.keepaliveParams = cc.dopts.copts.KeepaliveParams

	if err = cc.initAuthority(); err != nil {
		return nil, err
	}

	
	
	cc.channelzRegistration(target)
	channelz.Infof(logger, cc.channelz, "parsed dial target is: %#v", cc.parsedTarget)
	channelz.Infof(logger, cc.channelz, "Channel authority set to %q", cc.authority)

	cc.csMgr = newConnectivityStateManager(cc.ctx, cc.channelz)
	cc.pickerWrapper = newPickerWrapper(cc.dopts.copts.StatsHandlers)

	cc.metricsRecorderList = stats.NewMetricsRecorderList(cc.dopts.copts.StatsHandlers)

	cc.initIdleStateLocked() 
	cc.idlenessMgr = idle.NewManager((*idler)(cc), cc.dopts.idleTimeout)

	return cc, nil
}




func Dial(target string, opts ...DialOption) (*ClientConn, error) {
	return DialContext(context.Background(), target, opts...)
}












func DialContext(ctx context.Context, target string, opts ...DialOption) (conn *ClientConn, err error) {
	
	
	
	
	
	
	
	opts = append([]DialOption{withDefaultScheme("passthrough"), WithLocalDNSResolution()}, opts...)
	cc, err := NewClient(target, opts...)
	if err != nil {
		return nil, err
	}

	
	
	
	defer func() {
		if err != nil {
			cc.Close()
		}
	}()

	
	if err := cc.idlenessMgr.ExitIdleMode(); err != nil {
		return nil, err
	}

	
	if !cc.dopts.block {
		return cc, nil
	}

	if cc.dopts.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cc.dopts.timeout)
		defer cancel()
	}
	defer func() {
		select {
		case <-ctx.Done():
			switch {
			case ctx.Err() == err:
				conn = nil
			case err == nil || !cc.dopts.returnLastError:
				conn, err = nil, ctx.Err()
			default:
				conn, err = nil, fmt.Errorf("%v: %v", ctx.Err(), err)
			}
		default:
		}
	}()

	
	for {
		s := cc.GetState()
		if s == connectivity.Idle {
			cc.Connect()
		}
		if s == connectivity.Ready {
			return cc, nil
		} else if cc.dopts.copts.FailOnNonTempDialError && s == connectivity.TransientFailure {
			if err = cc.connectionError(); err != nil {
				terr, ok := err.(interface {
					Temporary() bool
				})
				if ok && !terr.Temporary() {
					return nil, err
				}
			}
		}
		if !cc.WaitForStateChange(ctx, s) {
			
			if err = cc.connectionError(); err != nil && cc.dopts.returnLastError {
				return nil, err
			}
			return nil, ctx.Err()
		}
	}
}



func (cc *ClientConn) addTraceEvent(msg string) {
	ted := &channelz.TraceEvent{
		Desc:     fmt.Sprintf("Channel %s", msg),
		Severity: channelz.CtInfo,
	}
	if cc.dopts.channelzParent != nil {
		ted.Parent = &channelz.TraceEvent{
			Desc:     fmt.Sprintf("Nested channel(id:%d) %s", cc.channelz.ID, msg),
			Severity: channelz.CtInfo,
		}
	}
	channelz.AddTraceEvent(logger, cc.channelz, 0, ted)
}

type idler ClientConn

func (i *idler) EnterIdleMode() {
	(*ClientConn)(i).enterIdleMode()
}

func (i *idler) ExitIdleMode() error {
	return (*ClientConn)(i).exitIdleMode()
}




func (cc *ClientConn) exitIdleMode() (err error) {
	cc.mu.Lock()
	if cc.conns == nil {
		cc.mu.Unlock()
		return errConnClosing
	}
	cc.mu.Unlock()

	
	
	
	if err := cc.resolverWrapper.start(); err != nil {
		return err
	}

	cc.addTraceEvent("exiting idle mode")
	return nil
}


func (cc *ClientConn) initIdleStateLocked() {
	cc.resolverWrapper = newCCResolverWrapper(cc)
	cc.balancerWrapper = newCCBalancerWrapper(cc)
	cc.firstResolveEvent = grpcsync.NewEvent()
	
	
	
	cc.conns = make(map[*addrConn]struct{})
}




func (cc *ClientConn) enterIdleMode() {
	cc.mu.Lock()

	if cc.conns == nil {
		cc.mu.Unlock()
		return
	}

	conns := cc.conns

	rWrapper := cc.resolverWrapper
	rWrapper.close()
	cc.pickerWrapper.reset()
	bWrapper := cc.balancerWrapper
	bWrapper.close()
	cc.csMgr.updateState(connectivity.Idle)
	cc.addTraceEvent("entering idle mode")

	cc.initIdleStateLocked()

	cc.mu.Unlock()

	
	<-rWrapper.serializer.Done()
	<-bWrapper.serializer.Done()

	
	for ac := range conns {
		ac.tearDown(errConnIdling)
	}
}












func (cc *ClientConn) validateTransportCredentials() error {
	if cc.dopts.copts.TransportCredentials == nil && cc.dopts.copts.CredsBundle == nil {
		return errNoTransportSecurity
	}
	if cc.dopts.copts.TransportCredentials != nil && cc.dopts.copts.CredsBundle != nil {
		return errTransportCredsAndBundle
	}
	if cc.dopts.copts.CredsBundle != nil && cc.dopts.copts.CredsBundle.TransportCredentials() == nil {
		return errNoTransportCredsInBundle
	}
	transportCreds := cc.dopts.copts.TransportCredentials
	if transportCreds == nil {
		transportCreds = cc.dopts.copts.CredsBundle.TransportCredentials()
	}
	if transportCreds.Info().SecurityProtocol == "insecure" {
		for _, cd := range cc.dopts.copts.PerRPCCredentials {
			if cd.RequireTransportSecurity() {
				return errTransportCredentialsMissing
			}
		}
	}
	return nil
}








func (cc *ClientConn) channelzRegistration(target string) {
	parentChannel, _ := cc.dopts.channelzParent.(*channelz.Channel)
	cc.channelz = channelz.RegisterChannel(parentChannel, target)
	cc.addTraceEvent("created")
}


func chainUnaryClientInterceptors(cc *ClientConn) {
	interceptors := cc.dopts.chainUnaryInts
	
	
	if cc.dopts.unaryInt != nil {
		interceptors = append([]UnaryClientInterceptor{cc.dopts.unaryInt}, interceptors...)
	}
	var chainedInt UnaryClientInterceptor
	if len(interceptors) == 0 {
		chainedInt = nil
	} else if len(interceptors) == 1 {
		chainedInt = interceptors[0]
	} else {
		chainedInt = func(ctx context.Context, method string, req, reply any, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error {
			return interceptors[0](ctx, method, req, reply, cc, getChainUnaryInvoker(interceptors, 0, invoker), opts...)
		}
	}
	cc.dopts.unaryInt = chainedInt
}


func getChainUnaryInvoker(interceptors []UnaryClientInterceptor, curr int, finalInvoker UnaryInvoker) UnaryInvoker {
	if curr == len(interceptors)-1 {
		return finalInvoker
	}
	return func(ctx context.Context, method string, req, reply any, cc *ClientConn, opts ...CallOption) error {
		return interceptors[curr+1](ctx, method, req, reply, cc, getChainUnaryInvoker(interceptors, curr+1, finalInvoker), opts...)
	}
}


func chainStreamClientInterceptors(cc *ClientConn) {
	interceptors := cc.dopts.chainStreamInts
	
	
	if cc.dopts.streamInt != nil {
		interceptors = append([]StreamClientInterceptor{cc.dopts.streamInt}, interceptors...)
	}
	var chainedInt StreamClientInterceptor
	if len(interceptors) == 0 {
		chainedInt = nil
	} else if len(interceptors) == 1 {
		chainedInt = interceptors[0]
	} else {
		chainedInt = func(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, streamer Streamer, opts ...CallOption) (ClientStream, error) {
			return interceptors[0](ctx, desc, cc, method, getChainStreamer(interceptors, 0, streamer), opts...)
		}
	}
	cc.dopts.streamInt = chainedInt
}


func getChainStreamer(interceptors []StreamClientInterceptor, curr int, finalStreamer Streamer) Streamer {
	if curr == len(interceptors)-1 {
		return finalStreamer
	}
	return func(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, opts ...CallOption) (ClientStream, error) {
		return interceptors[curr+1](ctx, desc, cc, method, getChainStreamer(interceptors, curr+1, finalStreamer), opts...)
	}
}



func newConnectivityStateManager(ctx context.Context, channel *channelz.Channel) *connectivityStateManager {
	return &connectivityStateManager{
		channelz: channel,
		pubSub:   grpcsync.NewPubSub(ctx),
	}
}







type connectivityStateManager struct {
	mu         sync.Mutex
	state      connectivity.State
	notifyChan chan struct{}
	channelz   *channelz.Channel
	pubSub     *grpcsync.PubSub
}




func (csm *connectivityStateManager) updateState(state connectivity.State) {
	csm.mu.Lock()
	defer csm.mu.Unlock()
	if csm.state == connectivity.Shutdown {
		return
	}
	if csm.state == state {
		return
	}
	csm.state = state
	csm.channelz.ChannelMetrics.State.Store(&state)
	csm.pubSub.Publish(state)

	channelz.Infof(logger, csm.channelz, "Channel Connectivity change to %v", state)
	if csm.notifyChan != nil {
		
		close(csm.notifyChan)
		csm.notifyChan = nil
	}
}

func (csm *connectivityStateManager) getState() connectivity.State {
	csm.mu.Lock()
	defer csm.mu.Unlock()
	return csm.state
}

func (csm *connectivityStateManager) getNotifyChan() <-chan struct{} {
	csm.mu.Lock()
	defer csm.mu.Unlock()
	if csm.notifyChan == nil {
		csm.notifyChan = make(chan struct{})
	}
	return csm.notifyChan
}




type ClientConnInterface interface {
	
	
	Invoke(ctx context.Context, method string, args any, reply any, opts ...CallOption) error
	
	NewStream(ctx context.Context, desc *StreamDesc, method string, opts ...CallOption) (ClientStream, error)
}


var _ ClientConnInterface = (*ClientConn)(nil)













type ClientConn struct {
	ctx    context.Context    
	cancel context.CancelFunc 

	
	target              string            
	parsedTarget        resolver.Target   
	authority           string            
	dopts               dialOptions       
	channelz            *channelz.Channel 
	resolverBuilder     resolver.Builder  
	idlenessMgr         *idle.Manager
	metricsRecorderList *stats.MetricsRecorderList

	
	
	csMgr              *connectivityStateManager
	pickerWrapper      *pickerWrapper
	safeConfigSelector iresolver.SafeConfigSelector
	retryThrottler     atomic.Value 

	
	
	mu              sync.RWMutex
	resolverWrapper *ccResolverWrapper         
	balancerWrapper *ccBalancerWrapper         
	sc              *ServiceConfig             
	conns           map[*addrConn]struct{}     
	keepaliveParams keepalive.ClientParameters 
	
	
	
	
	
	firstResolveEvent *grpcsync.Event

	lceMu               sync.Mutex 
	lastConnectionError error
}



func (cc *ClientConn) WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool {
	ch := cc.csMgr.getNotifyChan()
	if cc.csMgr.getState() != sourceState {
		return true
	}
	select {
	case <-ctx.Done():
		return false
	case <-ch:
		return true
	}
}


func (cc *ClientConn) GetState() connectivity.State {
	return cc.csMgr.getState()
}









func (cc *ClientConn) Connect() {
	if err := cc.idlenessMgr.ExitIdleMode(); err != nil {
		cc.addTraceEvent(err.Error())
		return
	}
	
	
	cc.mu.Lock()
	cc.balancerWrapper.exitIdle()
	cc.mu.Unlock()
}




func (cc *ClientConn) waitForResolvedAddrs(ctx context.Context) error {
	
	
	if cc.firstResolveEvent.HasFired() {
		return nil
	}
	select {
	case <-cc.firstResolveEvent.Done():
		return nil
	case <-ctx.Done():
		return status.FromContextError(ctx.Err()).Err()
	case <-cc.ctx.Done():
		return ErrClientConnClosing
	}
}

var emptyServiceConfig *ServiceConfig

func init() {
	cfg := parseServiceConfig("{}", defaultMaxCallAttempts)
	if cfg.Err != nil {
		panic(fmt.Sprintf("impossible error parsing empty service config: %v", cfg.Err))
	}
	emptyServiceConfig = cfg.Config.(*ServiceConfig)

	internal.SubscribeToConnectivityStateChanges = func(cc *ClientConn, s grpcsync.Subscriber) func() {
		return cc.csMgr.pubSub.Subscribe(s)
	}
	internal.EnterIdleModeForTesting = func(cc *ClientConn) {
		cc.idlenessMgr.EnterIdleModeForTesting()
	}
	internal.ExitIdleModeForTesting = func(cc *ClientConn) error {
		return cc.idlenessMgr.ExitIdleMode()
	}
}

func (cc *ClientConn) maybeApplyDefaultServiceConfig() {
	if cc.sc != nil {
		cc.applyServiceConfigAndBalancer(cc.sc, nil)
		return
	}
	if cc.dopts.defaultServiceConfig != nil {
		cc.applyServiceConfigAndBalancer(cc.dopts.defaultServiceConfig, &defaultConfigSelector{cc.dopts.defaultServiceConfig})
	} else {
		cc.applyServiceConfigAndBalancer(emptyServiceConfig, &defaultConfigSelector{emptyServiceConfig})
	}
}

func (cc *ClientConn) updateResolverStateAndUnlock(s resolver.State, err error) error {
	defer cc.firstResolveEvent.Fire()
	
	
	
	if cc.conns == nil {
		cc.mu.Unlock()
		return nil
	}

	if err != nil {
		
		
		
		cc.maybeApplyDefaultServiceConfig()

		cc.balancerWrapper.resolverError(err)

		
		cc.mu.Unlock()
		return balancer.ErrBadResolverState
	}

	var ret error
	if cc.dopts.disableServiceConfig {
		channelz.Infof(logger, cc.channelz, "ignoring service config from resolver (%v) and applying the default because service config is disabled", s.ServiceConfig)
		cc.maybeApplyDefaultServiceConfig()
	} else if s.ServiceConfig == nil {
		cc.maybeApplyDefaultServiceConfig()
		
		
	} else {
		if sc, ok := s.ServiceConfig.Config.(*ServiceConfig); s.ServiceConfig.Err == nil && ok {
			configSelector := iresolver.GetConfigSelector(s)
			if configSelector != nil {
				if len(s.ServiceConfig.Config.(*ServiceConfig).Methods) != 0 {
					channelz.Infof(logger, cc.channelz, "method configs in service config will be ignored due to presence of config selector")
				}
			} else {
				configSelector = &defaultConfigSelector{sc}
			}
			cc.applyServiceConfigAndBalancer(sc, configSelector)
		} else {
			ret = balancer.ErrBadResolverState
			if cc.sc == nil {
				
				
				cc.applyFailingLBLocked(s.ServiceConfig)
				cc.mu.Unlock()
				return ret
			}
		}
	}

	balCfg := cc.sc.lbConfig
	bw := cc.balancerWrapper
	cc.mu.Unlock()

	uccsErr := bw.updateClientConnState(&balancer.ClientConnState{ResolverState: s, BalancerConfig: balCfg})
	if ret == nil {
		ret = uccsErr 
		
	}
	return ret
}







func (cc *ClientConn) applyFailingLBLocked(sc *serviceconfig.ParseResult) {
	var err error
	if sc.Err != nil {
		err = status.Errorf(codes.Unavailable, "error parsing service config: %v", sc.Err)
	} else {
		err = status.Errorf(codes.Unavailable, "illegal service config type: %T", sc.Config)
	}
	cc.safeConfigSelector.UpdateConfigSelector(&defaultConfigSelector{nil})
	cc.pickerWrapper.updatePicker(base.NewErrPicker(err))
	cc.csMgr.updateState(connectivity.TransientFailure)
}



func copyAddresses(in []resolver.Address) []resolver.Address {
	out := make([]resolver.Address, len(in))
	copy(out, in)
	return out
}




func (cc *ClientConn) newAddrConnLocked(addrs []resolver.Address, opts balancer.NewSubConnOptions) (*addrConn, error) {
	if cc.conns == nil {
		return nil, ErrClientConnClosing
	}

	ac := &addrConn{
		state:        connectivity.Idle,
		cc:           cc,
		addrs:        copyAddresses(addrs),
		scopts:       opts,
		dopts:        cc.dopts,
		channelz:     channelz.RegisterSubChannel(cc.channelz, ""),
		resetBackoff: make(chan struct{}),
	}
	ac.ctx, ac.cancel = context.WithCancel(cc.ctx)
	
	
	ac.channelz.ChannelMetrics.Target.Store(&addrs[0].Addr)

	channelz.AddTraceEvent(logger, ac.channelz, 0, &channelz.TraceEvent{
		Desc:     "Subchannel created",
		Severity: channelz.CtInfo,
		Parent: &channelz.TraceEvent{
			Desc:     fmt.Sprintf("Subchannel(id:%d) created", ac.channelz.ID),
			Severity: channelz.CtInfo,
		},
	})

	
	cc.conns[ac] = struct{}{}
	return ac, nil
}



func (cc *ClientConn) removeAddrConn(ac *addrConn, err error) {
	cc.mu.Lock()
	if cc.conns == nil {
		cc.mu.Unlock()
		return
	}
	delete(cc.conns, ac)
	cc.mu.Unlock()
	ac.tearDown(err)
}


func (cc *ClientConn) Target() string {
	return cc.target
}








func (cc *ClientConn) CanonicalTarget() string {
	return cc.parsedTarget.String()
}

func (cc *ClientConn) incrCallsStarted() {
	cc.channelz.ChannelMetrics.CallsStarted.Add(1)
	cc.channelz.ChannelMetrics.LastCallStartedTimestamp.Store(time.Now().UnixNano())
}

func (cc *ClientConn) incrCallsSucceeded() {
	cc.channelz.ChannelMetrics.CallsSucceeded.Add(1)
}

func (cc *ClientConn) incrCallsFailed() {
	cc.channelz.ChannelMetrics.CallsFailed.Add(1)
}




func (ac *addrConn) connect() error {
	ac.mu.Lock()
	if ac.state == connectivity.Shutdown {
		if logger.V(2) {
			logger.Infof("connect called on shutdown addrConn; ignoring.")
		}
		ac.mu.Unlock()
		return errConnClosing
	}
	if ac.state != connectivity.Idle {
		if logger.V(2) {
			logger.Infof("connect called on addrConn in non-idle state (%v); ignoring.", ac.state)
		}
		ac.mu.Unlock()
		return nil
	}

	ac.resetTransportAndUnlock()
	return nil
}





func equalAddressIgnoringBalAttributes(a, b *resolver.Address) bool {
	return a.Addr == b.Addr && a.ServerName == b.ServerName &&
		a.Attributes.Equal(b.Attributes) &&
		a.Metadata == b.Metadata
}

func equalAddressesIgnoringBalAttributes(a, b []resolver.Address) bool {
	return slices.EqualFunc(a, b, func(a, b resolver.Address) bool { return equalAddressIgnoringBalAttributes(&a, &b) })
}



func (ac *addrConn) updateAddrs(addrs []resolver.Address) {
	addrs = copyAddresses(addrs)
	limit := len(addrs)
	if limit > 5 {
		limit = 5
	}
	channelz.Infof(logger, ac.channelz, "addrConn: updateAddrs addrs (%d of %d): %v", limit, len(addrs), addrs[:limit])

	ac.mu.Lock()
	if equalAddressesIgnoringBalAttributes(ac.addrs, addrs) {
		ac.mu.Unlock()
		return
	}

	ac.addrs = addrs

	if ac.state == connectivity.Shutdown ||
		ac.state == connectivity.TransientFailure ||
		ac.state == connectivity.Idle {
		
		ac.mu.Unlock()
		return
	}

	if ac.state == connectivity.Ready {
		
		for _, a := range addrs {
			a.ServerName = ac.cc.getServerName(a)
			if equalAddressIgnoringBalAttributes(&a, &ac.curAddr) {
				
				
				ac.mu.Unlock()
				return
			}
		}
	}

	
	

	ac.cancel()
	ac.ctx, ac.cancel = context.WithCancel(ac.cc.ctx)

	
	
	if ac.transport != nil {
		defer ac.transport.GracefulClose()
		ac.transport = nil
	}

	if len(addrs) == 0 {
		ac.updateConnectivityState(connectivity.Idle, nil)
	}

	
	
	go ac.resetTransportAndUnlock()
}











func (cc *ClientConn) getServerName(addr resolver.Address) string {
	if cc.dopts.authority != "" {
		return cc.dopts.authority
	}
	if addr.ServerName != "" {
		return addr.ServerName
	}
	return cc.authority
}

func getMethodConfig(sc *ServiceConfig, method string) MethodConfig {
	if sc == nil {
		return MethodConfig{}
	}
	if m, ok := sc.Methods[method]; ok {
		return m
	}
	i := strings.LastIndex(method, "/")
	if m, ok := sc.Methods[method[:i+1]]; ok {
		return m
	}
	return sc.Methods[""]
}









func (cc *ClientConn) GetMethodConfig(method string) MethodConfig {
	
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	return getMethodConfig(cc.sc, method)
}

func (cc *ClientConn) healthCheckConfig() *healthCheckConfig {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	if cc.sc == nil {
		return nil
	}
	return cc.sc.healthCheckConfig
}

func (cc *ClientConn) getTransport(ctx context.Context, failfast bool, method string) (transport.ClientTransport, balancer.PickResult, error) {
	return cc.pickerWrapper.pick(ctx, failfast, balancer.PickInfo{
		Ctx:            ctx,
		FullMethodName: method,
	})
}

func (cc *ClientConn) applyServiceConfigAndBalancer(sc *ServiceConfig, configSelector iresolver.ConfigSelector) {
	if sc == nil {
		
		return
	}
	cc.sc = sc
	if configSelector != nil {
		cc.safeConfigSelector.UpdateConfigSelector(configSelector)
	}

	if cc.sc.retryThrottling != nil {
		newThrottler := &retryThrottler{
			tokens: cc.sc.retryThrottling.MaxTokens,
			max:    cc.sc.retryThrottling.MaxTokens,
			thresh: cc.sc.retryThrottling.MaxTokens / 2,
			ratio:  cc.sc.retryThrottling.TokenRatio,
		}
		cc.retryThrottler.Store(newThrottler)
	} else {
		cc.retryThrottler.Store((*retryThrottler)(nil))
	}
}

func (cc *ClientConn) resolveNow(o resolver.ResolveNowOptions) {
	cc.mu.RLock()
	cc.resolverWrapper.resolveNow(o)
	cc.mu.RUnlock()
}

func (cc *ClientConn) resolveNowLocked(o resolver.ResolveNowOptions) {
	cc.resolverWrapper.resolveNow(o)
}














func (cc *ClientConn) ResetConnectBackoff() {
	cc.mu.Lock()
	conns := cc.conns
	cc.mu.Unlock()
	for ac := range conns {
		ac.resetConnectBackoff()
	}
}


func (cc *ClientConn) Close() error {
	defer func() {
		cc.cancel()
		<-cc.csMgr.pubSub.Done()
	}()

	
	
	cc.idlenessMgr.Close()

	cc.mu.Lock()
	if cc.conns == nil {
		cc.mu.Unlock()
		return ErrClientConnClosing
	}

	conns := cc.conns
	cc.conns = nil
	cc.csMgr.updateState(connectivity.Shutdown)

	
	
	cc.mu.Unlock()

	cc.resolverWrapper.close()
	
	
	cc.pickerWrapper.close()
	cc.balancerWrapper.close()

	<-cc.resolverWrapper.serializer.Done()
	<-cc.balancerWrapper.serializer.Done()
	var wg sync.WaitGroup
	for ac := range conns {
		wg.Add(1)
		go func(ac *addrConn) {
			defer wg.Done()
			ac.tearDown(ErrClientConnClosing)
		}(ac)
	}
	wg.Wait()
	cc.addTraceEvent("deleted")
	
	
	
	channelz.RemoveEntry(cc.channelz.ID)

	return nil
}


type addrConn struct {
	ctx    context.Context
	cancel context.CancelFunc

	cc     *ClientConn
	dopts  dialOptions
	acbw   *acBalancerWrapper
	scopts balancer.NewSubConnOptions

	
	
	
	
	transport transport.ClientTransport 

	
	
	
	
	mu      sync.Mutex
	curAddr resolver.Address   
	addrs   []resolver.Address 

	
	state connectivity.State

	backoffIdx   int 
	resetBackoff chan struct{}

	channelz *channelz.SubChannel
}


func (ac *addrConn) updateConnectivityState(s connectivity.State, lastErr error) {
	if ac.state == s {
		return
	}
	ac.state = s
	ac.channelz.ChannelMetrics.State.Store(&s)
	if lastErr == nil {
		channelz.Infof(logger, ac.channelz, "Subchannel Connectivity change to %v", s)
	} else {
		channelz.Infof(logger, ac.channelz, "Subchannel Connectivity change to %v, last error: %s", s, lastErr)
	}
	ac.acbw.updateState(s, ac.curAddr, lastErr)
}



func (ac *addrConn) adjustParams(r transport.GoAwayReason) {
	switch r {
	case transport.GoAwayTooManyPings:
		v := 2 * ac.dopts.copts.KeepaliveParams.Time
		ac.cc.mu.Lock()
		if v > ac.cc.keepaliveParams.Time {
			ac.cc.keepaliveParams.Time = v
		}
		ac.cc.mu.Unlock()
	}
}




func (ac *addrConn) resetTransportAndUnlock() {
	acCtx := ac.ctx
	if acCtx.Err() != nil {
		ac.mu.Unlock()
		return
	}

	addrs := ac.addrs
	backoffFor := ac.dopts.bs.Backoff(ac.backoffIdx)
	
	dialDuration := minConnectTimeout
	if ac.dopts.minConnectTimeout != nil {
		dialDuration = ac.dopts.minConnectTimeout()
	}

	if dialDuration < backoffFor {
		
		dialDuration = backoffFor
	}
	
	
	
	
	
	
	connectDeadline := time.Now().Add(dialDuration)

	ac.updateConnectivityState(connectivity.Connecting, nil)
	ac.mu.Unlock()

	if err := ac.tryAllAddrs(acCtx, addrs, connectDeadline); err != nil {
		
		
		ac.cc.resolveNow(resolver.ResolveNowOptions{})
		ac.mu.Lock()
		if acCtx.Err() != nil {
			
			ac.mu.Unlock()
			return
		}
		
		
		ac.updateConnectivityState(connectivity.TransientFailure, err)

		
		b := ac.resetBackoff
		ac.mu.Unlock()

		timer := time.NewTimer(backoffFor)
		select {
		case <-timer.C:
			ac.mu.Lock()
			ac.backoffIdx++
			ac.mu.Unlock()
		case <-b:
			timer.Stop()
		case <-acCtx.Done():
			timer.Stop()
			return
		}

		ac.mu.Lock()
		if acCtx.Err() == nil {
			ac.updateConnectivityState(connectivity.Idle, err)
		}
		ac.mu.Unlock()
		return
	}
	
	ac.mu.Lock()
	ac.backoffIdx = 0
	ac.mu.Unlock()
}




func (ac *addrConn) tryAllAddrs(ctx context.Context, addrs []resolver.Address, connectDeadline time.Time) error {
	var firstConnErr error
	for _, addr := range addrs {
		ac.channelz.ChannelMetrics.Target.Store(&addr.Addr)
		if ctx.Err() != nil {
			return errConnClosing
		}
		ac.mu.Lock()

		ac.cc.mu.RLock()
		ac.dopts.copts.KeepaliveParams = ac.cc.keepaliveParams
		ac.cc.mu.RUnlock()

		copts := ac.dopts.copts
		if ac.scopts.CredsBundle != nil {
			copts.CredsBundle = ac.scopts.CredsBundle
		}
		ac.mu.Unlock()

		channelz.Infof(logger, ac.channelz, "Subchannel picks a new address %q to connect", addr.Addr)

		err := ac.createTransport(ctx, addr, copts, connectDeadline)
		if err == nil {
			return nil
		}
		if firstConnErr == nil {
			firstConnErr = err
		}
		ac.cc.updateConnectionError(err)
	}

	
	return firstConnErr
}




func (ac *addrConn) createTransport(ctx context.Context, addr resolver.Address, copts transport.ConnectOptions, connectDeadline time.Time) error {
	addr.ServerName = ac.cc.getServerName(addr)
	hctx, hcancel := context.WithCancel(ctx)

	onClose := func(r transport.GoAwayReason) {
		ac.mu.Lock()
		defer ac.mu.Unlock()
		
		ac.adjustParams(r)
		if ctx.Err() != nil {
			
			
			
			
			return
		}
		hcancel()
		if ac.transport == nil {
			
			
			
			
			return
		}
		ac.transport = nil
		
		ac.cc.resolveNow(resolver.ResolveNowOptions{})
		
		
		ac.updateConnectivityState(connectivity.Idle, nil)
	}

	connectCtx, cancel := context.WithDeadline(ctx, connectDeadline)
	defer cancel()
	copts.ChannelzParent = ac.channelz

	newTr, err := transport.NewHTTP2Client(connectCtx, ac.cc.ctx, addr, copts, onClose)
	if err != nil {
		if logger.V(2) {
			logger.Infof("Creating new client transport to %q: %v", addr, err)
		}
		
		hcancel()
		channelz.Warningf(logger, ac.channelz, "grpc: addrConn.createTransport failed to connect to %s. Err: %v", addr, err)
		return err
	}

	ac.mu.Lock()
	defer ac.mu.Unlock()
	if ctx.Err() != nil {
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		go newTr.Close(transport.ErrConnClosing)
		return nil
	}
	if hctx.Err() != nil {
		
		
		
		ac.updateConnectivityState(connectivity.Idle, nil)
		return nil
	}
	ac.curAddr = addr
	ac.transport = newTr
	ac.startHealthCheck(hctx) 
	return nil
}













func (ac *addrConn) startHealthCheck(ctx context.Context) {
	var healthcheckManagingState bool
	defer func() {
		if !healthcheckManagingState {
			ac.updateConnectivityState(connectivity.Ready, nil)
		}
	}()

	if ac.cc.dopts.disableHealthCheck {
		return
	}
	healthCheckConfig := ac.cc.healthCheckConfig()
	if healthCheckConfig == nil {
		return
	}
	if !ac.scopts.HealthCheckEnabled {
		return
	}
	healthCheckFunc := internal.HealthCheckFunc
	if healthCheckFunc == nil {
		
		
		
		channelz.Error(logger, ac.channelz, "Health check is requested but health check function is not set.")
		return
	}

	healthcheckManagingState = true

	
	currentTr := ac.transport
	newStream := func(method string) (any, error) {
		ac.mu.Lock()
		if ac.transport != currentTr {
			ac.mu.Unlock()
			return nil, status.Error(codes.Canceled, "the provided transport is no longer valid to use")
		}
		ac.mu.Unlock()
		return newNonRetryClientStream(ctx, &StreamDesc{ServerStreams: true}, method, currentTr, ac)
	}
	setConnectivityState := func(s connectivity.State, lastErr error) {
		ac.mu.Lock()
		defer ac.mu.Unlock()
		if ac.transport != currentTr {
			return
		}
		ac.updateConnectivityState(s, lastErr)
	}
	
	go func() {
		err := healthCheckFunc(ctx, newStream, setConnectivityState, healthCheckConfig.ServiceName)
		if err != nil {
			if status.Code(err) == codes.Unimplemented {
				channelz.Error(logger, ac.channelz, "Subchannel health check is unimplemented at server side, thus health check is disabled")
			} else {
				channelz.Errorf(logger, ac.channelz, "Health checking failed: %v", err)
			}
		}
	}()
}

func (ac *addrConn) resetConnectBackoff() {
	ac.mu.Lock()
	close(ac.resetBackoff)
	ac.backoffIdx = 0
	ac.resetBackoff = make(chan struct{})
	ac.mu.Unlock()
}


func (ac *addrConn) getReadyTransport() transport.ClientTransport {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	if ac.state == connectivity.Ready {
		return ac.transport
	}
	return nil
}





func (ac *addrConn) tearDown(err error) {
	ac.mu.Lock()
	if ac.state == connectivity.Shutdown {
		ac.mu.Unlock()
		return
	}
	curTr := ac.transport
	ac.transport = nil
	
	
	ac.updateConnectivityState(connectivity.Shutdown, nil)
	ac.cancel()
	ac.curAddr = resolver.Address{}

	channelz.AddTraceEvent(logger, ac.channelz, 0, &channelz.TraceEvent{
		Desc:     "Subchannel deleted",
		Severity: channelz.CtInfo,
		Parent: &channelz.TraceEvent{
			Desc:     fmt.Sprintf("Subchannel(id:%d) deleted", ac.channelz.ID),
			Severity: channelz.CtInfo,
		},
	})
	
	
	
	channelz.RemoveEntry(ac.channelz.ID)
	ac.mu.Unlock()

	
	
	if curTr != nil {
		if err == errConnDrain {
			
			
			
			
			
			
			curTr.GracefulClose()
		} else {
			
			
			
			
			
			
			
			curTr.Close(err)
		}
	}
}

type retryThrottler struct {
	max    float64
	thresh float64
	ratio  float64

	mu     sync.Mutex
	tokens float64 
}




func (rt *retryThrottler) throttle() bool {
	if rt == nil {
		return false
	}
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.tokens--
	if rt.tokens < 0 {
		rt.tokens = 0
	}
	return rt.tokens <= rt.thresh
}

func (rt *retryThrottler) successfulRPC() {
	if rt == nil {
		return
	}
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.tokens += rt.ratio
	if rt.tokens > rt.max {
		rt.tokens = rt.max
	}
}

func (ac *addrConn) incrCallsStarted() {
	ac.channelz.ChannelMetrics.CallsStarted.Add(1)
	ac.channelz.ChannelMetrics.LastCallStartedTimestamp.Store(time.Now().UnixNano())
}

func (ac *addrConn) incrCallsSucceeded() {
	ac.channelz.ChannelMetrics.CallsSucceeded.Add(1)
}

func (ac *addrConn) incrCallsFailed() {
	ac.channelz.ChannelMetrics.CallsFailed.Add(1)
}






var ErrClientConnTimeout = errors.New("grpc: timed out when dialing")




func (cc *ClientConn) getResolver(scheme string) resolver.Builder {
	for _, rb := range cc.dopts.resolvers {
		if scheme == rb.Scheme() {
			return rb
		}
	}
	return resolver.Get(scheme)
}

func (cc *ClientConn) updateConnectionError(err error) {
	cc.lceMu.Lock()
	cc.lastConnectionError = err
	cc.lceMu.Unlock()
}

func (cc *ClientConn) connectionError() error {
	cc.lceMu.Lock()
	defer cc.lceMu.Unlock()
	return cc.lastConnectionError
}








func (cc *ClientConn) initParsedTargetAndResolverBuilder() error {
	logger.Infof("original dial target is: %q", cc.target)

	var rb resolver.Builder
	parsedTarget, err := parseTarget(cc.target)
	if err == nil {
		rb = cc.getResolver(parsedTarget.URL.Scheme)
		if rb != nil {
			cc.parsedTarget = parsedTarget
			cc.resolverBuilder = rb
			return nil
		}
	}

	
	
	
	
	
	defScheme := cc.dopts.defaultScheme
	if internal.UserSetDefaultScheme {
		defScheme = resolver.GetDefaultScheme()
	}

	canonicalTarget := defScheme + ":///" + cc.target

	parsedTarget, err = parseTarget(canonicalTarget)
	if err != nil {
		return err
	}
	rb = cc.getResolver(parsedTarget.URL.Scheme)
	if rb == nil {
		return fmt.Errorf("could not get resolver for default scheme: %q", parsedTarget.URL.Scheme)
	}
	cc.parsedTarget = parsedTarget
	cc.resolverBuilder = rb
	return nil
}




func parseTarget(target string) (resolver.Target, error) {
	u, err := url.Parse(target)
	if err != nil {
		return resolver.Target{}, err
	}

	return resolver.Target{URL: *u}, nil
}



func encodeAuthority(authority string) string {
	const upperhex = "0123456789ABCDEF"

	
	
	
	shouldEscape := func(c byte) bool {
		
		if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
			return false
		}
		switch c {
		case '-', '_', '.', '~': 
			return false
		case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=': 
			return false
		case ':', '[', ']', '@': 
			return false
		}
		
		return true
	}

	hexCount := 0
	for i := 0; i < len(authority); i++ {
		c := authority[i]
		if shouldEscape(c) {
			hexCount++
		}
	}

	if hexCount == 0 {
		return authority
	}

	required := len(authority) + 2*hexCount
	t := make([]byte, required)

	j := 0
	
	for i := 0; i < len(authority); i++ {
		switch c := authority[i]; {
		case shouldEscape(c):
			t[j] = '%'
			t[j+1] = upperhex[c>>4]
			t[j+2] = upperhex[c&15]
			j += 3
		default:
			t[j] = authority[i]
			j++
		}
	}
	return string(t)
}












func (cc *ClientConn) initAuthority() error {
	dopts := cc.dopts
	
	
	
	
	
	
	
	
	
	
	authorityFromCreds := ""
	if creds := dopts.copts.TransportCredentials; creds != nil && creds.Info().ServerName != "" {
		authorityFromCreds = creds.Info().ServerName
	}
	authorityFromDialOption := dopts.authority
	if (authorityFromCreds != "" && authorityFromDialOption != "") && authorityFromCreds != authorityFromDialOption {
		return fmt.Errorf("ClientConn's authority from transport creds %q and dial option %q don't match", authorityFromCreds, authorityFromDialOption)
	}

	endpoint := cc.parsedTarget.Endpoint()
	if authorityFromDialOption != "" {
		cc.authority = authorityFromDialOption
	} else if authorityFromCreds != "" {
		cc.authority = authorityFromCreds
	} else if auth, ok := cc.resolverBuilder.(resolver.AuthorityOverrider); ok {
		cc.authority = auth.OverrideAuthority(cc.parsedTarget)
	} else if strings.HasPrefix(endpoint, ":") {
		cc.authority = "localhost" + endpoint
	} else {
		cc.authority = encodeAuthority(endpoint)
	}
	return nil
}
