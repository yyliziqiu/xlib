package xpgsql

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/yyliziqiu/xlib/xutil"
)

var dbs map[string]*sql.DB

func Init(configs ...Config) error {
	dbs = make(map[string]*sql.DB, len(configs))
	for _, config := range configs {
		db, err := New(config)
		if err != nil {
			Finally()
			return err
		}
		dbs[xutil.IES(config.Id, DefaultId)] = db
	}
	return nil
}

func New(config Config) (*sql.DB, error) {
	config = config.WithDefault()

	db, err := sql.Open("postgres", config.DSN)
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

func NewWithDSN(dsn string) (*sql.DB, error) {
	return New(Config{DSN: dsn})
}

func GetDB(id string) *sql.DB {
	return dbs[id]
}

func GetDefaultDB() *sql.DB {
	return GetDB(DefaultId)
}
