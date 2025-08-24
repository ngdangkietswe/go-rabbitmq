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

const (
	DevelopmentEnv = "development"
)

func NewAppConfig(path string) {
	loadConfig(path)
}

func loadConfig(path string) {
	viper.AutomaticEnv()

	var env string

	if os.Getenv("ENV") != "" {
		env = os.Getenv("ENV")
	} else {
		env = DevelopmentEnv
	}

	if env == DevelopmentEnv {
		viper.AddConfigPath(path)
		viper.SetConfigName("app")
		viper.SetConfigType("env")

		log.Printf("Loading development configuration from %s", path)

		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("failed to read config file: %w", err))
		}

		settings := viper.AllSettings()
		jsonBytes, _ := json.MarshalIndent(settings, "", "  ")

		fmt.Println(string(jsonBytes))
	}

	return
}

func Get(key string, defaultValue interface{}) interface{} {
	if viper.IsSet(key) {
		switch defaultValue.(type) {
		case int:
			return viper.GetInt(key)
		case string:
			return viper.GetString(key)
		case bool:
			return viper.GetBool(key)
		default:
			log.Printf("Unsupported type for key %s", key)
			return defaultValue
		}
	} else {
		log.Printf("Key %s not found", key)
		return defaultValue
	}
}

func GetInt(key string, defaultValue int) int {
	return Get(key, defaultValue).(int)
}

func GetString(key string, defaultValue string) string {
	return Get(key, defaultValue).(string)
}

func GetBool(key string, defaultValue bool) bool {
	return Get(key, defaultValue).(bool)
}
