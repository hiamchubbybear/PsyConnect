package settings

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the logging service
type Config struct {
	Kafka      KafkaConfig
	Loki       LokiConfig
	FileLogger FileLoggerConfig
	Server     ServerConfig
}

type KafkaConfig struct {
	BootstrapServers string
	GroupID          string
	Topics           []string
	AutoOffsetReset  string
}

type LokiConfig struct {
	Enabled   bool
	URL       string
	Username  string
	Password  string
	BatchSize int
	Timeout   int // seconds
}

type FileLoggerConfig struct {
	Enabled    bool
	AppLogPath string
	ErrLogPath string
	MaxSize    int // megabytes
	MaxBackups int
	MaxAge     int // days
	Compress   bool
}

type ServerConfig struct {
	Port int
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	return &Config{
		Kafka: KafkaConfig{
			BootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "127.0.0.1:19092"),
			GroupID:          getEnv("KAFKA_GROUP_ID", "logging-service"),
			Topics: []string{
				"logging-service",
				"identity-service",
				"profile-service",
				"friends-service",
				"consultation-service",
				"chat-service",
				"notification-service",
			},
			AutoOffsetReset: getEnv("KAFKA_AUTO_OFFSET_RESET", "earliest"),
		},
		Loki: LokiConfig{
			Enabled:   getEnvAsBool("LOKI_ENABLED", true),
			URL:       getEnv("LOKI_URL", "http://localhost:3100"),
			Username:  getEnv("LOKI_USERNAME", ""),
			Password:  getEnv("LOKI_PASSWORD", ""),
			BatchSize: getEnvAsInt("LOKI_BATCH_SIZE", 100),
			Timeout:   getEnvAsInt("LOKI_TIMEOUT", 10),
		},
		FileLogger: FileLoggerConfig{
			Enabled:    getEnvAsBool("FILE_LOGGER_ENABLED", true),
			AppLogPath: getEnv("APP_LOG_PATH", "./logs/app.log"),
			ErrLogPath: getEnv("ERROR_LOG_PATH", "./logs/error.log"),
			MaxSize:    getEnvAsInt("LOG_MAX_SIZE", 15),
			MaxBackups: getEnvAsInt("LOG_MAX_BACKUPS", 10),
			MaxAge:     getEnvAsInt("LOG_MAX_AGE", 120),
			Compress:   getEnvAsBool("LOG_COMPRESS", true),
		},
		Server: ServerConfig{
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func (c *Config) Validate() error {
	if c.Kafka.BootstrapServers == "" {
		return fmt.Errorf("KAFKA_BOOTSTRAP_SERVERS is required")
	}
	if len(c.Kafka.Topics) == 0 {
		return fmt.Errorf("at least one Kafka topic is required")
	}
	if c.Loki.Enabled && c.Loki.URL == "" {
		return fmt.Errorf("LOKI_URL is required when Loki is enabled")
	}
	return nil
}
