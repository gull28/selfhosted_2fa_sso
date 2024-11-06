package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Port        string
		Environment string
	}
	Database struct {
		URL            string
		MaxConnections int
	}
	JWT struct {
		Secret     string
		Expiration int
	}
	Logging struct {
		Level string
		File  string
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	log.Printf("Loaded configuration: %+v", config)
	return &config, nil
}
