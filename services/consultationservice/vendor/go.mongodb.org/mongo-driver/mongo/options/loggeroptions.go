





package options

import (
	"go.mongodb.org/mongo-driver/internal/logger"
)


type LogLevel int

const (
	
	
	LogLevelInfo LogLevel = LogLevel(logger.LevelInfo)

	
	
	
	LogLevelDebug LogLevel = LogLevel(logger.LevelDebug)
)



type LogComponent int

const (
	
	LogComponentAll LogComponent = LogComponent(logger.ComponentAll)

	
	LogComponentCommand LogComponent = LogComponent(logger.ComponentCommand)

	
	LogComponentTopology LogComponent = LogComponent(logger.ComponentTopology)

	
	LogComponentServerSelection LogComponent = LogComponent(logger.ComponentServerSelection)

	
	LogComponentConnection LogComponent = LogComponent(logger.ComponentConnection)
)



type LogSink interface {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Info(level int, message string, keysAndValues ...interface{})

	
	Error(err error, message string, keysAndValues ...interface{})
}


type LoggerOptions struct {
	
	
	
	ComponentLevels map[LogComponent]LogLevel

	
	
	Sink LogSink

	
	
	
	MaxDocumentLength uint
}


func Logger() *LoggerOptions {
	return &LoggerOptions{
		ComponentLevels: map[LogComponent]LogLevel{},
	}
}


func (opts *LoggerOptions) SetComponentLevel(component LogComponent, level LogLevel) *LoggerOptions {
	opts.ComponentLevels[component] = level

	return opts
}


func (opts *LoggerOptions) SetMaxDocumentLength(maxDocumentLength uint) *LoggerOptions {
	opts.MaxDocumentLength = maxDocumentLength

	return opts
}


func (opts *LoggerOptions) SetSink(sink LogSink) *LoggerOptions {
	opts.Sink = sink

	return opts
}
