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
	Id   string `json:"id"`
	DSN  string `json:"dsn"`
	Type string `json:"type"`

	// optional
	MaxOpenConns    int           `json:"maxOpenConns"`
	MaxIdleConns    int           `json:"maxIdleConns"`
	ConnMaxLifetime time.Duration `json:"connMaxLifetime"`
	ConnMaxIdleTime time.Duration `json:"connMaxIdleTime"`
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
	return c
}
