package producer

import (
	"context"
	"time"
	configs "wallet-service/internal/config"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Producer struct {
	producer *kafka.Writer
}

func New(cfg *configs.Config) *Producer {
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.KafkaCfg.Brokers,
		Topic:   cfg.KafkaCfg.Topic,
	})

	return &Producer{
		producer: producer,
	}
}

func (p *Producer) Produce(cfg *configs.Config, ctx context.Context, value []byte) error {
	logrus.Info("Writing a messages")
	err := p.producer.WriteMessages(ctx, kafka.Message{
		Value: value,
		Time:  time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
