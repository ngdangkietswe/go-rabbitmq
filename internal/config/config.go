/**
 * Author : ngdangkietswe
 * Since  : 8/13/2025
 */

package config

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

type AppConfig struct {
	Env         string `mapstructure:"ENV"`
	HttpPort    string `mapstructure:"PORT"`
	RabbitMQUrl string `mapstructure:"RABBITMQ_URL"`
}

const (
	DevelopmentEnv = "development"
)

func NewAppConfig(path string) AppConfig {
	config, err := loadConfig(path)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return config
}

func (cfg *AppConfig) GetHttpPort() string {
	if cfg.HttpPort == "" {
		return "3000"
	}
	return cfg.HttpPort
}

func (cfg *AppConfig) GetRabbitMQUrl() string {
	if cfg.RabbitMQUrl == "" {
		return "amqp://admin:admin123@localhost:5672/"
	}
	return cfg.RabbitMQUrl
}

func loadConfig(path string) (config AppConfig, err error) {
	viper.AutomaticEnv()
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	var env string

	if os.Getenv("ENV") != "" {
		env = os.Getenv("ENV")
	} else {
		env = DevelopmentEnv
	}

	if env == DevelopmentEnv {
		log.Printf("Loading development configuration from %s", path)

		if err := viper.ReadInConfig(); err != nil {
			return config, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if env == DevelopmentEnv {
		fmt.Println("Current configuration settings on development environment:")
		settings := viper.AllSettings()
		jsonBytes, _ := json.MarshalIndent(settings, "", "  ")
		fmt.Println(string(jsonBytes))
	}

	return
}
