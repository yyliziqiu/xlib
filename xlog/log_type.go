package xlog

import "time"

const (
	TextFormatterName = "text"
	JSONFormatterName = "json"
)

type Config struct {
	Console         bool
	Name            string
	Path            string
	Level           string
	MaxAge          time.Duration
	RotationTime    time.Duration
	RotationLevel   int
	Formatter       string
	EnableCaller    bool
	TimestampFormat string
}

func (c *Config) Default() {
	if c.Name == "" {
		c.Name = "app"
	}
	if c.Path == "" {
		c.Path = "logs"
	}
	if c.Level == "" {
		c.Level = "debug"
	}
	if c.MaxAge == 0 {
		c.MaxAge = 7 * 24 * time.Hour
	}
	if c.RotationTime == 0 {
		c.RotationTime = 24 * time.Hour
	}
	if c.Formatter == "" {
		c.Formatter = TextFormatterName
	}
	if c.TimestampFormat == "" {
		c.TimestampFormat = "2006-01-02 15:04:05"
	}
}
