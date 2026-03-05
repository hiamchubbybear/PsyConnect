

package grpc

import (
	"context"
	"net"
	"net/url"
	"time"

	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/channelz"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/internal"
	internalbackoff "google.golang.org/grpc/internal/backoff"
	"google.golang.org/grpc/internal/binarylog"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/mem"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/stats"
)

const (
	
	defaultMaxCallAttempts = 5
)

func init() {
	internal.AddGlobalDialOptions = func(opt ...DialOption) {
		globalDialOptions = append(globalDialOptions, opt...)
	}
	internal.ClearGlobalDialOptions = func() {
		globalDialOptions = nil
	}
	internal.AddGlobalPerTargetDialOptions = func(opt any) {
		if ptdo, ok := opt.(perTargetDialOption); ok {
			globalPerTargetDialOptions = append(globalPerTargetDialOptions, ptdo)
		}
	}
	internal.ClearGlobalPerTargetDialOptions = func() {
		globalPerTargetDialOptions = nil
	}
	internal.WithBinaryLogger = withBinaryLogger
	internal.JoinDialOptions = newJoinDialOption
	internal.DisableGlobalDialOptions = newDisableGlobalDialOptions
	internal.WithBufferPool = withBufferPool
}



type dialOptions struct {
	unaryInt  UnaryClientInterceptor
	streamInt StreamClientInterceptor

	chainUnaryInts  []UnaryClientInterceptor
	chainStreamInts []StreamClientInterceptor

	compressorV0                Compressor
	dc                          Decompressor
	bs                          internalbackoff.Strategy
	block                       bool
	returnLastError             bool
	timeout                     time.Duration
	authority                   string
	binaryLogger                binarylog.Logger
	copts                       transport.ConnectOptions
	callOptions                 []CallOption
	channelzParent              channelz.Identifier
	disableServiceConfig        bool
	disableRetry                bool
	disableHealthCheck          bool
	minConnectTimeout           func() time.Duration
	defaultServiceConfig        *ServiceConfig 
	defaultServiceConfigRawJSON *string
	resolvers                   []resolver.Builder
	idleTimeout                 time.Duration
	defaultScheme               string
	maxCallAttempts             int
	enableLocalDNSResolution    bool 
	useProxy                    bool 
}


type DialOption interface {
	apply(*dialOptions)
}

var globalDialOptions []DialOption







type perTargetDialOption interface {
	
	DialOptionForTarget(parsedTarget url.URL) DialOption
}

var globalPerTargetDialOptions []perTargetDialOption








type EmptyDialOption struct{}

func (EmptyDialOption) apply(*dialOptions) {}

type disableGlobalDialOptions struct{}

func (disableGlobalDialOptions) apply(*dialOptions) {}



func newDisableGlobalDialOptions() DialOption {
	return &disableGlobalDialOptions{}
}



type funcDialOption struct {
	f func(*dialOptions)
}

func (fdo *funcDialOption) apply(do *dialOptions) {
	fdo.f(do)
}

func newFuncDialOption(f func(*dialOptions)) *funcDialOption {
	return &funcDialOption{
		f: f,
	}
}

type joinDialOption struct {
	opts []DialOption
}

func (jdo *joinDialOption) apply(do *dialOptions) {
	for _, opt := range jdo.opts {
		opt.apply(do)
	}
}

func newJoinDialOption(opts ...DialOption) DialOption {
	return &joinDialOption{opts: opts}
}









func WithSharedWriteBuffer(val bool) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.SharedWriteBuffer = val
	})
}







func WithWriteBufferSize(s int) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.WriteBufferSize = s
	})
}







func WithReadBufferSize(s int) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.ReadBufferSize = s
	})
}




func WithInitialWindowSize(s int32) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.InitialWindowSize = s
	})
}




func WithInitialConnWindowSize(s int32) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.InitialConnWindowSize = s
	})
}






func WithMaxMsgSize(s int) DialOption {
	return WithDefaultCallOptions(MaxCallRecvMsgSize(s))
}



func WithDefaultCallOptions(cos ...CallOption) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.callOptions = append(o.callOptions, cos...)
	})
}






func WithCodec(c Codec) DialOption {
	return WithDefaultCallOptions(CallCustomCodec(c))
}






func WithCompressor(cp Compressor) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.compressorV0 = cp
	})
}











func WithDecompressor(dc Decompressor) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.dc = dc
	})
}









func WithConnectParams(p ConnectParams) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.bs = internalbackoff.Exponential{Config: p.Backoff}
		o.minConnectTimeout = func() time.Duration {
			return p.MinConnectTimeout
		}
	})
}





func WithBackoffMaxDelay(md time.Duration) DialOption {
	return WithBackoffConfig(BackoffConfig{MaxDelay: md})
}





func WithBackoffConfig(b BackoffConfig) DialOption {
	bc := backoff.DefaultConfig
	bc.MaxDelay = b.MaxDelay
	return withBackoff(internalbackoff.Exponential{Config: bc})
}





