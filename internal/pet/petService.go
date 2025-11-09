package pet

import (
	"UltimateDesktopPet/internal/database"
	pp "UltimateDesktopPet/pkg/print"
	"context"
	"errors"

	"gorm.io/gorm"
)

type PetMeta struct {
	Controller *database.BaseController[Pet]
	Pet        *Pet
}

func init() {
	p := newPetController(nil)
	database.RegisterSchema(database.Pets, p)
	pp.Assert(pp.Pet, "pet init complete")
}

func newPetController(model **Pet) *database.BaseController[Pet] {
	return &database.BaseController[Pet]{Model: model}
}

func (p *PetMeta) Service(c context.Context, db *gorm.DB) {
	p.petServiceInit(db)

	go p.Pet.periodicallyUpdateStates(c)
	go p.Pet.periodicallyPrintStatus(c)
}

func (p *PetMeta) Shutdown(db *gorm.DB) {
	err := p.Controller.Create(db)
	if err != nil {
		pp.Warn(pp.Pet, "failed to save pet state: %v", err)
	} else {
		pp.Info(pp.Pet, "pet state saved successfully")
	}
	pp.Assert(pp.Pet, "pet service stopped")
}

func (p *PetMeta) petServiceInit(db *gorm.DB) {
	var err error
	p.Pet = &Pet{}
	p.Controller = newPetController(&p.Pet)

	p.Pet, err = p.Controller.ReadFirst(db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// create a default pet when none exists
			defaultPet := &Pet{
				Water:  WaterMax,
				Hunger: HungerMax,
				Health: HealthMax,
				Mood:   MoodMax,
				Energy: EnergyMax,
				Money:  0,
			}
			if err := db.Create(defaultPet).Error; err != nil {
				pp.Fatal(pp.Pet, "failed to create default pet: %v", err)
			}
			pp.Info(pp.Pet, "created default pet")
			p.Pet = defaultPet
			return
		}
		pp.Fatal(pp.Pet, "Read first pet entry failed: %v", err)
	}
}
