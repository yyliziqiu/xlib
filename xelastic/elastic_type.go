package xelastic

const DefaultId = "default"

type Config struct {
	// must
	Hosts    []string
	Username string
	Password string

	// optional
	Id             string
	EnableLog      bool
	EnableLogTrace bool
	LogName        string
}

func (c Config) WithDefault() Config {
	if c.Id == "" {
		c.Id = DefaultId
	}
	return c
}
