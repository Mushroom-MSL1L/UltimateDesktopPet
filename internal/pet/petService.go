package pet

import (
	"UltimateDesktopPet/internal/attributes"
	pp "UltimateDesktopPet/pkg/print"
	"fmt"
	"time"
)

/* PerformActivity to be modified */
func (p *PetMeta) PerformActivity(name string, attr attributes.Attributes, durationMinutes int16) error {
	if p.Status == Acting {
		return fmt.Errorf("PerformActivity: Pet already acting %s", p.StatusDetail)
	}
	p.Status = Acting
	p.StatusDetail = name
	p.DoneSignal = make(chan struct{})

	starting := time.Now()
	duration := time.Duration(durationMinutes) * time.Minute
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-p.DoneSignal:
				p.Status = Idle
				p.StatusDetail = ""
				return
			case <-ticker.C:
				elapsed := time.Since(starting)
				pp.Info(pp.Pet, "Acting %s for %v", name, elapsed)
				if elapsed >= duration {
					p.Status = Idle
					p.StatusDetail = ""
					pp.Info(pp.Pet, "Activity %s completed", name)
					return
				}
			}
		}
	}()
	return nil
}

func (p *PetMeta) StopActivity() {
	if p.Status != Acting {
		return
	}
	p.DoneSignal <- struct{}{}
	close(p.DoneSignal)
}

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
