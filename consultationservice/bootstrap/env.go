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
}

func LoadEnv() *Env {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("../")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		log.Printf("Warning: No .env file found, using environment variables only")
	}
	var env Env
	if err := v.Unmarshal(&env); err != nil {
		log.Fatal("Failed to unmarshal environment variables:", err)
	}

	return &env
}
