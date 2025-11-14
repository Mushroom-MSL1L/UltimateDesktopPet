package items

import (
	"UltimateDesktopPet/internal/database"
	pp "UltimateDesktopPet/pkg/print"
	"context"
	"fmt"
)

type ItemsMeta struct {
	Controller *database.BaseController[Item]
	ImagePath  string
	DB         database.DB
	Item       *Item
}

func init() {
	i := newItemsController(nil)
	database.RegisterSchema(database.StaticAssets, i)
	pp.Assert(pp.Items, "items init complete")
}

func newItemsController(model **Item) *database.BaseController[Item] {
	return &database.BaseController[Item]{Model: model}
}

func (i *ItemsMeta) Service(c context.Context) {
	i.Item = &Item{}
	i.Controller = newItemsController(&i.Item)
}

func (i *ItemsMeta) Shutdown() {
	i.DB.CloseDB()
	pp.Assert(pp.Items, "item service stopped")
}

func (i *ItemsMeta) LoadAll() (*[]Item, error) {
	items, err := i.Controller.ReadAll(i.DB.GetDB())
	if items == nil {
		return nil, fmt.Errorf("LoadAll return nil")
	}
	return items, err
}
