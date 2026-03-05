








package pickfirstleaf

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"sync"
	"time"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/pickfirst/internal"
	"google.golang.org/grpc/connectivity"
	expstats "google.golang.org/grpc/experimental/stats"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/envconfig"
	internalgrpclog "google.golang.org/grpc/internal/grpclog"
	"google.golang.org/grpc/internal/pretty"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

func init() {
	if envconfig.NewPickFirstEnabled {
		
		Name = "pick_first"
	}
	balancer.Register(pickfirstBuilder{})
}

type (
	
	
	enableHealthListenerKeyType struct{}
	
	
	
	
	
	
	managedByPickfirstKeyType struct{}
)

var (
	logger = grpclog.Component("pick-first-leaf-lb")
	
	
	
	Name                 = "pick_first_leaf"
	disconnectionsMetric = expstats.RegisterInt64Count(expstats.MetricDescriptor{
		Name:        "grpc.lb.pick_first.disconnections",
		Description: "EXPERIMENTAL. Number of times the selected subchannel becomes disconnected.",
		Unit:        "disconnection",
		Labels:      []string{"grpc.target"},
		Default:     false,
	})
	connectionAttemptsSucceededMetric = expstats.RegisterInt64Count(expstats.MetricDescriptor{
		Name:        "grpc.lb.pick_first.connection_attempts_succeeded",
		Description: "EXPERIMENTAL. Number of successful connection attempts.",
		Unit:        "attempt",
		Labels:      []string{"grpc.target"},
		Default:     false,
	})
	connectionAttemptsFailedMetric = expstats.RegisterInt64Count(expstats.MetricDescriptor{
		Name:        "grpc.lb.pick_first.connection_attempts_failed",
		Description: "EXPERIMENTAL. Number of failed connection attempts.",
		Unit:        "attempt",
		Labels:      []string{"grpc.target"},
		Default:     false,
	})
)

const (
	
	logPrefix = "[pick-first-leaf-lb %p] "
	
	
	connectionDelayInterval = 250 * time.Millisecond
)

type ipAddrFamily int

const (
	
	
	ipAddrFamilyUnknown ipAddrFamily = iota
	ipAddrFamilyV4
	ipAddrFamilyV6
)

type pickfirstBuilder struct{}

func (pickfirstBuilder) Build(cc balancer.ClientConn, bo balancer.BuildOptions) balancer.Balancer {
	b := &pickfirstBalancer{
		cc:              cc,
		target:          bo.Target.String(),
		metricsRecorder: cc.MetricsRecorder(),

		subConns:              resolver.NewAddressMap(),
		state:                 connectivity.Connecting,
		cancelConnectionTimer: func() {},
	}
	b.logger = internalgrpclog.NewPrefixLogger(logger, fmt.Sprintf(logPrefix, b))
	return b
}

func (b pickfirstBuilder) Name() string {
	return Name
}

func (pickfirstBuilder) ParseConfig(js json.RawMessage) (serviceconfig.LoadBalancingConfig, error) {
	var cfg pfConfig
	if err := json.Unmarshal(js, &cfg); err != nil {
		return nil, fmt.Errorf("pickfirst: unable to unmarshal LB policy config: %s, error: %v", string(js), err)
	}
	return cfg, nil
}



func EnableHealthListener(state resolver.State) resolver.State {
	state.Attributes = state.Attributes.WithValue(enableHealthListenerKeyType{}, true)
	return state
}








func IsManagedByPickfirst(addr resolver.Address) bool {
	return addr.BalancerAttributes.Value(managedByPickfirstKeyType{}) != nil
}

type pfConfig struct {
	serviceconfig.LoadBalancingConfig `json:"-"`

	
	
	
	ShuffleAddressList bool `json:"shuffleAddressList"`
}



type scData struct {
	
	
	subConn balancer.SubConn
	addr    resolver.Address

	rawConnectivityState connectivity.State
	
	
	effectiveState              connectivity.State
	lastErr                     error
	connectionFailedInFirstPass bool
}

func (b *pickfirstBalancer) newSCData(addr resolver.Address) (*scData, error) {
	addr.BalancerAttributes = addr.BalancerAttributes.WithValue(managedByPickfirstKeyType{}, true)
	sd := &scData{
		rawConnectivityState: connectivity.Idle,
		effectiveState:       connectivity.Idle,
		addr:                 addr,
	}
	sc, err := b.cc.NewSubConn([]resolver.Address{addr}, balancer.NewSubConnOptions{
		StateListener: func(state balancer.SubConnState) {
			b.updateSubConnState(sd, state)
		},
	})
	if err != nil {
		return nil, err
	}
	sd.subConn = sc
	return sd, nil
}

