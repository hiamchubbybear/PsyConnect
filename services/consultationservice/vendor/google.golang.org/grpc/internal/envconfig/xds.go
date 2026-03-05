

package envconfig

import (
	"os"
)

const (
	
	
	
	
	
	XDSBootstrapFileNameEnv = "GRPC_XDS_BOOTSTRAP"
	
	
	
	
	
	XDSBootstrapFileContentEnv = "GRPC_XDS_BOOTSTRAP_CONFIG"
)

var (
	
	
	
	
	
	XDSBootstrapFileName = os.Getenv(XDSBootstrapFileNameEnv)
	
	
	
	
	
	XDSBootstrapFileContent = os.Getenv(XDSBootstrapFileContentEnv)

	
	C2PResolverTestOnlyTrafficDirectorURI = os.Getenv("GRPC_TEST_ONLY_GOOGLE_C2P_RESOLVER_TRAFFIC_DIRECTOR_URI")

	
	
	XDSDualstackEndpointsEnabled = boolFromEnv("GRPC_EXPERIMENTAL_XDS_DUALSTACK_ENDPOINTS", true)

	
	
	
	
	XDSSystemRootCertsEnabled = boolFromEnv("GRPC_EXPERIMENTAL_XDS_SYSTEM_ROOT_CERTS", false)
)
