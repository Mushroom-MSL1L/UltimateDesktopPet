package items

import (
	"fmt"

	pp "github.com/Mushroom-MSL1L/UltimateDesktopPet/pkg/print"
)

func (i *ItemsMeta) LoadAllItems() ([]Item, error) {
	items, err := i.Controller.ReadAll(i.DB.GetDB())
	if items == nil {
		pp.Warn(pp.Items, "LoadAllItems: failed to load all items: %v", err)
		return nil, fmt.Errorf("LoadAll return nil")
	}

	validItems := make([]Item, 0, len(*items))
	for _, item := range *items {
		validItems = append(validItems, item)
	}
	return validItems, nil
}

func (i *ItemsMeta) LoadItemByID(id uint) (Item, error) {
	item, err := i.Controller.Read(i.DB.GetDB(), id)
	if err != nil {
		return Item{}, err
	}
	return *item, nil
}

func (i *ItemsMeta) LoadItemFrameByID(id uint) (string, error) {
	item, err := i.Controller.Read(i.DB.GetDB(), id)
	if err != nil {
		return "", err
	}
	frame, err := i.ST.LoadFrameFromDir(item.Path)
	if err != nil {
		pp.Warn(pp.Items, "LoadItemFrameByID: failed to load frame for item ID %d: %v", id, err)
		return "", err
	}
	return frame, nil

}
