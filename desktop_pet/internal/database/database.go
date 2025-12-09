package database

import (
	"context"
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
