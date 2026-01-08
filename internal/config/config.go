// internal/config.config.go
package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	DatabaseURL          string        `mapstructure:"DATABASE_URL" validate:"required"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS" validate:"required"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY" validate:"required,len=32"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION" validate:"required"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION" validate:"required"`
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
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return cfg, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}
