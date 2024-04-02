package xkafka

import (
	"errors"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	_configs   map[string]Config
	_consumers map[string]*kafka.Consumer
	_producers map[string]*kafka.Producer
)

func Init(configs ...Config) error {
	_configs = make(map[string]Config, 16)
	for _, config := range configs {
		_configs[config.Id] = config.Default()
	}

	_consumers = make(map[string]*kafka.Consumer, 8)
	_producers = make(map[string]*kafka.Producer, 8)
	for _, config := range _configs {
		if !config.Auto {
			continue
		}
		switch config.Role {
		case RoleConsumer:
			consumer, err := NewConsumer(config)
			if err != nil {
				Finally()
				return err
			}
			_consumers[config.Id] = consumer
		case RoleProducer:
			producer, err := NewProducer(config)
			if err != nil {
				Finally()
				return err
			}
			_producers[config.Id] = producer
		default:
			return errors.New("not support kafka role")
		}
	}

	return nil
}

func Finally() {
	for _, consumer := range _consumers {
		_ = consumer.Close()
	}
	for _, producer := range _producers {
		producer.Close()
	}
}

func GetConfig(id string) Config {
	return _configs[id]
}

func GetDefaultConfig() Config {
	return GetConfig(DefaultId)
}

func GetConsumer(id string) *kafka.Consumer {
	return _consumers[id]
}

func GetDefaultConsumer() *kafka.Consumer {
	return GetConsumer(DefaultId)
}

func GetProducer(id string) *kafka.Producer {
	return _producers[id]
}

func GetDefaultProducer() *kafka.Producer {
	return GetProducer(DefaultId)
}
