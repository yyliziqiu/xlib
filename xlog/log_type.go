package xlog

import "time"

const (
	TextFormatterName = "text"
	JSONFormatterName = "json"
)

type Config struct {
	Console         bool
	Path            string
	Name            string
	Level           string
	MaxAge          time.Duration
	RotationTime    time.Duration
	RotationLevel   int
	Formatter       string
	EnableCaller    bool
	TimestampFormat string
}

func (c Config) Default() Config {
	if c.Path == "" {
		c.Console = true
	}
	if c.Name == "" {
		c.Name = "app"
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
		c.TimestampFormat = time.DateTime
	}
	return c
}
