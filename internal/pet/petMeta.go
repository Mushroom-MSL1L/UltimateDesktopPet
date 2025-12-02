package pet

import (
	"UltimateDesktopPet/internal/activities"
	"UltimateDesktopPet/internal/attributes"
	"UltimateDesktopPet/internal/database"
	"UltimateDesktopPet/internal/items"
	"UltimateDesktopPet/pkg/file"
	pp "UltimateDesktopPet/pkg/print"
	"context"
	"errors"

	"gorm.io/gorm"
)

type PetStatus int8

const (
	Idle PetStatus = iota
	Acting

	petStaticAssetPath      = "assets/petImages"
	petDefaultImageFolder   = "default"
	defaultPetAnimationType = "stand"
)

type PetMeta struct {
	Controller    *database.BaseController[Pet]
	DB            database.DB
	Pet           *Pet
	ST            file.SpriteTool
	ItemUsing     *items.ItemsMeta
	ActivityDoing *activities.ActivityMeta
	Status        PetStatus
	StatusDetail  string
	DoneSignal    chan struct{}
}

func init() {
	p := newPetController(nil)
	database.RegisterSchema(database.Pets, p)
	pp.Assert(pp.Pet, "pet init complete")
}

func newPetController(model **Pet) *database.BaseController[Pet] {
	return &database.BaseController[Pet]{Model: model}
}

func NewPetMeta(newDB database.DB, itemMeta *items.ItemsMeta, activityMeta *activities.ActivityMeta) *PetMeta {
	p := &PetMeta{
		DB:            newDB,
		ItemUsing:     itemMeta,
		ActivityDoing: activityMeta,
	}

	var err error
	p.Pet = &Pet{}
	p.Controller = newPetController(&p.Pet)
	p.Status = Idle
	p.ST = file.NewSpriteTool(p.ST)
	p.ST.StaticAssetPath = petStaticAssetPath
	p.ST.DefaultImageFolder = petDefaultImageFolder
	db := p.DB.GetDB()

	p.Pet, err = p.Controller.ReadFirst(db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// create a default pet when none exists
			defaultPet := &Pet{
				Attributes: attributes.Attributes{
					Experience: 0,
					Money:      100000,
					Water:      Max,
					Hunger:     Max,
					Health:     Max,
					Mood:       Max,
					Energy:     Max,
				},
			}
			if err := db.Create(defaultPet).Error; err != nil {
				pp.Fatal(pp.Pet, "failed to create default pet: %v", err)
			}
			pp.Info(pp.Pet, "created default pet")
			p.Pet = defaultPet
			return p
		}
		pp.Fatal(pp.Pet, "Read first pet entry failed: %v", err)
	}
	return p
}

func (p *PetMeta) Service(c context.Context) {
	go p.Pet.periodicallyUpdateStates(c)
	go p.Pet.periodicallyPrintStatus(c)
}

func (p *PetMeta) Shutdown() {
	p.Pet.Lock()
	p.storePet()
	p.Pet.Unlock()

	p.DB.CloseDB()
	pp.Assert(pp.Pet, "pet service stopped")
}

func (p *PetMeta) storePet() {
	/* Caller should lock the Pet struct */
	db := p.DB.GetDB()
	err := p.Controller.Create(db)

	if err != nil {
		pp.Warn(pp.Pet, "failed to save pet state: %v", err)
	} else {
		pp.Info(pp.Pet, "pet state saved successfully")
	}
}

func (p *PetMeta) GetPetStatus() Pet {
	return (*p.Pet).getStatus()
}

func (p *PetMeta) UpdateStatus(attr attributes.Attributes) {
	/* Caller should lock the Pet struct */
	(*p.Pet).updateStatus(attr)
	p.storePet()
}
