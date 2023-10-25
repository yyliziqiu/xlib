package xorm

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/yyliziqiu/xlib/xlog"
	"github.com/yyliziqiu/xlib/xsql"
)

var (
	globalLogger *logrus.Logger

	cfs map[string]Config
	dbs map[string]*gorm.DB
)

func SetGlobalLogger(logger *logrus.Logger) {
	globalLogger = logger
}

func Init(configs ...Config) error {
	cfs = make(map[string]Config, len(configs))
	for _, config := range configs {
		config = config.WithDefault()
		cfs[config.Id] = config
	}

	dbs = make(map[string]*gorm.DB, len(cfs))
	for _, cf := range cfs {
		db, err := New(cf)
		if err != nil {
			return err
		}
		dbs[cf.Id] = db
	}

	return nil
}

func New(config Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch config.Type {
	case DBTypeMySQL:
		db, err = gorm.Open(mysql.Open(config.DSN), config.GORMConfig())
	case DBTypePostgres:
		db, err = gorm.Open(postgres.Open(config.DSN), config.GORMConfig())
	default:
		return nil, errors.New("not support db type")
	}
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxLifetime)

	return db, nil
}

func NewWithDSN(typ string, dsn string) (*gorm.DB, error) {
	return New(Config{DSN: dsn, Type: typ})
}

func NewWithSQLDB(id string) (*gorm.DB, error) {
	config := xsql.GetConfig(id)
	db := xsql.GetDB(id)

	gormConfig := newGORMConfigBySQLConfig(config)

	switch config.Type {
	case xsql.DBTypeMySQL:
		return gorm.Open(mysql.New(mysql.Config{Conn: db}), gormConfig)
	case xsql.DBTypePostgres:
		return gorm.Open(postgres.New(postgres.Config{Conn: db}), gormConfig)
	default:
		return nil, errors.New("not support db type")
	}
}

func newGORMConfigBySQLConfig(c xsql.Config) *gorm.Config {
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

func GetDB(id string) *gorm.DB {
	return dbs[id]
}

func GetDefaultDB() *gorm.DB {
	return GetDB(DefaultId)
}
