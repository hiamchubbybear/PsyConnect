





package logger

import "strings"




const DiffToInfo = 1






type Level int

const (
	
	LevelOff Level = iota

	
	
	LevelInfo

	
	
	
	LevelDebug
)

const (
	levelLiteralOff       = "off"
	levelLiteralEmergency = "emergency"
	levelLiteralAlert     = "alert"
	levelLiteralCritical  = "critical"
	levelLiteralError     = "error"
	levelLiteralWarning   = "warning"
	levelLiteralNotice    = "notice"
	levelLiteralInfo      = "info"
	levelLiteralDebug     = "debug"
	levelLiteralTrace     = "trace"
)

var LevelLiteralMap = map[string]Level{
	levelLiteralOff:       LevelOff,
	levelLiteralEmergency: LevelInfo,
	levelLiteralAlert:     LevelInfo,
	levelLiteralCritical:  LevelInfo,
	levelLiteralError:     LevelInfo,
	levelLiteralWarning:   LevelInfo,
	levelLiteralNotice:    LevelInfo,
	levelLiteralInfo:      LevelInfo,
	levelLiteralDebug:     LevelDebug,
	levelLiteralTrace:     LevelDebug,
}




func ParseLevel(str string) Level {
	for literal, level := range LevelLiteralMap {
		if strings.EqualFold(literal, str) {
			return level
		}
	}

	return LevelOff
}
