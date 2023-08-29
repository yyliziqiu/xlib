package xlog

import "time"

type Config struct {
	Console    bool          `json:"console"`
	Path       string        `json:"path"`
	Name       string        `json:"name"`
	Level      string        `json:"level"`
	MaxAge     time.Duration `json:"maxAge"`
	RotateTime time.Duration `json:"rotateTime"`
}

func (c Config) WithDefault() Config {
	if c.Path == "" {
		c.Path = "logs"
	}
	if c.Level == "" {
		c.Level = "debug"
	}
	if c.MaxAge == 0 {
		c.MaxAge = 15 * 24 * time.Hour
	}
	if c.RotateTime == 0 {
		c.RotateTime = 24 * time.Hour
	}
	return c
}
