package xkafka

const (
	DefaultId = "default"

	RoleConsumer = "consumer"
	RoleProducer = "producer"
)

type Config struct {
	Id   string
	Role string

	// common
	BootstrapServers string
	SecurityProtocol string
	SaslUsername     string
	SaslPassword     string
	SaslMechanism    string
	SslCaLocation    string

	// producer
	RequiredAcks int
	Topic        string

	// consumer
	Topics            []string
	GroupId           string
	OffsetReset       string
	PollInterval      int
	SessionTimeout    int
	HeartbeatInterval int
	FetchMax          int
	PartitionFetchMax int
}

func (c Config) WithDefault() Config {
	if c.Id == "" {
		c.Id = DefaultId
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
	return c
}

func (c Config) GetRole() string {
	if c.Role == RoleProducer {
		return RoleProducer
	}
	return RoleConsumer
}
