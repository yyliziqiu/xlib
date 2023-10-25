package xkafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	cfs       map[string]Config
	consumers map[string]*kafka.Consumer
	producers map[string]*kafka.Producer
)

func Init(autoNew bool, configs ...Config) error {
	cfs = make(map[string]Config, len(configs))
	for _, config := range configs {
		config = config.WithDefault()
		cfs[config.Id] = config
	}

	if !autoNew {
		return nil
	}

	consumers = make(map[string]*kafka.Consumer, 8)
	producers = make(map[string]*kafka.Producer, 8)
	for _, cf := range cfs {
		switch cf.GetRole() {
		case RoleConsumer:
			consumer, err := NewConsumer(cf)
			if err != nil {
				Finally()
				return err
			}
			consumers[cf.Id] = consumer
		case RoleProducer:
			producer, err := NewProducer(cf)
			if err != nil {
				Finally()
				return err
			}
			producers[cf.Id] = producer
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

func GetDefaultConfig() Config {
	return GetConfig(DefaultId)
}

func GetConfig(id string) Config {
	return cfs[id]
}

func GetDefaultTopic() string {
	return GetTopic(DefaultId)
}

func GetTopic(id string) string {
	return cfs[id].Topic
}

func GetDefaultConsumer() *kafka.Consumer {
	return GetConsumer(DefaultId)
}

func GetConsumer(id string) *kafka.Consumer {
	return consumers[id]
}

func GetDefaultProducer() *kafka.Producer {
	return GetProducer(DefaultId)
}

func GetProducer(id string) *kafka.Producer {
	return producers[id]
}
