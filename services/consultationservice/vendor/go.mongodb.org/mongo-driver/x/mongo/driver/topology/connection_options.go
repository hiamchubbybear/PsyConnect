





package topology

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/httputil"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/ocsp"
)


type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}


type DialerFunc func(ctx context.Context, network, address string) (net.Conn, error)


func (df DialerFunc) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return df(ctx, network, address)
}





var DefaultDialer Dialer = &net.Dialer{}




type Handshaker = driver.Handshaker


type generationNumberFn func(serviceID *primitive.ObjectID) uint64

type connectionConfig struct {
	connectTimeout           time.Duration
	dialer                   Dialer
	handshaker               Handshaker
	idleTimeout              time.Duration
	cmdMonitor               *event.CommandMonitor
	readTimeout              time.Duration
	writeTimeout             time.Duration
	tlsConfig                *tls.Config
	httpClient               *http.Client
	compressors              []string
	zlibLevel                *int
	zstdLevel                *int
	ocspCache                ocsp.Cache
	disableOCSPEndpointCheck bool
	tlsConnectionSource      tlsConnectionSource
	loadBalanced             bool
	getGenerationFn          generationNumberFn
}

func newConnectionConfig(opts ...ConnectionOption) *connectionConfig {
	cfg := &connectionConfig{
		connectTimeout:      30 * time.Second,
		dialer:              nil,
		tlsConnectionSource: defaultTLSConnectionSource,
		httpClient:          httputil.DefaultHTTPClient,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(cfg)
	}

	if cfg.dialer == nil {
		
		
		cfg.dialer = &net.Dialer{}
	}

	return cfg
}


type ConnectionOption func(*connectionConfig)

func withTLSConnectionSource(fn func(tlsConnectionSource) tlsConnectionSource) ConnectionOption {
	return func(c *connectionConfig) {
		c.tlsConnectionSource = fn(c.tlsConnectionSource)
	}
}


func WithCompressors(fn func([]string) []string) ConnectionOption {
	return func(c *connectionConfig) {
		c.compressors = fn(c.compressors)
	}
}



func WithConnectTimeout(fn func(time.Duration) time.Duration) ConnectionOption {
	return func(c *connectionConfig) {
		c.connectTimeout = fn(c.connectTimeout)
	}
}


func WithDialer(fn func(Dialer) Dialer) ConnectionOption {
	return func(c *connectionConfig) {
		c.dialer = fn(c.dialer)
	}
}



func WithHandshaker(fn func(Handshaker) Handshaker) ConnectionOption {
	return func(c *connectionConfig) {
		c.handshaker = fn(c.handshaker)
	}
}


func WithIdleTimeout(fn func(time.Duration) time.Duration) ConnectionOption {
	return func(c *connectionConfig) {
		c.idleTimeout = fn(c.idleTimeout)
	}
}


func WithReadTimeout(fn func(time.Duration) time.Duration) ConnectionOption {
	return func(c *connectionConfig) {
		c.readTimeout = fn(c.readTimeout)
	}
}


func WithWriteTimeout(fn func(time.Duration) time.Duration) ConnectionOption {
	return func(c *connectionConfig) {
		c.writeTimeout = fn(c.writeTimeout)
	}
}


func WithTLSConfig(fn func(*tls.Config) *tls.Config) ConnectionOption {
	return func(c *connectionConfig) {
		c.tlsConfig = fn(c.tlsConfig)
	}
}


func WithHTTPClient(fn func(*http.Client) *http.Client) ConnectionOption {
	return func(c *connectionConfig) {
		c.httpClient = fn(c.httpClient)
	}
}


func WithMonitor(fn func(*event.CommandMonitor) *event.CommandMonitor) ConnectionOption {
	return func(c *connectionConfig) {
		c.cmdMonitor = fn(c.cmdMonitor)
	}
}


func WithZlibLevel(fn func(*int) *int) ConnectionOption {
	return func(c *connectionConfig) {
		c.zlibLevel = fn(c.zlibLevel)
	}
}


func WithZstdLevel(fn func(*int) *int) ConnectionOption {
	return func(c *connectionConfig) {
		c.zstdLevel = fn(c.zstdLevel)
	}
}


func WithOCSPCache(fn func(ocsp.Cache) ocsp.Cache) ConnectionOption {
	return func(c *connectionConfig) {
		c.ocspCache = fn(c.ocspCache)
	}
}




func WithDisableOCSPEndpointCheck(fn func(bool) bool) ConnectionOption {
	return func(c *connectionConfig) {
		c.disableOCSPEndpointCheck = fn(c.disableOCSPEndpointCheck)
	}
}


func WithConnectionLoadBalanced(fn func(bool) bool) ConnectionOption {
	return func(c *connectionConfig) {
		c.loadBalanced = fn(c.loadBalanced)
	}
}

func withGenerationNumberFn(fn func(generationNumberFn) generationNumberFn) ConnectionOption {
	return func(c *connectionConfig) {
		c.getGenerationFn = fn(c.getGenerationFn)
	}
}
