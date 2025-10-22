package database

import (
	"log"

	"gorm.io/gorm"
)

type BaseCRUD[T any] struct {
	Model *T
}

func (b *BaseCRUD[T]) InitTable(db *gorm.DB) {
	if err := db.AutoMigrate(b.Model); err != nil {
		log.Fatalf("Init table %T error: %v", b.Model, err)
	}
}

func (b *BaseCRUD[T]) Create(db *gorm.DB) error {
	result := db.Model(b.Model).Create(b.Model)
	return result.Error
}

func (b *BaseCRUD[T]) Read(db *gorm.DB, id uint) (*T, error) {
	var data T
	result := db.Model(b.Model).First(&data, id)
	return &data, result.Error
}
