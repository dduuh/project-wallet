package main

import (
	"context"
	"wallet-service/internal/repository"
	"wallet-service/internal/service"
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

	// init migrations
	migrats, err := postgresql.NewMigrations(cfg.PostgreSQL(), "file://migrations")
	if err != nil {
		logrus.Panicf("Migrations error: %v\n", err)
	}

	// Up migrations
	err = migrats.Up()
	if err != nil {
		logrus.Panicf("Migrations UP is failed: %v\n", err)
	}

	// init repository
	repo := repository.NewUsersRepository(psql)

	// init service
	service := service.NewService(repo)

	// init kafka consumer
	consumer := consumer.New(cfg.KafkaCfg.Brokers, cfg.KafkaCfg.Topic, cfg.KafkaCfg.GroupID)

	// start consuming
	if err := consumer.Consume(ctx); err != nil {
		logrus.Panicf("Consumer error: %v\n", err)
	}

}
