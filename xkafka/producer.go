package xkafka

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewProducer(config Config) (*kafka.Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers":     config.BootstrapServers,
		"request.required.acks": config.RequestRequiredAcks,
	}

	if config.SecurityProtocol != "" {
		_ = conf.SetKey("security.protocol", config.SecurityProtocol)
		switch config.SecurityProtocol {
		case "plaintext":
		case "sasl_plaintext":
			_ = conf.SetKey("sasl.username", config.SaslUsername)
			_ = conf.SetKey("sasl.password", config.SaslPassword)
			_ = conf.SetKey("sasl.mechanism", config.SaslMechanism)
		case "sasl_ssl":
			_ = conf.SetKey("sasl.username", config.SaslUsername)
			_ = conf.SetKey("sasl.password", config.SaslPassword)
			_ = conf.SetKey("sasl.mechanism", config.SaslMechanism)
			_ = conf.SetKey("ssl.ca.location", config.SslCaLocation)
		}
	}

	producer, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("create producer failed [%v]", err)
	}

	go func(deliveredCallback func(partition kafka.TopicPartition), deliverFailedCallback func(partition kafka.TopicPartition)) {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error == nil {
					if deliveredCallback != nil {
						deliveredCallback(ev.TopicPartition)
					}
				} else {
					if deliverFailedCallback != nil {
						deliverFailedCallback(ev.TopicPartition)
					}
				}
			}
		}
	}(config.DeliveredCallback, config.DeliverFailedCallback)

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
