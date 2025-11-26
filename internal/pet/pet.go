package pet

import (
	"UltimateDesktopPet/pkg/attributes.go"
	"UltimateDesktopPet/pkg/math"
	pp "UltimateDesktopPet/pkg/print"
	"context"
	"sync"
	"time"
)

type Pet struct {
	sync.Mutex
	ID                    uint `json:"id" gorm:"primaryKey;autoIncrement"`
	attributes.Attributes `json:"attributes" gorm:"embedded"`
}

const (
	Max int16 = 100
	Min int16 = 0
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
	pp.Info(pp.Pet, "Pet Status - Expr: %d Water: %d, Hunger: %d, Health: %d, Mood: %d, Energy: %d, Money: %d",
		p.Experience, p.Water, p.Hunger, p.Health, p.Mood, p.Energy, p.Money)
	p.Unlock()
}

func (p *Pet) getStatus() Pet {
	/* Caller should lock the Pet struct */
	return Pet{
		ID: p.ID,
		Attributes: attributes.Attributes{
			Experience: p.Experience,
			Water:      p.Water,
			Hunger:     p.Hunger,
			Health:     p.Health,
			Mood:       p.Mood,
			Energy:     p.Energy,
			Money:      p.Money,
		},
	}
}

func (p *Pet) updateStatus(attr attributes.Attributes) {
	/* Caller should lock the Pet struct */
	p.Experience = p.Experience + attr.Experience
	p.Water = math.InRange(p.Water+attr.Water, Max, Min)
	p.Hunger = math.InRange(p.Hunger+attr.Hunger, Max, Min)
	p.Health = math.InRange(p.Health+attr.Health, Max, Min)
	p.Mood = math.InRange(p.Mood+attr.Mood, Max, Min)
	p.Energy = math.InRange(p.Energy+attr.Energy, Max, Min)
	p.Money = p.Money + attr.Money
}
