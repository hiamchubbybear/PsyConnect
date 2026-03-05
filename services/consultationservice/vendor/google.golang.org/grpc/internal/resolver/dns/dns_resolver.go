



package dns

import (
	"context"
	"encoding/json"
	"fmt"
	rand "math/rand/v2"
	"net"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	grpclbstate "google.golang.org/grpc/balancer/grpclb/state"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/backoff"
	"google.golang.org/grpc/internal/envconfig"
	"google.golang.org/grpc/internal/resolver/dns/internal"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

var (
	
	
	EnableSRVLookups = false

	
	
	MinResolutionInterval = 30 * time.Second

	
	
	
	
	
	ResolvingTimeout = 30 * time.Second

	logger = grpclog.Component("dns")
)

func init() {
	resolver.Register(NewBuilder())
	internal.TimeAfterFunc = time.After
	internal.TimeNowFunc = time.Now
	internal.TimeUntilFunc = time.Until
	internal.NewNetResolver = newNetResolver
	internal.AddressDialer = addressDialer
}

const (
	defaultPort       = "443"
	defaultDNSSvrPort = "53"
	golang            = "GO"
	
	
	txtPrefix = "_grpc_config."
	
	
	txtAttribute = "grpc_config="
)

var addressDialer = func(address string) func(context.Context, string, string) (net.Conn, error) {
	return func(ctx context.Context, network, _ string) (net.Conn, error) {
		var dialer net.Dialer
		return dialer.DialContext(ctx, network, address)
	}
}

var newNetResolver = func(authority string) (internal.NetResolver, error) {
	if authority == "" {
		return net.DefaultResolver, nil
	}

	host, port, err := parseTarget(authority, defaultDNSSvrPort)
	if err != nil {
		return nil, err
	}

	authorityWithPort := net.JoinHostPort(host, port)

	return &net.Resolver{
		PreferGo: true,
		Dial:     internal.AddressDialer(authorityWithPort),
	}, nil
}


func NewBuilder() resolver.Builder {
	return &dnsBuilder{}
}

type dnsBuilder struct{}



func (b *dnsBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	host, port, err := parseTarget(target.Endpoint(), defaultPort)
	if err != nil {
		return nil, err
	}

	
	if ipAddr, err := formatIP(host); err == nil {
		addr := []resolver.Address{{Addr: ipAddr + ":" + port}}
		cc.UpdateState(resolver.State{Addresses: addr})
		return deadResolver{}, nil
	}

	
	ctx, cancel := context.WithCancel(context.Background())
	d := &dnsResolver{
		host:                 host,
		port:                 port,
		ctx:                  ctx,
		cancel:               cancel,
		cc:                   cc,
		rn:                   make(chan struct{}, 1),
		disableServiceConfig: opts.DisableServiceConfig,
	}

	d.resolver, err = internal.NewNetResolver(target.URL.Host)
	if err != nil {
		return nil, err
	}

	d.wg.Add(1)
	go d.watcher()
	return d, nil
}


func (b *dnsBuilder) Scheme() string {
	return "dns"
}


type deadResolver struct{}

func (deadResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (deadResolver) Close() {}


type dnsResolver struct {
	host     string
	port     string
	resolver internal.NetResolver
	ctx      context.Context
	cancel   context.CancelFunc
	cc       resolver.ClientConn
	
	
	rn chan struct{}
	
	
	
	
	
	
	
	wg                   sync.WaitGroup
	disableServiceConfig bool
}



func (d *dnsResolver) ResolveNow(resolver.ResolveNowOptions) {
	select {
	case d.rn <- struct{}{}:
	default:
	}
}


func (d *dnsResolver) Close() {
	d.cancel()
	d.wg.Wait()
}

func (d *dnsResolver) watcher() {
	defer d.wg.Done()
	backoffIndex := 1
	for {
		state, err := d.lookup()
		if err != nil {
			
			d.cc.ReportError(err)
		} else {
			err = d.cc.UpdateState(*state)
		}

		var nextResolutionTime time.Time
		if err == nil {
			
			
			backoffIndex = 1
			nextResolutionTime = internal.TimeNowFunc().Add(MinResolutionInterval)
			select {
			case <-d.ctx.Done():
				return
			case <-d.rn:
			}
		} else {
			
			
			nextResolutionTime = internal.TimeNowFunc().Add(backoff.DefaultExponential.Backoff(backoffIndex))
			backoffIndex++
		}
		select {
		case <-d.ctx.Done():
			return
		case <-internal.TimeAfterFunc(internal.TimeUntilFunc(nextResolutionTime)):
		}
	}
}

func (d *dnsResolver) lookupSRV(ctx context.Context) ([]resolver.Address, error) {
	
	
	if !EnableSRVLookups || d.host == "metadata.google.internal." {
		return nil, nil
	}
	var newAddrs []resolver.Address
	_, srvs, err := d.resolver.LookupSRV(ctx, "grpclb", "tcp", d.host)
	if err != nil {
		err = handleDNSError(err, "SRV") 
		return nil, err
	}
	for _, s := range srvs {
		lbAddrs, err := d.resolver.LookupHost(ctx, s.Target)
		if err != nil {
			err = handleDNSError(err, "A") 
			if err == nil {
				
				
				continue
			}
			return nil, err
		}
		for _, a := range lbAddrs {
			ip, err := formatIP(a)
			if err != nil {
				return nil, fmt.Errorf("dns: error parsing A record IP address %v: %v", a, err)
			}
			addr := ip + ":" + strconv.Itoa(int(s.Port))
			newAddrs = append(newAddrs, resolver.Address{Addr: addr, ServerName: s.Target})
		}
	}
	return newAddrs, nil
}

func handleDNSError(err error, lookupType string) error {
	dnsErr, ok := err.(*net.DNSError)
	if ok && !dnsErr.IsTimeout && !dnsErr.IsTemporary {
		
		
		
		return nil
	}
	if err != nil {
		err = fmt.Errorf("dns: %v record lookup error: %v", lookupType, err)
		logger.Info(err)
	}
	return err
}

func (d *dnsResolver) lookupTXT(ctx context.Context) *serviceconfig.ParseResult {
	ss, err := d.resolver.LookupTXT(ctx, txtPrefix+d.host)
	if err != nil {
		if envconfig.TXTErrIgnore {
			return nil
		}
		if err = handleDNSError(err, "TXT"); err != nil {
			return &serviceconfig.ParseResult{Err: err}
		}
		return nil
	}
	var res string
	for _, s := range ss {
		res += s
	}

	
	
	if !strings.HasPrefix(res, txtAttribute) {
		logger.Warningf("dns: TXT record %v missing %v attribute", res, txtAttribute)
		
		
		return nil
	}
	sc := canaryingSC(strings.TrimPrefix(res, txtAttribute))
	return d.cc.ParseServiceConfig(sc)
}

func (d *dnsResolver) lookupHost(ctx context.Context) ([]resolver.Address, error) {
	addrs, err := d.resolver.LookupHost(ctx, d.host)
	if err != nil {
		err = handleDNSError(err, "A")
		return nil, err
	}
	newAddrs := make([]resolver.Address, 0, len(addrs))
	for _, a := range addrs {
		ip, err := formatIP(a)
		if err != nil {
			return nil, fmt.Errorf("dns: error parsing A record IP address %v: %v", a, err)
		}
		addr := ip + ":" + d.port
		newAddrs = append(newAddrs, resolver.Address{Addr: addr})
	}
	return newAddrs, nil
}

func (d *dnsResolver) lookup() (*resolver.State, error) {
	ctx, cancel := context.WithTimeout(d.ctx, ResolvingTimeout)
	defer cancel()
	srv, srvErr := d.lookupSRV(ctx)
	addrs, hostErr := d.lookupHost(ctx)
	if hostErr != nil && (srvErr != nil || len(srv) == 0) {
		return nil, hostErr
	}

	state := resolver.State{Addresses: addrs}
	if len(srv) > 0 {
		state = grpclbstate.Set(state, &grpclbstate.State{BalancerAddresses: srv})
	}
	if !d.disableServiceConfig {
		state.ServiceConfig = d.lookupTXT(ctx)
	}
	return &state, nil
}





func formatIP(addr string) (string, error) {
	ip, err := netip.ParseAddr(addr)
	if err != nil {
		return "", err
	}
	if ip.Is4() {
		return addr, nil
	}
	return "[" + addr + "]", nil
}










func parseTarget(target, defaultPort string) (host, port string, err error) {
	if target == "" {
		return "", "", internal.ErrMissingAddr
	}
	if _, err := netip.ParseAddr(target); err == nil {
		
		return target, defaultPort, nil
	}
	if host, port, err = net.SplitHostPort(target); err == nil {
		if port == "" {
			
			
			return "", "", internal.ErrEndsWithColon
		}
		
		if host == "" {
			
			
			host = "localhost"
		}
		return host, port, nil
	}
	if host, port, err = net.SplitHostPort(target + ":" + defaultPort); err == nil {
		
		return host, port, nil
	}
	return "", "", fmt.Errorf("invalid target address %v, error info: %v", target, err)
}

type rawChoice struct {
	ClientLanguage *[]string        `json:"clientLanguage,omitempty"`
	Percentage     *int             `json:"percentage,omitempty"`
	ClientHostName *[]string        `json:"clientHostName,omitempty"`
	ServiceConfig  *json.RawMessage `json:"serviceConfig,omitempty"`
}

func containsString(a *[]string, b string) bool {
	if a == nil {
		return true
	}
	for _, c := range *a {
		if c == b {
			return true
		}
	}
	return false
}

func chosenByPercentage(a *int) bool {
	if a == nil {
		return true
	}
	return rand.IntN(100)+1 <= *a
}

func canaryingSC(js string) string {
	if js == "" {
		return ""
	}
	var rcs []rawChoice
	err := json.Unmarshal([]byte(js), &rcs)
	if err != nil {
		logger.Warningf("dns: error parsing service config json: %v", err)
		return ""
	}
	cliHostname, err := os.Hostname()
	if err != nil {
		logger.Warningf("dns: error getting client hostname: %v", err)
		return ""
	}
	var sc string
	for _, c := range rcs {
		if !containsString(c.ClientLanguage, golang) ||
			!chosenByPercentage(c.Percentage) ||
			!containsString(c.ClientHostName, cliHostname) ||
			c.ServiceConfig == nil {
			continue
		}
		sc = string(*c.ServiceConfig)
		break
	}
	return sc
}
