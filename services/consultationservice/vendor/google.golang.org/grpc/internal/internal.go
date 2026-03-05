




package internal

import (
	"context"
	"time"

	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/serviceconfig"
)

var (
	
	HealthCheckFunc HealthChecker
	
	
	
	RegisterClientHealthCheckListener any 
	
	BalancerUnregister func(name string)
	
	
	KeepaliveMinPingTime = 10 * time.Second
	
	
	
	KeepaliveMinServerPingTime = time.Second
	
	ParseServiceConfig any 
	
	
	
	
	EqualServiceConfigForTesting func(a, b serviceconfig.Config) bool
	
	
	
	
	GetCertificateProviderBuilder any 
	
	
	
	GetXDSHandshakeInfoForTesting any 
	
	
	
	GetServerCredentials any 
	
	
	MetricsRecorderForServer any 
	
	
	
	
	
	CanonicalString any 
	
	
	IsRegisteredMethod any 
	
	ServerFromContext any 
	
	
	
	
	
	
	AddGlobalServerOptions any 
	
	
	
	
	
	ClearGlobalServerOptions func()
	
	
	
	
	
	
	AddGlobalDialOptions any 
	
	
	
	
	
	
	DisableGlobalDialOptions any 
	
	
	
	
	
	ClearGlobalDialOptions func()

	
	
	AddGlobalPerTargetDialOptions any 
	
	
	ClearGlobalPerTargetDialOptions func()

	
	
	JoinDialOptions any 
	
	
	JoinServerOptions any 

	
	
	
	
	
	WithBinaryLogger any 
	
	
	
	
	
	BinaryLogger any 

	
	
	SubscribeToConnectivityStateChanges any 

	
	
	
	
	
	
	
	
	
	NewXDSResolverWithConfigForTesting any 

	
	
	
	
	
	
	
	
	
	
	
	
	NewXDSResolverWithPoolForTesting any 

	
	
	
	
	
	
	
	
	
	
	
	
	NewXDSResolverWithClientForTesting any 

	
	
	
	
	
	RegisterRLSClusterSpecifierPluginForTesting func()

	
	
	
	
	
	
	UnregisterRLSClusterSpecifierPluginForTesting func()

	
	
	
	
	RegisterRBACHTTPFilterForTesting func()

	
	
	
	
	
	
	UnregisterRBACHTTPFilterForTesting func()

	
	ORCAAllowAnyMinReportingInterval any 

	
	
	GRPCResolverSchemeExtraMetadata = "xds"

	
	EnterIdleModeForTesting any 

	
	ExitIdleModeForTesting any 

	
	
	ChannelzTurnOffForTesting func()

	
	
	TriggerXDSResourceNotFoundForTesting any 

	
	
	FromOutgoingContextRaw any 

	
	
	UserSetDefaultScheme = false

	
	
	ConnectedAddress any 

	
	SetConnectedAddress any 

	
	
	
	SnapshotMetricRegistryForTesting func() func()

	
	
	SetDefaultBufferPoolForTesting any 

	
	
	SetBufferPoolingThresholdForTesting any 
)










type HealthChecker func(ctx context.Context, newStream func(string) (any, error), setConnectivityState func(connectivity.State, error), serviceName string) error

const (
	
	CredsBundleModeFallback = "fallback"
	
	
	CredsBundleModeBalancer = "balancer"
	
	
	CredsBundleModeBackendFromBalancer = "backend-from-balancer"
)





const RLSLoadBalancingPolicyName = "rls_experimental"



type EnforceSubConnEmbedding interface {
	enforceSubConnEmbedding()
}



type EnforceClientConnEmbedding interface {
	enforceClientConnEmbedding()
}
