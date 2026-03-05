



package state

import (
	"google.golang.org/grpc/resolver"
)


type keyType string

const key = keyType("grpc.grpclb.state")


type State struct {
	
	
	BalancerAddresses []resolver.Address
}



func Set(state resolver.State, s *State) resolver.State {
	state.Attributes = state.Attributes.WithValue(key, s)
	return state
}



func Get(state resolver.State) *State {
	s, _ := state.Attributes.Value(key).(*State)
	return s
}
