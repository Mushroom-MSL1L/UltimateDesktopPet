package database

import (
	"gorm.io/gorm"
)

type DBID int

const (
	Activities DBID = iota
	Items
	Pets
)

type Schema interface {
	InitTable(db *gorm.DB)
}

var registerSchemas = make(map[DBID][]Schema)

func RegisterSchema(id DBID, s Schema) {
	registerSchemas[id] = append(registerSchemas[id], s)
}

func (d *DB) loadSchemas(id DBID) {
	for _, schema := range registerSchemas[id] {
		schema.InitTable(d.db)
	}
}
