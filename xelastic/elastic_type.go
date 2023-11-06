package xelastic

const DefaultId = "default"

type Config struct {
	Id        string   // optional
	Hosts     []string // must
	Username  string   // must
	Password  string   // must
	EnableLog bool     // optional
}

func (c *Config) Default() {
	if c.Id == "" {
		c.Id = DefaultId
	}
}
