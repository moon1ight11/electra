package config

import "time"

type Config struct {
	Environment string         `mapstructure:"environment"`
	DataBase    DataBaseConfig `mapstructure:"database"`
	JWT         JWTConfig      `mapstructure:"jwt"`
	Logger      LoggerConfig   `mapstructure:"logger"`
	Server      ServerConfig   `mapstructure:"server"`
}

// конфигурация бд
type DataBaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"db_name"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_life_time"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	MigrationsDir   string        `mapstructure:"migrations_dir"`
}

// конфигурация jwt
type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
}

// конфигурация логгера
type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"filepath"`
	MaxSize  int    `mapstructure:"maxsize"`
}

// конфигурация сервера
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}
