




package roundrobin

import (
	"fmt"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/endpointsharding"
	"google.golang.org/grpc/balancer/pickfirst/pickfirstleaf"
	"google.golang.org/grpc/grpclog"
	internalgrpclog "google.golang.org/grpc/internal/grpclog"
)


const Name = "round_robin"

var logger = grpclog.Component("roundrobin")

func init() {
	balancer.Register(builder{})
}

type builder struct{}

func (bb builder) Name() string {
	return Name
}

func (bb builder) Build(cc balancer.ClientConn, opts balancer.BuildOptions) balancer.Balancer {
	childBuilder := balancer.Get(pickfirstleaf.Name).Build
	bal := &rrBalancer{
		cc:       cc,
		Balancer: endpointsharding.NewBalancer(cc, opts, childBuilder, endpointsharding.Options{}),
	}
	bal.logger = internalgrpclog.NewPrefixLogger(logger, fmt.Sprintf("[%p] ", bal))
	bal.logger.Infof("Created")
	return bal
}

type rrBalancer struct {
	balancer.Balancer
	cc     balancer.ClientConn
	logger *internalgrpclog.PrefixLogger
}

func (b *rrBalancer) UpdateClientConnState(ccs balancer.ClientConnState) error {
	return b.Balancer.UpdateClientConnState(balancer.ClientConnState{
		
		
		ResolverState: pickfirstleaf.EnableHealthListener(ccs.ResolverState),
	})
}

func (b *rrBalancer) ExitIdle() {
	
	if ei, ok := b.Balancer.(balancer.ExitIdler); ok {
		ei.ExitIdle()
	}
}
