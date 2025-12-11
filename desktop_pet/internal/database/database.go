package database

import (
	"context"
	"os"
	"time"

	pkg "github.com/Mushroom-MSL1L/UltimateDesktopPet/desktop_pet/pkg/file"
	pp "github.com/Mushroom-MSL1L/UltimateDesktopPet/desktop_pet/pkg/print"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	db     *gorm.DB
	dbFile string
}

func (d *DB) InitDB(c context.Context, dbFile string, id DBID) {
	pp.Info(pp.DB, "DB initializing")
	d.setdbFile(dbFile)
	d.connectDB()
	d.loadSchemas(id)
	pp.Assert(pp.DB, "DB working")
}

func (d *DB) setdbFile(dbFile string) {
	pkg.SafeCreateFile(dbFile)
	d.dbFile = dbFile
}

func (d *DB) connectDB() {
	var err error
	for retry := 1; retry <= 5; retry++ {
		pp.Info(pp.DB, "DB %ded trying", retry)
		d.db, err = gorm.Open(sqlite.Open(d.dbFile), &gorm.Config{})
		if err == nil {
			pp.Assert(pp.DB, "DB connected")
			return
		}
		waitTime := time.Second
		time.Sleep(waitTime)
		pp.Info(pp.DB, "DB %ded not response, wait %s seconds", waitTime)
	}
	if err != nil {
		pp.Fatal(pp.DB, "failed to open database after retries: %v", err)
	}
}

func (d *DB) GetDB() *gorm.DB {
	return d.db
}

func (d *DB) LoadSQLFileIfEmpty(sqlFile string) {
	pp.Info(pp.DB, "Checking all tables before loading SQL file")

	var tableNames []string
	if err := d.db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';").Scan(&tableNames).Error; err != nil {
		pp.Fatal(pp.DB, "failed to list tables: %v", err)
		return
	}
	for _, table := range tableNames {
		var count int64
		if err := d.db.Table(table).Count(&count).Error; err != nil {
			pp.Fatal(pp.DB, "failed to count table %s: %v", table, err)
			return
		}
		if count > 0 {
			pp.Info(pp.DB, "Table %s has %d rows, skipping SQL load", table, count)
			return
		}
	}

	pp.Info(pp.DB, "All tables empty, loading SQL file: %s", sqlFile)
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		pp.Fatal(pp.DB, "failed to read SQL file: %v", err)
		return
	}
	sqlStatements := string(sqlBytes)

	result := d.db.Exec(sqlStatements)
	if result.Error != nil {
		pp.Fatal(pp.DB, "failed to execute SQL statements: %v", result.Error)
		return
	}
	pp.Assert(pp.DB, "SQL file loaded successfully")
}

func (d *DB) CloseDB() {
	if d.db == nil {
		pp.Info(pp.DB, "DB already closed")
		return
	}
	sqlDB, err := d.db.DB()
	if err != nil {
		pp.Fatal(pp.DB, "failed to get sql DB: %v", err)
		return
	}
	sqlDB.Close()
	pp.Assert(pp.DB, "DB closed")
}
