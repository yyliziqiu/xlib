package xkafka

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewProducer(config Config) (*kafka.Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers":     config.BootstrapServers,
		"request.required.acks": config.RequiredAcks,
	}

	producer, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("create producer failed [%v]", err)
	}

	return producer, nil
}

func NewProducerById(id string) (*kafka.Producer, error) {
	config, ok := configs[id]
	if !ok {
		return nil, fmt.Errorf("kafka config [%s] not exist", id)
	}
	return NewProducer(config)
}

func NewDefaultProducer() (*kafka.Producer, error) {
	return NewProducerById(DefaultId)
}

func Produce(producer *kafka.Producer, topic string, message []byte) error {
	return producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
}

func ProduceObject(producer *kafka.Producer, topic string, object interface{}) error {
	message, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return Produce(producer, topic, message)
}
