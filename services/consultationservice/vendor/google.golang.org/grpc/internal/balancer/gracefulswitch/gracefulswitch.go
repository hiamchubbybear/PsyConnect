


package gracefulswitch

import (
	"errors"
	"fmt"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/resolver"
)

var errBalancerClosed = errors.New("gracefulSwitchBalancer is closed")
var _ balancer.Balancer = (*Balancer)(nil)


func NewBalancer(cc balancer.ClientConn, opts balancer.BuildOptions) *Balancer {
	return &Balancer{
		cc:    cc,
		bOpts: opts,
	}
}



type Balancer struct {
	bOpts balancer.BuildOptions
	cc    balancer.ClientConn

	
	
	
	
	
	
	
	
	
	mu              sync.Mutex
	balancerCurrent *balancerWrapper
	balancerPending *balancerWrapper
	closed          bool 

	
	
	
	
	
	currentMu sync.Mutex
}



func (gsb *Balancer) swap() {
	gsb.cc.UpdateState(gsb.balancerPending.lastState)
	cur := gsb.balancerCurrent
	gsb.balancerCurrent = gsb.balancerPending
	gsb.balancerPending = nil
	go func() {
		gsb.currentMu.Lock()
		defer gsb.currentMu.Unlock()
		cur.Close()
	}()
}



func (gsb *Balancer) balancerCurrentOrPending(bw *balancerWrapper) bool {
	return bw == gsb.balancerCurrent || bw == gsb.balancerPending
}









func (gsb *Balancer) SwitchTo(builder balancer.Builder) error {
	_, err := gsb.switchTo(builder)
	return err
}

func (gsb *Balancer) switchTo(builder balancer.Builder) (*balancerWrapper, error) {
	gsb.mu.Lock()
	if gsb.closed {
		gsb.mu.Unlock()
		return nil, errBalancerClosed
	}
	bw := &balancerWrapper{
		ClientConn: gsb.cc,
		builder:    builder,
		gsb:        gsb,
		lastState: balancer.State{
			ConnectivityState: connectivity.Connecting,
			Picker:            base.NewErrPicker(balancer.ErrNoSubConnAvailable),
		},
		subconns: make(map[balancer.SubConn]bool),
	}
	balToClose := gsb.balancerPending 
	if gsb.balancerCurrent == nil {
		gsb.balancerCurrent = bw
	} else {
		gsb.balancerPending = bw
	}
	gsb.mu.Unlock()
	balToClose.Close()
	
	
	newBalancer := builder.Build(bw, gsb.bOpts)
	if newBalancer == nil {
		
		
		gsb.mu.Lock()
		if gsb.balancerPending != nil {
			gsb.balancerPending = nil
		} else {
			gsb.balancerCurrent = nil
		}
		gsb.mu.Unlock()
		return nil, balancer.ErrBadResolverState
	}

	
	
	
	
	
	bw.Balancer = newBalancer
	return bw, nil
}


func (gsb *Balancer) latestBalancer() *balancerWrapper {
	gsb.mu.Lock()
	defer gsb.mu.Unlock()
	if gsb.balancerPending != nil {
		return gsb.balancerPending
	}
	return gsb.balancerCurrent
}







func (gsb *Balancer) UpdateClientConnState(state balancer.ClientConnState) error {
	
	balToUpdate := gsb.latestBalancer()
	gsbCfg, ok := state.BalancerConfig.(*lbConfig)
	if ok {
		
		if balToUpdate == nil || gsbCfg.childBuilder.Name() != balToUpdate.builder.Name() {
			var err error
			balToUpdate, err = gsb.switchTo(gsbCfg.childBuilder)
			if err != nil {
				return fmt.Errorf("could not switch to new child balancer: %w", err)
			}
		}
		
		state.BalancerConfig = gsbCfg.childConfig
	}

	if balToUpdate == nil {
		return errBalancerClosed
	}

	
	
	
	return balToUpdate.UpdateClientConnState(state)
}


func (gsb *Balancer) ResolverError(err error) {
	
	balToUpdate := gsb.latestBalancer()
	if balToUpdate == nil {
		gsb.cc.UpdateState(balancer.State{
			ConnectivityState: connectivity.TransientFailure,
			Picker:            base.NewErrPicker(err),
		})
		return
	}
	
	
	
	balToUpdate.ResolverError(err)
}





func (gsb *Balancer) ExitIdle() {
	balToUpdate := gsb.latestBalancer()
	if balToUpdate == nil {
		return
	}
	
	
	
	if ei, ok := balToUpdate.Balancer.(balancer.ExitIdler); ok {
		ei.ExitIdle()
		return
	}
	gsb.mu.Lock()
	defer gsb.mu.Unlock()
	for sc := range balToUpdate.subconns {
		sc.Connect()
	}
}


