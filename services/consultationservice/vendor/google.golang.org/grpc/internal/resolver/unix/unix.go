


package unix

import (
	"fmt"

	"google.golang.org/grpc/internal/transport/networktype"
	"google.golang.org/grpc/resolver"
)

const unixScheme = "unix"
const unixAbstractScheme = "unix-abstract"

type builder struct {
	scheme string
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	if target.URL.Host != "" {
		return nil, fmt.Errorf("invalid (non-empty) authority: %v", target.URL.Host)
	}

	
	
	
	
	
	endpoint := target.URL.Path
	if endpoint == "" {
		endpoint = target.URL.Opaque
	}
	addr := resolver.Address{Addr: endpoint}
	if b.scheme == unixAbstractScheme {
		
		
		addr.Addr = "@" + addr.Addr
	}
	cc.UpdateState(resolver.State{Addresses: []resolver.Address{networktype.Set(addr, "unix")}})
	return &nopResolver{}, nil
}

func (b *builder) Scheme() string {
	return b.scheme
}

func (b *builder) OverrideAuthority(resolver.Target) string {
	return "localhost"
}

type nopResolver struct {
}

func (*nopResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (*nopResolver) Close() {}

func init() {
	resolver.Register(&builder{scheme: unixScheme})
	resolver.Register(&builder{scheme: unixAbstractScheme})
}
