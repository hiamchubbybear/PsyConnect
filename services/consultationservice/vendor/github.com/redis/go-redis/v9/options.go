package redis

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9/auth"
	"github.com/redis/go-redis/v9/internal/pool"
	"github.com/redis/go-redis/v9/internal/proto"
)


type Limiter interface {
	
	
	
	Allow() error
	
	
	ReportResult(result error)
}


type Options struct {

	
	
	
	Network string

	
	Addr string

	
	ClientName string

	
	
	Dialer func(ctx context.Context, network, addr string) (net.Conn, error)

	
	OnConnect func(ctx context.Context, cn *Conn) error

	
	
	
	Protocol int

	
	
	
	Username string

	
	
	
	
	Password string

	
	
	CredentialsProvider func() (username string, password string)

	
	
	
	
	CredentialsProviderContext func(ctx context.Context) (username string, password string, err error)

	
	
	
	
	
	
	StreamingCredentialsProvider auth.StreamingCredentialsProvider

	
	DB int

	
	
	
	
	MaxRetries int

	
	
	
	
	MinRetryBackoff time.Duration

	
	
	
	MaxRetryBackoff time.Duration

	
	
	
	DialTimeout time.Duration

	
	
	
	
	
	
	
	ReadTimeout time.Duration

	
	
	
	
	
	
	
	WriteTimeout time.Duration

	
	
	ContextTimeoutEnabled bool

	
	
	
	
	
	ReadBufferSize int

	
	
	
	
	
	WriteBufferSize int

	
	
	
	
	
	
	
	PoolFIFO bool

	
	
	
	
	
	
	PoolSize int

	
	
	
	
	PoolTimeout time.Duration

	
	
	
	
	MinIdleConns int

	
	
	
	
	MaxIdleConns int

	
	
	
	MaxActiveConns int

	
	
	
	
	
	
	
	
	ConnMaxIdleTime time.Duration

	
	
	
	
	
	
	ConnMaxLifetime time.Duration

	
	TLSConfig *tls.Config

	
	Limiter Limiter

	
	readOnly bool

	
	
	
	
	
	DisableIndentity bool

	
	
	
	DisableIdentity bool

	
	
	IdentitySuffix string

	
	
	UnstableResp3 bool
}

func (opt *Options) init() {
	if opt.Addr == "" {
		opt.Addr = "localhost:6379"
	}
	if opt.Network == "" {
		if strings.HasPrefix(opt.Addr, "/") {
			opt.Network = "unix"
		} else {
			opt.Network = "tcp"
		}
	}
	if opt.Protocol < 2 {
		opt.Protocol = 3
	}
	if opt.DialTimeout == 0 {
		opt.DialTimeout = 5 * time.Second
	}
	if opt.Dialer == nil {
		opt.Dialer = NewDialer(opt)
	}
	if opt.PoolSize == 0 {
		opt.PoolSize = 10 * runtime.GOMAXPROCS(0)
	}
	if opt.ReadBufferSize == 0 {
		opt.ReadBufferSize = proto.DefaultBufferSize
	}
	if opt.WriteBufferSize == 0 {
		opt.WriteBufferSize = proto.DefaultBufferSize
	}
	switch opt.ReadTimeout {
	case -2:
		opt.ReadTimeout = -1
	case -1:
		opt.ReadTimeout = 0
	case 0:
		opt.ReadTimeout = 3 * time.Second
	}
	switch opt.WriteTimeout {
	case -2:
		opt.WriteTimeout = -1
	case -1:
		opt.WriteTimeout = 0
	case 0:
		opt.WriteTimeout = opt.ReadTimeout
	}
	if opt.PoolTimeout == 0 {
		if opt.ReadTimeout > 0 {
			opt.PoolTimeout = opt.ReadTimeout + time.Second
		} else {
			opt.PoolTimeout = 30 * time.Second
		}
	}
	if opt.ConnMaxIdleTime == 0 {
		opt.ConnMaxIdleTime = 30 * time.Minute
	}

	switch opt.MaxRetries {
	case -1:
		opt.MaxRetries = 0
	case 0:
		opt.MaxRetries = 3
	}
	switch opt.MinRetryBackoff {
	case -1:
		opt.MinRetryBackoff = 0
	case 0:
		opt.MinRetryBackoff = 8 * time.Millisecond
	}
	switch opt.MaxRetryBackoff {
	case -1:
		opt.MaxRetryBackoff = 0
	case 0:
		opt.MaxRetryBackoff = 512 * time.Millisecond
	}
}

