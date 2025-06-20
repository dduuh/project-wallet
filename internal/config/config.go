package configs

import (
	"fmt"
	"net"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		HTTP     HTTPConfig
		Postgres PostgreSQLConfig
		Kafka    KafkaConfig
	}

	HTTPConfig struct {
		Port         string        `envconfig:"HTTP_PORT" default:"8080"`
		ReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"10s"`
		WriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"10s"`
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
		Brokers []string `envconfig:"KAFKA_BROKERS" default:"localhost:9094"`
		GroupID string   `envconfig:"KAFKA_GROUP_ID" default:"wallet_users"`
		Topic   string   `envconfig:"KAFKA_TOPIC" default:"users"`
	}
)

func Init() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process all configs: %w", err)
	}

	return &cfg, nil
}

func (c *Config) PostgreSQL() string {
	hostPort := net.JoinHostPort(c.Postgres.Host, c.Postgres.Port)

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		c.Postgres.User,
		c.Postgres.Password,
		hostPort,
		c.Postgres.DBName,
		c.Postgres.SSLMode)
}
