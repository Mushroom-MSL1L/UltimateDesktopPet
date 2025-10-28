package database

import (
	pp "UltimateDesktopPet/pkg/print"

	"gorm.io/gorm"
)

type BaseController[T any] struct {
	Model *T
}

func (b *BaseController[T]) InitTable(db *gorm.DB) {
	if err := db.AutoMigrate(b.Model); err != nil {
		pp.Fatal(pp.DB, "Init table %T error: %v", b.Model, err)
	}
}

func (b *BaseController[T]) Create(db *gorm.DB) error {
	result := db.Model(b.Model).Create(b.Model)
	return result.Error
}

func (b *BaseController[T]) Read(db *gorm.DB, id uint) (*T, error) {
	var data T
	result := db.Model(b.Model).First(&data, id)
	return &data, result.Error
}

func (b *BaseController[T]) ReadFirst(db *gorm.DB) (*T, error) {
	var data T
	result := db.Model(b.Model).First(&data)
	return &data, result.Error
}
