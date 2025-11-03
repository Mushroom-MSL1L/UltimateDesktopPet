package pet

import (
	pp "UltimateDesktopPet/pkg/print"
	"context"
	"sync"
	"time"
)

type Pet struct {
	sync.Mutex
	ID     uint  `json:"id" gorm:"primaryKey;autoIncrement"`
	Water  int16 `json:"water"`
	Hunger int16 `json:"hunger"`
	Health int16 `json:"health"`
	Mood   int16 `json:"mood"`
	Energy int16 `json:"energy"`
	Money  int16 `json:"money"`
}

const (
	Max       int16 = 100
	WaterMax        = Max
	HungerMax       = Max
	HealthMax       = Max
	MoodMax         = Max
	EnergyMax       = Max

	Min       int16 = 0
	WaterMin        = Min
	HungerMin       = Min
	HealthMin       = Min
	MoodMin         = Min
	EnergyMin       = Min
)

func (p *Pet) periodicallyUpdateStates(c context.Context) {
	const waitTime = 5 * time.Second
	for {
		select {
		case <-c.Done():
			return
		case <-time.After(waitTime):
			p.Lock()
			p.Water = max(WaterMin, p.Water-1)
			p.Hunger = max(HungerMin, p.Hunger-1)
			p.Mood = max(MoodMin, p.Mood-1)
			p.Energy = min(EnergyMax, p.Energy+1)
			p.Unlock()
		}
	}
}

func (p *Pet) periodicallyPrintStatus(c context.Context) {
	const waitTime = 2 * time.Second
	for {
		select {
		case <-c.Done():
			return
		case <-time.After(waitTime):
			p.printStatus()
		}
	}
}

func (p *Pet) printStatus() {
	p.Lock()
	pp.Info(pp.Pet, "Pet Status - Water: %d, Hunger: %d, Health: %d, Mood: %d, Energy: %d, Money: %d",
		p.Water, p.Hunger, p.Health, p.Mood, p.Energy, p.Money)
	p.Unlock()
}
