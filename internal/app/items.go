package app

import (
	"UltimateDesktopPet/internal/items"
	"fmt"

	pp "UltimateDesktopPet/pkg/print"
)

func (a *App) LoadAllItems() ([]items.Item, error) {
	items, err := a.itemsMeta.LoadAll()
	if err != nil {
		pp.Warn(pp.App, "LoadAllItems: failed to load all items: %v", err)
		return nil, err
	}
	return items, err
}

func (a *App) LoadItemFrameByID(id uint) (string, error) {
	frame, err := a.itemsMeta.LoadFrameByID(id)
	if err != nil {
		pp.Warn(pp.App, "LoadItemFrameByID: failed to load frame for item ID %d: %v", id, err)
		return "", err
	}
	return frame, nil
}

func (a *App) UseItem(id uint) error {
	item, err := a.itemsMeta.LoadByID(id)
	if err != nil {
		pp.Warn(pp.App, "UseItem: failed to load item by ID %d: %v", id, err)
		return err
	}
	pp.Info(pp.App, "UseItem : id %d, name %s", id, item.Name)

	oldPet := a.petMeta.GetPetStatus()
	if (oldPet.Money + item.MoneyCost) < 0 {
		errMessage := fmt.Sprintf("UseItem: pet has not enough money to use %s", item.Name)
		pp.Warn(pp.App, errMessage)
		return fmt.Errorf(errMessage)
	}

	a.petMeta.UpdateStatus(item.Experience, item.Water, item.Hunger, item.Health, item.Mood, item.Energy, item.MoneyCost)
	pp.Info(pp.App, "UseItem: pet use \"%s\" and update status", item.Name)
	return nil
}
