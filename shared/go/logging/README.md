# 📚 PsyConnect Shared Logging Library

Shared Go logging library cho tất cả PsyConnect microservices.

## 🎯 Features

- ✅ Structured logging với JSON format
- ✅ Kafka integration cho centralized logging
- ✅ Support nhiều log levels (DEBUG, INFO, WARN, ERROR, FATAL, AUDIT)
- ✅ Helper methods cho HTTP requests và user actions
- ✅ Automatic metadata extraction
- ✅ Non-blocking Kafka producer
- ✅ Graceful shutdown

## 📦 Installation

```bash
# In your service directory
go get github.com/psyconnect/shared/logging
```

## 🚀 Quick Start

### **Basic Usage**

```go
package main

import (
    "github.com/psyconnect/shared/logging"
)

func main() {
    // Create logger
    logger, err := logging.NewKafkaLogger(logging.Config{
        Brokers:     "localhost:9092",
        ServiceName: "my-service",
        Topic:       "logging-service",
        Environment: "development",
        Version:     "1.0.0",
    })
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    // Log messages
    logger.Info("Service started", nil)

    logger.Error("Something went wrong", map[string]interface{}{
        "error": "connection timeout",
        "userId": "user-123",
    })
}
```

### **HTTP Request Logging**

```go
// Log HTTP requests automatically
logger.LogHTTPRequest(
    "GET",              // method
    "/api/users/123",   // path
    200,                // status code
    45,                 // duration in ms
    map[string]interface{}{
        "userId": "user-123",
        "ip": "192.168.1.1",
    },
)
```

### **User Action Logging**

```go
// Log user actions for audit trail
logger.LogAction(
    "CREATE_POST",      // action
    "user-123",         // userId
    "User created a new post",
    map[string]interface{}{
        "postId": "post-456",
        "title": "My First Post",
    },
)
```

## 📝 Log Levels

| Level | Usage | Example |
|-------|-------|---------|
| **DEBUG** | Development debugging | `logger.Debug("Variable value", map[string]interface{}{"value": x})` |
| **INFO** | General information | `logger.Info("Service started", nil)` |
| **WARN** | Warning messages | `logger.Warn("High memory usage", map[string]interface{}{"usage": "85%"})` |
| **ERROR** | Error messages | `logger.Error("Failed to connect", map[string]interface{}{"error": err.Error()})` |
| **FATAL** | Fatal errors | `logger.Fatal("Cannot start service", map[string]interface{}{"error": err.Error()})` |
| **AUDIT** | Audit trail | `logger.Audit("User logged in", map[string]interface{}{"userId": "123"})` |

## 🔧 Configuration

```go
config := logging.Config{
    Brokers:     "localhost:9092",      // Kafka brokers
    ServiceName: "consultation-service", // Your service name
    Topic:       "logging-service",      // Kafka topic
    Environment: "production",           // Environment (dev/staging/prod)
    Version:     "1.0.0",               // Service version
}
```

## 📊 Log Event Structure

```json
{
  "timestamp": "2025-11-26T02:44:00Z",
  "level": "INFO",
  "service": "consultation-service",
  "message": "Consultation created",
  "traceId": "abc123",
  "userId": "user-123",
  "action": "CREATE_CONSULTATION",
  "metadata": {
    "consultationId": "consult-456",
    "doctorId": "doctor-789"
  },
  "method": "POST",
  "path": "/api/consultations",
  "statusCode": 201,
  "duration": 45,
  "ip": "192.168.1.1",
  "userAgent": "Mozilla/5.0...",
  "environment": "production",
  "version": "1.0.0",
  "hostname": "server-01"
}
```

## 🎯 Best Practices

### **1. Always Close Logger**

```go
logger, _ := logging.NewKafkaLogger(config)
defer logger.Close() // Ensures messages are flushed
```

### **2. Use Appropriate Log Levels**

```go
// ❌ Don't
logger.Info("Error occurred", map[string]interface{}{"error": err.Error()})

// ✅ Do
logger.Error("Failed to process request", map[string]interface{}{"error": err.Error()})
```

### **3. Include Context**

```go
// ❌ Don't
logger.Error("Failed", nil)

// ✅ Do
logger.Error("Failed to create consultation", map[string]interface{}{
    "userId": userID,
    "error": err.Error(),
    "consultationId": consultationID,
})
```

### **4. Use Helper Methods**

```go
// ❌ Don't
logger.Info("HTTP request", map[string]interface{}{
    "method": "GET",
    "path": "/api/users",
    "statusCode": 200,
})

// ✅ Do
logger.LogHTTPRequest("GET", "/api/users", 200, 45, nil)
```

## 🔗 Integration Examples

### **Gin Framework**

```go
func LoggingMiddleware(logger *logging.KafkaLogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        c.Next()

        duration := time.Since(start).Milliseconds()

        logger.LogHTTPRequest(
            c.Request.Method,
            c.Request.URL.Path,
            c.Writer.Status(),
            duration,
            map[string]interface{}{
                "ip": c.ClientIP(),
                "userAgent": c.Request.UserAgent(),
                "userId": c.GetString("userId"), // if available
            },
        )
    }
}
```

### **Error Handling**

```go
func handleError(logger *logging.KafkaLogger, err error, context map[string]interface{}) {
    if context == nil {
        context = make(map[string]interface{})
    }
    context["error"] = err.Error()

    logger.Error("Operation failed", context)
}
```

## 📚 API Reference

### **Methods**

| Method | Description |
|--------|-------------|
| `NewKafkaLogger(config)` | Create new logger instance |
| `Debug(message, fields)` | Log debug message |
| `Info(message, fields)` | Log info message |
| `Warn(message, fields)` | Log warning message |
| `Error(message, fields)` | Log error message |
| `Fatal(message, fields)` | Log fatal message |
| `Audit(message, fields)` | Log audit message |
| `LogHTTPRequest(...)` | Log HTTP request |
| `LogAction(...)` | Log user action |
| `Close()` | Close logger and flush messages |

### **Fields**

Common fields you can include:

- `traceId` - Distributed tracing ID
- `userId` - User identifier
- `action` - Action being performed
- `error` - Error message
- `stackTrace` - Stack trace
- `method` - HTTP method
- `path` - HTTP path
- `statusCode` - HTTP status code
- `duration` - Request duration (ms)
- `ip` - Client IP address
- `userAgent` - User agent string

## 🧪 Testing

```go
// Mock logger for testing
type MockLogger struct {
    logs []string
}

func (m *MockLogger) Info(message string, fields map[string]interface{}) {
    m.logs = append(m.logs, message)
}
```

## 📖 Examples

See `/examples` directory for complete examples:
- Basic logging
- HTTP middleware
- Error handling
- User action tracking

## 🤝 Contributing

When adding new features:
1. Update this README
2. Add examples
3. Test with real Kafka
4. Update version in go.mod

## 📄 License

Internal use only - PsyConnect Platform

---

**Happy Logging! 📝**
