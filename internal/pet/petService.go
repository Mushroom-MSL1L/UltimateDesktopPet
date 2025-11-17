package pet

import (
	"UltimateDesktopPet/internal/database"
	"UltimateDesktopPet/pkg/file"
	pp "UltimateDesktopPet/pkg/print"
	"context"
	"errors"
	"fmt"
	"time"

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
	Controller   *database.BaseController[Pet]
	DB           database.DB
	Pet          *Pet
	ST           file.SpriteTool
	Status       PetStatus
	StatusDetail string
	DoneSignal   chan struct{}
}

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
				Experience: 0,
				Money:      100000,
				Water:      Max,
				Hunger:     Max,
				Health:     Max,
				Mood:       Max,
				Energy:     Max,
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

func (p *PetMeta) UpdateStatus(expr int, water, hunger, health, mood, energy int16, money int) {
	(*p.Pet).updateStatus(expr, water, hunger, health, mood, energy, money)
	p.storePet()
}

func (p *PetMeta) PerformActivity(name string, expr int, water, hunger, health, mood, energy int16, money int, durationMinutes int16) error {
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
