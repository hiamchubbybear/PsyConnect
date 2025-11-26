package filelogger

import (
	"github.com/loggingservice/pkg/models"
	"github.com/loggingservice/pkg/settings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// FileLogger handles writing logs to files
type FileLogger struct {
	config    settings.FileLoggerConfig
	appLogger *zap.Logger
	errLogger *zap.Logger
}

// NewFileLogger creates a new file logger
func NewFileLogger(config settings.FileLoggerConfig) *FileLogger {
	if !config.Enabled {
		return &FileLogger{
			config:    config,
			appLogger: zap.NewNop(),
			errLogger: zap.NewNop(),
		}
	}

	// Encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		MessageKey:     "msg",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// App log writer (all logs)
	appLogWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.AppLogPath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
		LocalTime:  true,
	})

	// Error log writer (errors only)
	errLogWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.ErrLogPath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
		LocalTime:  true,
	})

	// Create app logger (all levels)
	appCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		appLogWriter,
		zapcore.DebugLevel,
	)
	appLogger := zap.New(appCore)

	// Create error logger (errors only)
	errCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		errLogWriter,
		zapcore.ErrorLevel,
	)
	errLogger := zap.New(errCore)

	return &FileLogger{
		config:    config,
		appLogger: appLogger,
		errLogger: errLogger,
	}
}

// Write writes a log event to file
func (f *FileLogger) Write(event *models.LogEvent) error {
	if !f.config.Enabled {
		return nil
	}

	// Build zap fields from event
	fields := []zap.Field{
		zap.String("service", event.Service),
		zap.String("timestamp", event.Timestamp),
	}

	if event.TraceID != "" {
		fields = append(fields, zap.String("traceId", event.TraceID))
	}
	if event.UserID != "" {
		fields = append(fields, zap.String("userId", event.UserID))
	}
	if event.Action != "" {
		fields = append(fields, zap.String("action", event.Action))
	}
	if event.Metadata != nil {
		fields = append(fields, zap.Any("metadata", event.Metadata))
	}
	if event.Error != "" {
		fields = append(fields, zap.String("error", event.Error))
	}
	if event.StackTrace != "" {
		fields = append(fields, zap.String("stackTrace", event.StackTrace))
	}
	if event.Method != "" {
		fields = append(fields, zap.String("method", event.Method))
	}
	if event.Path != "" {
		fields = append(fields, zap.String("path", event.Path))
	}
	if event.StatusCode > 0 {
		fields = append(fields, zap.Int("statusCode", event.StatusCode))
	}
	if event.Duration > 0 {
		fields = append(fields, zap.Int64("duration", event.Duration))
	}

	// Write to appropriate logger based on level
	switch models.LogLevel(event.Level) {
	case models.LevelDebug:
		f.appLogger.Debug(event.Message, fields...)
	case models.LevelInfo, models.LevelLog, models.LevelAudit:
		f.appLogger.Info(event.Message, fields...)
	case models.LevelWarn:
		f.appLogger.Warn(event.Message, fields...)
		f.errLogger.Warn(event.Message, fields...)
	case models.LevelError, models.LevelFatal:
		f.appLogger.Error(event.Message, fields...)
		f.errLogger.Error(event.Message, fields...)
	default:
		f.appLogger.Info(event.Message, fields...)
	}

	return nil
}

// Sync flushes any buffered log entries
func (f *FileLogger) Sync() error {
	if !f.config.Enabled {
		return nil
	}
	_ = f.appLogger.Sync()
	_ = f.errLogger.Sync()
	return nil
}

// Close closes the file logger
func (f *FileLogger) Close() error {
	return f.Sync()
}
