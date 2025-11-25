package items

import (
	"UltimateDesktopPet/internal/database"
	"UltimateDesktopPet/pkg/file"
	pp "UltimateDesktopPet/pkg/print"
	"fmt"
)

type ItemsMeta struct {
	Controller *database.BaseController[Item]
	DB         database.DB
	Item       *Item
	ST         file.SpriteTool
}

const itemStaticAssetPath = "assets/itemImages"
const itemDefaultImageFolder = "default"

func init() {
	i := newItemsController(nil)
	database.RegisterSchema(database.StaticAssets, i)
	pp.Assert(pp.Items, "items init complete")
}

func newItemsController(model **Item) *database.BaseController[Item] {
	return &database.BaseController[Item]{Model: model}
}

func NewItemMeta() *ItemsMeta {
	i := &ItemsMeta{}
	i.Item = &Item{}
	i.Controller = newItemsController(&i.Item)
	i.ST = file.NewSpriteTool(i.ST)
	i.ST.StaticAssetPath = itemStaticAssetPath
	i.ST.DefaultImageFolder = itemDefaultImageFolder
	return i
}

func (i *ItemsMeta) Shutdown() {
	i.DB.CloseDB()
	pp.Assert(pp.Items, "item service stopped")
}

func (i *ItemsMeta) LoadAll() ([]Item, error) {
	items, err := i.Controller.ReadAll(i.DB.GetDB())
	if items == nil {
		return nil, fmt.Errorf("LoadAll return nil")
	}

	validItems := make([]Item, 0, len(*items))
	for _, item := range *items {
		validItems = append(validItems, item)
	}
	return validItems, err
}

func (i *ItemsMeta) LoadByID(id uint) (Item, error) {
	item, err := i.Controller.Read(i.DB.GetDB(), id)
	if err != nil {
		return Item{}, err
	}
	return *item, nil
}

func (i *ItemsMeta) LoadFrameByID(id uint) (string, error) {
	item, err := i.Controller.Read(i.DB.GetDB(), id)
	if err != nil {
		return "", err
	}
	return i.ST.LoadFrameFromDir(item.Path)
}
