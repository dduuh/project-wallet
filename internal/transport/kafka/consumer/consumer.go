package consumer

import (
	"context"
	configs "wallet-service/internal/config"
	"wallet-service/internal/domain"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	kf   *kafka.Reader
	repo usersDb
}

type usersDb interface {
	UpsertUser(ctx context.Context, user domain.User) error
}

func New(cfg configs.Config, repo usersDb) *Consumer {
	kf := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.KafkaCfg.Brokers,
		GroupID: cfg.KafkaCfg.GroupID,
		Topic:   cfg.KafkaCfg.Topic,
	})

	return &Consumer{
		kf: kf,
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	for {
		msg, err := c.kf.ReadMessage(ctx)
		if err != nil {
			return err
		}

		logrus.Printf("topic: %s message: %v", msg.Topic, msg.Value)
	}
}

func (c *Consumer) Close() error {
	return c.kf.Close()
}
