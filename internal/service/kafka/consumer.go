package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
)

type ConsumerConfig struct {
	BufferSize uint `yaml:"buffer_size"`
}

type Consumer struct {
	client *kafka.Reader
	config ConsumerConfig
	logger log.Logger

	cancel context.CancelFunc
	done   chan struct{}

	messages chan *models.KafkaMessage
}

func NewConsumer(client *kafka.Reader, config ConsumerConfig, logger log.Logger) *Consumer {
	return &Consumer{
		client: client,
		config: config,
		logger: logger.With(log.ComponentKey, "Kafka consumer"),
	}
}

func (c *Consumer) Start() error {
	ctx, cancel := context.WithCancel(context.Background())

	c.cancel = cancel
	c.done = make(chan struct{})

	c.messages = make(chan *models.KafkaMessage, c.config.BufferSize)

	go c.run(ctx)

	return nil
}

func (c *Consumer) Stop(ctx context.Context) error {
	c.cancel()

	select {
	case <-c.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Consumer) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.logger.WithError(ctx.Err()).Info("kafka consumer has been closed")
			close(c.messages)
			close(c.done)

		default:
			msg, err := c.client.ReadMessage(ctx)
			if err != nil {
				c.logger.
					WithError(err).
					Error("failed to read message from kafka")
				continue
			}

			c.logger.With("kafka_message", msg).Debug("recieved message from Kafka")

			c.messages <- &models.KafkaMessage{
				Key:     msg.Key,
				Message: msg.Value,
			}
		}
	}
}

func (c *Consumer) GetMessageChan() <-chan *models.KafkaMessage {
	return c.messages
}
