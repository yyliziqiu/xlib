package xsql

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	globalLogger *logrus.Logger

	cfs  map[string]Config
	sqls map[string]*sql.DB
	orms map[string]*gorm.DB
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

	sqls = make(map[string]*sql.DB, len(cfs))
	orms = make(map[string]*gorm.DB, len(cfs))

	for _, cf := range cfs {
		db, err := NewSQL(cf)
		if err != nil {
			Finally()
			return err
		}
		sqls[cf.Id] = db

		if !cf.EnableORM {
			continue
		}

		orm, err := NewORM(cf, db)
		if err != nil {
			Finally()
			return err
		}
		orms[cf.Id] = orm
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
	switch config.Type {
	case DBTypeMySQL:
		return gorm.Open(mysql.New(mysql.Config{Conn: db}), config.GORMConfig())
	case DBTypePostgres:
		return gorm.Open(postgres.New(postgres.Config{Conn: db}), config.GORMConfig())
	default:
		return nil, errors.New("not support db type")
	}
}

func Finally() {
	for _, db := range sqls {
		_ = db.Close()
	}
}

func GetConfig(id string) Config {
	return cfs[id]
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
