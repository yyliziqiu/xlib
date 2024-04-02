package xdb

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/yyliziqiu/xlib/xlog"
)

type TablesMigration struct {
	DB       *gorm.DB
	Once     []schema.Tabler
	Cron     []schema.Tabler
	Interval time.Duration
}

func MigrateTablesBatch(ctx context.Context, migrations func() []TablesMigration) (err error) {
	for _, migration := range migrations() {
		err = MigrateTables(ctx, migration)
		if err != nil {
			return err
		}
	}
	return nil
}

func MigrateTables(ctx context.Context, migration TablesMigration) (err error) {
	db := migration.DB.Set("gorm:table_options", "ENGINE=InnoDB")

	err = migrateTables(db, migration.Once)
	if err != nil {
		return fmt.Errorf("migrate DB error [%v]", err)
	}

	if len(migration.Cron) == 0 {
		return nil
	}

	err = migrateTables(db, migration.Cron)
	if err != nil {
		return fmt.Errorf("migrate tables error [%v]", err)
	}

	go runCronMigrateTables(ctx, migration.Interval, db, migration.Cron)

	return nil
}

func migrateTables(db *gorm.DB, tables []schema.Tabler) error {
	for _, table := range tables {
		err := db.Table(table.TableName()).Migrator().AutoMigrate(&table)
		if err != nil {
			return fmt.Errorf("create table [%s] failed [%v]", table.TableName(), err)
		}
	}
	return nil
}

func runCronMigrateTables(ctx context.Context, interval time.Duration, db *gorm.DB, tables []schema.Tabler) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			err := migrateTables(db, tables)
			if err != nil {
				xlog.Errorf("Migrate table failed, error: %v.", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

type RecordMigration interface {
	Exist() (bool, error)
	Create() error
}

func MigrateRecords(migrations []RecordMigration) error {
	for _, migration := range migrations {
		exist, err := migration.Exist()
		if err != nil {
			return err
		}
		if exist {
			continue
		}
		err = migration.Create()
		if err != nil {
			return err
		}
	}
	return nil
}
