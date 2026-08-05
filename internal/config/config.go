package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// загрузка конфигурации
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Config file not found: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Переопределяем из переменных окружения
	if envSecret := os.Getenv("JWT_SECRET"); envSecret != "" {
		cfg.JWT.Secret = envSecret
	}
	if envDBPassword := os.Getenv("DB_PASSWORD"); envDBPassword != "" {
		cfg.DataBase.Password = envDBPassword
	}

	return &cfg, nil
}

// дефолтные значения полей конфигурации
func setDefaults() {
	viper.SetDefault("environment", "development")

	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.cors_origin", "http://localhost:3000")

	viper.SetDefault("jwt.expiration", "24h")
}
