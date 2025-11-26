package logging

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// LogLevel represents the severity of a log event
type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	LOG   LogLevel = "LOG"
	AUDIT LogLevel = "AUDIT"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
)

// LogEvent represents a structured log event
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
	StackTrace  string                 `json:"stackTrace,omitempty"`
	Method      string                 `json:"method,omitempty"`
	Path        string                 `json:"path,omitempty"`
	StatusCode  int                    `json:"statusCode,omitempty"`
	Duration    int64                  `json:"duration,omitempty"`
	IP          string                 `json:"ip,omitempty"`
	UserAgent   string                 `json:"userAgent,omitempty"`
	Environment string                 `json:"environment,omitempty"`
	Version     string                 `json:"version,omitempty"`
	Hostname    string                 `json:"hostname,omitempty"`
}

// KafkaLogger handles logging to Kafka
type KafkaLogger struct {
	producer    *kafka.Producer
	serviceName string
	topic       string
	environment string
	version     string
	hostname    string
}

// Config holds configuration for KafkaLogger
type Config struct {
	Brokers     string
	ServiceName string
	Topic       string
	Environment string
	Version     string
}

// NewKafkaLogger creates a new Kafka logger instance
func NewKafkaLogger(config Config) (*KafkaLogger, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.Brokers,
		"client.id":         config.ServiceName,
		"acks":              "1",
		"retries":           3,
		"max.in.flight.requests.per.connection": 5,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	// Get hostname
	hostname := "unknown"
	// You can add hostname detection here if needed

	logger := &KafkaLogger{
		producer:    p,
		serviceName: config.ServiceName,
		topic:       config.Topic,
		environment: config.Environment,
		version:     config.Version,
		hostname:    hostname,
	}

	// Start delivery report handler
	go logger.handleDeliveryReports()

	return logger, nil
}

// handleDeliveryReports handles Kafka delivery reports in background
func (l *KafkaLogger) handleDeliveryReports() {
	for e := range l.producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition.Error)
			}
		}
	}
}

// Log sends a log event to Kafka
func (l *KafkaLogger) Log(level LogLevel, message string, fields map[string]interface{}) {
	event := LogEvent{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Level:       string(level),
		Service:     l.serviceName,
		Message:     message,
		Environment: l.environment,
		Version:     l.version,
		Hostname:    l.hostname,
		Metadata:    make(map[string]interface{}),
	}

	// Extract known fields
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
		if stackTrace, ok := fields["stackTrace"].(string); ok {
			event.StackTrace = stackTrace
			delete(fields, "stackTrace")
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
		if userAgent, ok := fields["userAgent"].(string); ok {
			event.UserAgent = userAgent
			delete(fields, "userAgent")
		}

		// Remaining fields go to metadata
		for k, v := range fields {
			event.Metadata[k] = v
		}
	}

	// Marshal to JSON
	data, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("Failed to marshal log event: %v\n", err)
		return
	}

	// Send to Kafka (non-blocking)
	err = l.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &l.topic,
			Partition: kafka.PartitionAny,
		},
		Value: data,
	}, nil)

	if err != nil {
		fmt.Printf("Failed to produce message: %v\n", err)
	}
}

// Debug logs a debug message
func (l *KafkaLogger) Debug(message string, fields map[string]interface{}) {
	l.Log(DEBUG, message, fields)
}

// Info logs an info message
func (l *KafkaLogger) Info(message string, fields map[string]interface{}) {
	l.Log(INFO, message, fields)
}

// Warn logs a warning message
func (l *KafkaLogger) Warn(message string, fields map[string]interface{}) {
	l.Log(WARN, message, fields)
}

// Error logs an error message
func (l *KafkaLogger) Error(message string, fields map[string]interface{}) {
	l.Log(ERROR, message, fields)
}

// Fatal logs a fatal message
func (l *KafkaLogger) Fatal(message string, fields map[string]interface{}) {
	l.Log(FATAL, message, fields)
}

// Audit logs an audit message
func (l *KafkaLogger) Audit(message string, fields map[string]interface{}) {
	l.Log(AUDIT, message, fields)
}

// LogHTTPRequest logs an HTTP request
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
	l.Log(level, message, fields)
}

// LogAction logs a user action
func (l *KafkaLogger) LogAction(action, userID, message string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["action"] = action
	fields["userId"] = userID

	l.Log(AUDIT, message, fields)
}

// Close flushes and closes the Kafka producer
func (l *KafkaLogger) Close() {
	// Flush remaining messages (wait up to 5 seconds)
	l.producer.Flush(5000)
	l.producer.Close()
}