func (gsb *Balancer) updateSubConnState(sc balancer.SubConn, state balancer.SubConnState, cb func(balancer.SubConnState)) {
	gsb.currentMu.Lock()
	defer gsb.currentMu.Unlock()
	gsb.mu.Lock()
	
	
	
	var balToUpdate *balancerWrapper
	if gsb.balancerCurrent != nil && gsb.balancerCurrent.subconns[sc] {
		balToUpdate = gsb.balancerCurrent
	} else if gsb.balancerPending != nil && gsb.balancerPending.subconns[sc] {
		balToUpdate = gsb.balancerPending
	}
	if balToUpdate == nil {
		
		
		gsb.mu.Unlock()
		return
	}
	if state.ConnectivityState == connectivity.Shutdown {
		delete(balToUpdate.subconns, sc)
	}
	gsb.mu.Unlock()
	if cb != nil {
		cb(state)
	} else {
		balToUpdate.UpdateSubConnState(sc, state)
	}
}


func (gsb *Balancer) UpdateSubConnState(sc balancer.SubConn, state balancer.SubConnState) {
	gsb.updateSubConnState(sc, state, nil)
}


func (gsb *Balancer) Close() {
	gsb.mu.Lock()
	gsb.closed = true
	currentBalancerToClose := gsb.balancerCurrent
	gsb.balancerCurrent = nil
	pendingBalancerToClose := gsb.balancerPending
	gsb.balancerPending = nil
	gsb.mu.Unlock()

	currentBalancerToClose.Close()
	pendingBalancerToClose.Close()
}










type balancerWrapper struct {
	balancer.ClientConn
	balancer.Balancer
	gsb     *Balancer
	builder balancer.Builder

	lastState balancer.State
	subconns  map[balancer.SubConn]bool 
}





func (bw *balancerWrapper) Close() {
	
	if bw == nil {
		return
	}
	
	
	
	
	bw.Balancer.Close()
	bw.gsb.mu.Lock()
	for sc := range bw.subconns {
		sc.Shutdown()
	}
	bw.gsb.mu.Unlock()
}

func (bw *balancerWrapper) UpdateState(state balancer.State) {
	
	
	
	bw.gsb.mu.Lock()
	defer bw.gsb.mu.Unlock()
	bw.lastState = state

	if !bw.gsb.balancerCurrentOrPending(bw) {
		return
	}

	if bw == bw.gsb.balancerCurrent {
		
		
		
		
		
		if state.ConnectivityState != connectivity.Ready && bw.gsb.balancerPending != nil {
			bw.gsb.swap()
			return
		}
		
		
		
		
		
		
		bw.gsb.cc.UpdateState(state)
		return
	}
	
	
	
	
	
	if state.ConnectivityState != connectivity.Connecting || bw.gsb.balancerCurrent.lastState.ConnectivityState != connectivity.Ready {
		bw.gsb.swap()
	}
}

func (bw *balancerWrapper) NewSubConn(addrs []resolver.Address, opts balancer.NewSubConnOptions) (balancer.SubConn, error) {
	bw.gsb.mu.Lock()
	if !bw.gsb.balancerCurrentOrPending(bw) {
		bw.gsb.mu.Unlock()
		return nil, fmt.Errorf("%T at address %p that called NewSubConn is deleted", bw, bw)
	}
	bw.gsb.mu.Unlock()

	var sc balancer.SubConn
	oldListener := opts.StateListener
	opts.StateListener = func(state balancer.SubConnState) { bw.gsb.updateSubConnState(sc, state, oldListener) }
	sc, err := bw.gsb.cc.NewSubConn(addrs, opts)
	if err != nil {
		return nil, err
	}
	bw.gsb.mu.Lock()
	if !bw.gsb.balancerCurrentOrPending(bw) { 
		sc.Shutdown()
		bw.gsb.mu.Unlock()
		return nil, fmt.Errorf("%T at address %p that called NewSubConn is deleted", bw, bw)
	}
	bw.subconns[sc] = true
	bw.gsb.mu.Unlock()
	return sc, nil
}

func (bw *balancerWrapper) ResolveNow(opts resolver.ResolveNowOptions) {
	
	
	if bw != bw.gsb.latestBalancer() {
		return
	}
	bw.gsb.cc.ResolveNow(opts)
}

func (bw *balancerWrapper) RemoveSubConn(sc balancer.SubConn) {
	
	
	sc.Shutdown()
}

func (bw *balancerWrapper) UpdateAddresses(sc balancer.SubConn, addrs []resolver.Address) {
	bw.gsb.mu.Lock()
	if !bw.gsb.balancerCurrentOrPending(bw) {
		bw.gsb.mu.Unlock()
		return
	}
	bw.gsb.mu.Unlock()
	bw.gsb.cc.UpdateAddresses(sc, addrs)
}
