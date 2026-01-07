// internal/config.config.go
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	DatabaseURL       string
	HTTPServerAddress string
}

// LoadConfig reads configuration from app.env and environment variables
func LoadConfig(path string) (Config, error) {
	var cfg Config

	viper.SetConfigName("app") // app.env → "app"
	viper.SetConfigType("env") // file type is .env
	viper.AddConfigPath(path)  // where to look for the file
	viper.AutomaticEnv()       // read from OS environment variables

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("failed to read config: %w", err)
	}

	// Map env values to struct
	cfg.DatabaseURL = viper.GetString("DATABASE_URL")
	cfg.HTTPServerAddress = viper.GetString("HTTP_SERVER_ADDRESS")

	return cfg, nil
}
