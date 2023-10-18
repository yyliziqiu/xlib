package xorm

import (
	"fmt"
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
	Id   string `json:"id"`
	DSN  string `json:"dsn"`
	Type string `json:"type"`

	// optional
	MaxOpenConns                 int           `json:"max_open_conns"`
	MaxIdleConns                 int           `json:"max_idle_conns"`
	ConnMaxLifetime              time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime              time.Duration `json:"conn_max_idle_time"`
	LogEnabled                   bool          `json:"log_enabled"`
	LogName                      string        `json:"log_name"`
	LogLevel                     int           `json:"log_level"`
	LogSlowThreshold             time.Duration `json:"log_slow_threshold"`
	LogParameterizedQueries      bool          `json:"log_parameterized_queries"`
	LogIgnoreRecordNotFoundError bool          `json:"log_ignore_record_not_found_error"`
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
	if c.LogName == "" {
		c.LogName = fmt.Sprintf("gorm-%s-", c.Id)
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
	if !c.LogEnabled {
		return &gorm.Config{}
	}

	var lg *logrus.Logger
	if globalLogger != nil {
		lg = globalLogger
	} else {
		lg = xlog.MustNewLoggerByName(c.LogName)
	}

	return &gorm.Config{Logger: logger.New(lg, logger.Config{
		LogLevel:                  logger.LogLevel(c.LogLevel),    // Log level
		SlowThreshold:             c.LogSlowThreshold,             // Slow SQL threshold
		ParameterizedQueries:      c.LogParameterizedQueries,      // Don't include params in the SQL log
		IgnoreRecordNotFoundError: c.LogIgnoreRecordNotFoundError, // Ignore ErrRecordNotFound error for logger
	})}
}
