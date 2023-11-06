package xkafka

import (
	"errors"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	configs   map[string]Config
	consumers map[string]*kafka.Consumer
	producers map[string]*kafka.Producer
)

func Init(autoNew bool, cfs ...Config) error {
	configs = make(map[string]Config, 16)
	for _, config := range cfs {
		config.Default()
		configs[config.ID] = config
	}

	if !autoNew {
		return nil
	}

	consumers = make(map[string]*kafka.Consumer, 16)
	producers = make(map[string]*kafka.Producer, 16)
	for _, config := range configs {
		switch config.Role {
		case RoleConsumer:
			consumer, err := NewConsumer(config)
			if err != nil {
				Finally()
				return err
			}
			consumers[config.ID] = consumer
		case RoleProducer:
			producer, err := NewProducer(config)
			if err != nil {
				Finally()
				return err
			}
			producers[config.ID] = producer
		default:
			return errors.New("not support kafka role")
		}
	}

	return nil
}

func Finally() {
	for _, consumer := range consumers {
		_ = consumer.Close()
	}
	for _, producer := range producers {
		producer.Close()
	}
}

func GetConfig(id string) Config {
	return configs[id]
}

func GetDefaultConfig() Config {
	return GetConfig(DefaultID)
}

func GetTopic(id string) string {
	return configs[id].Topic
}

func GetDefaultTopic() string {
	return GetTopic(DefaultID)
}

func GetConsumer(id string) *kafka.Consumer {
	return consumers[id]
}

func GetDefaultConsumer() *kafka.Consumer {
	return GetConsumer(DefaultID)
}

func GetProducer(id string) *kafka.Producer {
	return producers[id]
}

func GetDefaultProducer() *kafka.Producer {
	return GetProducer(DefaultID)
}
