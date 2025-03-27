package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	Port string `mapstructure:"SERVICE_PORT"`
	Addr string `mapstructure:"SERVICE_HOST"`
}

func LoadEnv() *Env {
	v := viper.New()
	v.SetConfigFile(".env")
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
