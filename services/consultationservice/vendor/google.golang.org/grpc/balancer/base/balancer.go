

package base

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

var logger = grpclog.Component("balancer")

type baseBuilder struct {
	name          string
	pickerBuilder PickerBuilder
	config        Config
}

func (bb *baseBuilder) Build(cc balancer.ClientConn, _ balancer.BuildOptions) balancer.Balancer {
	bal := &baseBalancer{
		cc:            cc,
		pickerBuilder: bb.pickerBuilder,

		subConns: resolver.NewAddressMap(),
		scStates: make(map[balancer.SubConn]connectivity.State),
		csEvltr:  &balancer.ConnectivityStateEvaluator{},
		config:   bb.config,
		state:    connectivity.Connecting,
	}
	
	
	
	bal.picker = NewErrPicker(balancer.ErrNoSubConnAvailable)
	return bal
}

func (bb *baseBuilder) Name() string {
	return bb.name
}

type baseBalancer struct {
	cc            balancer.ClientConn
	pickerBuilder PickerBuilder

	csEvltr *balancer.ConnectivityStateEvaluator
	state   connectivity.State

	subConns *resolver.AddressMap
	scStates map[balancer.SubConn]connectivity.State
	picker   balancer.Picker
	config   Config

	resolverErr error 
	connErr     error 
}

func (b *baseBalancer) ResolverError(err error) {
	b.resolverErr = err
	if b.subConns.Len() == 0 {
		b.state = connectivity.TransientFailure
	}

	if b.state != connectivity.TransientFailure {
		
		
		return
	}
	b.regeneratePicker()
	b.cc.UpdateState(balancer.State{
		ConnectivityState: b.state,
		Picker:            b.picker,
	})
}

func (b *baseBalancer) UpdateClientConnState(s balancer.ClientConnState) error {
	
	if logger.V(2) {
		logger.Info("base.baseBalancer: got new ClientConn state: ", s)
	}
	
	b.resolverErr = nil
	
	addrsSet := resolver.NewAddressMap()
	for _, a := range s.ResolverState.Addresses {
		addrsSet.Set(a, nil)
		if _, ok := b.subConns.Get(a); !ok {
			
			var sc balancer.SubConn
			opts := balancer.NewSubConnOptions{
				HealthCheckEnabled: b.config.HealthCheck,
				StateListener:      func(scs balancer.SubConnState) { b.updateSubConnState(sc, scs) },
			}
			sc, err := b.cc.NewSubConn([]resolver.Address{a}, opts)
			if err != nil {
				logger.Warningf("base.baseBalancer: failed to create new SubConn: %v", err)
				continue
			}
			b.subConns.Set(a, sc)
			b.scStates[sc] = connectivity.Idle
			b.csEvltr.RecordTransition(connectivity.Shutdown, connectivity.Idle)
			sc.Connect()
		}
	}
	for _, a := range b.subConns.Keys() {
		sci, _ := b.subConns.Get(a)
		sc := sci.(balancer.SubConn)
		
		if _, ok := addrsSet.Get(a); !ok {
			sc.Shutdown()
			b.subConns.Delete(a)
			
			
		}
	}
	
	
	
	
	if len(s.ResolverState.Addresses) == 0 {
		b.ResolverError(errors.New("produced zero addresses"))
		return balancer.ErrBadResolverState
	}

	b.regeneratePicker()
	b.cc.UpdateState(balancer.State{ConnectivityState: b.state, Picker: b.picker})
	return nil
}



func (b *baseBalancer) mergeErrors() error {
	
	
	if b.connErr == nil {
		return fmt.Errorf("last resolver error: %v", b.resolverErr)
	}
	if b.resolverErr == nil {
		return fmt.Errorf("last connection error: %v", b.connErr)
	}
	return fmt.Errorf("last connection error: %v; last resolver error: %v", b.connErr, b.resolverErr)
}





func (b *baseBalancer) regeneratePicker() {
	if b.state == connectivity.TransientFailure {
		b.picker = NewErrPicker(b.mergeErrors())
		return
	}
	readySCs := make(map[balancer.SubConn]SubConnInfo)

	
	for _, addr := range b.subConns.Keys() {
		sci, _ := b.subConns.Get(addr)
		sc := sci.(balancer.SubConn)
		if st, ok := b.scStates[sc]; ok && st == connectivity.Ready {
			readySCs[sc] = SubConnInfo{Address: addr}
		}
	}
	b.picker = b.pickerBuilder.Build(PickerBuildInfo{ReadySCs: readySCs})
}


func (b *baseBalancer) UpdateSubConnState(sc balancer.SubConn, state balancer.SubConnState) {
	logger.Errorf("base.baseBalancer: UpdateSubConnState(%v, %+v) called unexpectedly", sc, state)
}

func (b *baseBalancer) updateSubConnState(sc balancer.SubConn, state balancer.SubConnState) {
	s := state.ConnectivityState
	if logger.V(2) {
		logger.Infof("base.baseBalancer: handle SubConn state change: %p, %v", sc, s)
	}
	oldS, ok := b.scStates[sc]
	if !ok {
		if logger.V(2) {
			logger.Infof("base.baseBalancer: got state changes for an unknown SubConn: %p, %v", sc, s)
		}
		return
	}
	if oldS == connectivity.TransientFailure &&
		(s == connectivity.Connecting || s == connectivity.Idle) {
		
		
		
		if s == connectivity.Idle {
			sc.Connect()
		}
		return
	}
	b.scStates[sc] = s
	switch s {
	case connectivity.Idle:
		sc.Connect()
	case connectivity.Shutdown:
		
		
		delete(b.scStates, sc)
	case connectivity.TransientFailure:
		
		b.connErr = state.ConnectionError
	}

	b.state = b.csEvltr.RecordTransition(oldS, s)

	
	
	
	
	if (s == connectivity.Ready) != (oldS == connectivity.Ready) ||
		b.state == connectivity.TransientFailure {
		b.regeneratePicker()
	}
	b.cc.UpdateState(balancer.State{ConnectivityState: b.state, Picker: b.picker})
}



func (b *baseBalancer) Close() {
}



func (b *baseBalancer) ExitIdle() {
}


func NewErrPicker(err error) balancer.Picker {
	return &errPicker{err: err}
}




var NewErrPickerV2 = NewErrPicker

type errPicker struct {
	err error 
}

func (p *errPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	return balancer.PickResult{}, p.err
}
