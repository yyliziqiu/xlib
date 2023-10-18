package xlog

import "time"

const textFormatter = "text"
const jsonFormatter = "json"

type Config struct {
	Console              bool
	Name                 string
	Path                 string
	Formatter            string
	Level                string
	MaxAge               time.Duration
	RotationTime         time.Duration
	DisableLevelRotation bool
}

func (c Config) WithDefault() Config {
	if c.Name == "" {
		c.Name = "app"
	}
	if c.Path == "" {
		c.Path = "logs"
	}
	if c.Formatter == "" {
		c.Formatter = textFormatter
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
	return c
}
