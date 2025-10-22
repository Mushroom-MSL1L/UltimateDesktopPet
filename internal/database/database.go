package database

import (
	"log"
	"time"

	pkg "UltimateDesktopPet/pkg/file"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	db     *gorm.DB
	dbFile string
}

func (d *DB) InitDB(dbFile string) {
	log.Println("database init")
	d.setdbFile(dbFile)
	d.connectDB()
	d.loadSchema()
}

func (d *DB) setdbFile(dbFile string) {
	pkg.SafeCreateFile(dbFile)
	d.dbFile = dbFile
}

func (d *DB) connectDB() {
	var err error
	for retry := 0; retry < 5; retry++ {
		d.db, err = gorm.Open(sqlite.Open(d.dbFile), &gorm.Config{})
		if err == nil {
			return
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Fatalf("failed to open database after retries: %v", err)
	}
}

func (d *DB) CloseDB() {
	if d.db == nil {
		return
	}
	sqlDB, err := d.db.DB()
	if err != nil {
		log.Printf("failed to get sql DB: %v", err)
		return
	}
	sqlDB.Close()
}
