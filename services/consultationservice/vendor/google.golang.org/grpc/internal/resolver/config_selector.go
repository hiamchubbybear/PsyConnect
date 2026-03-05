


package resolver

import (
	"context"
	"sync"

	"google.golang.org/grpc/internal/serviceconfig"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
)


type ConfigSelector interface {
	
	
	
	SelectConfig(RPCInfo) (*RPCConfig, error)
}


type RPCInfo struct {
	
	
	
	Context context.Context
	Method  string 
}


type RPCConfig struct {
	
	
	Context      context.Context
	MethodConfig serviceconfig.MethodConfig 
	OnCommitted  func()                     
	Interceptor  ClientInterceptor
}



type ClientStream interface {
	
	
	Header() (metadata.MD, error)
	
	
	
	Trailer() metadata.MD
	
	
	
	CloseSend() error
	
	
	
	
	Context() context.Context
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	SendMsg(m any) error
	
	
	
	
	
	
	
	
	RecvMsg(m any) error
}


type ClientInterceptor interface {
	
	
	
	
	
	
	
	
	NewStream(ctx context.Context, ri RPCInfo, done func(), newStream func(ctx context.Context, done func()) (ClientStream, error)) (ClientStream, error)
}


type ServerInterceptor interface {
	
	
	
	AllowRPC(ctx context.Context) error 
}

type csKeyType string

const csKey = csKeyType("grpc.internal.resolver.configSelector")



func SetConfigSelector(state resolver.State, cs ConfigSelector) resolver.State {
	state.Attributes = state.Attributes.WithValue(csKey, cs)
	return state
}



func GetConfigSelector(state resolver.State) ConfigSelector {
	cs, _ := state.Attributes.Value(csKey).(ConfigSelector)
	return cs
}




type SafeConfigSelector struct {
	mu sync.RWMutex
	cs ConfigSelector
}



func (scs *SafeConfigSelector) UpdateConfigSelector(cs ConfigSelector) {
	scs.mu.Lock()
	defer scs.mu.Unlock()
	scs.cs = cs
}


func (scs *SafeConfigSelector) SelectConfig(r RPCInfo) (*RPCConfig, error) {
	scs.mu.RLock()
	defer scs.mu.RUnlock()
	return scs.cs.SelectConfig(r)
}