func withBackoff(bs internalbackoff.Strategy) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.bs = bs
	})
}










func WithBlock() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.block = true
	})
}











func WithReturnConnectionError() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.block = true
		o.returnLastError = true
	})
}










func WithInsecure() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.TransportCredentials = insecure.NewCredentials()
	})
}








func WithNoProxy() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.useProxy = false
	})
}










func WithLocalDNSResolution() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.enableLocalDNSResolution = true
	})
}




func WithTransportCredentials(creds credentials.TransportCredentials) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.TransportCredentials = creds
	})
}



func WithPerRPCCredentials(creds credentials.PerRPCCredentials) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.PerRPCCredentials = append(o.copts.PerRPCCredentials, creds)
	})
}









func WithCredentialsBundle(b credentials.Bundle) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.CredsBundle = b
	})
}






func WithTimeout(d time.Duration) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.timeout = d
	})
}





















func WithContextDialer(f func(context.Context, string) (net.Conn, error)) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.Dialer = f
	})
}








func WithDialer(f func(string, time.Duration) (net.Conn, error)) DialOption {
	return WithContextDialer(
		func(ctx context.Context, addr string) (net.Conn, error) {
			if deadline, ok := ctx.Deadline(); ok {
				return f(addr, time.Until(deadline))
			}
			return f(addr, 0)
		})
}



func WithStatsHandler(h stats.Handler) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		if h == nil {
			logger.Error("ignoring nil parameter in grpc.WithStatsHandler ClientOption")
			
			
			return
		}
		o.copts.StatsHandlers = append(o.copts.StatsHandlers, h)
	})
}



func withBinaryLogger(bl binarylog.Logger) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.binaryLogger = bl
	})
}















func FailOnNonTempDialError(f bool) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.FailOnNonTempDialError = f
	})
}



func WithUserAgent(s string) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.UserAgent = s + " " + grpcUA
	})
}





func WithKeepaliveParams(kp keepalive.ClientParameters) DialOption {
	if kp.Time < internal.KeepaliveMinPingTime {
		logger.Warningf("Adjusting keepalive ping interval to minimum period of %v", internal.KeepaliveMinPingTime)
		kp.Time = internal.KeepaliveMinPingTime
	}
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.KeepaliveParams = kp
	})
}



func WithUnaryInterceptor(f UnaryClientInterceptor) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.unaryInt = f
	})
}






func WithChainUnaryInterceptor(interceptors ...UnaryClientInterceptor) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.chainUnaryInts = append(o.chainUnaryInts, interceptors...)
	})
}



func WithStreamInterceptor(f StreamClientInterceptor) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.streamInt = f
	})
}






func WithChainStreamInterceptor(interceptors ...StreamClientInterceptor) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.chainStreamInts = append(o.chainStreamInts, interceptors...)
	})
}



func WithAuthority(a string) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.authority = a
	})
}









func WithChannelzParentID(c channelz.Identifier) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.channelzParent = c
	})
}







func WithDisableServiceConfig() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.disableServiceConfig = true
	})
}














func WithDefaultServiceConfig(s string) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.defaultServiceConfigRawJSON = &s
	})
}





func WithDisableRetry() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.disableRetry = true
	})
}



type MaxHeaderListSizeDialOption struct {
	MaxHeaderListSize uint32
}

func (o MaxHeaderListSizeDialOption) apply(do *dialOptions) {
	do.copts.MaxHeaderListSize = &o.MaxHeaderListSize
}



func WithMaxHeaderListSize(s uint32) DialOption {
	return MaxHeaderListSizeDialOption{
		MaxHeaderListSize: s,
	}
}








func WithDisableHealthCheck() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.disableHealthCheck = true
	})
}

func defaultDialOptions() dialOptions {
	return dialOptions{
		copts: transport.ConnectOptions{
			ReadBufferSize:  defaultReadBufSize,
			WriteBufferSize: defaultWriteBufSize,
			UserAgent:       grpcUA,
			BufferPool:      mem.DefaultBufferPool(),
		},
		bs:                       internalbackoff.DefaultExponential,
		idleTimeout:              30 * time.Minute,
		defaultScheme:            "dns",
		maxCallAttempts:          defaultMaxCallAttempts,
		useProxy:                 true,
		enableLocalDNSResolution: false,
	}
}






func withMinConnectDeadline(f func() time.Duration) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.minConnectTimeout = f
	})
}



func withDefaultScheme(s string) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.defaultScheme = s
	})
}










func WithResolvers(rs ...resolver.Builder) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.resolvers = append(o.resolvers, rs...)
	})
}















func WithIdleTimeout(d time.Duration) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.idleTimeout = d
	})
}









func WithMaxCallAttempts(n int) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		if n < 2 {
			n = defaultMaxCallAttempts
		}
		o.maxCallAttempts = n
	})
}

func withBufferPool(bufferPool mem.BufferPool) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.BufferPool = bufferPool
	})
}
