package models

type KafkaMessage struct {
	Key     []byte
	Message []byte
}
