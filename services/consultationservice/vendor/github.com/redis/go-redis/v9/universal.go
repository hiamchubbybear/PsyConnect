package redis

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/redis/go-redis/v9/auth"
)



type UniversalOptions struct {
	
	
	Addrs []string

	
	ClientName string

	
	
	DB int

	

	Dialer    func(ctx context.Context, network, addr string) (net.Conn, error)
	OnConnect func(ctx context.Context, cn *Conn) error

	Protocol int
	Username string
	Password string
	
	
	CredentialsProvider func() (username string, password string)

	
	
	
	
	CredentialsProviderContext func(ctx context.Context) (username string, password string, err error)

	
	
	
	
	
	
	StreamingCredentialsProvider auth.StreamingCredentialsProvider

	SentinelUsername string
	SentinelPassword string

	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration

	DialTimeout           time.Duration
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
	ContextTimeoutEnabled bool

	
	PoolFIFO bool

	PoolSize        int
	PoolTimeout     time.Duration
	MinIdleConns    int
	MaxIdleConns    int
	MaxActiveConns  int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration

	TLSConfig *tls.Config

	

	MaxRedirects   int
	ReadOnly       bool
	RouteByLatency bool
	RouteRandomly  bool

	
	
	MasterName string

	
	
	
	
	
	DisableIndentity bool

	
	
	
	DisableIdentity bool

	IdentitySuffix string
	UnstableResp3  bool

	
	IsClusterMode bool
}


func (o *UniversalOptions) Cluster() *ClusterOptions {
	if len(o.Addrs) == 0 {
		o.Addrs = []string{"127.0.0.1:6379"}
	}

	return &ClusterOptions{
		Addrs:      o.Addrs,
		ClientName: o.ClientName,
		Dialer:     o.Dialer,
		OnConnect:  o.OnConnect,

		Protocol:                     o.Protocol,
		Username:                     o.Username,
		Password:                     o.Password,
		CredentialsProvider:          o.CredentialsProvider,
		CredentialsProviderContext:   o.CredentialsProviderContext,
		StreamingCredentialsProvider: o.StreamingCredentialsProvider,

		MaxRedirects:   o.MaxRedirects,
		ReadOnly:       o.ReadOnly,
		RouteByLatency: o.RouteByLatency,
		RouteRandomly:  o.RouteRandomly,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:           o.DialTimeout,
		ReadTimeout:           o.ReadTimeout,
		WriteTimeout:          o.WriteTimeout,
		ContextTimeoutEnabled: o.ContextTimeoutEnabled,

		PoolFIFO: o.PoolFIFO,

		PoolSize:        o.PoolSize,
		PoolTimeout:     o.PoolTimeout,
		MinIdleConns:    o.MinIdleConns,
		MaxIdleConns:    o.MaxIdleConns,
		MaxActiveConns:  o.MaxActiveConns,
		ConnMaxIdleTime: o.ConnMaxIdleTime,
		ConnMaxLifetime: o.ConnMaxLifetime,

		TLSConfig: o.TLSConfig,

		DisableIdentity:  o.DisableIdentity,
		DisableIndentity: o.DisableIndentity,
		IdentitySuffix:   o.IdentitySuffix,
		UnstableResp3:    o.UnstableResp3,
	}
}


