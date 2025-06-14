package main

import (
	"context"
	"encoding/json"
	"wallet-service/internal/domain"
	"wallet-service/internal/repository"
	"wallet-service/internal/transport/kafka/consumer"
	"wallet-service/internal/transport/kafka/producer"

	"github.com/sirupsen/logrus"

	configs "wallet-service/internal/config"
	postgresql "wallet-service/internal/repository/psql"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
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
		if err.Error() == "no change" {
			logrus.Info("No migrations to apply.")
		} else {
			logrus.Panicf("Migrations error: %v\n", err)
		}
	}

	// init repository
	repo := repository.NewUsersRepository(psql.Database())

	// init producer
	producer := producer.New(cfg)

	user := &domain.User{
		Id:        1,
		BlockedAt: nil,
		DeletedAt: nil,
	}

	userData, err := json.Marshal(user)
	if err != nil {
		logrus.Panicf("JSON Marshal error: %v\n", err)
	}

	if err := producer.Produce(ctx, userData); err != nil {
		logrus.Panicf("Producer error: %v\n", err)
	}

	// init kafka consumer
	consumer := consumer.New(cfg, repo)

	// start consuming
	if err := consumer.Consume(ctx); err != nil {
		logrus.Panicf("Consumer error: %v\n", err)
	}
}
