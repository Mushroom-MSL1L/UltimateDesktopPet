package pet

import (
	"UltimateDesktopPet/internal/database"
	"UltimateDesktopPet/pkg/file"
	pp "UltimateDesktopPet/pkg/print"
	"context"
	"errors"

	"gorm.io/gorm"
)

type PetMeta struct {
	Controller *database.BaseController[Pet]
	DB         database.DB
	Pet        *Pet
	ST         file.SpriteTool
}

const petStaticAssetPath = "assets/petImages"
const petDefaultImageFolder = "default"
const defaultPetAnimationType = "stand"

func init() {
	p := newPetController(nil)
	database.RegisterSchema(database.Pets, p)
	pp.Assert(pp.Pet, "pet init complete")
}

func newPetController(model **Pet) *database.BaseController[Pet] {
	return &database.BaseController[Pet]{Model: model}
}

func (p *PetMeta) Service(c context.Context) {
	p.petServiceInit()

	go p.Pet.periodicallyUpdateStates(c)
	go p.Pet.periodicallyPrintStatus(c)
}

func (p *PetMeta) Shutdown() {
	p.storePet()
	p.DB.CloseDB()
	pp.Assert(pp.Pet, "pet service stopped")
}

func (p *PetMeta) petServiceInit() {
	var err error
	p.Pet = &Pet{}
	p.Controller = newPetController(&p.Pet)
	p.ST = file.NewSpriteTool(p.ST)
	p.ST.StaticAssetPath = petStaticAssetPath
	p.ST.DefaultImageFolder = petDefaultImageFolder
	db := p.DB.GetDB()

	p.Pet, err = p.Controller.ReadFirst(db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// create a default pet when none exists
			defaultPet := &Pet{
				Money:  0,
				Water:  Max,
				Hunger: Max,
				Health: Max,
				Mood:   Max,
				Energy: Max,
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

func (p *PetMeta) storePet() {
	(*p.Pet).Lock()
	db := p.DB.GetDB()
	err := p.Controller.Create(db)

	if err != nil {
		pp.Warn(pp.Pet, "failed to save pet state: %v", err)
	} else {
		pp.Info(pp.Pet, "pet state saved successfully")
	}
	(*p.Pet).Unlock()
}

func (p *PetMeta) LoadDefaultFrames() ([]string, error) {
	return p.ST.LoadFramesFromDir(defaultPetAnimationType)
}

func (p *PetMeta) GetPetStatus() Pet {
	return (*p.Pet).getStatus()
}

func (p *PetMeta) UpdateStatus(water, hunger, health, mood, energy int16, money int) {
	(*p.Pet).updateStatus(water, hunger, health, mood, energy, money)
	p.storePet()
}
