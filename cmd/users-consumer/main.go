package main

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	configs "wallet-service/internal/config"
	"wallet-service/internal/repository"
	postgresql "wallet-service/internal/repository/psql"
	"wallet-service/internal/transport/kafka/consumer"
)

func main() {
	ctx := context.Background()

	cfg, err := configs.Init()
	if err != nil {
		logrus.Panicf("Config error: %v\n", err)
	}

	psql, err := postgresql.New(cfg)
	if err != nil {
		logrus.Panicf("Postgres error: %v\n", err)
	}

	if err := psql.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logrus.Info("No migrations to apply.")
		} else {
			logrus.Panicf("Migrations error: %v\n", err)
		}
	}

	repo := repository.NewUsersRepository(psql.Database())

	consumer := consumer.New(cfg, repo)

	if err := consumer.Consume(ctx); err != nil {
		logrus.Panicf("Consumer error: %v\n", err)
	}
}
