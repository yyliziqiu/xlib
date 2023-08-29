package xelastic

import (
	"fmt"
)

const DefaultId = "default"

type Config struct {
	Id         string   `json:"id"`
	Hosts      []string `json:"hosts"`
	User       string   `json:"user"`
	Password   string   `json:"password"`
	LogEnabled bool     `json:"log_enabled"`
	LogTrace   bool     `json:"log_trace"`
	LogName    string   `json:"log_name"`
}

func (c Config) WithDefault() Config {
	if c.Id == "" {
		c.Id = DefaultId
	}
	if c.LogName == "" {
		c.LogName = fmt.Sprintf("elastic-%s-", c.Id)
	}
	return c
}
