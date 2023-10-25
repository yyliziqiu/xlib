package xsql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var (
	cfs map[string]Config
	dbs map[string]*sql.DB
)

func Init(configs ...Config) error {
	cfs = make(map[string]Config, len(configs))
	for _, config := range configs {
		config = config.WithDefault()
		cfs[config.Id] = config
	}

	dbs = make(map[string]*sql.DB, len(cfs))
	for _, cf := range cfs {
		db, err := New(cf)
		if err != nil {
			Finally()
			return err
		}
		dbs[cf.Id] = db
	}

	return nil
}

func New(config Config) (*sql.DB, error) {
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

func Finally() {
	for _, db := range dbs {
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
	return dbs[id]
}

func GetDefaultDB() *sql.DB {
	return GetDB(DefaultId)
}
