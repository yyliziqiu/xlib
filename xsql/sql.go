package xsql

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/yyliziqiu/xlib/xlog"
)

var (
	logger     *logrus.Logger
	loggerOnce sync.Once
)

func SetLogger(lgg *logrus.Logger) {
	logger = lgg
}

func GetLogger() *logrus.Logger {
	if logger != nil {
		return logger
	}
	loggerOnce.Do(func() {
		if logger == nil {
			logger = xlog.NewWithNameMust("gorm")
		}
	})
	return logger
}

var (
	configs map[string]Config
	sqls    map[string]*sql.DB
	orms    map[string]*gorm.DB
)

func Init(cfs ...Config) error {
	configs = make(map[string]Config, 16)
	for _, config := range cfs {
		config.Default()
		configs[config.Id] = config
	}

	sqls = make(map[string]*sql.DB, 16)
	orms = make(map[string]*gorm.DB, 16)
	for _, config := range configs {
		db, err := NewSQL(config)
		if err != nil {
			Finally()
			return err
		}
		sqls[config.Id] = db

		if !config.EnableORM {
			continue
		}

		orm, err := NewORM(config, db)
		if err != nil {
			Finally()
			return err
		}
		orms[config.Id] = orm
	}

	return nil
}

func NewSQL(config Config) (*sql.DB, error) {
	db, err := sql.Open(config.Type, config.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxLifetime)

	return db, nil
}

func NewORM(config Config, db *sql.DB) (*gorm.DB, error) {
	if db == nil {
		var err error
		db, err = NewSQL(config)
		if err != nil {
			return nil, err
		}
	}

	switch config.Type {
	case DBTypeMySQL:
		return gorm.Open(mysql.New(mysql.Config{Conn: db}), config.GORMConfig())
	case DBTypePostgres:
		return gorm.Open(postgres.New(postgres.Config{Conn: db}), config.GORMConfig())
	default:
		return nil, fmt.Errorf("not support db type [%s]", config.Type)
	}
}

func Finally() {
	for _, db := range sqls {
		_ = db.Close()
	}
}

func GetConfig(id string) Config {
	return configs[id]
}

func GetDefaultConfig() Config {
	return GetConfig(DefaultId)
}

func GetDB(id string) *sql.DB {
	return sqls[id]
}

func GetDefaultDB() *sql.DB {
	return GetDB(DefaultId)
}

func GetORM(id string) *gorm.DB {
	return orms[id]
}

func GetDefaultORM() *gorm.DB {
	return GetORM(DefaultId)
}
