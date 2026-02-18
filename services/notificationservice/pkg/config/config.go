package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// SMTP Configuration
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	// Kafka Configuration
	KafkaBrokers []string
	KafkaGroupID string

	// Database Configuration
	DatabaseDSN string

	// Firebase Configuration
	FirebaseCredentials string

	// Server Configuration
	ServerPort  string
	Environment string
}

func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Parse SMTP Port
	smtpPort, err := strconv.Atoi(getEnv("SMPT_PORT", "587"))
	if err != nil {
		log.Fatal("Invalid SMTP port:", err)
	}

	// Parse Kafka Brokers
	kafkaBrokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:19092"), ",")

	// Build MySQL DSN
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "notification_db")

	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	return &Config{
		// SMTP
		SMTPHost:     getEnv("SMPT_HOST", "smtp.gmail.com"),
		SMTPPort:     smtpPort,
		SMTPUsername: getEnv("SMPT_MAIL", ""),
		SMTPPassword: getEnv("SMPT_APP_PASS", ""),
		SMTPFrom:     getEnv("SMPT_MAIL", "noreply@psyconnect.dev"),

		// Kafka
		KafkaBrokers: kafkaBrokers,
		KafkaGroupID: getEnv("KAFKA_GROUP_ID", "notification-group"),

		// Database
		DatabaseDSN: getEnv("DATABASE_DSN", dsn),

		// Firebase
		FirebaseCredentials: getEnv("FIREBASE_CREDENTIALS", ""),

		// Server
		ServerPort:  getEnv("PORT", "8082"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