func (opt *Options) clone() *Options {
	clone := *opt
	return &clone
}



func NewDialer(opt *Options) func(context.Context, string, string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		netDialer := &net.Dialer{
			Timeout:   opt.DialTimeout,
			KeepAlive: 5 * time.Minute,
		}
		if opt.TLSConfig == nil {
			return netDialer.DialContext(ctx, network, addr)
		}
		return tls.DialWithDialer(netDialer, network, addr, opt.TLSConfig)
	}
}






































func ParseURL(redisURL string) (*Options, error) {
	u, err := url.Parse(redisURL)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "redis", "rediss":
		return setupTCPConn(u)
	case "unix":
		return setupUnixConn(u)
	default:
		return nil, fmt.Errorf("redis: invalid URL scheme: %s", u.Scheme)
	}
}

func setupTCPConn(u *url.URL) (*Options, error) {
	o := &Options{Network: "tcp"}

	o.Username, o.Password = getUserPassword(u)

	h, p := getHostPortWithDefaults(u)
	o.Addr = net.JoinHostPort(h, p)

	f := strings.FieldsFunc(u.Path, func(r rune) bool {
		return r == '/'
	})
	switch len(f) {
	case 0:
		o.DB = 0
	case 1:
		var err error
		if o.DB, err = strconv.Atoi(f[0]); err != nil {
			return nil, fmt.Errorf("redis: invalid database number: %q", f[0])
		}
	default:
		return nil, fmt.Errorf("redis: invalid URL path: %s", u.Path)
	}

	if u.Scheme == "rediss" {
		o.TLSConfig = &tls.Config{
			ServerName: h,
			MinVersion: tls.VersionTLS12,
		}
	}

	return setupConnParams(u, o)
}




func getHostPortWithDefaults(u *url.URL) (string, string) {
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host
	}
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}
	return host, port
}

func setupUnixConn(u *url.URL) (*Options, error) {
	o := &Options{
		Network: "unix",
	}

	if strings.TrimSpace(u.Path) == "" { 
		return nil, errors.New("redis: empty unix socket path")
	}
	o.Addr = u.Path
	o.Username, o.Password = getUserPassword(u)
	return setupConnParams(u, o)
}

type queryOptions struct {
	q   url.Values
	err error
}

func (o *queryOptions) has(name string) bool {
	return len(o.q[name]) > 0
}

func (o *queryOptions) string(name string) string {
	vs := o.q[name]
	if len(vs) == 0 {
		return ""
	}
	delete(o.q, name) 
	return vs[len(vs)-1]
}

func (o *queryOptions) strings(name string) []string {
	vs := o.q[name]
	delete(o.q, name)
	return vs
}

func (o *queryOptions) int(name string) int {
	s := o.string(name)
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err == nil {
		return i
	}
	if o.err == nil {
		o.err = fmt.Errorf("redis: invalid %s number: %s", name, err)
	}
	return 0
}

func (o *queryOptions) duration(name string) time.Duration {
	s := o.string(name)
	if s == "" {
		return 0
	}
	
	if i, err := strconv.Atoi(s); err == nil {
		if i <= 0 {
			
			return -1
		}
		return time.Duration(i) * time.Second
	}
	dur, err := time.ParseDuration(s)
	if err == nil {
		return dur
	}
	if o.err == nil {
		o.err = fmt.Errorf("redis: invalid %s duration: %w", name, err)
	}
	return 0
}