func (o *UniversalOptions) Failover() *FailoverOptions {
	if len(o.Addrs) == 0 {
		o.Addrs = []string{"127.0.0.1:26379"}
	}

	return &FailoverOptions{
		SentinelAddrs: o.Addrs,
		MasterName:    o.MasterName,
		ClientName:    o.ClientName,

		Dialer:    o.Dialer,
		OnConnect: o.OnConnect,

		DB:                           o.DB,
		Protocol:                     o.Protocol,
		Username:                     o.Username,
		Password:                     o.Password,
		CredentialsProvider:          o.CredentialsProvider,
		CredentialsProviderContext:   o.CredentialsProviderContext,
		StreamingCredentialsProvider: o.StreamingCredentialsProvider,

		SentinelUsername: o.SentinelUsername,
		SentinelPassword: o.SentinelPassword,

		RouteByLatency: o.RouteByLatency,
		RouteRandomly:  o.RouteRandomly,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:           o.DialTimeout,
		ReadTimeout:           o.ReadTimeout,
		WriteTimeout:          o.WriteTimeout,
		ContextTimeoutEnabled: o.ContextTimeoutEnabled,

		PoolFIFO:        o.PoolFIFO,
		PoolSize:        o.PoolSize,
		PoolTimeout:     o.PoolTimeout,
		MinIdleConns:    o.MinIdleConns,
		MaxIdleConns:    o.MaxIdleConns,
		MaxActiveConns:  o.MaxActiveConns,
		ConnMaxIdleTime: o.ConnMaxIdleTime,
		ConnMaxLifetime: o.ConnMaxLifetime,

		TLSConfig: o.TLSConfig,

		ReplicaOnly: o.ReadOnly,

		DisableIdentity:  o.DisableIdentity,
		DisableIndentity: o.DisableIndentity,
		IdentitySuffix:   o.IdentitySuffix,
		UnstableResp3:    o.UnstableResp3,
	}
}


func (o *UniversalOptions) Simple() *Options {
	addr := "127.0.0.1:6379"
	if len(o.Addrs) > 0 {
		addr = o.Addrs[0]
	}

	return &Options{
		Addr:       addr,
		ClientName: o.ClientName,
		Dialer:     o.Dialer,
		OnConnect:  o.OnConnect,

		DB:                           o.DB,
		Protocol:                     o.Protocol,
		Username:                     o.Username,
		Password:                     o.Password,
		CredentialsProvider:          o.CredentialsProvider,
		CredentialsProviderContext:   o.CredentialsProviderContext,
		StreamingCredentialsProvider: o.StreamingCredentialsProvider,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:           o.DialTimeout,
		ReadTimeout:           o.ReadTimeout,
		WriteTimeout:          o.WriteTimeout,
		ContextTimeoutEnabled: o.ContextTimeoutEnabled,

		PoolFIFO:        o.PoolFIFO,
		PoolSize:        o.PoolSize,
		PoolTimeout:     o.PoolTimeout,
		MinIdleConns:    o.MinIdleConns,
		MaxIdleConns:    o.MaxIdleConns,
		MaxActiveConns:  o.MaxActiveConns,
		ConnMaxIdleTime: o.ConnMaxIdleTime,
		ConnMaxLifetime: o.ConnMaxLifetime,

		TLSConfig: o.TLSConfig,

		DisableIdentity:  o.DisableIdentity,
		DisableIndentity: o.DisableIndentity,
		IdentitySuffix:   o.IdentitySuffix,
		UnstableResp3:    o.UnstableResp3,
	}
}







type UniversalClient interface {
	Cmdable
	AddHook(Hook)
	Watch(ctx context.Context, fn func(*Tx) error, keys ...string) error
	Do(ctx context.Context, args ...interface{}) *Cmd
	Process(ctx context.Context, cmd Cmder) error
	Subscribe(ctx context.Context, channels ...string) *PubSub
	PSubscribe(ctx context.Context, channels ...string) *PubSub
	SSubscribe(ctx context.Context, channels ...string) *PubSub
	Close() error
	PoolStats() *PoolStats
}

var (
	_ UniversalClient = (*Client)(nil)
	_ UniversalClient = (*ClusterClient)(nil)
	_ UniversalClient = (*Ring)(nil)
)











func NewUniversalClient(opts *UniversalOptions) UniversalClient {
	if opts == nil {
		panic("redis: NewUniversalClient nil options")
	}

	switch {
	case opts.MasterName != "" && (opts.RouteByLatency || opts.RouteRandomly || opts.IsClusterMode):
		return NewFailoverClusterClient(opts.Failover())
	case opts.MasterName != "":
		return NewFailoverClient(opts.Failover())
	case len(opts.Addrs) > 1 || opts.IsClusterMode:
		return NewClusterClient(opts.Cluster())
	default:
		return NewClient(opts.Simple())
	}
}