type pickfirstBalancer struct {
	
	
	logger          *internalgrpclog.PrefixLogger
	cc              balancer.ClientConn
	target          string
	metricsRecorder expstats.MetricsRecorder 

	
	
	
	mu sync.Mutex
	
	
	state connectivity.State
	
	subConns              *resolver.AddressMap
	addressList           addressList
	firstPass             bool
	numTF                 int
	cancelConnectionTimer func()
	healthCheckingEnabled bool
}



func (b *pickfirstBalancer) ResolverError(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.resolverErrorLocked(err)
}

func (b *pickfirstBalancer) resolverErrorLocked(err error) {
	if b.logger.V(2) {
		b.logger.Infof("Received error from the name resolver: %v", err)
	}

	
	
	
	if b.state != connectivity.TransientFailure && b.addressList.size() > 0 {
		if b.logger.V(2) {
			b.logger.Infof("Ignoring resolver error because balancer is using a previous good update.")
		}
		return
	}

	b.updateBalancerState(balancer.State{
		ConnectivityState: connectivity.TransientFailure,
		Picker:            &picker{err: fmt.Errorf("name resolver error: %v", err)},
	})
}

func (b *pickfirstBalancer) UpdateClientConnState(state balancer.ClientConnState) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.cancelConnectionTimer()
	if len(state.ResolverState.Addresses) == 0 && len(state.ResolverState.Endpoints) == 0 {
		
		
		b.closeSubConnsLocked()
		b.addressList.updateAddrs(nil)
		b.resolverErrorLocked(errors.New("produced zero addresses"))
		return balancer.ErrBadResolverState
	}
	b.healthCheckingEnabled = state.ResolverState.Attributes.Value(enableHealthListenerKeyType{}) != nil
	cfg, ok := state.BalancerConfig.(pfConfig)
	if state.BalancerConfig != nil && !ok {
		return fmt.Errorf("pickfirst: received illegal BalancerConfig (type %T): %v: %w", state.BalancerConfig, state.BalancerConfig, balancer.ErrBadResolverState)
	}

	if b.logger.V(2) {
		b.logger.Infof("Received new config %s, resolver state %s", pretty.ToJSON(cfg), pretty.ToJSON(state.ResolverState))
	}

	var newAddrs []resolver.Address
	if endpoints := state.ResolverState.Endpoints; len(endpoints) != 0 {
		
		
		
		if cfg.ShuffleAddressList {
			endpoints = append([]resolver.Endpoint{}, endpoints...)
			internal.RandShuffle(len(endpoints), func(i, j int) { endpoints[i], endpoints[j] = endpoints[j], endpoints[i] })
		}

		
		
		for _, endpoint := range endpoints {
			newAddrs = append(newAddrs, endpoint.Addresses...)
		}
	} else {
		
		
		
		
		
		
		newAddrs = state.ResolverState.Addresses
		if cfg.ShuffleAddressList {
			newAddrs = append([]resolver.Address{}, newAddrs...)
			internal.RandShuffle(len(endpoints), func(i, j int) { endpoints[i], endpoints[j] = endpoints[j], endpoints[i] })
		}
	}

	
	
	
	
	
	newAddrs = deDupAddresses(newAddrs)
	newAddrs = interleaveAddresses(newAddrs)

	prevAddr := b.addressList.currentAddress()
	prevSCData, found := b.subConns.Get(prevAddr)
	prevAddrsCount := b.addressList.size()
	isPrevRawConnectivityStateReady := found && prevSCData.(*scData).rawConnectivityState == connectivity.Ready
	b.addressList.updateAddrs(newAddrs)

	
	
	if isPrevRawConnectivityStateReady && b.addressList.seekTo(prevAddr) {
		return nil
	}

	b.reconcileSubConnsLocked(newAddrs)
	
	
	
	
	
	
	
	if isPrevRawConnectivityStateReady || b.state == connectivity.Connecting || prevAddrsCount == 0 {
		
		b.forceUpdateConcludedStateLocked(balancer.State{
			ConnectivityState: connectivity.Connecting,
			Picker:            &picker{err: balancer.ErrNoSubConnAvailable},
		})
		b.startFirstPassLocked()
	} else if b.state == connectivity.TransientFailure {
		
		
		b.startFirstPassLocked()
	}
	return nil
}



func (b *pickfirstBalancer) UpdateSubConnState(subConn balancer.SubConn, state balancer.SubConnState) {
	b.logger.Errorf("UpdateSubConnState(%v, %+v) called unexpectedly", subConn, state)
}

