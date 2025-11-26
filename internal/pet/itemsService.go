package pet

import (
	"fmt"

	pp "UltimateDesktopPet/pkg/print"
)

func (p *PetMeta) UseItemByID(itemID uint) error {
	item, err := p.ItemUsing.LoadItemByID(itemID)
	if err != nil {
		pp.Warn(pp.Items, "UseItem: failed to load item by ID %d: %v", itemID, err)
		return err
	}
	pp.Info(pp.Items, "UseItem : id %d, name %s", itemID, item.Name)

	p.Pet.Lock()
	oldPet := p.GetPetStatus()
	if (oldPet.Money + item.Money) < 0 {
		errMessage := fmt.Sprintf("UseItem: pet has not enough money to use %s", item.Name)
		pp.Warn(pp.Pet, errMessage)
		p.Pet.Unlock()

		pp.Warn(pp.Items, "UseItemByID: failed to use item ID %d: %v", itemID, err)
		return fmt.Errorf(errMessage)
	}
	p.UpdateStatus(item.Attributes)
	p.Pet.Unlock()
	pp.Info(pp.Pet, "UseItem: pet use \"%s\" and update status", item.Name)

	return nil
}
