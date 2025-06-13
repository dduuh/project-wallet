package producer

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	producer *kafka.Writer
}

func New(brokers []string, topic string) *Producer {
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})

	return &Producer{
		producer: producer,
	}
}

func (p *Producer) Produce(ctx context.Context, topic string, value []byte) error {
	err := p.producer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
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