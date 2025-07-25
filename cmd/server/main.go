package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	configs "wallet-service/internal/config"
	"wallet-service/internal/repository"
	postgresql "wallet-service/internal/repository/psql"
	"wallet-service/internal/service"
	"wallet-service/internal/transport/rest"
)

func main() {
	ctx := context.Background()

	cfg, err := configs.Init()
	if err != nil {
		logrus.Panicf("Configs error: %v\n", err)
	}

	psql, err := postgresql.New(cfg)
	if err != nil {
		logrus.Panicf("PostgreSQL error: %v\n", err)
	}

	if err := psql.Up(); err != nil {
		if err.Error() == "no change" {
			logrus.Info("no migrations to apply")
		} else {
			logrus.Panicf("Migrations error: %v\n", err)
		}
	}

	repo := repository.NewUsersRepository(psql.Database())
	walletRepo := repository.NewWalletRepository(psql.Database())
	services := service.New(repo, walletRepo)
	server := rest.New(services, repo)

	logrus.Infof("HTTP Server started on port %s\n", cfg.HTTP.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Run(ctx, cfg, server.InitRoutes()); err != nil {
			logrus.Panicf("HTTP Server error: %v\n", err)
		}
	}()

	<-quit

	if err := psql.Close(); err != nil {
		logrus.Panicf("PostgreSQL Close error: %v\n", err)
	}
}
