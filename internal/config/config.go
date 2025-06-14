package configs

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		HTTPCfg       HTTPConfig
		PostgreSQLCfg PostgreSQLConfig
		KafkaCfg      KafkaConfig
	}

	HTTPConfig struct {
		Host string `envconfig:"HTTP_HOST" default:"localhost"`
		Port string `envconfig:"HTTP_PORT" default:"8080"`
	}

	PostgreSQLConfig struct {
		Host     string `envconfig:"POSTGRES_HOST" default:"localhost"`
		Port     string `envconfig:"POSTGRES_PORT" default:"5432"`
		User     string `envconfig:"POSTGRES_USER" default:"postgres"`
		Password string `envconfig:"POSTGRES_PASSWORD" default:"David3410"`
		DBName   string `envconfig:"POSTGRES_DBNAME" default:"postgres"`
		SSLMode  string `envconfig:"POSTGRES_SSLMODE" default:"disable"`
	}

	KafkaConfig struct {
		Brokers []string `envconfig:"KAFKA_BROKERS" default:"localhost:9092"`
		GroupID string   `envconfig:"KAFKA_GROUP_ID" default:"wallet_users"`
		Topic   string   `envconfig:"KAFKA_TOPIC" default:"users"`
	}
)

func Init() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) PostgreSQL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.PostgreSQLCfg.User,
		c.PostgreSQLCfg.Password,
		c.PostgreSQLCfg.Host,
		c.PostgreSQLCfg.Port,
		c.PostgreSQLCfg.DBName,
		c.PostgreSQLCfg.SSLMode)
}
