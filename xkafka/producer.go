package xkafka

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewDefaultProducer() (*kafka.Producer, error) {
	return NewProducerById(DefaultId)
}

func NewProducerById(id string) (*kafka.Producer, error) {
	config, ok := cs[id]
	if !ok {
		return nil, fmt.Errorf("kafka ID: %s not exist")
	}
	return NewProducer(config)
}

func NewProducer(config Config) (*kafka.Producer, error) {
	config = config.WithDefault()

	conf := &kafka.ConfigMap{
		"bootstrap.servers":     config.BootstrapServers,
		"request.required.acks": config.RequiredAcks,
	}

	producer, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("create producer error: %v", err)
	}

	return producer, nil
}

func SendObject(producer *kafka.Producer, topic string, object interface{}) error {
	message, err := json.Marshal(object)
	if err != nil {
		return err
	}

	return producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
}

func SendMessage(producer *kafka.Producer, topic string, message []byte) error {
	return producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
}
