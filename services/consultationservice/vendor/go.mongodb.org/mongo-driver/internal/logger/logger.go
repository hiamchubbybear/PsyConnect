







package logger

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)



const DefaultMaxDocumentLength = 1000




const TruncationSuffix = "..."

const logSinkPathEnvVar = "MONGODB_LOG_PATH"
const maxDocumentLengthEnvVar = "MONGODB_LOG_MAX_DOCUMENT_LENGTH"



type LogSink interface {
	
	
	Info(level int, msg string, keysAndValues ...interface{})

	
	Error(err error, msg string, keysAndValues ...interface{})
}


type Logger struct {
	ComponentLevels   map[Component]Level 
	Sink              LogSink             
	MaxDocumentLength uint                
	logFile           *os.File            
}





func New(sink LogSink, maxDocLen uint, compLevels map[Component]Level) (*Logger, error) {
	logger := &Logger{
		ComponentLevels:   selectComponentLevels(compLevels),
		MaxDocumentLength: selectMaxDocumentLength(maxDocLen),
	}

	sink, logFile, err := selectLogSink(sink)
	if err != nil {
		return nil, err
	}

	logger.Sink = sink
	logger.logFile = logFile

	return logger, nil
}


func (logger *Logger) Close() error {
	if logger.logFile != nil {
		return logger.logFile.Close()
	}

	return nil
}










func (logger *Logger) LevelComponentEnabled(level Level, component Component) bool {
	if level == LevelOff {
		return false
	}

	if logger.ComponentLevels == nil {
		return false
	}

	return logger.ComponentLevels[component] >= level ||
		logger.ComponentLevels[ComponentAll] >= level
}











func (logger *Logger) Print(level Level, component Component, msg string, keysAndValues ...interface{}) {
	
	
	if !logger.LevelComponentEnabled(level, component) {
		return
	}

	
	if logger.Sink == nil {
		return
	}

	logger.Sink.Info(int(level)-DiffToInfo, msg, keysAndValues...)
}




func (logger *Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	if logger.Sink == nil {
		return
	}

	logger.Sink.Error(err, msg, keysAndValues...)
}







func selectMaxDocumentLength(maxDocLen uint) uint {
	if maxDocLen != 0 {
		return maxDocLen
	}

	maxDocLenEnv := os.Getenv(maxDocumentLengthEnvVar)
	if maxDocLenEnv != "" {
		maxDocLenEnvInt, err := strconv.ParseUint(maxDocLenEnv, 10, 32)
		if err == nil {
			return uint(maxDocLenEnvInt)
		}
	}

	return DefaultMaxDocumentLength
}

const (
	logSinkPathStdout = "stdout"
	logSinkPathStderr = "stderr"
)




func selectLogSink(sink LogSink) (LogSink, *os.File, error) {
	if sink != nil {
		return sink, nil, nil
	}

	path := os.Getenv(logSinkPathEnvVar)
	lowerPath := strings.ToLower(path)

	if lowerPath == string(logSinkPathStderr) {
		return NewIOSink(os.Stderr), nil, nil
	}

	if lowerPath == string(logSinkPathStdout) {
		return NewIOSink(os.Stdout), nil, nil
	}

	if path != "" {
		logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to open log file: %w", err)
		}

		return NewIOSink(logFile), logFile, nil
	}

	return NewIOSink(os.Stderr), nil, nil
}




func selectComponentLevels(componentLevels map[Component]Level) map[Component]Level {
	selected := make(map[Component]Level)

	
	var globalEnvLevel *Level
	if all := os.Getenv(mongoDBLogAllEnvVar); all != "" {
		level := ParseLevel(all)
		globalEnvLevel = &level
	}

	for envVar, component := range componentEnvVarMap {
		
		if _, ok := componentLevels[component]; ok {
			selected[component] = componentLevels[component]

			continue
		}

		
		
		
		if globalEnvLevel != nil {
			selected[component] = *globalEnvLevel

			continue
		}

		
		
		selected[component] = ParseLevel(os.Getenv(envVar))
	}

	return selected
}




func truncate(str string, width uint) string {
	if width == 0 {
		return ""
	}

	if len(str) <= int(width) {
		return str
	}

	
	newStr := str[:width]

	
	
	if newStr[len(newStr)-1]&0xC0 == 0xC0 {
		return newStr[:len(newStr)-1] + TruncationSuffix
	}

	
	
	if newStr[len(newStr)-1]&0xC0 == 0x80 {
		for i := len(newStr) - 1; i >= 0; i-- {
			if newStr[i]&0xC0 == 0xC0 {
				return newStr[:i] + TruncationSuffix
			}
		}
	}

	return newStr + TruncationSuffix
}



func FormatMessage(msg string, width uint) string {
	if len(msg) == 0 {
		return "{}"
	}

	return truncate(msg, width)
}
