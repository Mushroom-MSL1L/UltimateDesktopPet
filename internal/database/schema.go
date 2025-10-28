package database

import (
	"gorm.io/gorm"
)

type Schema interface {
	InitTable(db *gorm.DB)
}

var registerSchemas []Schema

func RegisterSchema(s Schema) {
	registerSchemas = append(registerSchemas, s)
}

func (d *DB) loadSchemas() {
	for _, schema := range registerSchemas {
		schema.InitTable(d.db)
	}
}
