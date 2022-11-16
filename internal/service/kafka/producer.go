package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	client *kafka.Writer
}

func NewProducer(client *kafka.Writer) *Producer {
	return &Producer{
		client: client,
	}
}

func (p *Producer) SendMessage(ctx context.Context, key []byte, value []byte) error {
	err := p.client.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to send the message to kafka: %w", err)
	}

	return nil
}
