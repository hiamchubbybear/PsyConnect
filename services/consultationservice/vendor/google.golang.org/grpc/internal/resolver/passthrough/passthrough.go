



package passthrough

import (
	"errors"

	"google.golang.org/grpc/resolver"
)

const scheme = "passthrough"

type passthroughBuilder struct{}

func (*passthroughBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	if target.Endpoint() == "" && opts.Dialer == nil {
		return nil, errors.New("passthrough: received empty target in Build()")
	}
	r := &passthroughResolver{
		target: target,
		cc:     cc,
	}
	r.start()
	return r, nil
}

func (*passthroughBuilder) Scheme() string {
	return scheme
}

type passthroughResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
}

func (r *passthroughResolver) start() {
	r.cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: r.target.Endpoint()}}})
}

func (*passthroughResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (*passthroughResolver) Close() {}

func init() {
	resolver.Register(&passthroughBuilder{})
}
