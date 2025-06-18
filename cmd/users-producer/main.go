package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"
	configs "wallet-service/internal/config"
	"wallet-service/internal/transport/kafka/producer"
	"wallet-service/pkg/generator"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := configs.Init()
	if err != nil {
		logrus.Panicf("Config error: %v\n", err)
	}

	producer := producer.New(cfg)

	defer func() {
		if err := producer.Close(); err != nil {
			logrus.Panicf("Close producer error: %v\n", err)
		}
	}()

	for {
		fakeUser, err := generator.GenerateUser()
		if err != nil {
			logrus.Panicf("Error to create: %v\n", err)
		}

		if err := producer.Produce(cfg, ctx, fakeUser); err != nil {
			logrus.Panicf("Producer error: %v\n", err)
		}

		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return
		}
	}
}
