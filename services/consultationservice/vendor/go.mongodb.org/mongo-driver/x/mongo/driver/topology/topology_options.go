





package topology

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth"
	"go.mongodb.org/mongo-driver/x/mongo/driver/ocsp"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

const defaultServerSelectionTimeout = 30 * time.Second


type Config struct {
	Mode                   MonitorMode
	ReplicaSetName         string
	SeedList               []string
	ServerOpts             []ServerOption
	URI                    string
	ServerSelectionTimeout time.Duration
	ServerMonitor          *event.ServerMonitor
	SRVMaxHosts            int
	SRVServiceName         string
	LoadBalanced           bool
	logger                 *logger.Logger
}


func ConvertToDriverAPIOptions(s *options.ServerAPIOptions) *driver.ServerAPIOptions {
	driverOpts := driver.NewServerAPIOptions(string(s.ServerAPIVersion))
	if s.Strict != nil {
		driverOpts.SetStrict(*s.Strict)
	}
	if s.DeprecationErrors != nil {
		driverOpts.SetDeprecationErrors(*s.DeprecationErrors)
	}
	return driverOpts
}

func newLogger(opts *options.LoggerOptions) (*logger.Logger, error) {
	if opts == nil {
		opts = options.Logger()
	}

	componentLevels := make(map[logger.Component]logger.Level)
	for component, level := range opts.ComponentLevels {
		componentLevels[logger.Component(component)] = logger.Level(level)
	}

	log, err := logger.New(opts.Sink, opts.MaxDocumentLength, componentLevels)
	if err != nil {
		return nil, fmt.Errorf("error creating logger: %w", err)
	}

	return log, nil
}



func convertOIDCArgs(args *driver.OIDCArgs) *options.OIDCArgs {
	if args == nil {
		return nil
	}
	return &options.OIDCArgs{
		Version:      args.Version,
		IDPInfo:      (*options.IDPInfo)(args.IDPInfo),
		RefreshToken: args.RefreshToken,
	}
}



func ConvertCreds(cred *options.Credential) *driver.Cred {
	if cred == nil {
		return nil
	}

	var oidcMachineCallback auth.OIDCCallback
	if cred.OIDCMachineCallback != nil {
		oidcMachineCallback = func(ctx context.Context, args *driver.OIDCArgs) (*driver.OIDCCredential, error) {
			cred, err := cred.OIDCMachineCallback(ctx, convertOIDCArgs(args))
			return (*driver.OIDCCredential)(cred), err
		}
	}

	var oidcHumanCallback auth.OIDCCallback
	if cred.OIDCHumanCallback != nil {
		oidcHumanCallback = func(ctx context.Context, args *driver.OIDCArgs) (*driver.OIDCCredential, error) {
			cred, err := cred.OIDCHumanCallback(ctx, convertOIDCArgs(args))
			return (*driver.OIDCCredential)(cred), err
		}
	}

	return &auth.Cred{
		Source:              cred.AuthSource,
		Username:            cred.Username,
		Password:            cred.Password,
		PasswordSet:         cred.PasswordSet,
		Props:               cred.AuthMechanismProperties,
		OIDCMachineCallback: oidcMachineCallback,
		OIDCHumanCallback:   oidcHumanCallback,
	}
}



func NewConfig(co *options.ClientOptions, clock *session.ClusterClock) (*Config, error) {
	var authenticator driver.Authenticator
	var err error
	if co.Auth != nil {
		authenticator, err = auth.CreateAuthenticator(
			co.Auth.AuthMechanism,
			ConvertCreds(co.Auth),
			co.HTTPClient,
		)
		if err != nil {
			return nil, fmt.Errorf("error creating authenticator: %w", err)
		}
	}
	return NewConfigWithAuthenticator(co, clock, authenticator)
}





