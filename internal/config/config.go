package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/config/")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("config file not found, using defaults and env: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	if env := os.Getenv("JWT_SECRET"); env != "" {
		cfg.JWT.Secret = env
	}
	if env := os.Getenv("DB_HOST"); env != "" {
		cfg.DataBase.Host = env
	}
	if env := os.Getenv("DB_PORT"); env != "" {
		fmt.Sscanf(env, "%d", &cfg.DataBase.Port)
	}
	if env := os.Getenv("DB_USER"); env != "" {
		cfg.DataBase.User = env
	}
	if env := os.Getenv("DB_PASSWORD"); env != "" {
		cfg.DataBase.Password = env
	}
	if env := os.Getenv("DB_NAME"); env != "" {
		cfg.DataBase.DBName = env
	}

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("environment", "development")

	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.cors_origin", "http://localhost:3000")

	viper.SetDefault("jwt.expiration", "24h")
}