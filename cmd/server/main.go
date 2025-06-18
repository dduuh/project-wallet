package main

import (
	configs "wallet-service/internal/config"
	"wallet-service/internal/repository"
	postgresql "wallet-service/internal/repository/psql"
	"wallet-service/internal/service"
	"wallet-service/internal/transport/rest"

	"github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
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
	handlers := rest.New(services, repo)

	logrus.Infof("HTTP Server started on port %s\n", cfg.HTTPCfg.Port)

	if err := handlers.Run(cfg, handlers.InitRoutes()); err != nil {
		logrus.Panicf("HTTP Server error: %v\n", err)
	}
}
