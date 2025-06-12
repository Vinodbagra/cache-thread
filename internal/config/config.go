package config

import (
	"time"

	"github.com/Vinodbagra/cache-thread/internal/constants"
	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	Port        int    `mapstructure:"PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`
	Debug       bool   `mapstructure:"DEBUG"`

	// Cache Configuration
	CacheMaxSize int           `mapstructure:"CACHE_MAX_SIZE"`
	CacheTTL     time.Duration `mapstructure:"CACHE_TTL"`
}

func InitializeAppConfig() error {
	viper.SetConfigName(".env") // allow directly reading from .env file
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("internal/config")
	viper.AddConfigPath("/")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return constants.ErrLoadConfig
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		return constants.ErrParseConfig
	}

	// check required fields
	if AppConfig.Port == 0 || AppConfig.Environment == "" {
		return constants.ErrEmptyVar
	}

	// Set default cache values if not provided
	if AppConfig.CacheMaxSize == 0 {
		AppConfig.CacheMaxSize = 1000 // Default max size
	}
	if AppConfig.CacheTTL == 0 {
		AppConfig.CacheTTL = 30 * time.Minute // Default TTL
	}

	// Database validation (only if environment requires it)
	switch AppConfig.Environment {
	case constants.EnvironmentDevelopment:
		// Database is optional for development
	case constants.EnvironmentProduction:
		// Database is optional for production
	}

	return nil
}
