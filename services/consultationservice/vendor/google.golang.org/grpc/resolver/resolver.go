



package resolver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/experimental/stats"
	"google.golang.org/grpc/internal"
	"google.golang.org/grpc/serviceconfig"
)

var (
	
	m = make(map[string]Builder)
	
	defaultScheme = "passthrough"
)










func Register(b Builder) {
	m[b.Scheme()] = b
}




func Get(scheme string) Builder {
	if b, ok := m[scheme]; ok {
		return b
	}
	return nil
}







func SetDefaultScheme(scheme string) {
	defaultScheme = scheme
	internal.UserSetDefaultScheme = true
}



func GetDefaultScheme() string {
	return defaultScheme
}







type Address struct {
	
	Addr string

	
	
	
	
	
	
	
	
	ServerName string

	
	
	Attributes *attributes.Attributes

	
	
	
	
	
	
	BalancerAttributes *attributes.Attributes

	
	
	
	
	Metadata any
}







func (a Address) Equal(o Address) bool {
	return a.Addr == o.Addr && a.ServerName == o.ServerName &&
		a.Attributes.Equal(o.Attributes) &&
		a.BalancerAttributes.Equal(o.BalancerAttributes) &&
		a.Metadata == o.Metadata
}


func (a Address) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("{Addr: %q, ", a.Addr))
	sb.WriteString(fmt.Sprintf("ServerName: %q, ", a.ServerName))
	if a.Attributes != nil {
		sb.WriteString(fmt.Sprintf("Attributes: %v, ", a.Attributes.String()))
	}
	if a.BalancerAttributes != nil {
		sb.WriteString(fmt.Sprintf("BalancerAttributes: %v", a.BalancerAttributes.String()))
	}
	sb.WriteString("}")
	return sb.String()
}



type BuildOptions struct {
	
	
	DisableServiceConfig bool
	
	
	
	
	
	DialCreds credentials.TransportCredentials
	
	
	
	
	
	CredsBundle credentials.Bundle
	
	
	
	
	
	Dialer func(context.Context, string) (net.Conn, error)
	
	
	Authority string
	
	MetricsRecorder stats.MetricsRecorder
}



type Endpoint struct {
	
	Addresses []Address

	
	
	Attributes *attributes.Attributes
}


type State struct {
	
	
	
	
	
	
	
	
	
	Addresses []Address

	
	
	
	
	
	Endpoints []Endpoint

	
	
	
	ServiceConfig *serviceconfig.ParseResult

	
	
	Attributes *attributes.Attributes
}








type ClientConn interface {
	
	
	
	
	
	
	
	
	
	
	UpdateState(State) error
	
	
	
	ReportError(error)
	
	
	
	
	
	NewAddress(addresses []Address)
	
	
	ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult
}











type Target struct {
	
	
	
	
	URL url.URL
}



func (t Target) Endpoint() string {
	endpoint := t.URL.Path
	if endpoint == "" {
		endpoint = t.URL.Opaque
	}
	
	
	
	
	
	
	
	
	return strings.TrimPrefix(endpoint, "/")
}


func (t Target) String() string {
	return t.URL.Scheme + "://" + t.URL.Host + "/" + t.Endpoint()
}


type Builder interface {
	
	
	
	
	Build(target Target, cc ClientConn, opts BuildOptions) (Resolver, error)
	
	
	
	
	Scheme() string
}


type ResolveNowOptions struct{}



type Resolver interface {
	
	
	
	
	ResolveNow(ResolveNowOptions)
	
	Close()
}




type AuthorityOverrider interface {
	
	
	
	OverrideAuthority(Target) string
}





func ValidateEndpoints(endpoints []Endpoint) error {
	if len(endpoints) == 0 {
		return errors.New("endpoints list is empty")
	}

	for _, endpoint := range endpoints {
		for range endpoint.Addresses {
			return nil
		}
	}
	return errors.New("endpoints list contains no addresses")
}
