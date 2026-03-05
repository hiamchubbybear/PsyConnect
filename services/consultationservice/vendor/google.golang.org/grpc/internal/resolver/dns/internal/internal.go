


package internal

import (
	"context"
	"errors"
	"net"
	"time"
)




type NetResolver interface {
	LookupHost(ctx context.Context, host string) (addrs []string, err error)
	LookupSRV(ctx context.Context, service, proto, name string) (cname string, addrs []*net.SRV, err error)
	LookupTXT(ctx context.Context, name string) (txts []string, err error)
}

var (
	
	
	ErrMissingAddr = errors.New("dns resolver: missing address")

	
	
	
	
	
	ErrEndsWithColon = errors.New("dns resolver: missing port after port-separator colon")
)


var (
	
	
	
	
	TimeAfterFunc func(time.Duration) <-chan time.Time

	
	
	
	TimeNowFunc func() time.Time

	
	
	
	
	TimeUntilFunc func(time.Time) time.Duration

	
	NewNetResolver func(string) (NetResolver, error)

	
	
	
	AddressDialer func(address string) func(context.Context, string, string) (net.Conn, error)
)
