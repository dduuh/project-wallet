package main

import (
	"context"
	"wallet-service/internal/repository"
	"wallet-service/internal/transport/kafka/consumer"

	"github.com/sirupsen/logrus"

	configs "wallet-service/internal/config"
	postgresql "wallet-service/internal/repository/psql"
)

func main() {
	ctx := context.Background()

	// init configs
	cfg, err := configs.Init()
	if err != nil {
		logrus.Panicf("Config error: %v\n", err)
	}

	// init db
	psql, err := postgresql.New(cfg)
	if err != nil {
		logrus.Panicf("Postgres error: %v\n", err)
	}

	// up migrations
	if err := psql.Up(); err != nil {
		logrus.Panicf("Migrations error: %v\n", err)
	}

	// init repository
	repo := repository.NewUsersRepository(psql.Database())

	// init kafka consumer
	consumer := consumer.New(*cfg, repo)

	// start consuming
	if err := consumer.Consume(ctx); err != nil {
		logrus.Panicf("Consumer error: %v\n", err)
	}

}
