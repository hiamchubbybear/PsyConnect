



package stdr

import (
	"log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
)


var globalVerbosity int






func SetVerbosity(v int) int {
	old := globalVerbosity
	globalVerbosity = v
	return old
}






func New(std StdLogger) logr.Logger {
	return NewWithOptions(std, Options{})
}



func NewWithOptions(std StdLogger, opts Options) logr.Logger {
	if std == nil {
		
		std = log.New(os.Stderr, "", log.LstdFlags)
	}

	if opts.Depth < 0 {
		opts.Depth = 0
	}

	fopts := funcr.Options{
		LogCaller: funcr.MessageClass(opts.LogCaller),
	}

	sl := &logger{
		Formatter: funcr.NewFormatter(fopts),
		std:       std,
	}

	
	sl.Formatter.AddCallDepth(1 + opts.Depth)

	return logr.New(sl)
}


type Options struct {
	
	
	
	
	Depth int

	
	
	LogCaller MessageClass

	
}


type MessageClass int

const (
	
	None MessageClass = iota
	
	All
	
	Info
	
	Error
)



type StdLogger interface {
	
	Output(calldepth int, logline string) error
}

type logger struct {
	funcr.Formatter
	std StdLogger
}

var _ logr.LogSink = &logger{}
var _ logr.CallDepthLogSink = &logger{}

func (l logger) Enabled(level int) bool {
	return globalVerbosity >= level
}

func (l logger) Info(level int, msg string, kvList ...interface{}) {
	prefix, args := l.FormatInfo(level, msg, kvList)
	if prefix != "" {
		args = prefix + ": " + args
	}
	_ = l.std.Output(l.Formatter.GetDepth()+1, args)
}

func (l logger) Error(err error, msg string, kvList ...interface{}) {
	prefix, args := l.FormatError(err, msg, kvList)
	if prefix != "" {
		args = prefix + ": " + args
	}
	_ = l.std.Output(l.Formatter.GetDepth()+1, args)
}

func (l logger) WithName(name string) logr.LogSink {
	l.Formatter.AddName(name)
	return &l
}

func (l logger) WithValues(kvList ...interface{}) logr.LogSink {
	l.Formatter.AddValues(kvList)
	return &l
}

func (l logger) WithCallDepth(depth int) logr.LogSink {
	l.Formatter.AddCallDepth(depth)
	return &l
}





type Underlier interface {
	GetUnderlying() StdLogger
}



func (l logger) GetUnderlying() StdLogger {
	return l.std
}
