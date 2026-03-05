

package envconfig

import "os"

const (
	envObservabilityConfig     = "GRPC_GCP_OBSERVABILITY_CONFIG"
	envObservabilityConfigFile = "GRPC_GCP_OBSERVABILITY_CONFIG_FILE"
)

var (
	
	
	
	
	
	ObservabilityConfig = os.Getenv(envObservabilityConfig)
	
	
	
	
	
	
	ObservabilityConfigFile = os.Getenv(envObservabilityConfigFile)
)
