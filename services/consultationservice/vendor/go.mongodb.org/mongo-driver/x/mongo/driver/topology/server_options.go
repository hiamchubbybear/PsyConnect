





package topology

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

var defaultRegistry = bson.NewRegistryBuilder().Build()

type serverConfig struct {
	clock                *session.ClusterClock
	compressionOpts      []string
	connectionOpts       []ConnectionOption
	appname              string
	heartbeatInterval    time.Duration
	heartbeatTimeout     time.Duration
	serverMonitoringMode string
	serverMonitor        *event.ServerMonitor
	registry             *bsoncodec.Registry
	monitoringDisabled   bool
	serverAPI            *driver.ServerAPIOptions
	loadBalanced         bool

	
	maxConns             uint64
	minConns             uint64
	maxConnecting        uint64
	poolMonitor          *event.PoolMonitor
	logger               *logger.Logger
	poolMaxIdleTime      time.Duration
	poolMaintainInterval time.Duration
}

func newServerConfig(opts ...ServerOption) *serverConfig {
	cfg := &serverConfig{
		heartbeatInterval: 10 * time.Second,
		heartbeatTimeout:  10 * time.Second,
		registry:          defaultRegistry,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(cfg)
	}

	return cfg
}


type ServerOption func(*serverConfig)



func ServerAPIFromServerOptions(opts []ServerOption) *driver.ServerAPIOptions {
	return newServerConfig(opts...).serverAPI
}

func withMonitoringDisabled(fn func(bool) bool) ServerOption {
	return func(cfg *serverConfig) {
		cfg.monitoringDisabled = fn(cfg.monitoringDisabled)
	}
}


func WithConnectionOptions(fn func(...ConnectionOption) []ConnectionOption) ServerOption {
	return func(cfg *serverConfig) {
		cfg.connectionOpts = fn(cfg.connectionOpts...)
	}
}


func WithCompressionOptions(fn func(...string) []string) ServerOption {
	return func(cfg *serverConfig) {
		cfg.compressionOpts = fn(cfg.compressionOpts...)
	}
}


func WithServerAppName(fn func(string) string) ServerOption {
	return func(cfg *serverConfig) {
		cfg.appname = fn(cfg.appname)
	}
}


func WithHeartbeatInterval(fn func(time.Duration) time.Duration) ServerOption {
	return func(cfg *serverConfig) {
		cfg.heartbeatInterval = fn(cfg.heartbeatInterval)
	}
}



func WithHeartbeatTimeout(fn func(time.Duration) time.Duration) ServerOption {
	return func(cfg *serverConfig) {
		cfg.heartbeatTimeout = fn(cfg.heartbeatTimeout)
	}
}



func WithMaxConnections(fn func(uint64) uint64) ServerOption {
	return func(cfg *serverConfig) {
		cfg.maxConns = fn(cfg.maxConns)
	}
}




func WithMinConnections(fn func(uint64) uint64) ServerOption {
	return func(cfg *serverConfig) {
		cfg.minConns = fn(cfg.minConns)
	}
}




func WithMaxConnecting(fn func(uint64) uint64) ServerOption {
	return func(cfg *serverConfig) {
		cfg.maxConnecting = fn(cfg.maxConnecting)
	}
}




func WithConnectionPoolMaxIdleTime(fn func(time.Duration) time.Duration) ServerOption {
	return func(cfg *serverConfig) {
		cfg.poolMaxIdleTime = fn(cfg.poolMaxIdleTime)
	}
}



func WithConnectionPoolMaintainInterval(fn func(time.Duration) time.Duration) ServerOption {
	return func(cfg *serverConfig) {
		cfg.poolMaintainInterval = fn(cfg.poolMaintainInterval)
	}
}


func WithConnectionPoolMonitor(fn func(*event.PoolMonitor) *event.PoolMonitor) ServerOption {
	return func(cfg *serverConfig) {
		cfg.poolMonitor = fn(cfg.poolMonitor)
	}
}


func WithServerMonitor(fn func(*event.ServerMonitor) *event.ServerMonitor) ServerOption {
	return func(cfg *serverConfig) {
		cfg.serverMonitor = fn(cfg.serverMonitor)
	}
}


func WithClock(fn func(clock *session.ClusterClock) *session.ClusterClock) ServerOption {
	return func(cfg *serverConfig) {
		cfg.clock = fn(cfg.clock)
	}
}



func WithRegistry(fn func(*bsoncodec.Registry) *bsoncodec.Registry) ServerOption {
	return func(cfg *serverConfig) {
		cfg.registry = fn(cfg.registry)
	}
}


func WithServerAPI(fn func(serverAPI *driver.ServerAPIOptions) *driver.ServerAPIOptions) ServerOption {
	return func(cfg *serverConfig) {
		cfg.serverAPI = fn(cfg.serverAPI)
	}
}


func WithServerLoadBalanced(fn func(bool) bool) ServerOption {
	return func(cfg *serverConfig) {
		cfg.loadBalanced = fn(cfg.loadBalanced)
	}
}


func withLogger(fn func() *logger.Logger) ServerOption {
	return func(cfg *serverConfig) {
		cfg.logger = fn()
	}
}



func withServerMonitoringMode(mode *string) ServerOption {
	return func(cfg *serverConfig) {
		if mode != nil {
			cfg.serverMonitoringMode = *mode

			return
		}

		cfg.serverMonitoringMode = connstring.ServerMonitoringModeAuto
	}
}
