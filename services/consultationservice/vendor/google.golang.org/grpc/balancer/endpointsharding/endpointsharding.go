








package endpointsharding

import (
	"errors"
	rand "math/rand/v2"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/resolver"
)



type ChildState struct {
	Endpoint resolver.Endpoint
	State    balancer.State

	
	
	Balancer balancer.ExitIdler
}



type Options struct {
	
	
	
	
	
	DisableAutoReconnect bool
}



type ChildBuilderFunc func(cc balancer.ClientConn, opts balancer.BuildOptions) balancer.Balancer




func NewBalancer(cc balancer.ClientConn, opts balancer.BuildOptions, childBuilder ChildBuilderFunc, esOpts Options) balancer.Balancer {
	es := &endpointSharding{
		cc:           cc,
		bOpts:        opts,
		esOpts:       esOpts,
		childBuilder: childBuilder,
	}
	es.children.Store(resolver.NewEndpointMap())
	return es
}




type endpointSharding struct {
	cc           balancer.ClientConn
	bOpts        balancer.BuildOptions
	esOpts       Options
	childBuilder ChildBuilderFunc

	
	
	
	childMu  sync.Mutex
	children atomic.Pointer[resolver.EndpointMap] 

	
	
	
	inhibitChildUpdates atomic.Bool

	
	
	
	
	mu sync.Mutex
}







func (es *endpointSharding) UpdateClientConnState(state balancer.ClientConnState) error {
	es.childMu.Lock()
	defer es.childMu.Unlock()

	es.inhibitChildUpdates.Store(true)
	defer func() {
		es.inhibitChildUpdates.Store(false)
		es.updateState()
	}()
	var ret error

	children := es.children.Load()
	newChildren := resolver.NewEndpointMap()

	
	for _, endpoint := range state.ResolverState.Endpoints {
		if _, ok := newChildren.Get(endpoint); ok {
			
			
			continue
		}
		var childBalancer *balancerWrapper
		if val, ok := children.Get(endpoint); ok {
			childBalancer = val.(*balancerWrapper)
			
			es.mu.Lock()
			childBalancer.childState.Endpoint = endpoint
			es.mu.Unlock()
		} else {
			childBalancer = &balancerWrapper{
				childState: ChildState{Endpoint: endpoint},
				ClientConn: es.cc,
				es:         es,
			}
			childBalancer.childState.Balancer = childBalancer
			childBalancer.child = es.childBuilder(childBalancer, es.bOpts)
		}
		newChildren.Set(endpoint, childBalancer)
		if err := childBalancer.updateClientConnStateLocked(balancer.ClientConnState{
			BalancerConfig: state.BalancerConfig,
			ResolverState: resolver.State{
				Endpoints:  []resolver.Endpoint{endpoint},
				Attributes: state.ResolverState.Attributes,
			},
		}); err != nil && ret == nil {
			
			
			
			
			ret = err
		}
	}
	
	for _, e := range children.Keys() {
		child, _ := children.Get(e)
		if _, ok := newChildren.Get(e); !ok {
			child.(*balancerWrapper).closeLocked()
		}
	}
	es.children.Store(newChildren)
	if newChildren.Len() == 0 {
		return balancer.ErrBadResolverState
	}
	return ret
}




func (es *endpointSharding) ResolverError(err error) {
	es.childMu.Lock()
	defer es.childMu.Unlock()
	es.inhibitChildUpdates.Store(true)
	defer func() {
		es.inhibitChildUpdates.Store(false)
		es.updateState()
	}()
	children := es.children.Load()
	for _, child := range children.Values() {
		child.(*balancerWrapper).resolverErrorLocked(err)
	}
}

func (es *endpointSharding) UpdateSubConnState(balancer.SubConn, balancer.SubConnState) {
	
}

func (es *endpointSharding) Close() {
	es.childMu.Lock()
	defer es.childMu.Unlock()
	children := es.children.Load()
	for _, child := range children.Values() {
		child.(*balancerWrapper).closeLocked()
	}
}




func (es *endpointSharding) updateState() {
	if es.inhibitChildUpdates.Load() {
		return
	}
	var readyPickers, connectingPickers, idlePickers, transientFailurePickers []balancer.Picker

	es.mu.Lock()
	defer es.mu.Unlock()

	children := es.children.Load()
	childStates := make([]ChildState, 0, children.Len())

	for _, child := range children.Values() {
		bw := child.(*balancerWrapper)
		childState := bw.childState
		childStates = append(childStates, childState)
		childPicker := childState.State.Picker
		switch childState.State.ConnectivityState {
		case connectivity.Ready:
			readyPickers = append(readyPickers, childPicker)
		case connectivity.Connecting:
			connectingPickers = append(connectingPickers, childPicker)
		case connectivity.Idle:
			idlePickers = append(idlePickers, childPicker)
		case connectivity.TransientFailure:
			transientFailurePickers = append(transientFailurePickers, childPicker)
			
		}
	}

	
	
	
	var aggState connectivity.State
	var pickers []balancer.Picker
	if len(readyPickers) >= 1 {
		aggState = connectivity.Ready
		pickers = readyPickers
	} else if len(connectingPickers) >= 1 {
		aggState = connectivity.Connecting
		pickers = connectingPickers
	} else if len(idlePickers) >= 1 {
		aggState = connectivity.Idle
		pickers = idlePickers
	} else if len(transientFailurePickers) >= 1 {
		aggState = connectivity.TransientFailure
		pickers = transientFailurePickers
	} else {
		aggState = connectivity.TransientFailure
		pickers = []balancer.Picker{base.NewErrPicker(errors.New("no children to pick from"))}
	} 
	p := &pickerWithChildStates{
		pickers:     pickers,
		childStates: childStates,
		next:        uint32(rand.IntN(len(pickers))),
	}
	es.cc.UpdateState(balancer.State{
		ConnectivityState: aggState,
		Picker:            p,
	})
}




type pickerWithChildStates struct {
	pickers     []balancer.Picker
	childStates []ChildState
	next        uint32
}

func (p *pickerWithChildStates) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	nextIndex := atomic.AddUint32(&p.next, 1)
	picker := p.pickers[nextIndex%uint32(len(p.pickers))]
	return picker.Pick(info)
}



func ChildStatesFromPicker(picker balancer.Picker) []ChildState {
	p, ok := picker.(*pickerWithChildStates)
	if !ok {
		return nil
	}
	return p.childStates
}



type balancerWrapper struct {
	
	

	
	
	child               balancer.Balancer
	balancer.ClientConn 

	es *endpointSharding

	

	childState ChildState
	isClosed   bool
}

func (bw *balancerWrapper) UpdateState(state balancer.State) {
	bw.es.mu.Lock()
	bw.childState.State = state
	bw.es.mu.Unlock()
	if state.ConnectivityState == connectivity.Idle && !bw.es.esOpts.DisableAutoReconnect {
		bw.ExitIdle()
	}
	bw.es.updateState()
}



func (bw *balancerWrapper) ExitIdle() {
	if ei, ok := bw.child.(balancer.ExitIdler); ok {
		go func() {
			bw.es.childMu.Lock()
			if !bw.isClosed {
				ei.ExitIdle()
			}
			bw.es.childMu.Unlock()
		}()
	}
}




func (bw *balancerWrapper) updateClientConnStateLocked(ccs balancer.ClientConnState) error {
	return bw.child.UpdateClientConnState(ccs)
}



func (bw *balancerWrapper) closeLocked() {
	bw.child.Close()
	bw.isClosed = true
}

func (bw *balancerWrapper) resolverErrorLocked(err error) {
	bw.child.ResolverError(err)
}
