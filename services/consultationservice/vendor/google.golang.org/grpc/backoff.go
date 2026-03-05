




package grpc

import (
	"time"

	"google.golang.org/grpc/backoff"
)





var DefaultBackoffConfig = BackoffConfig{
	MaxDelay: 120 * time.Second,
}




type BackoffConfig struct {
	
	MaxDelay time.Duration
}










type ConnectParams struct {
	
	Backoff backoff.Config
	
	
	MinConnectTimeout time.Duration
}