func (b *pickfirstBalancer) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.closeSubConnsLocked()
	b.cancelConnectionTimer()
	b.state = connectivity.Shutdown
}




func (b *pickfirstBalancer) ExitIdle() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == connectivity.Idle {
		b.startFirstPassLocked()
	}
}

func (b *pickfirstBalancer) startFirstPassLocked() {
	b.firstPass = true
	b.numTF = 0
	
	for _, sd := range b.subConns.Values() {
		sd.(*scData).connectionFailedInFirstPass = false
	}
	b.requestConnectionLocked()
}

func (b *pickfirstBalancer) closeSubConnsLocked() {
	for _, sd := range b.subConns.Values() {
		sd.(*scData).subConn.Shutdown()
	}
	b.subConns = resolver.NewAddressMap()
}


func deDupAddresses(addrs []resolver.Address) []resolver.Address {
	seenAddrs := resolver.NewAddressMap()
	retAddrs := []resolver.Address{}

	for _, addr := range addrs {
		if _, ok := seenAddrs.Get(addr); ok {
			continue
		}
		retAddrs = append(retAddrs, addr)
	}
	return retAddrs
}












func interleaveAddresses(addrs []resolver.Address) []resolver.Address {
	familyAddrsMap := map[ipAddrFamily][]resolver.Address{}
	interleavingOrder := []ipAddrFamily{}
	for _, addr := range addrs {
		family := addressFamily(addr.Addr)
		if _, found := familyAddrsMap[family]; !found {
			interleavingOrder = append(interleavingOrder, family)
		}
		familyAddrsMap[family] = append(familyAddrsMap[family], addr)
	}

	interleavedAddrs := make([]resolver.Address, 0, len(addrs))

	for curFamilyIdx := 0; len(interleavedAddrs) < len(addrs); curFamilyIdx = (curFamilyIdx + 1) % len(interleavingOrder) {
		
		
		
		family := interleavingOrder[curFamilyIdx]
		remainingMembers := familyAddrsMap[family]
		if len(remainingMembers) > 0 {
			interleavedAddrs = append(interleavedAddrs, remainingMembers[0])
			familyAddrsMap[family] = remainingMembers[1:]
		}
	}

	return interleavedAddrs
}






func addressFamily(address string) ipAddrFamily {
	
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return ipAddrFamilyUnknown
	}
	ip, err := netip.ParseAddr(host)
	if err != nil {
		return ipAddrFamilyUnknown
	}
	switch {
	case ip.Is4() || ip.Is4In6():
		return ipAddrFamilyV4
	case ip.Is6():
		return ipAddrFamilyV6
	default:
		return ipAddrFamilyUnknown
	}
}










func (b *pickfirstBalancer) reconcileSubConnsLocked(newAddrs []resolver.Address) {
	newAddrsMap := resolver.NewAddressMap()
	for _, addr := range newAddrs {
		newAddrsMap.Set(addr, true)
	}

	for _, oldAddr := range b.subConns.Keys() {
		if _, ok := newAddrsMap.Get(oldAddr); ok {
			continue
		}
		val, _ := b.subConns.Get(oldAddr)
		val.(*scData).subConn.Shutdown()
		b.subConns.Delete(oldAddr)
	}
}



func (b *pickfirstBalancer) shutdownRemainingLocked(selected *scData) {
	b.cancelConnectionTimer()
	for _, v := range b.subConns.Values() {
		sd := v.(*scData)
		if sd.subConn != selected.subConn {
			sd.subConn.Shutdown()
		}
	}
	b.subConns = resolver.NewAddressMap()
	b.subConns.Set(selected.addr, selected)
}





func (b *pickfirstBalancer) requestConnectionLocked() {
	if !b.addressList.isValid() {
		return
	}
	var lastErr error
	for valid := true; valid; valid = b.addressList.increment() {
		curAddr := b.addressList.currentAddress()
		sd, ok := b.subConns.Get(curAddr)
		if !ok {
			var err error
			
			
			sd, err = b.newSCData(curAddr)
			if err != nil {
				
				
				if b.logger.V(2) {
					b.logger.Infof("Failed to create a subConn for address %v: %v", curAddr.String(), err)
				}
				
				return
			}
			b.subConns.Set(curAddr, sd)
		}

		scd := sd.(*scData)
		switch scd.rawConnectivityState {
		case connectivity.Idle:
			scd.subConn.Connect()
			b.scheduleNextConnectionLocked()
			return
		case connectivity.TransientFailure:
			
			
			
			scd.connectionFailedInFirstPass = true
			lastErr = scd.lastErr
			continue
		case connectivity.Connecting:
			
			
			b.scheduleNextConnectionLocked()
			return
		default:
			b.logger.Errorf("SubConn with unexpected state %v present in SubConns map.", scd.rawConnectivityState)
			return

		}
	}

	
	
	b.endFirstPassIfPossibleLocked(lastErr)
}

