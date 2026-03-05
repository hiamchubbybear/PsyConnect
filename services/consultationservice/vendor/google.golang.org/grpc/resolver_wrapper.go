

package grpc

import (
	"context"
	"strings"
	"sync"

	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/grpcsync"
	"google.golang.org/grpc/internal/pretty"
	"google.golang.org/grpc/internal/resolver/delegatingresolver"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)



type ccResolverWrapper struct {
	
	
	cc                  *ClientConn
	ignoreServiceConfig bool
	serializer          *grpcsync.CallbackSerializer
	serializerCancel    context.CancelFunc

	resolver resolver.Resolver 

	
	
	mu       sync.Mutex
	curState resolver.State
	closed   bool
}



func newCCResolverWrapper(cc *ClientConn) *ccResolverWrapper {
	ctx, cancel := context.WithCancel(cc.ctx)
	return &ccResolverWrapper{
		cc:                  cc,
		ignoreServiceConfig: cc.dopts.disableServiceConfig,
		serializer:          grpcsync.NewCallbackSerializer(ctx),
		serializerCancel:    cancel,
	}
}




func (ccr *ccResolverWrapper) start() error {
	errCh := make(chan error)
	ccr.serializer.TrySchedule(func(ctx context.Context) {
		if ctx.Err() != nil {
			return
		}
		opts := resolver.BuildOptions{
			DisableServiceConfig: ccr.cc.dopts.disableServiceConfig,
			DialCreds:            ccr.cc.dopts.copts.TransportCredentials,
			CredsBundle:          ccr.cc.dopts.copts.CredsBundle,
			Dialer:               ccr.cc.dopts.copts.Dialer,
			Authority:            ccr.cc.authority,
			MetricsRecorder:      ccr.cc.metricsRecorderList,
		}
		var err error
		
		
		
		
		
		if ccr.cc.dopts.copts.Dialer != nil || !ccr.cc.dopts.useProxy {
			ccr.resolver, err = ccr.cc.resolverBuilder.Build(ccr.cc.parsedTarget, ccr, opts)
		} else {
			ccr.resolver, err = delegatingresolver.New(ccr.cc.parsedTarget, ccr, opts, ccr.cc.resolverBuilder, ccr.cc.dopts.enableLocalDNSResolution)
		}
		errCh <- err
	})
	return <-errCh
}

func (ccr *ccResolverWrapper) resolveNow(o resolver.ResolveNowOptions) {
	ccr.serializer.TrySchedule(func(ctx context.Context) {
		if ctx.Err() != nil || ccr.resolver == nil {
			return
		}
		ccr.resolver.ResolveNow(o)
	})
}




func (ccr *ccResolverWrapper) close() {
	channelz.Info(logger, ccr.cc.channelz, "Closing the name resolver")
	ccr.mu.Lock()
	ccr.closed = true
	ccr.mu.Unlock()

	ccr.serializer.TrySchedule(func(context.Context) {
		if ccr.resolver == nil {
			return
		}
		ccr.resolver.Close()
		ccr.resolver = nil
	})
	ccr.serializerCancel()
}



func (ccr *ccResolverWrapper) UpdateState(s resolver.State) error {
	ccr.cc.mu.Lock()
	ccr.mu.Lock()
	if ccr.closed {
		ccr.mu.Unlock()
		ccr.cc.mu.Unlock()
		return nil
	}
	if s.Endpoints == nil {
		s.Endpoints = make([]resolver.Endpoint, 0, len(s.Addresses))
		for _, a := range s.Addresses {
			ep := resolver.Endpoint{Addresses: []resolver.Address{a}, Attributes: a.BalancerAttributes}
			ep.Addresses[0].BalancerAttributes = nil
			s.Endpoints = append(s.Endpoints, ep)
		}
	}
	ccr.addChannelzTraceEvent(s)
	ccr.curState = s
	ccr.mu.Unlock()
	return ccr.cc.updateResolverStateAndUnlock(s, nil)
}



func (ccr *ccResolverWrapper) ReportError(err error) {
	ccr.cc.mu.Lock()
	ccr.mu.Lock()
	if ccr.closed {
		ccr.mu.Unlock()
		ccr.cc.mu.Unlock()
		return
	}
	ccr.mu.Unlock()
	channelz.Warningf(logger, ccr.cc.channelz, "ccResolverWrapper: reporting error to cc: %v", err)
	ccr.cc.updateResolverStateAndUnlock(resolver.State{}, err)
}



func (ccr *ccResolverWrapper) NewAddress(addrs []resolver.Address) {
	ccr.cc.mu.Lock()
	ccr.mu.Lock()
	if ccr.closed {
		ccr.mu.Unlock()
		ccr.cc.mu.Unlock()
		return
	}
	s := resolver.State{Addresses: addrs, ServiceConfig: ccr.curState.ServiceConfig}
	ccr.addChannelzTraceEvent(s)
	ccr.curState = s
	ccr.mu.Unlock()
	ccr.cc.updateResolverStateAndUnlock(s, nil)
}



func (ccr *ccResolverWrapper) ParseServiceConfig(scJSON string) *serviceconfig.ParseResult {
	return parseServiceConfig(scJSON, ccr.cc.dopts.maxCallAttempts)
}



func (ccr *ccResolverWrapper) addChannelzTraceEvent(s resolver.State) {
	if !logger.V(0) && !channelz.IsOn() {
		return
	}
	var updates []string
	var oldSC, newSC *ServiceConfig
	var oldOK, newOK bool
	if ccr.curState.ServiceConfig != nil {
		oldSC, oldOK = ccr.curState.ServiceConfig.Config.(*ServiceConfig)
	}
	if s.ServiceConfig != nil {
		newSC, newOK = s.ServiceConfig.Config.(*ServiceConfig)
	}
	if oldOK != newOK || (oldOK && newOK && oldSC.rawJSONString != newSC.rawJSONString) {
		updates = append(updates, "service config updated")
	}
	if len(ccr.curState.Addresses) > 0 && len(s.Addresses) == 0 {
		updates = append(updates, "resolver returned an empty address list")
	} else if len(ccr.curState.Addresses) == 0 && len(s.Addresses) > 0 {
		updates = append(updates, "resolver returned new addresses")
	}
	channelz.Infof(logger, ccr.cc.channelz, "Resolver state updated: %s (%v)", pretty.ToJSON(s), strings.Join(updates, "; "))
}