func NewConfigWithAuthenticator(
	co *options.ClientOptions,
	clock *session.ClusterClock,
	authenticator driver.Authenticator,
) (*Config, error) {
	var serverAPI *driver.ServerAPIOptions

	if err := co.Validate(); err != nil {
		return nil, err
	}

	var connOpts []ConnectionOption
	var serverOpts []ServerOption

	cfgp := &Config{}

	
	cfgp.ServerSelectionTimeout = defaultServerSelectionTimeout

	
	cfgp.SeedList = []string{"localhost:27017"}

	

	
	
	if co.ServerAPIOptions != nil {
		serverAPI = ConvertToDriverAPIOptions(co.ServerAPIOptions)
		serverOpts = append(serverOpts, WithServerAPI(func(*driver.ServerAPIOptions) *driver.ServerAPIOptions {
			return serverAPI
		}))
	}

	cfgp.URI = co.GetURI()

	if co.SRVServiceName != nil {
		cfgp.SRVServiceName = *co.SRVServiceName
	}

	if co.SRVMaxHosts != nil {
		cfgp.SRVMaxHosts = *co.SRVMaxHosts
	}

	
	var appName string
	if co.AppName != nil {
		appName = *co.AppName

		serverOpts = append(serverOpts, WithServerAppName(func(string) string {
			return appName
		}))
	}
	
	var comps []string
	if len(co.Compressors) > 0 {
		comps = co.Compressors

		connOpts = append(connOpts, WithCompressors(
			func(compressors []string) []string {
				return append(compressors, comps...)
			},
		))

		for _, comp := range comps {
			switch comp {
			case "zlib":
				connOpts = append(connOpts, WithZlibLevel(func(*int) *int {
					return co.ZlibLevel
				}))
			case "zstd":
				connOpts = append(connOpts, WithZstdLevel(func(*int) *int {
					return co.ZstdLevel
				}))
			}
		}

		serverOpts = append(serverOpts, WithCompressionOptions(
			func(opts ...string) []string { return append(opts, comps...) },
		))
	}

	var loadBalanced bool
	if co.LoadBalanced != nil {
		loadBalanced = *co.LoadBalanced
	}

	
	var handshaker func(driver.Handshaker) driver.Handshaker
	if authenticator != nil {
		handshakeOpts := &auth.HandshakeOptions{
			AppName:       appName,
			Authenticator: authenticator,
			Compressors:   comps,
			ServerAPI:     serverAPI,
			LoadBalanced:  loadBalanced,
			ClusterClock:  clock,
		}

		if co.Auth.AuthMechanism == "" {
			
			handshakeOpts.DBUser = co.Auth.AuthSource + "." + co.Auth.Username
		}
		if co.AuthenticateToAnything != nil && *co.AuthenticateToAnything {
			
			handshakeOpts.PerformAuthentication = func(description.Server) bool {
				return true
			}
		}

		handshaker = func(driver.Handshaker) driver.Handshaker {
			return auth.Handshaker(nil, handshakeOpts)
		}
	} else {
		handshaker = func(driver.Handshaker) driver.Handshaker {
			return operation.NewHello().
				AppName(appName).
				Compressors(comps).
				ClusterClock(clock).
				ServerAPI(serverAPI).
				LoadBalanced(loadBalanced)
		}
	}

	connOpts = append(connOpts, WithHandshaker(handshaker))
	
	if co.ConnectTimeout != nil {
		serverOpts = append(serverOpts, WithHeartbeatTimeout(
			func(time.Duration) time.Duration { return *co.ConnectTimeout },
		))
		connOpts = append(connOpts, WithConnectTimeout(
			func(time.Duration) time.Duration { return *co.ConnectTimeout },
		))
	}
	
	if co.Dialer != nil {
		connOpts = append(connOpts, WithDialer(
			func(Dialer) Dialer { return co.Dialer },
		))
	}
	
	if co.Direct != nil && *co.Direct {
		cfgp.Mode = SingleMode
	}

	
	if co.HeartbeatInterval != nil {
		serverOpts = append(serverOpts, WithHeartbeatInterval(
			func(time.Duration) time.Duration { return *co.HeartbeatInterval },
		))
	}
	
	cfgp.SeedList = []string{"localhost:27017"} 
	if len(co.Hosts) > 0 {
		cfgp.SeedList = co.Hosts
	}

	
	if co.MaxConnIdleTime != nil {
		serverOpts = append(serverOpts, WithConnectionPoolMaxIdleTime(
			func(time.Duration) time.Duration { return *co.MaxConnIdleTime },
		))
	}
	
	if co.MaxPoolSize != nil {
		serverOpts = append(
			serverOpts,
			WithMaxConnections(func(uint64) uint64 { return *co.MaxPoolSize }),
		)
	}
	
	if co.MinPoolSize != nil {
		serverOpts = append(
			serverOpts,
			WithMinConnections(func(uint64) uint64 { return *co.MinPoolSize }),
		)
	}
	
	if co.MaxConnecting != nil {
		serverOpts = append(
			serverOpts,
			WithMaxConnecting(func(uint64) uint64 { return *co.MaxConnecting }),
		)
	}
	
	if co.PoolMonitor != nil {
		serverOpts = append(
			serverOpts,
			WithConnectionPoolMonitor(func(*event.PoolMonitor) *event.PoolMonitor { return co.PoolMonitor }),
		)
	}
	
	if co.Monitor != nil {
		connOpts = append(connOpts, WithMonitor(
			func(*event.CommandMonitor) *event.CommandMonitor { return co.Monitor },
		))
	}
	
	if co.ServerMonitor != nil {
		serverOpts = append(
			serverOpts,
			WithServerMonitor(func(*event.ServerMonitor) *event.ServerMonitor { return co.ServerMonitor }),
		)
		cfgp.ServerMonitor = co.ServerMonitor
	}
	
	if co.ReplicaSet != nil {
		cfgp.ReplicaSetName = *co.ReplicaSet
	}
	
	if co.ServerSelectionTimeout != nil {
		cfgp.ServerSelectionTimeout = *co.ServerSelectionTimeout
	}
	
	if co.SocketTimeout != nil {
		connOpts = append(
			connOpts,
			WithReadTimeout(func(time.Duration) time.Duration { return *co.SocketTimeout }),
			WithWriteTimeout(func(time.Duration) time.Duration { return *co.SocketTimeout }),
		)
	}
	
	if co.TLSConfig != nil {
		connOpts = append(connOpts, WithTLSConfig(
			func(*tls.Config) *tls.Config {
				return co.TLSConfig
			},
		))
	}

	
	if co.HTTPClient != nil {
		connOpts = append(connOpts, WithHTTPClient(
			func(*http.Client) *http.Client {
				return co.HTTPClient
			},
		))
	}

	
	ocspCache := ocsp.NewCache()
	connOpts = append(
		connOpts,
		WithOCSPCache(func(ocsp.Cache) ocsp.Cache { return ocspCache }),
	)

	
	if co.DisableOCSPEndpointCheck != nil {
		connOpts = append(
			connOpts,
			WithDisableOCSPEndpointCheck(func(bool) bool { return *co.DisableOCSPEndpointCheck }),
		)
	}

	
	if co.LoadBalanced != nil {
		cfgp.LoadBalanced = *co.LoadBalanced

		serverOpts = append(
			serverOpts,
			WithServerLoadBalanced(func(bool) bool { return *co.LoadBalanced }),
		)
		connOpts = append(
			connOpts,
			WithConnectionLoadBalanced(func(bool) bool { return *co.LoadBalanced }),
		)
	}

	lgr, err := newLogger(co.LoggerOptions)
	if err != nil {
		return nil, err
	}

	serverOpts = append(
		serverOpts,
		withLogger(func() *logger.Logger { return lgr }),
		withServerMonitoringMode(co.ServerMonitoringMode),
	)

	cfgp.logger = lgr

	serverOpts = append(
		serverOpts,
		WithClock(func(*session.ClusterClock) *session.ClusterClock { return clock }),
		WithConnectionOptions(func(...ConnectionOption) []ConnectionOption { return connOpts }))

	cfgp.ServerOpts = serverOpts

	return cfgp, nil
}
