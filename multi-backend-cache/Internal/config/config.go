package config

import (
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	IsTenantBased   bool     `mapstructure:"IsTenantBased"`
	NumberOfTenants string   `mapstructure:"NumberOfTenants"`
	TenantIDs       []string `mapstructure:"TenantIDs"`
	DefaultTTL      int      `mapstructure:"defaultTTL"`
	CacheSystems    []string `mapstructure:"CacheSystems"`
}

var AppConfig Config

func LoadConfig(configFile string) {
	absPath, err := filepath.Abs(configFile)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}
	log.Printf("Using config file: %s", absPath)

	viper.SetConfigFile(absPath)
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatalf("Failed to unmarshal config file: %v", err)
	}
}
