package pet

import (
	"UltimateDesktopPet/pkg/math"
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
	Money  int `json:"money"`
}

const (
	Max       int16 = 100
	Min       int16 = 0
)

func (p *Pet) periodicallyUpdateStates(c context.Context) {
	const waitTime = 5 * time.Second
	for {
		select {
		case <-c.Done():
			return
		case <-time.After(waitTime):
			p.Lock()
			p.Water = max(Min, p.Water-1)
			p.Hunger = max(Min, p.Hunger-1)
			p.Mood = max(Min, p.Mood-1)
			p.Energy = min(Max, p.Energy+1)
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

func (p *Pet) getStatus() Pet {
	p.Lock()
	defer p.Unlock()
	return Pet{
		Water: p.Water,
		Hunger: p.Hunger,
		Health: p.Health,
		Mood: p.Mood,
		Energy: p.Energy,
		Money: p.Money,
	}
}

func (p *Pet) updateStatus(water, hunger, health, mood, energy int16, money int) {
	p.Lock()
	defer p.Unlock()
	p.Water  = math.InRange(p.Water+water, Max, Min)
	p.Hunger = math.InRange(p.Hunger+hunger, Max, Min)
	p.Health = math.InRange(p.Health+health, Max, Min)
	p.Mood   = math.InRange(p.Mood+mood, Max, Min)
	p.Energy = math.InRange(p.Energy+energy, Max, Min)
	p.Money = p.Money+money
} 