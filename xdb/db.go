package xdb

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	_configs map[string]Config
	_sqlDBS  map[string]*sql.DB
	_ormDBS  map[string]*gorm.DB
)

func Init(configs ...Config) error {
	_configs = make(map[string]Config, 8)
	for _, config := range configs {
		_configs[config.Id] = config.Default()
	}

	_sqlDBS = make(map[string]*sql.DB, 8)
	_ormDBS = make(map[string]*gorm.DB, 8)
	for _, config := range _configs {
		db, err := NewSQL(config)
		if err != nil {
			Finally()
			return err
		}
		_sqlDBS[config.Id] = db

		if !config.EnableORM {
			continue
		}

		orm, err := NewORM(config, db)
		if err != nil {
			Finally()
			return err
		}
		_ormDBS[config.Id] = orm
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
		return gorm.Open(mysql.New(mysql.Config{Conn: db}), config.ORMConfig())
	case DBTypePostgres:
		return gorm.Open(postgres.New(postgres.Config{Conn: db}), config.ORMConfig())
	default:
		return nil, fmt.Errorf("not support db type [%s]", config.Type)
	}
}

func Finally() {
	for _, db := range _sqlDBS {
		_ = db.Close()
	}
}

func GetConfig(id string) Config {
	return _configs[id]
}

func GetDefaultConfig() Config {
	return GetConfig(DefaultId)
}

func GetDB(id string) *sql.DB {
	return _sqlDBS[id]
}

func GetDefaultDB() *sql.DB {
	return GetDB(DefaultId)
}

func GetORMDB(id string) *gorm.DB {
	return _ormDBS[id]
}

func GetDefaultORMDB() *gorm.DB {
	return GetORMDB(DefaultId)
}
