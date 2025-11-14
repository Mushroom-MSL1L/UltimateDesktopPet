package items

import (
	"UltimateDesktopPet/internal/database"
	"UltimateDesktopPet/pkg/file"
	pp "UltimateDesktopPet/pkg/print"
	"context"
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

func (i *ItemsMeta) Service(c context.Context) {
	i.Item = &Item{}
	i.Controller = newItemsController(&i.Item)
	i.ST = file.NewSpriteTool(i.ST)
	i.ST.StaticAssetPath = itemStaticAssetPath
	i.ST.DefaultImageFolder = itemDefaultImageFolder
}

func (i *ItemsMeta) Shutdown() {
	i.DB.CloseDB()
	pp.Assert(pp.Items, "item service stopped")
}

func (i *ItemsMeta) LoadAll() ([]ItemWithFrame, error) {
	items, err := i.Controller.ReadAll(i.DB.GetDB())
	if items == nil {
		return nil, fmt.Errorf("LoadAll return nil")
	}

	validItems := make([]ItemWithFrame, 0, len(*items))
	for _, item := range *items {
		frame, err := i.ST.LoadFrameFromDir(item.Path)
		if err != nil {
			pp.Warn(pp.Items, "%v", err)
			continue
		}
		var temp ItemWithFrame
		temp.Frame = frame
		temp.Item = item
		validItems = append(validItems, temp)
	}
	return validItems, err
}
