package database

import (
	pp "UltimateDesktopPet/pkg/print"
	"errors"
	"reflect"

	"gorm.io/gorm"
)

type BaseController[T any] struct {
	Model **T
}

func (b *BaseController[T]) InitTable(db *gorm.DB) {
	if err := db.AutoMigrate(b.Model); err != nil {
		pp.Fatal(pp.DB, "Init table %T error: %v", b.Model, err)
	}
}

func (b *BaseController[T]) Create(db *gorm.DB) error {
	if b.Model == nil || *b.Model == nil {
		return errors.New("model not set")
	}
	reflect.ValueOf(*b.Model).Elem().FieldByName("ID").SetUint(0)
	result := db.Model(*b.Model).Create(*b.Model)
	return result.Error
}

func (b *BaseController[T]) Read(db *gorm.DB, id uint) (*T, error) {
	var instance T
	if err := db.First(&instance, id).Error; err != nil {
		return nil, err
	}
	*b.Model = &instance
	return *b.Model, nil

}

func (b *BaseController[T]) ReadFirst(db *gorm.DB) (*T, error) {
	var instance T
	if err := db.Order("id DESC").First(&instance).Error; err != nil {
		return nil, err
	}
	*b.Model = &instance
	return *b.Model, nil
}
