package main

import (
	"context"
	"wallet-service/internal/repository"
	"wallet-service/internal/transport/kafka/consumer"

	"github.com/sirupsen/logrus"

	configs "wallet-service/internal/config"
	postgresql "wallet-service/internal/repository/psql"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
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
		if err.Error() == "no change" {
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
