package xkafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewConsumer(config Config) (*kafka.Consumer, error) {
	var conf = &kafka.ConfigMap{
		"bootstrap.servers":         config.BootstrapServers,
		"group.id":                  config.GroupId,
		"auto.offset.reset":         config.AutoOffsetReset,
		"max.poll.interval.ms":      config.MaxPollIntervalMS,
		"session.timeout.ms":        config.SessionTimeoutMS,
		"heartbeat.interval.ms":     config.HeartbeatIntervalMS,
		"fetch.max.bytes":           config.FetchMaxBytes,
		"max.partition.fetch.bytes": config.MaxPartitionFetchBytes,
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

	consumer, err := kafka.NewConsumer(conf)
	if err != nil {
		return nil, fmt.Errorf("create consumer failed [%v]", err)
	}

	err = consumer.SubscribeTopics(config.Topics, nil)
	if err != nil {
		return nil, fmt.Errorf("subscribe topic failed [%v]", err)
	}

	return consumer, nil
}

func NewConsumerByConfigId(id string) (*kafka.Consumer, error) {
	config, ok := configs[id]
	if !ok {
		return nil, fmt.Errorf("kafka config [%s] not exist", id)
	}
	return NewConsumer(config)
}

func NewDefaultConsumer() (*kafka.Consumer, error) {
	return NewConsumerByConfigId(DefaultId)
}

// func consume(consumer *kafka.Consumer) {
// 	for {
// 		select {
// 		case <-quit:
// 			log.Info("[runConsumeKafkaMessage] Quit.")
// 			return
// 		default:
// 			msg, err := consumer.ReadMessage(-1)
// 			if err != nil {
// 				continue
// 			}
// 			// to do something
// 		}
// 	}
// }
