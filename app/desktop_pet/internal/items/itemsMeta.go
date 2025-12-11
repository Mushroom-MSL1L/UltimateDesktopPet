package items

import (
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/database"
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/file"
	pp "github.com/Mushroom-MSL1L/UltimateDesktopPet/pkg/print"
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
