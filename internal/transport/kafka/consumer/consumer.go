package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	configs "wallet-service/internal/config"
	"wallet-service/internal/domain"
)

type Consumer struct {
	kf   *kafka.Reader
	repo usersDb
}

type usersDb interface {
	UpsertUser(ctx context.Context, user domain.User) error
	GetUser(ctx context.Context, user uuid.UUID) (domain.User, error)
}

func New(cfg *configs.Config, repo usersDb) *Consumer {
	kf := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		GroupID: cfg.Kafka.GroupID,
		Topic:   cfg.Kafka.Topic,
	})

	return &Consumer{
		kf:   kf,
		repo: repo,
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	logrus.Info("Consuming messages...")

	for {
		msg, err := c.kf.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("failed to consume a messages: %w", err)
		}

		var user domain.User
		if err := json.Unmarshal(msg.Value, &user); err != nil {
			return fmt.Errorf("failed to unmarshal domain.User: %w", err)
		}

		if err := c.repo.UpsertUser(ctx, user); err != nil {
			return fmt.Errorf("failed to create or update the user: %w", err)
		}

		logrus.Printf("topic: %s message: %s", msg.Topic, string(msg.Value))
	}
}

func (c *Consumer) Close() error {
	if err := c.kf.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka consumer: %w", err)
	}

	return nil
}
