


package envconfig

import (
	"os"
	"strconv"
	"strings"
)

var (
	
	TXTErrIgnore = boolFromEnv("GRPC_GO_IGNORE_TXT_ERRORS", true)
	
	
	
	
	RingHashCap = uint64FromEnv("GRPC_RING_HASH_CAP", 4096, 1, 8*1024*1024)
	
	
	
	LeastRequestLB = boolFromEnv("GRPC_EXPERIMENTAL_ENABLE_LEAST_REQUEST", false)
	
	
	ALTSMaxConcurrentHandshakes = uint64FromEnv("GRPC_ALTS_MAX_CONCURRENT_HANDSHAKES", 100, 1, 100)
	
	
	
	
	
	EnforceALPNEnabled = boolFromEnv("GRPC_ENFORCE_ALPN_ENABLED", true)
	
	
	
	XDSFallbackSupport = boolFromEnv("GRPC_EXPERIMENTAL_XDS_FALLBACK", true)
	
	
	
	
	NewPickFirstEnabled = boolFromEnv("GRPC_EXPERIMENTAL_ENABLE_NEW_PICK_FIRST", false)
)

func boolFromEnv(envVar string, def bool) bool {
	if def {
		
		return !strings.EqualFold(os.Getenv(envVar), "false")
	}
	
	return strings.EqualFold(os.Getenv(envVar), "true")
}

func uint64FromEnv(envVar string, def, min, max uint64) uint64 {
	v, err := strconv.ParseUint(os.Getenv(envVar), 10, 64)
	if err != nil {
		return def
	}
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
