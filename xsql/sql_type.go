package xsql

import (
	"time"
)

const (
	DefaultId = "default"

	DBTypeMySQL    = "mysql"
	DBTypePostgres = "postgres"
)

type Config struct {
	// must
	DSN string

	// optional
	Id              string
	Type            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration

	// only valid when use gorm
	EnableLog                    bool
	LogName                      string
	LogLevel                     int
	LogSlowThreshold             time.Duration
	LogParameterizedQueries      bool
	LogIgnoreRecordNotFoundError bool
}

func (c Config) WithDefault() Config {
	if c.Id == "" {
		c.Id = DefaultId
	}
	if c.Type == "" {
		c.Type = DBTypeMySQL
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 50
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 10
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = time.Hour
	}
	if c.ConnMaxIdleTime == 0 {
		c.ConnMaxLifetime = 10 * time.Minute
	}
	if c.LogLevel == 0 {
		c.LogLevel = 4
	}
	if c.LogSlowThreshold == 0 {
		c.LogSlowThreshold = 15 * time.Second
	}
	return c
}
