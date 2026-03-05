package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)


type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
	AUDIT LogLevel = "AUDIT"
)


type LogEvent struct {
	Timestamp   string                 `json:"timestamp"`
	Level       string                 `json:"level"`
	Service     string                 `json:"service"`
	Message     string                 `json:"message"`
	TraceID     string                 `json:"traceId,omitempty"`
	UserID      string                 `json:"userId,omitempty"`
	Action      string                 `json:"action,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Method      string                 `json:"method,omitempty"`
	Path        string                 `json:"path,omitempty"`
	StatusCode  int                    `json:"statusCode,omitempty"`
	Duration    int64                  `json:"duration,omitempty"`
	IP          string                 `json:"ip,omitempty"`
	Environment string                 `json:"environment,omitempty"`
	Version     string                 `json:"version,omitempty"`
}


type KafkaLogger struct {
	writer      *kafka.Writer
	serviceName string
	environment string
	version     string
}


type Config struct {
	Brokers     []string
	Topic       string
	ServiceName string
	Environment string
	Version     string
}


func NewKafkaLogger(config Config) *KafkaLogger {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		Topic:        config.Topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        true, 
	}

	return &KafkaLogger{
		writer:      writer,
		serviceName: config.ServiceName,
		environment: config.Environment,
		version:     config.Version,
	}
}


func (l *KafkaLogger) log(level LogLevel, message string, fields map[string]interface{}) {
	event := LogEvent{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Level:       string(level),
		Service:     l.serviceName,
		Message:     message,
		Environment: l.environment,
		Version:     l.version,
		Metadata:    make(map[string]interface{}),
	}

	
	if fields != nil {
		if traceID, ok := fields["traceId"].(string); ok {
			event.TraceID = traceID
			delete(fields, "traceId")
		}
		if userID, ok := fields["userId"].(string); ok {
			event.UserID = userID
			delete(fields, "userId")
		}
		if action, ok := fields["action"].(string); ok {
			event.Action = action
			delete(fields, "action")
		}
		if err, ok := fields["error"].(string); ok {
			event.Error = err
			delete(fields, "error")
		}
		if method, ok := fields["method"].(string); ok {
			event.Method = method
			delete(fields, "method")
		}
		if path, ok := fields["path"].(string); ok {
			event.Path = path
			delete(fields, "path")
		}
		if statusCode, ok := fields["statusCode"].(int); ok {
			event.StatusCode = statusCode
			delete(fields, "statusCode")
		}
		if duration, ok := fields["duration"].(int64); ok {
			event.Duration = duration
			delete(fields, "duration")
		}
		if ip, ok := fields["ip"].(string); ok {
			event.IP = ip
			delete(fields, "ip")
		}

		
		for k, v := range fields {
			event.Metadata[k] = v
		}
	}

	
	data, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("Failed to marshal log event: %v\n", err)
		return
	}

	
	err = l.writer.WriteMessages(context.Background(), kafka.Message{
		Value: data,
	})

	if err != nil {
		
		
		fmt.Printf("⚠️ Kafka Log Error: %v (service continues)\n", err)
	}
}


func (l *KafkaLogger) Debug(message string, fields map[string]interface{}) {
	l.log(DEBUG, message, fields)
}


func (l *KafkaLogger) Info(message string, fields map[string]interface{}) {
	l.log(INFO, message, fields)
}


func (l *KafkaLogger) Warn(message string, fields map[string]interface{}) {
	l.log(WARN, message, fields)
}


func (l *KafkaLogger) Error(message string, fields map[string]interface{}) {
	l.log(ERROR, message, fields)
}


func (l *KafkaLogger) Fatal(message string, fields map[string]interface{}) {
	l.log(FATAL, message, fields)
	fmt.Printf("🔴 FATAL: %s\n", message)
}


func (l *KafkaLogger) Audit(message string, fields map[string]interface{}) {
	l.log(AUDIT, message, fields)
}


func (l *KafkaLogger) LogHTTPRequest(method, path string, statusCode int, duration int64, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["method"] = method
	fields["path"] = path
	fields["statusCode"] = statusCode
	fields["duration"] = duration

	level := INFO
	if statusCode >= 500 {
		level = ERROR
	} else if statusCode >= 400 {
		level = WARN
	}

	message := fmt.Sprintf("%s %s - %d (%dms)", method, path, statusCode, duration)
	l.log(level, message, fields)
}


func (l *KafkaLogger) LogAction(action, userID, message string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["action"] = action
	fields["userId"] = userID

	l.log(AUDIT, message, fields)
}


func (l *KafkaLogger) Close() error {
	return l.writer.Close()
}
