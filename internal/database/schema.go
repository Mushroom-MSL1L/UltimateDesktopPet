package database

import "gorm.io/gorm"

type Schema interface {
	InitTable(db *gorm.DB)
}

var registerSchema []Schema

func RegisterSchema(s Schema) {
	registerSchema = append(registerSchema, s)
}

func (d *DB) loadSchema() {
	for _, schema := range registerSchema {
		schema.InitTable(d.db)
	}
}
