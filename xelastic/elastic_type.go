package xelastic

const DefaultID = "default"

type Config struct {
	ID        string   // optional
	Hosts     []string // must
	Username  string   // must
	Password  string   // must
	EnableLog bool     // optional
}

func (c *Config) Default() {
	if c.ID == "" {
		c.ID = DefaultID
	}
}
