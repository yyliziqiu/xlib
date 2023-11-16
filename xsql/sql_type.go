package xsql

import (
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const (
	DefaultId = "default"

	DBTypeMySQL    = "mysql"
	DBTypePostgres = "postgres"
)

type Config struct {
	Id              string        // optional
	Type            string        // optional
	DSN             string        // must
	MaxOpenConns    int           // optional
	MaxIdleConns    int           // optional
	ConnMaxLifetime time.Duration // optional
	ConnMaxIdleTime time.Duration // optional

	// only valid when use gorm
	EnableORM                    bool          // optional
	LogLevel                     int           // optional
	LogSlowThreshold             time.Duration // optional
	LogParameterizedQueries      bool          // optional
	LogIgnoreRecordNotFoundError bool          // optional
}

func (c *Config) Default() {
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
		c.LogLevel = 1
	}
	if c.LogSlowThreshold == 0 {
		c.LogSlowThreshold = 3 * time.Second
	}
}

func (c Config) GORMConfig() *gorm.Config {
	return &gorm.Config{Logger: gormlogger.New(GetLogger(), gormlogger.Config{
		LogLevel:                  gormlogger.LogLevel(c.LogLevel), // Log level
		SlowThreshold:             c.LogSlowThreshold,              // Slow SQL threshold
		ParameterizedQueries:      c.LogParameterizedQueries,       // Don't include params in the SQL log
		IgnoreRecordNotFoundError: c.LogIgnoreRecordNotFoundError,  // Ignore ErrRecordNotFound error for logger
	})}
}