func (b *pickfirstBalancer) scheduleNextConnectionLocked() {
	b.cancelConnectionTimer()
	if !b.addressList.hasNext() {
		return
	}
	curAddr := b.addressList.currentAddress()
	cancelled := false 
	closeFn := internal.TimeAfterFunc(connectionDelayInterval, func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		
		if cancelled {
			return
		}
		if b.logger.V(2) {
			b.logger.Infof("Happy Eyeballs timer expired while waiting for connection to %q.", curAddr.Addr)
		}
		if b.addressList.increment() {
			b.requestConnectionLocked()
		}
	})
	
	
	b.cancelConnectionTimer = sync.OnceFunc(func() {
		cancelled = true
		closeFn()
	})
}

func (b *pickfirstBalancer) updateSubConnState(sd *scData, newState balancer.SubConnState) {
	b.mu.Lock()
	defer b.mu.Unlock()
	oldState := sd.rawConnectivityState
	sd.rawConnectivityState = newState.ConnectivityState
	
	
	
	
	if !b.isActiveSCData(sd) {
		return
	}
	if newState.ConnectivityState == connectivity.Shutdown {
		sd.effectiveState = connectivity.Shutdown
		return
	}

	
	if newState.ConnectivityState == connectivity.TransientFailure {
		sd.connectionFailedInFirstPass = true
		connectionAttemptsFailedMetric.Record(b.metricsRecorder, 1, b.target)
	}

	if newState.ConnectivityState == connectivity.Ready {
		connectionAttemptsSucceededMetric.Record(b.metricsRecorder, 1, b.target)
		b.shutdownRemainingLocked(sd)
		if !b.addressList.seekTo(sd.addr) {
			
			
			b.logger.Errorf("Address %q not found address list in  %v", sd.addr, b.addressList.addresses)
			return
		}
		if !b.healthCheckingEnabled {
			if b.logger.V(2) {
				b.logger.Infof("SubConn %p reported connectivity state READY and the health listener is disabled. Transitioning SubConn to READY.", sd.subConn)
			}

			sd.effectiveState = connectivity.Ready
			b.updateBalancerState(balancer.State{
				ConnectivityState: connectivity.Ready,
				Picker:            &picker{result: balancer.PickResult{SubConn: sd.subConn}},
			})
			return
		}
		if b.logger.V(2) {
			b.logger.Infof("SubConn %p reported connectivity state READY. Registering health listener.", sd.subConn)
		}
		
		
		sd.effectiveState = connectivity.Connecting
		b.updateBalancerState(balancer.State{
			ConnectivityState: connectivity.Connecting,
			Picker:            &picker{err: balancer.ErrNoSubConnAvailable},
		})
		sd.subConn.RegisterHealthListener(func(scs balancer.SubConnState) {
			b.updateSubConnHealthState(sd, scs)
		})
		return
	}

	
	
	
	
	
	
	
	
	if oldState == connectivity.Ready || (oldState == connectivity.Connecting && newState.ConnectivityState == connectivity.Idle) {
		
		
		b.shutdownRemainingLocked(sd)
		sd.effectiveState = newState.ConnectivityState
		
		
		if oldState == connectivity.Connecting {
			
			
			
			connectionAttemptsSucceededMetric.Record(b.metricsRecorder, 1, b.target)
		}
		disconnectionsMetric.Record(b.metricsRecorder, 1, b.target)
		b.addressList.reset()
		b.updateBalancerState(balancer.State{
			ConnectivityState: connectivity.Idle,
			Picker:            &idlePicker{exitIdle: sync.OnceFunc(b.ExitIdle)},
		})
		return
	}

	if b.firstPass {
		switch newState.ConnectivityState {
		case connectivity.Connecting:
			
			
			
			if sd.effectiveState != connectivity.TransientFailure {
				sd.effectiveState = connectivity.Connecting
				b.updateBalancerState(balancer.State{
					ConnectivityState: connectivity.Connecting,
					Picker:            &picker{err: balancer.ErrNoSubConnAvailable},
				})
			}
		case connectivity.TransientFailure:
			sd.lastErr = newState.ConnectionError
			sd.effectiveState = connectivity.TransientFailure
			
			
			
			

			if curAddr := b.addressList.currentAddress(); equalAddressIgnoringBalAttributes(&curAddr, &sd.addr) {
				b.cancelConnectionTimer()
				if b.addressList.increment() {
					b.requestConnectionLocked()
					return
				}
			}

			
			
			b.endFirstPassIfPossibleLocked(newState.ConnectionError)
		}
		return
	}

	
	switch newState.ConnectivityState {
	case connectivity.TransientFailure:
		b.numTF = (b.numTF + 1) % b.subConns.Len()
		sd.lastErr = newState.ConnectionError
		if b.numTF%b.subConns.Len() == 0 {
			b.updateBalancerState(balancer.State{
				ConnectivityState: connectivity.TransientFailure,
				Picker:            &picker{err: newState.ConnectionError},
			})
		}
		
		
		
		
	case connectivity.Idle:
		sd.subConn.Connect()
	}
}



