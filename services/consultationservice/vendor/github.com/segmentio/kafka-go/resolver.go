package kafka

import (
	"context"
	"net"
)



type Resolver interface {
	
	
	LookupHost(ctx context.Context, host string) (addrs []string, err error)
}










type BrokerResolver interface {
	
	LookupBrokerIPAddr(ctx context.Context, broker Broker) ([]net.IPAddr, error)
}




func NewBrokerResolver(r *net.Resolver) BrokerResolver {
	return brokerResolver{r}
}

type brokerResolver struct {
	*net.Resolver
}

func (r brokerResolver) LookupBrokerIPAddr(ctx context.Context, broker Broker) ([]net.IPAddr, error) {
	ipAddrs, err := r.LookupIPAddr(ctx, broker.Host)
	if err != nil {
		return nil, err
	}

	if len(ipAddrs) == 0 {
		return nil, &net.DNSError{
			Err:         "no addresses were returned by the resolver",
			Name:        broker.Host,
			IsTemporary: true,
			IsNotFound:  true,
		}
	}

	return ipAddrs, nil
}
