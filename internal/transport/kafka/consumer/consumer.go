package consumer

import (
	"context"
	"wallet-service/internal/domain"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	kf *kafka.Reader
}

type usersDb interface {
	UpsertUser(ctx context.Context, user domain.User) error
}

func New(brokers []string, topic, groupId string) *Consumer {
	kf := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupId,
		Topic:   topic,
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
