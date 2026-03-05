

package balancer

import "google.golang.org/grpc/connectivity"





type ConnectivityStateEvaluator struct {
	numReady            uint64 
	numConnecting       uint64 
	numTransientFailure uint64 
	numIdle             uint64 
}










func (cse *ConnectivityStateEvaluator) RecordTransition(oldState, newState connectivity.State) connectivity.State {
	
	for idx, state := range []connectivity.State{oldState, newState} {
		updateVal := 2*uint64(idx) - 1 
		switch state {
		case connectivity.Ready:
			cse.numReady += updateVal
		case connectivity.Connecting:
			cse.numConnecting += updateVal
		case connectivity.TransientFailure:
			cse.numTransientFailure += updateVal
		case connectivity.Idle:
			cse.numIdle += updateVal
		}
	}
	return cse.CurrentState()
}


func (cse *ConnectivityStateEvaluator) CurrentState() connectivity.State {
	
	if cse.numReady > 0 {
		return connectivity.Ready
	}
	if cse.numConnecting > 0 {
		return connectivity.Connecting
	}
	if cse.numIdle > 0 {
		return connectivity.Idle
	}
	return connectivity.TransientFailure
}
