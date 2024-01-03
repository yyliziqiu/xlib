package xsql

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Migration struct {
	DB       *gorm.DB
	Once     []schema.Tabler
	Cron     []schema.Tabler
	Interval time.Duration
}

func Migrates(ctx context.Context, migrations func() []Migration) (err error) {
	for _, migration := range migrations() {
		err = Migrate(ctx, migration)
		if err != nil {
			return err
		}
	}
	return nil
}

func Migrate(ctx context.Context, migration Migration) (err error) {
	db := migration.DB.Set("gorm:table_options", "ENGINE=InnoDB")

	err = migrateTables(db, migration.Once)
	if err != nil {
		logger.Errorf("Migrate DB failed, error: %s.", err)
		return err
	}

	if len(migration.Cron) == 0 {
		return nil
	}

	err = migrateTables(db, migration.Cron)
	if err != nil {
		logger.Errorf("Migrate DB failed, error: %s.", err)
		return err
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
				logger.Errorf("Migrate DB failed, error: %s.", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
