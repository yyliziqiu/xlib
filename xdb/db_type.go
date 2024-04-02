package xdb

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	ormlogger "gorm.io/gorm/logger"

	"github.com/yyliziqiu/xlib/xlog"
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
	EnableORM                       bool           // optional
	ORMLogger                       *logrus.Logger // optional
	ORMLogLevel                     int            // optional
	ORMLogSlowThreshold             time.Duration  // optional
	ORMLogParameterizedQueries      bool           // optional
	ORMLogIgnoreRecordNotFoundError bool           // optional
}

func (c Config) Default() Config {
	if c.Id == "" {
		c.Id = DefaultId
	}
	if c.Type == "" {
		c.Type = DBTypeMySQL
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 20
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 10
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = time.Hour
	}
	if c.ConnMaxIdleTime == 0 {
		c.ConnMaxLifetime = 30 * time.Minute
	}
	if c.ORMLogLevel == 0 {
		c.ORMLogLevel = 1
	}
	if c.ORMLogSlowThreshold == 0 {
		c.ORMLogSlowThreshold = 5 * time.Second
	}
	if c.ORMLogLevel > 1 && c.ORMLogger == nil {
		c.ORMLogger = xlog.NewWithNameMust("orm-" + c.Id)
	}
	return c
}

func (c Config) ORMConfig() *gorm.Config {
	return &gorm.Config{
		Logger: ormlogger.New(
			c.ORMLogger,
			ormlogger.Config{
				LogLevel:                  ormlogger.LogLevel(c.ORMLogLevel),
				SlowThreshold:             c.ORMLogSlowThreshold,
				ParameterizedQueries:      c.ORMLogParameterizedQueries,
				IgnoreRecordNotFoundError: c.ORMLogIgnoreRecordNotFoundError,
			},
		),
	}
}
