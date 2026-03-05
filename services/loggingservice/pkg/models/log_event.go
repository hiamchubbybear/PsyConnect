package models

import (
	"encoding/json"
	"time"
)


type LogEvent struct {
	
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Service   string `json:"service"`
	Message   string `json:"message"`

	
	TraceID string `json:"traceId,omitempty"`
	SpanID  string `json:"spanId,omitempty"`

	
	UserID   string `json:"userId,omitempty"`
	Username string `json:"username,omitempty"`

	
	Action     string                 `json:"action,omitempty"`
	TargetID   string                 `json:"targetId,omitempty"`
	TargetType string                 `json:"targetType,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`

	
	Error      string `json:"error,omitempty"`
	StackTrace string `json:"stackTrace,omitempty"`

	
	Method     string `json:"method,omitempty"`
	Path       string `json:"path,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
	Duration   int64  `json:"duration,omitempty"` 
	IP         string `json:"ip,omitempty"`
	UserAgent  string `json:"userAgent,omitempty"`

	
	Environment string `json:"environment,omitempty"`
	Version     string `json:"version,omitempty"`
	Hostname    string `json:"hostname,omitempty"`
}


type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelLog   LogLevel = "LOG"
	LevelAudit LogLevel = "AUDIT"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
	LevelFatal LogLevel = "FATAL"
)


func ParseLogEvent(data []byte) (*LogEvent, error) {
	var event LogEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	
	if event.Timestamp == "" {
		event.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
	}

	
	if event.Level == "" {
		event.Level = string(LevelInfo)
	}

	return &event, nil
}


func (e *LogEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}


func (e *LogEvent) ToJSONString() (string, error) {
	bytes, err := e.ToJSON()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}


func (e *LogEvent) IsError() bool {
	level := LogLevel(e.Level)
	return level == LevelError || level == LevelFatal
}


func (e *LogEvent) GetLabels() map[string]string {
	labels := map[string]string{
		"level":   e.Level,
		"service": e.Service,
	}

	if e.Environment != "" {
		labels["environment"] = e.Environment
	}

	if e.UserID != "" {
		labels["user_id"] = e.UserID
	}

	if e.Action != "" {
		labels["action"] = e.Action
	}

	if e.TraceID != "" {
		labels["trace_id"] = e.TraceID
	}

	return labels
}


func (e *LogEvent) GetLogLine() string {
	jsonStr, err := e.ToJSONString()
	if err != nil {
		return e.Message
	}
	return jsonStr
}
