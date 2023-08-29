package xkafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/yyliziqiu/xlib/xutil"
)

var (
	cs        map[string]Config
	consumers map[string]*kafka.Consumer
	producers map[string]*kafka.Producer
)

func Initialize(requireNew bool, configs ...Config) error {
	cs = make(map[string]Config)
	for _, config := range configs {
		cs[config.Id] = config
	}

	if !requireNew {
		return nil
	}

	consumers = make(map[string]*kafka.Consumer, 8)
	producers = make(map[string]*kafka.Producer, 8)
	for _, config := range configs {
		switch config.GetRole() {
		case RoleConsumer:
			consumer, err := NewConsumer(config)
			if err != nil {
				Finally()
				return err
			}
			consumers[xutil.IES(config.Id, DefaultId)] = consumer
		case RoleProducer:
			producer, err := NewProducer(config)
			if err != nil {
				Finally()
				return err
			}
			producers[xutil.IES(config.Id, DefaultId)] = producer
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
	return cs[id]
}

func GetDefaultTopic() string {
	return GetTopic(DefaultId)
}

func GetTopic(id string) string {
	return cs[id].Topic
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
