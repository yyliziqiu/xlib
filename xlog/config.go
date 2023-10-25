package xlog

import "time"

const (
	textFormatter = "text"
	jsonFormatter = "json"

	fieldsFilename = "filename"
	fieldsFunction = "function"
	fieldsAll      = "all"
)

type Config struct {
	Console              bool
	Name                 string
	Path                 string
	Level                string
	MaxAge               time.Duration
	RotationTime         time.Duration
	Formatter            string
	EnableCaller         bool
	CallerFields         string
	TimestampFormat      string
	DisableLevelRotation bool
}

func (c Config) WithDefault() Config {
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
		c.Formatter = textFormatter
	}
	if c.CallerFields == "" {
		c.CallerFields = fieldsFunction
	}
	if c.TimestampFormat == "" {
		c.TimestampFormat = "2006-01-02 15:04:05"
	}
	return c
}
