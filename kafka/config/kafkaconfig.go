package config

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

type Config struct {
	Kafka struct {
		Broker            []string `mapstructure:"broker"`
		Topics            []string `mapstructure:"topics"`
		Partitions        int32    `mapstructure:"partitions"`
		ReplicationFactor int32    `mapstructure:"replication_factor"`
	}
	Server struct {
		Port int `mapstructure:"port"`
	}
}

var AppConfig Config

func LoadConfig(path string) Config {
	if path == "" {
		path = "../config"
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Can not read config: %v", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Can not parse config: %v", err)
	}

	log.Println("Load config success")
	return AppConfig
}