func (b *pickfirstBalancer) endFirstPassIfPossibleLocked(lastErr error) {
	
	if b.addressList.isValid() {
		return
	}
	
	
	for _, v := range b.subConns.Values() {
		sd := v.(*scData)
		if !sd.connectionFailedInFirstPass {
			return
		}
	}
	b.firstPass = false
	b.updateBalancerState(balancer.State{
		ConnectivityState: connectivity.TransientFailure,
		Picker:            &picker{err: lastErr},
	})
	
	for _, v := range b.subConns.Values() {
		sd := v.(*scData)
		if sd.rawConnectivityState == connectivity.Idle {
			sd.subConn.Connect()
		}
	}
}

func (b *pickfirstBalancer) isActiveSCData(sd *scData) bool {
	activeSD, found := b.subConns.Get(sd.addr)
	return found && activeSD == sd
}

func (b *pickfirstBalancer) updateSubConnHealthState(sd *scData, state balancer.SubConnState) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	
	
	
	if !b.isActiveSCData(sd) {
		return
	}
	sd.effectiveState = state.ConnectivityState
	switch state.ConnectivityState {
	case connectivity.Ready:
		b.updateBalancerState(balancer.State{
			ConnectivityState: connectivity.Ready,
			Picker:            &picker{result: balancer.PickResult{SubConn: sd.subConn}},
		})
	case connectivity.TransientFailure:
		b.updateBalancerState(balancer.State{
			ConnectivityState: connectivity.TransientFailure,
			Picker:            &picker{err: fmt.Errorf("pickfirst: health check failure: %v", state.ConnectionError)},
		})
	case connectivity.Connecting:
		b.updateBalancerState(balancer.State{
			ConnectivityState: connectivity.Connecting,
			Picker:            &picker{err: balancer.ErrNoSubConnAvailable},
		})
	default:
		b.logger.Errorf("Got unexpected health update for SubConn %p: %v", state)
	}
}




func (b *pickfirstBalancer) updateBalancerState(newState balancer.State) {
	
	
	
	if newState.ConnectivityState == b.state && b.state != connectivity.TransientFailure {
		return
	}
	b.forceUpdateConcludedStateLocked(newState)
}






func (b *pickfirstBalancer) forceUpdateConcludedStateLocked(newState balancer.State) {
	b.state = newState.ConnectivityState
	b.cc.UpdateState(newState)
}

type picker struct {
	result balancer.PickResult
	err    error
}

func (p *picker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	return p.result, p.err
}



type idlePicker struct {
	exitIdle func()
}

func (i *idlePicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	i.exitIdle()
	return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
}





type addressList struct {
	addresses []resolver.Address
	idx       int
}

func (al *addressList) isValid() bool {
	return al.idx < len(al.addresses)
}

func (al *addressList) size() int {
	return len(al.addresses)
}



func (al *addressList) increment() bool {
	if !al.isValid() {
		return false
	}
	al.idx++
	return al.idx < len(al.addresses)
}



func (al *addressList) currentAddress() resolver.Address {
	if !al.isValid() {
		return resolver.Address{}
	}
	return al.addresses[al.idx]
}

func (al *addressList) reset() {
	al.idx = 0
}

func (al *addressList) updateAddrs(addrs []resolver.Address) {
	al.addresses = addrs
	al.reset()
}



func (al *addressList) seekTo(needle resolver.Address) bool {
	for ai, addr := range al.addresses {
		if !equalAddressIgnoringBalAttributes(&addr, &needle) {
			continue
		}
		al.idx = ai
		return true
	}
	return false
}




func (al *addressList) hasNext() bool {
	if !al.isValid() {
		return false
	}
	return al.idx+1 < len(al.addresses)
}





func equalAddressIgnoringBalAttributes(a, b *resolver.Address) bool {
	return a.Addr == b.Addr && a.ServerName == b.ServerName &&
		a.Attributes.Equal(b.Attributes) &&
		a.Metadata == b.Metadata
}
