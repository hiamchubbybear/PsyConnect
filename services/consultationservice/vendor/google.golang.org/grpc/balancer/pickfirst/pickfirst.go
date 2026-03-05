


package pickfirst

import (
	"encoding/json"
	"errors"
	"fmt"
	rand "math/rand/v2"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/pickfirst/internal"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/envconfig"
	internalgrpclog "google.golang.org/grpc/internal/grpclog"
	"google.golang.org/grpc/internal/pretty"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"

	_ "google.golang.org/grpc/balancer/pickfirst/pickfirstleaf" 
)

func init() {
	if envconfig.NewPickFirstEnabled {
		return
	}
	balancer.Register(pickfirstBuilder{})
}

var logger = grpclog.Component("pick-first-lb")

const (
	
	Name      = "pick_first"
	logPrefix = "[pick-first-lb %p] "
)

type pickfirstBuilder struct{}

func (pickfirstBuilder) Build(cc balancer.ClientConn, _ balancer.BuildOptions) balancer.Balancer {
	b := &pickfirstBalancer{cc: cc}
	b.logger = internalgrpclog.NewPrefixLogger(logger, fmt.Sprintf(logPrefix, b))
	return b
}

func (pickfirstBuilder) Name() string {
	return Name
}

type pfConfig struct {
	serviceconfig.LoadBalancingConfig `json:"-"`

	
	
	
	ShuffleAddressList bool `json:"shuffleAddressList"`
}

func (pickfirstBuilder) ParseConfig(js json.RawMessage) (serviceconfig.LoadBalancingConfig, error) {
	var cfg pfConfig
	if err := json.Unmarshal(js, &cfg); err != nil {
		return nil, fmt.Errorf("pickfirst: unable to unmarshal LB policy config: %s, error: %v", string(js), err)
	}
	return cfg, nil
}

type pickfirstBalancer struct {
	logger  *internalgrpclog.PrefixLogger
	state   connectivity.State
	cc      balancer.ClientConn
	subConn balancer.SubConn
}

func (b *pickfirstBalancer) ResolverError(err error) {
	if b.logger.V(2) {
		b.logger.Infof("Received error from the name resolver: %v", err)
	}
	if b.subConn == nil {
		b.state = connectivity.TransientFailure
	}

	if b.state != connectivity.TransientFailure {
		
		
		return
	}
	b.cc.UpdateState(balancer.State{
		ConnectivityState: connectivity.TransientFailure,
		Picker:            &picker{err: fmt.Errorf("name resolver error: %v", err)},
	})
}


type Shuffler interface {
	ShuffleAddressListForTesting(n int, swap func(i, j int))
}



func ShuffleAddressListForTesting(n int, swap func(i, j int)) { rand.Shuffle(n, swap) }

func (b *pickfirstBalancer) UpdateClientConnState(state balancer.ClientConnState) error {
	if len(state.ResolverState.Addresses) == 0 && len(state.ResolverState.Endpoints) == 0 {
		
		
		if b.subConn != nil {
			
			
			b.subConn.Shutdown()
			b.subConn = nil
		}
		b.ResolverError(errors.New("produced zero addresses"))
		return balancer.ErrBadResolverState
	}
	
	
	cfg, ok := state.BalancerConfig.(pfConfig)
	if state.BalancerConfig != nil && !ok {
		return fmt.Errorf("pickfirst: received illegal BalancerConfig (type %T): %v", state.BalancerConfig, state.BalancerConfig)
	}

	if b.logger.V(2) {
		b.logger.Infof("Received new config %s, resolver state %s", pretty.ToJSON(cfg), pretty.ToJSON(state.ResolverState))
	}

	var addrs []resolver.Address
	if endpoints := state.ResolverState.Endpoints; len(endpoints) != 0 {
		
		
		
		if cfg.ShuffleAddressList {
			endpoints = append([]resolver.Endpoint{}, endpoints...)
			internal.RandShuffle(len(endpoints), func(i, j int) { endpoints[i], endpoints[j] = endpoints[j], endpoints[i] })
		}

		
		
		for _, endpoint := range endpoints {
			
			
			
			addrs = append(addrs, endpoint.Addresses...)
		}
	} else {
		
		
		
		
		
		
		addrs = state.ResolverState.Addresses
		if cfg.ShuffleAddressList {
			addrs = append([]resolver.Address{}, addrs...)
			rand.Shuffle(len(addrs), func(i, j int) { addrs[i], addrs[j] = addrs[j], addrs[i] })
		}
	}

	if b.subConn != nil {
		b.cc.UpdateAddresses(b.subConn, addrs)
		return nil
	}

	var subConn balancer.SubConn
	subConn, err := b.cc.NewSubConn(addrs, balancer.NewSubConnOptions{
		StateListener: func(state balancer.SubConnState) {
			b.updateSubConnState(subConn, state)
		},
	})
	if err != nil {
		if b.logger.V(2) {
			b.logger.Infof("Failed to create new SubConn: %v", err)
		}
		b.state = connectivity.TransientFailure
		b.cc.UpdateState(balancer.State{
			ConnectivityState: connectivity.TransientFailure,
			Picker:            &picker{err: fmt.Errorf("error creating connection: %v", err)},
		})
		return balancer.ErrBadResolverState
	}
	b.subConn = subConn
	b.state = connectivity.Idle
	b.cc.UpdateState(balancer.State{
		ConnectivityState: connectivity.Connecting,
		Picker:            &picker{err: balancer.ErrNoSubConnAvailable},
	})
	b.subConn.Connect()
	return nil
}



func (b *pickfirstBalancer) UpdateSubConnState(subConn balancer.SubConn, state balancer.SubConnState) {
	b.logger.Errorf("UpdateSubConnState(%v, %+v) called unexpectedly", subConn, state)
}

func (b *pickfirstBalancer) updateSubConnState(subConn balancer.SubConn, state balancer.SubConnState) {
	if b.logger.V(2) {
		b.logger.Infof("Received SubConn state update: %p, %+v", subConn, state)
	}
	if b.subConn != subConn {
		if b.logger.V(2) {
			b.logger.Infof("Ignored state change because subConn is not recognized")
		}
		return
	}
	if state.ConnectivityState == connectivity.Shutdown {
		b.subConn = nil
		return
	}

	switch state.ConnectivityState {
	case connectivity.Ready:
		b.cc.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            &picker{result: balancer.PickResult{SubConn: subConn}},
		})
	case connectivity.Connecting:
		if b.state == connectivity.TransientFailure {
			
			return
		}
		b.cc.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            &picker{err: balancer.ErrNoSubConnAvailable},
		})
	case connectivity.Idle:
		if b.state == connectivity.TransientFailure {
			
			
			b.subConn.Connect()
			return
		}
		b.cc.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            &idlePicker{subConn: subConn},
		})
	case connectivity.TransientFailure:
		b.cc.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            &picker{err: state.ConnectionError},
		})
	}
	b.state = state.ConnectivityState
}

func (b *pickfirstBalancer) Close() {
}

func (b *pickfirstBalancer) ExitIdle() {
	if b.subConn != nil && b.state == connectivity.Idle {
		b.subConn.Connect()
	}
}

type picker struct {
	result balancer.PickResult
	err    error
}

func (p *picker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	return p.result, p.err
}



type idlePicker struct {
	subConn balancer.SubConn
}

func (i *idlePicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	i.subConn.Connect()
	return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
}
