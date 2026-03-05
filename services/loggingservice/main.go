package main

import (
	"fmt"
	"os"

	"github.com/loggingservice/pkg/filelogger"
	"github.com/loggingservice/pkg/kafka"
	"github.com/loggingservice/pkg/loki"
	"github.com/loggingservice/pkg/models"
	"github.com/loggingservice/pkg/server"
	"github.com/loggingservice/pkg/settings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	
	consoleLogger := initConsoleLogger()
	defer consoleLogger.Sync()

	consoleLogger.Info("🚀 Starting PsyConnect Logging Service")

	
	config := settings.LoadConfig()
	if err := config.Validate(); err != nil {
		consoleLogger.Fatal("Invalid configuration", zap.Error(err))
	}

	consoleLogger.Info("Configuration loaded",
		zap.String("kafka_servers", config.Kafka.BootstrapServers),
		zap.Strings("kafka_topics", config.Kafka.Topics),
		zap.Bool("loki_enabled", config.Loki.Enabled),
		zap.String("loki_url", config.Loki.URL),
		zap.Bool("file_logger_enabled", config.FileLogger.Enabled),
	)

	
	var lokiClient *loki.Client
	if config.Loki.Enabled {
		lokiClient = loki.NewClient(config.Loki)
		lokiClient.Start()
		defer lokiClient.Stop()
		consoleLogger.Info("Loki client initialized")
	} else {
		consoleLogger.Warn("Loki client disabled")
	}

	fileLogger := filelogger.NewFileLogger(config.FileLogger)
	defer fileLogger.Close()
	if config.FileLogger.Enabled {
		consoleLogger.Info("File logger initialized",
			zap.String("app_log", config.FileLogger.AppLogPath),
			zap.String("error_log", config.FileLogger.ErrLogPath),
		)
	} else {
		consoleLogger.Warn("File logger disabled")
	}

	
	httpServer := server.NewServer(config.Server, consoleLogger)
	if err := httpServer.Start(); err != nil {
		consoleLogger.Fatal("Failed to start HTTP server", zap.Error(err))
	}
	defer httpServer.Stop()
	consoleLogger.Info("HTTP server started", zap.Int("port", config.Server.Port))

	
	kafkaConsumer, err := kafka.NewConsumer(config.Kafka, consoleLogger)
	if err != nil {
		consoleLogger.Fatal("Failed to create Kafka consumer", zap.Error(err))
	}
	defer kafkaConsumer.Close()

	
	kafkaConsumer.AddHandler(func(event *models.LogEvent) error {
		httpServer.IncrementMessages()

		
		consoleLogger.Debug("Processing log event",
			zap.String("service", event.Service),
			zap.String("level", event.Level),
			zap.String("message", event.Message),
		)

		
		if err := fileLogger.Write(event); err != nil {
			consoleLogger.Error("Failed to write to file", zap.Error(err))
			httpServer.IncrementErrors()
			return err
		}

		
		if lokiClient != nil {
			if err := lokiClient.Push(event); err != nil {
				consoleLogger.Error("Failed to push to Loki", zap.Error(err))
				httpServer.IncrementErrors()
				return err
			}
		}

		return nil
	})

	consoleLogger.Info("Kafka consumer initialized and handlers registered")
	consoleLogger.Info(" Logging service is ready to process events")
	consoleLogger.Info(" Health check available at http://localhost:" + fmt.Sprintf("%d", config.Server.Port) + "/health")

	
	if err := kafkaConsumer.Start(); err != nil {
		consoleLogger.Fatal("Kafka consumer error", zap.Error(err))
	}

	consoleLogger.Info("👋 Logging service shutting down gracefully")
}


func initConsoleLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	return zap.New(core, zap.AddCaller())
}
