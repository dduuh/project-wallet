package configs

import (
	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		HTTPCfg       HTTPConfig
		PostgreSQLCfg PostgreSQLConfig
	}

	HTTPConfig struct {
		Host string `envconfig:"HTTP_HOST" default:"localhost"`
		Port string `envconfig:"HTTP_PORT" default:"8080"`
	}

	PostgreSQLConfig struct {
		Host     string `envconfig:"POSTGRES_HOST" default:"localhost"`
		Port     string `envconfig:"POSTGRES_PORT" default:"5432"`
		User     string `envconfig:"POSTGRES_USER" default:"postgres"`
		Password string `envconfig:"POSTGRES_PASSWORD"`
		DBName   string `envconfig:"POSTGRES_DBNAME"`
		SSLMode  string `envconfig:"POSTGRES_SSLMODE" default:"disable"`
	}
)

func Init() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
