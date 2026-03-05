

package balancer

import (
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/internal"
	"google.golang.org/grpc/resolver"
)

























type SubConn interface {
	
	
	
	
	
	
	
	
	
	
	UpdateAddresses([]resolver.Address)
	
	Connect()
	
	
	
	
	
	
	
	GetOrBuildProducer(ProducerBuilder) (p Producer, close func())
	
	
	
	
	
	
	Shutdown()
	
	
	
	
	
	
	
	RegisterHealthListener(func(SubConnState))
	
	
	
	internal.EnforceSubConnEmbedding
}



type ProducerBuilder interface {
	
	
	
	
	
	
	
	Build(grpcClientConnInterface any) (p Producer, close func())
}


type SubConnState struct {
	
	ConnectivityState connectivity.State
	
	
	ConnectionError error
	
	
	connectedAddress resolver.Address
}



func connectedAddress(scs SubConnState) resolver.Address {
	return scs.connectedAddress
}


func setConnectedAddress(scs *SubConnState, addr resolver.Address) {
	scs.connectedAddress = addr
}





type Producer any
