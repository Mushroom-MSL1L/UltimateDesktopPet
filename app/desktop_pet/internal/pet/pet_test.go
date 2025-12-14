package pet

import (
	"testing"

	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/attributes"
)

func TestPet_updateStatus_ClampsFields(t *testing.T) {
	p := &Pet{
		Attributes: attributes.Attributes{
			Experience: 10,
			Water:      0,
			Hunger:     100,
			Health:     50,
			Mood:       1,
			Energy:     100,
			Money:      5,
		},
	}

	p.updateStatus(attributes.Attributes{
		Experience: 3,
		Water:      -5,
		Hunger:     10,
		Health:     -100,
		Mood:       -5,
		Energy:     10,
		Money:      7,
	})

	if p.Experience != 13 {
		t.Fatalf("Experience = %d, want %d", p.Experience, 13)
	}
	if p.Money != 12 {
		t.Fatalf("Money = %d, want %d", p.Money, 12)
	}

	if p.Water != Min {
		t.Fatalf("Water = %d, want %d", p.Water, Min)
	}
	if p.Hunger != Max {
		t.Fatalf("Hunger = %d, want %d", p.Hunger, Max)
	}
	if p.Health != Min {
		t.Fatalf("Health = %d, want %d", p.Health, Min)
	}
	if p.Mood != Min {
		t.Fatalf("Mood = %d, want %d", p.Mood, Min)
	}
	if p.Energy != Max {
		t.Fatalf("Energy = %d, want %d", p.Energy, Max)
	}
}

func TestPet_getStatus_ReturnsCopy(t *testing.T) {
	p := &Pet{
		ID: 42,
		Attributes: attributes.Attributes{
			Experience: 1,
			Water:      2,
			Hunger:     3,
			Health:     4,
			Mood:       5,
			Energy:     6,
			Money:      7,
		},
	}

	got := p.getStatus()
	if got.ID != 42 || got.Water != 2 || got.Money != 7 {
		t.Fatalf("getStatus = %#v, want ID=42 Water=2 Money=7", got)
	}

	got.Water = 99
	if p.Water != 2 {
		t.Fatalf("modifying returned value should not affect original: p.Water=%d", p.Water)
	}
}

