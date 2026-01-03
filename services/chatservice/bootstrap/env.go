package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	Port              string `mapstructure:"SERVICE_PORT"`
	Addr              string `mapstructure:"SERVICE_HOST"`
	KafkaAddr         string `mapstructure:"KAFKA_ADDRESS"`
	KafkaTopic        string `mapstructure:"KAFKA_TOPIC"`
	KafkaPart         string `mapstructure:"KAFKA_PARTITION"`
	GrpcAdd           string `mapstructure:"GRPC_ADDRESS"`
	NotificationTopic string `mapstructure:"KAFKA_NOTIFICATION_TOPIC"`
	JwtSecret         string `mapstructure:"JWT_SIGNER_KEY"`
	DatabaseMongoUri  string `mapstructure:"MONGO_URI"`
	DatabaseUser      string `mapstructure:"DB_USER"`
	DatabasePassword  string `mapstructure:"DB_PASS"`
	DatabaseHost      string `mapstructure:"DB_HOST"`
	DatabasePort      string `mapstructure:"DB_PORT"`
	DatabaseName      string `mapstructure:"DB_NAME"`
}

func LoadEnv() *Env {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	v.BindEnv("SERVICE_PORT")
	v.BindEnv("SERVICE_HOST")
	v.BindEnv("KAFKA_ADDRESS")
	v.BindEnv("KAFKA_TOPIC")
	v.BindEnv("KAFKA_PARTITION")
	v.BindEnv("KAFKA_NOTIFICATION_TOPIC")
	v.BindEnv("MONGO_URI")
	v.BindEnv("DB_USER")
	v.BindEnv("DB_PASS")
	v.BindEnv("DB_HOST")
	v.BindEnv("DB_PORT")
	v.BindEnv("DB_NAME")
	v.BindEnv("GRPC_ADDRESS")
	v.BindEnv("JWT_SIGNER_KEY")

	if err := v.ReadInConfig(); err != nil {
		log.Printf("Note: No .env file found, using system environment variables")
	}
	var env Env
	if err := v.Unmarshal(&env); err != nil {
		log.Fatal("Failed to unmarshal environment variables:", err)
	}
	return &env
}
