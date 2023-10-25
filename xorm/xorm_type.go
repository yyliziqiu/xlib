package xorm

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/yyliziqiu/xlib/xlog"
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
	Id                           string
	Type                         string
	MaxOpenConns                 int
	MaxIdleConns                 int
	ConnMaxLifetime              time.Duration
	ConnMaxIdleTime              time.Duration
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

func (c Config) GORMConfig() *gorm.Config {
	if !c.EnableLog {
		return &gorm.Config{}
	}

	var lgg *logrus.Logger
	if c.LogName != "" {
		lgg = xlog.MustNewLoggerByName(c.LogName)
	} else {
		if globalLogger == nil {
			globalLogger = xlog.MustNewLoggerByName("gorm")
		}
		lgg = globalLogger
	}

	return &gorm.Config{Logger: logger.New(lgg, logger.Config{
		LogLevel:                  logger.LogLevel(c.LogLevel),    // Log level
		SlowThreshold:             c.LogSlowThreshold,             // Slow SQL threshold
		ParameterizedQueries:      c.LogParameterizedQueries,      // Don't include params in the SQL log
		IgnoreRecordNotFoundError: c.LogIgnoreRecordNotFoundError, // Ignore ErrRecordNotFound error for logger
	})}
}
