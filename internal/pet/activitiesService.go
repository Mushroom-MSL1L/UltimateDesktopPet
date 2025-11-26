package pet

import (
	pp "UltimateDesktopPet/pkg/print"
	"fmt"
	"time"
)

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
