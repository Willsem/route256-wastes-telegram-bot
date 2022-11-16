package startup

import "github.com/segmentio/kafka-go"

type KafkaConfig struct {
	Brockers []string `yaml:"brockers"`
	Topic    string   `yaml:"topic"`
}

func NewKafkaProducer(config KafkaConfig) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(config.Brockers...),
		Topic:    config.Topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func NewKafkaConsumer(config KafkaConfig) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   config.Brockers,
		Topic:     config.Topic,
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})
}
