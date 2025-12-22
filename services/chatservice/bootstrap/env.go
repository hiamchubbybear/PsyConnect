package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	Port              string `mapstructure:"CHAT_SERVICE_PORT"`
	Addr              string `mapstructure:"SERVICE_HOST"`
	KafkaAddr         string `mapstructure:"KAFKA_ADDRESS"`
	KafkaTopic        string `mapstructure:"KAFKA_TOPIC"`
	KafkaPart         string `mapstructure:"KAFKA_PARTITION"`
	GrpcAdd           string `mapstructure:"GRPC_ADDRESS"`
	NotificationTopic string `mapstructure:"KAFKA_NOTIFICATION_TOPIC"`
	JwtSecret         string `mapstructure:"JWT_SIGNER_KEY"`
}

func LoadEnv() *Env {
	v := viper.New()
	v.SetConfigName(".env")
	v.AutomaticEnv()
	v.BindEnv("CHAT_SERVICE_PORT")
	v.BindEnv("SERVICE_HOST")
	v.BindEnv("KAFKA_ADDRESS")
	v.BindEnv("KAFKA_TOPIC")
	v.BindEnv("KAFKA_PARTITION")
	v.BindEnv("KAFKA_NOTIFICATION_TOPIC")
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
