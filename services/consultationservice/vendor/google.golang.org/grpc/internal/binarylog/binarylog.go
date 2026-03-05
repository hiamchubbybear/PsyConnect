



package binarylog

import (
	"fmt"
	"os"

	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/grpcutil"
)

var grpclogLogger = grpclog.Component("binarylog")






type Logger interface {
	GetMethodLogger(methodName string) MethodLogger
}





var binLogger Logger




func SetLogger(l Logger) {
	binLogger = l
}




func GetLogger() Logger {
	return binLogger
}







func GetMethodLogger(methodName string) MethodLogger {
	if binLogger == nil {
		return nil
	}
	return binLogger.GetMethodLogger(methodName)
}

func init() {
	const envStr = "GRPC_BINARY_LOG_FILTER"
	configStr := os.Getenv(envStr)
	binLogger = NewLoggerFromConfigString(configStr)
}



type MethodLoggerConfig struct {
	
	Header, Message uint64
}


type LoggerConfig struct {
	All      *MethodLoggerConfig
	Services map[string]*MethodLoggerConfig
	Methods  map[string]*MethodLoggerConfig

	Blacklist map[string]struct{}
}

type logger struct {
	config LoggerConfig
}


func NewLoggerFromConfig(config LoggerConfig) Logger {
	return &logger{config: config}
}



func newEmptyLogger() *logger {
	return &logger{}
}


func (l *logger) setDefaultMethodLogger(ml *MethodLoggerConfig) error {
	if l.config.All != nil {
		return fmt.Errorf("conflicting global rules found")
	}
	l.config.All = ml
	return nil
}




func (l *logger) setServiceMethodLogger(service string, ml *MethodLoggerConfig) error {
	if _, ok := l.config.Services[service]; ok {
		return fmt.Errorf("conflicting service rules for service %v found", service)
	}
	if l.config.Services == nil {
		l.config.Services = make(map[string]*MethodLoggerConfig)
	}
	l.config.Services[service] = ml
	return nil
}




func (l *logger) setMethodMethodLogger(method string, ml *MethodLoggerConfig) error {
	if _, ok := l.config.Blacklist[method]; ok {
		return fmt.Errorf("conflicting blacklist rules for method %v found", method)
	}
	if _, ok := l.config.Methods[method]; ok {
		return fmt.Errorf("conflicting method rules for method %v found", method)
	}
	if l.config.Methods == nil {
		l.config.Methods = make(map[string]*MethodLoggerConfig)
	}
	l.config.Methods[method] = ml
	return nil
}


func (l *logger) setBlacklist(method string) error {
	if _, ok := l.config.Blacklist[method]; ok {
		return fmt.Errorf("conflicting blacklist rules for method %v found", method)
	}
	if _, ok := l.config.Methods[method]; ok {
		return fmt.Errorf("conflicting method rules for method %v found", method)
	}
	if l.config.Blacklist == nil {
		l.config.Blacklist = make(map[string]struct{})
	}
	l.config.Blacklist[method] = struct{}{}
	return nil
}







func (l *logger) GetMethodLogger(methodName string) MethodLogger {
	s, m, err := grpcutil.ParseMethod(methodName)
	if err != nil {
		grpclogLogger.Infof("binarylogging: failed to parse %q: %v", methodName, err)
		return nil
	}
	if ml, ok := l.config.Methods[s+"/"+m]; ok {
		return NewTruncatingMethodLogger(ml.Header, ml.Message)
	}
	if _, ok := l.config.Blacklist[s+"/"+m]; ok {
		return nil
	}
	if ml, ok := l.config.Services[s]; ok {
		return NewTruncatingMethodLogger(ml.Header, ml.Message)
	}
	if l.config.All == nil {
		return nil
	}
	return NewTruncatingMethodLogger(l.config.All.Header, l.config.All.Message)
}
