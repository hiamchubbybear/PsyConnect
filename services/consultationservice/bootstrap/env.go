package bootstrap

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Env struct {
	Port                      string `mapstructure:"SERVICE_PORT"`
	Addr                      string `mapstructure:"SERVICE_HOST"`
	KafkaAddr                 string `mapstructure:"KAFKA_ADDRESS"`
	KafkaTopic                string `mapstructure:"KAFKA_TOPIC"`
	KafkaPart                 string `mapstructure:"KAFKA_PARTITION"`
	GrpcAdd                   string `mapstructure:"PROFILE_GRPC_ADDR"`
	NotificationTopic         string `mapstructure:"KAFKA_NOTIFICATION_TOPIC"`
	ConsultationRedisAddress  string `mapstructure:"CONSULTATION_REDIS_ADDRESS"`
	ConsultationProtocol      int    `mapstructure:"CONSULTATION_PROTOCOL"`
	ConsultationRedisPassword string `mapstructure:"CONSULTATION_REDIS_PASSWORD"`
	ConsultationDB            int    `mapstructure:"CONSULTATION_DB"`
	MongoURI                  string `mapstructure:"MONGO_URI"`
	DBName                    string `mapstructure:"DB_NAME"`
	Environment               string `mapstructure:"ENVIRONMENT"`
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
	v.BindEnv("GRPC_ADDRESS", "PROFILE_GRPC_ADDR")
	v.BindEnv("PROFILE_GRPC_ADDR")
	v.BindEnv("KAFKA_NOTIFICATION_TOPIC")
	v.BindEnv("CONSULTATION_REDIS_ADDRESS")
	v.BindEnv("CONSULTATION_PROTOCOL")
	v.BindEnv("CONSULTATION_REDIS_PASSWORD")
	v.BindEnv("CONSULTATION_DB")
	v.BindEnv("MONGO_URI")
	v.BindEnv("DB_NAME")
	v.BindEnv("ENVIRONMENT")

	if err := v.ReadInConfig(); err != nil {
		log.Printf("Note: No .env file found, using system environment variables")
	} else {
		log.Printf("Using config file: %s", v.ConfigFileUsed())
	}

	
	keysToSync := []string{"MONGO_URI", "DB_NAME", "ENVIRONMENT", "KAFKA_BROKERS", "PROFILE_GRPC_ADDR"}
	for _, k := range keysToSync {
		val := v.GetString(k)
		if val != "" && os.Getenv(k) == "" {
			os.Setenv(k, val)
		}
	}

	var env Env
	if err := v.Unmarshal(&env); err != nil {
		log.Fatal("Failed to unmarshal environment variables:", err)
	}

	return &env
}
