package xorm

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/yyliziqiu/xlib/xutil"
)

var globalLogger *logrus.Logger

func SetGlobalLogger(logger *logrus.Logger) {
	globalLogger = logger
}

var dbs map[string]*gorm.DB

func Initialize(configs ...Config) error {
	dbs = make(map[string]*gorm.DB, len(configs))
	for _, config := range configs {
		db, err := New(config)
		if err != nil {
			return err
		}
		dbs[xutil.IES(config.Id, DefaultId)] = db
	}
	return nil
}

func New(config Config) (*gorm.DB, error) {
	config = config.WithDefault()

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

func GetDB(id string) *gorm.DB {
	return dbs[id]
}

func GetDefaultDB() *gorm.DB {
	return GetDB(DefaultId)
}
