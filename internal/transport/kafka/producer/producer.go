package producer

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	configs "wallet-service/internal/config"
)

type Producer struct {
	producer *kafka.Writer
}

func New(cfg *configs.Config) *Producer {
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
	})

	return &Producer{
		producer: producer,
	}
}

func (p *Producer) Produce(cfg *configs.Config, ctx context.Context, value []byte) error {
	err := p.producer.WriteMessages(ctx, kafka.Message{
		Value: value,
		Time:  time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to produce a messages: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	if err := p.producer.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka producer: %w", err)
	}

	return nil
}
