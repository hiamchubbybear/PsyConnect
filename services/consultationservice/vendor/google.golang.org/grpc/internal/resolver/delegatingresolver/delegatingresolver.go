



package delegatingresolver

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/proxyattributes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

var (
	logger = grpclog.Component("delegating-resolver")
	
	HTTPSProxyFromEnvironment = http.ProxyFromEnvironment
)






type delegatingResolver struct {
	target         resolver.Target     
	cc             resolver.ClientConn 
	targetResolver resolver.Resolver   
	proxyResolver  resolver.Resolver   
	proxyURL       *url.URL            

	mu                  sync.Mutex         
	targetResolverState *resolver.State    
	proxyAddrs          []resolver.Address 
}


type nopResolver struct{}

func (nopResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (nopResolver) Close() {}









func proxyURLForTarget(address string) (*url.URL, error) {
	req := &http.Request{URL: &url.URL{
		Scheme: "https",
		Host:   address,
	}}
	return HTTPSProxyFromEnvironment(req)
}












func New(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions, targetResolverBuilder resolver.Builder, targetResolutionEnabled bool) (resolver.Resolver, error) {
	r := &delegatingResolver{
		target: target,
		cc:     cc,
	}

	var err error
	r.proxyURL, err = proxyURLForTarget(target.Endpoint())
	if err != nil {
		return nil, fmt.Errorf("delegating_resolver: failed to determine proxy URL for target %s: %v", target, err)
	}

	
	
	if r.proxyURL == nil {
		return targetResolverBuilder.Build(target, cc, opts)
	}

	if logger.V(2) {
		logger.Infof("Proxy URL detected : %s", r.proxyURL)
	}

	
	
	
	if target.URL.Scheme == "dns" && !targetResolutionEnabled {
		state := resolver.State{
			Addresses: []resolver.Address{{Addr: target.Endpoint()}},
			Endpoints: []resolver.Endpoint{{Addresses: []resolver.Address{{Addr: target.Endpoint()}}}},
		}
		r.targetResolverState = &state
	} else {
		wcc := &wrappingClientConn{
			stateListener: r.updateTargetResolverState,
			parent:        r,
		}
		if r.targetResolver, err = targetResolverBuilder.Build(target, wcc, opts); err != nil {
			return nil, fmt.Errorf("delegating_resolver: unable to build the resolver for target %s: %v", target, err)
		}
	}

	if r.proxyResolver, err = r.proxyURIResolver(opts); err != nil {
		return nil, fmt.Errorf("delegating_resolver: failed to build resolver for proxy URL %q: %v", r.proxyURL, err)
	}

	if r.targetResolver == nil {
		r.targetResolver = nopResolver{}
	}
	if r.proxyResolver == nil {
		r.proxyResolver = nopResolver{}
	}
	return r, nil
}




func (r *delegatingResolver) proxyURIResolver(opts resolver.BuildOptions) (resolver.Resolver, error) {
	proxyBuilder := resolver.Get("dns")
	if proxyBuilder == nil {
		panic("delegating_resolver: resolver for proxy not found for scheme dns")
	}
	url := *r.proxyURL
	url.Scheme = "dns"
	url.Path = "/" + r.proxyURL.Host
	url.Host = "" 

	proxyTarget := resolver.Target{URL: url}
	wcc := &wrappingClientConn{
		stateListener: r.updateProxyResolverState,
		parent:        r,
	}
	return proxyBuilder.Build(proxyTarget, wcc, opts)
}

func (r *delegatingResolver) ResolveNow(o resolver.ResolveNowOptions) {
	r.targetResolver.ResolveNow(o)
	r.proxyResolver.ResolveNow(o)
}

func (r *delegatingResolver) Close() {
	r.targetResolver.Close()
	r.targetResolver = nil

	r.proxyResolver.Close()
	r.proxyResolver = nil
}







func (r *delegatingResolver) updateClientConnStateLocked() error {
	if r.targetResolverState == nil || r.proxyAddrs == nil {
		return nil
	}

	curState := *r.targetResolverState
	
	
	
	
	
	
	
	
	
	var proxyAddr resolver.Address
	if len(r.proxyAddrs) == 1 {
		proxyAddr = r.proxyAddrs[0]
	} else {
		proxyAddr = resolver.Address{Addr: r.proxyURL.Host}
	}
	var addresses []resolver.Address
	for _, targetAddr := range (*r.targetResolverState).Addresses {
		addresses = append(addresses, proxyattributes.Set(proxyAddr, proxyattributes.Options{
			User:        r.proxyURL.User,
			ConnectAddr: targetAddr.Addr,
		}))
	}

	
	
	
	
	
	
	
	var endpoints []resolver.Endpoint
	for _, endpt := range (*r.targetResolverState).Endpoints {
		var addrs []resolver.Address
		for _, proxyAddr := range r.proxyAddrs {
			for _, targetAddr := range endpt.Addresses {
				addrs = append(addrs, proxyattributes.Set(proxyAddr, proxyattributes.Options{
					User:        r.proxyURL.User,
					ConnectAddr: targetAddr.Addr,
				}))
			}
		}
		endpoints = append(endpoints, resolver.Endpoint{Addresses: addrs})
	}
	
	
	
	
	curState.Addresses = addresses
	curState.Endpoints = endpoints
	return r.cc.UpdateState(curState)
}






func (r *delegatingResolver) updateProxyResolverState(state resolver.State) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if logger.V(2) {
		logger.Infof("Addresses received from proxy resolver: %s", state.Addresses)
	}
	if len(state.Endpoints) > 0 {
		
		
		r.proxyAddrs = make([]resolver.Address, 0, len(state.Endpoints))
		for _, endpoint := range state.Endpoints {
			r.proxyAddrs = append(r.proxyAddrs, endpoint.Addresses...)
		}
	} else if state.Addresses != nil {
		r.proxyAddrs = state.Addresses
	} else {
		r.proxyAddrs = []resolver.Address{} 
	}
	err := r.updateClientConnStateLocked()
	
	
	
	
	
	if err != nil {
		r.targetResolver.ResolveNow(resolver.ResolveNowOptions{})
	}
	return err
}






func (r *delegatingResolver) updateTargetResolverState(state resolver.State) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if logger.V(2) {
		logger.Infof("Addresses received from target resolver: %v", state.Addresses)
	}
	r.targetResolverState = &state
	err := r.updateClientConnStateLocked()
	if err != nil {
		r.proxyResolver.ResolveNow(resolver.ResolveNowOptions{})
	}
	return nil
}




type wrappingClientConn struct {
	
	stateListener func(state resolver.State) error
	parent        *delegatingResolver
}



func (wcc *wrappingClientConn) UpdateState(state resolver.State) error {
	return wcc.stateListener(state)
}


func (wcc *wrappingClientConn) ReportError(err error) {
	wcc.parent.cc.ReportError(err)
}



func (wcc *wrappingClientConn) NewAddress(addrs []resolver.Address) {
	wcc.UpdateState(resolver.State{Addresses: addrs})
}



func (wcc *wrappingClientConn) ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult {
	return wcc.parent.cc.ParseServiceConfig(serviceConfigJSON)
}
