



package balancer

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"strings"

	"google.golang.org/grpc/channelz"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	estats "google.golang.org/grpc/experimental/stats"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

var (
	
	m = make(map[string]Builder)

	logger = grpclog.Component("balancer")
)










func Register(b Builder) {
	name := strings.ToLower(b.Name())
	if name != b.Name() {
		
		
		
		logger.Warningf("Balancer registered with name %q. grpc-go will be switching to case sensitive balancer registries soon", b.Name())
	}
	m[name] = b
}





func unregisterForTesting(name string) {
	delete(m, name)
}

func init() {
	internal.BalancerUnregister = unregisterForTesting
	internal.ConnectedAddress = connectedAddress
	internal.SetConnectedAddress = setConnectedAddress
}




func Get(name string) Builder {
	if strings.ToLower(name) != name {
		
		
		
		logger.Warningf("Balancer retrieved for name %q. grpc-go will be switching to case sensitive balancer registries soon", name)
	}
	if b, ok := m[strings.ToLower(name)]; ok {
		return b
	}
	return nil
}


type NewSubConnOptions struct {
	
	
	
	
	
	
	CredsBundle credentials.Bundle
	
	
	HealthCheckEnabled bool
	
	
	
	
	StateListener func(SubConnState)
}


type State struct {
	
	
	ConnectivityState connectivity.State
	
	Picker Picker
}














type ClientConn interface {
	
	
	
	
	
	
	NewSubConn([]resolver.Address, NewSubConnOptions) (SubConn, error)
	
	
	
	
	RemoveSubConn(SubConn)
	
	
	
	
	
	
	
	
	
	UpdateAddresses(SubConn, []resolver.Address)

	
	
	
	
	
	UpdateState(State)

	
	ResolveNow(resolver.ResolveNowOptions)

	
	
	
	Target() string

	
	
	
	
	MetricsRecorder() estats.MetricsRecorder

	
	
	
	internal.EnforceClientConnEmbedding
}


type BuildOptions struct {
	
	
	
	DialCreds credentials.TransportCredentials
	
	
	
	CredsBundle credentials.Bundle
	
	
	
	Dialer func(context.Context, string) (net.Conn, error)
	
	
	
	
	Authority string
	
	ChannelzParent channelz.Identifier
	
	
	
	CustomUserAgent string
	
	
	
	Target resolver.Target
}


type Builder interface {
	
	Build(cc ClientConn, opts BuildOptions) Balancer
	
	
	Name() string
}


type ConfigParser interface {
	
	
	
	ParseConfig(LoadBalancingConfigJSON json.RawMessage) (serviceconfig.LoadBalancingConfig, error)
}


type PickInfo struct {
	
	
	FullMethodName string
	
	
	Ctx context.Context
}


type DoneInfo struct {
	
	Err error
	
	Trailer metadata.MD
	
	BytesSent bool
	
	BytesReceived bool
	
	
	
	
	ServerLoad any
}

var (
	
	
	ErrNoSubConnAvailable = errors.New("no SubConn is available")
	
	
	
	
	
	
	ErrTransientFailure = errors.New("all SubConns are in TransientFailure")
)


type PickResult struct {
	
	
	
	
	SubConn SubConn

	
	
	
	
	Done func(DoneInfo)

	
	
	
	
	
	
	Metadata metadata.MD
}






func TransientFailureError(e error) error { return e }






type Picker interface {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Pick(info PickInfo) (PickResult, error)
}









type Balancer interface {
	
	
	
	
	
	UpdateClientConnState(ClientConnState) error
	
	ResolverError(error)
	
	
	
	
	
	UpdateSubConnState(SubConn, SubConnState)
	
	
	
	Close()
}








type ExitIdler interface {
	
	
	
	ExitIdle()
}



type ClientConnState struct {
	ResolverState resolver.State
	
	
	BalancerConfig serviceconfig.LoadBalancingConfig
}



var ErrBadResolverState = errors.New("bad resolver state")
