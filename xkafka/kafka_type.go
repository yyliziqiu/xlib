package xkafka

const (
	DefaultID = "default"

	RoleConsumer = "consumer"
	RoleProducer = "producer"
)

type Config struct {
	ID   string // optional
	Role string // optional, default is consumer

	// common
	BootstrapServers string // must
	SecurityRequired bool   // optional
	SecurityProtocol string // optional
	SaslUsername     string // optional
	SaslPassword     string // optional
	SaslMechanism    string // optional
	SslCaLocation    string // optional

	// producer
	RequiredAcks int    // optional
	Topic        string // must

	// consumer
	Topics            []string // must
	GroupId           string   // must
	OffsetReset       string   // optional
	PollInterval      int      // optional
	SessionTimeout    int      // optional
	HeartbeatInterval int      // optional
	FetchMax          int      // optional
	PartitionFetchMax int      // optional
}

func (c *Config) Default() {
	if c.ID == "" {
		c.ID = DefaultID
	}
	if c.Role == "" {
		c.Role = RoleConsumer
	}
	if c.OffsetReset == "" {
		c.OffsetReset = "latest"
	}
	if c.PollInterval == 0 {
		c.PollInterval = 10000 // 10s
	}
	if c.SessionTimeout == 0 {
		c.SessionTimeout = 10000 // 10s
	}
	if c.HeartbeatInterval == 0 {
		c.HeartbeatInterval = 2000 // 3s
	}
	if c.FetchMax == 0 {
		c.FetchMax = 1024000 // 1M
	}
	if c.PartitionFetchMax == 0 {
		c.PartitionFetchMax = 512000 // 500K
	}
}
