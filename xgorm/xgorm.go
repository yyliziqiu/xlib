package xorm

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/yyliziqiu/xlib/xlog"
	"github.com/yyliziqiu/xlib/xutil"
)

var globalLogger *logrus.Logger

func SetGlobalLogger(logger *logrus.Logger) {
	globalLogger = logger
}

var dbs map[string]*gorm.DB

func Init(configs ...Config) error {
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

	db, err := gorm.Open(config.Type, config.DSN)
	if err != nil {
		return nil, err
	}

	db.LogMode(config.LogEnabled)
	if globalLogger != nil {
		db.SetLogger(globalLogger)
	} else {
		db.SetLogger(xlog.NewLoggerMust(config.LogName))
	}

	sqlDB := db.DB()
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
