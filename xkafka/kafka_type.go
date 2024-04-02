package xkafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	DefaultId = "default"

	RoleConsumer = "consumer"
	RoleProducer = "producer"
)

type Config struct {
	// common
	Id               string // optional
	Role             string // optional, default is consumer
	Auto             bool   // optional
	BootstrapServers string // must
	SecurityProtocol string // optional
	SaslUsername     string // optional
	SaslPassword     string // optional
	SaslMechanism    string // optional
	SslCaLocation    string // optional

	// producer
	Topic                 string                     // must
	RequestRequiredAcks   int                        // optional
	DeliveredCallback     func(kafka.TopicPartition) // optional
	DeliverFailedCallback func(kafka.TopicPartition) // optional

	// consumer
	Topics                 []string // must
	GroupId                string   // must
	AutoOffsetReset        string   // optional
	MaxPollIntervalMS      int      // optional
	SessionTimeoutMS       int      // optional
	HeartbeatIntervalMS    int      // optional
	FetchMaxBytes          int      // optional
	MaxPartitionFetchBytes int      // optional
}

func (c Config) Default() Config {
	if c.Id == "" {
		c.Id = DefaultId
	}
	if c.Role == "" {
		c.Role = RoleConsumer
	}
	if c.AutoOffsetReset == "" {
		c.AutoOffsetReset = "latest"
	}
	if c.MaxPollIntervalMS == 0 {
		c.MaxPollIntervalMS = 10000 // 10s
	}
	if c.SessionTimeoutMS == 0 {
		c.SessionTimeoutMS = 10000 // 10s
	}
	if c.HeartbeatIntervalMS == 0 {
		c.HeartbeatIntervalMS = 3000 // 3s
	}
	if c.FetchMaxBytes == 0 {
		c.FetchMaxBytes = 1024000 // 1M
	}
	if c.MaxPartitionFetchBytes == 0 {
		c.MaxPartitionFetchBytes = 512000 // 500K
	}
	return c
}
