package main

import (
	"wallet-service/internal/configs"
	"wallet-service/internal/repository"
	"wallet-service/internal/rest/handlers"
	"wallet-service/internal/server"
	"wallet-service/internal/service"
	"wallet-service/pkg/psql"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := configs.Init()
	if err != nil {
		logrus.Fatalf("Config error: %v\n", err)
	}

	psql, err := psql.NewPostgreSQL()
	if err != nil {
		logrus.Fatalf("Postgres error: %v\n", err)
	}

	repo := repository.NewRepository(psql)
	service := service.NewService(repo)
	handler := handlers.NewHandler(service)
	srv := server.NewServer()

	if err := srv.Run(cfg.HTTPCfg.Port, handler.InitRoutes()); err != nil {
		logrus.Fatalf("HTTP error: %v\n", err)
	}
}
