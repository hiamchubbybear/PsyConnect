



package dns

import (
	"time"

	"google.golang.org/grpc/internal/resolver/dns"
	"google.golang.org/grpc/resolver"
)













func SetResolvingTimeout(timeout time.Duration) {
	dns.ResolvingTimeout = timeout
}




func NewBuilder() resolver.Builder {
	return dns.NewBuilder()
}






func SetMinResolutionInterval(d time.Duration) {
	dns.MinResolutionInterval = d
}
