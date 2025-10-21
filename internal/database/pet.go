package database

import (
	"gorm.io/gorm"
)

type SystemAccounting struct {
	BeginTime string `json:"begin_time"`
}

type Pet struct {
	Water     string `json:"water"`
	Hunger    string `json:"hunger"`
	Money     string `json:"money"`
	Health    string `json:"health"`
	Happiness string `json:"happiness"`
}

func initPetTable(db *gorm.DB) {
	db.AutoMigrate(&SystemAccounting{})
	db.AutoMigrate(&Pet{})
}
