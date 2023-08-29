package xkafka

const (
	DefaultId = "default"

	RoleConsumer = "consumer"
	RoleProducer = "producer"
)

type Config struct {
	Id   string `json:"id"`
	Role string `json:"role"`

	// common
	BootstrapServers string `json:"bootstrapServers"`
	SecurityProtocol string `json:"securityProtocol"`
	SaslUsername     string `json:"saslUsername"`
	SaslPassword     string `json:"saslPassword"`
	SaslMechanism    string `json:"saslMechanism"`
	SslCaLocation    string `json:"sslCaLocation"`

	// producer
	RequiredAcks int    `json:"requiredAcks"`
	Topic        string `json:"topic"`

	// consumer
	Topics            []string `json:"topics"`
	GroupId           string   `json:"groupId"`
	OffsetReset       string   `json:"offsetReset"`
	PollInterval      int      `json:"pollInterval"`
	SessionTimeout    int      `json:"sessionTimeout"`
	HeartbeatInterval int      `json:"heartbeatInterval"`
	FetchMax          int      `json:"fetchMax"`
	PartitionFetchMax int      `json:"partitionFetchMax"`
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
