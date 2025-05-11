package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	ExchangePrivateKey string
	ETHHost            string
	ServerPort         string
}

// LoadConfig loads configuration from the given file path
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error reading config file: %w", err)
	}

	return &Config{
		ExchangePrivateKey: viper.GetString("ExchangePrivateKey"),
		ETHHost:            viper.GetString("ETHHost"),
		ServerPort:         viper.GetString("ServerPort"),
	}, nil
}