func (o *queryOptions) bool(name string) bool {
	switch s := o.string(name); s {
	case "true", "1":
		return true
	case "false", "0", "":
		return false
	default:
		if o.err == nil {
			o.err = fmt.Errorf("redis: invalid %s boolean: expected true/false/1/0 or an empty string, got %q", name, s)
		}
		return false
	}
}

func (o *queryOptions) remaining() []string {
	if len(o.q) == 0 {
		return nil
	}
	keys := make([]string, 0, len(o.q))
	for k := range o.q {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}


func setupConnParams(u *url.URL, o *Options) (*Options, error) {
	q := queryOptions{q: u.Query()}

	
	if tmp := q.string("db"); tmp != "" {
		db, err := strconv.Atoi(tmp)
		if err != nil {
			return nil, fmt.Errorf("redis: invalid database number: %w", err)
		}
		o.DB = db
	}

	o.Protocol = q.int("protocol")
	o.ClientName = q.string("client_name")
	o.MaxRetries = q.int("max_retries")
	o.MinRetryBackoff = q.duration("min_retry_backoff")
	o.MaxRetryBackoff = q.duration("max_retry_backoff")
	o.DialTimeout = q.duration("dial_timeout")
	o.ReadTimeout = q.duration("read_timeout")
	o.WriteTimeout = q.duration("write_timeout")
	o.PoolFIFO = q.bool("pool_fifo")
	o.PoolSize = q.int("pool_size")
	o.PoolTimeout = q.duration("pool_timeout")
	o.MinIdleConns = q.int("min_idle_conns")
	o.MaxIdleConns = q.int("max_idle_conns")
	o.MaxActiveConns = q.int("max_active_conns")
	if q.has("conn_max_idle_time") {
		o.ConnMaxIdleTime = q.duration("conn_max_idle_time")
	} else {
		o.ConnMaxIdleTime = q.duration("idle_timeout")
	}
	if q.has("conn_max_lifetime") {
		o.ConnMaxLifetime = q.duration("conn_max_lifetime")
	} else {
		o.ConnMaxLifetime = q.duration("max_conn_age")
	}
	if q.err != nil {
		return nil, q.err
	}
	if o.TLSConfig != nil && q.has("skip_verify") {
		o.TLSConfig.InsecureSkipVerify = q.bool("skip_verify")
	}

	
	if r := q.remaining(); len(r) > 0 {
		return nil, fmt.Errorf("redis: unexpected option: %s", strings.Join(r, ", "))
	}

	return o, nil
}

func getUserPassword(u *url.URL) (string, string) {
	var user, password string
	if u.User != nil {
		user = u.User.Username()
		if p, ok := u.User.Password(); ok {
			password = p
		}
	}
	return user, password
}

func newConnPool(
	opt *Options,
	dialer func(ctx context.Context, network, addr string) (net.Conn, error),
) *pool.ConnPool {
	return pool.NewConnPool(&pool.Options{
		Dialer: func(ctx context.Context) (net.Conn, error) {
			return dialer(ctx, opt.Network, opt.Addr)
		},
		PoolFIFO:        opt.PoolFIFO,
		PoolSize:        opt.PoolSize,
		PoolTimeout:     opt.PoolTimeout,
		DialTimeout:     opt.DialTimeout,
		MinIdleConns:    opt.MinIdleConns,
		MaxIdleConns:    opt.MaxIdleConns,
		MaxActiveConns:  opt.MaxActiveConns,
		ConnMaxIdleTime: opt.ConnMaxIdleTime,
		ConnMaxLifetime: opt.ConnMaxLifetime,
		ReadBufferSize:  opt.ReadBufferSize,
		WriteBufferSize: opt.WriteBufferSize,
	})
}
