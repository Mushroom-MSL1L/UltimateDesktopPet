package app

import (
	"UltimateDesktopPet/internal/activities"
	pp "UltimateDesktopPet/pkg/print"
)

func (a *App) LoadAllActivities() ([]activities.Activity, error) {
	activities, err := a.activityMeta.LoadAll()
	if err != nil {
		pp.Warn(pp.App, "LoadAllActivities: failed to load all activities: %v", err)
		return nil, err
	}
	return activities, err
}

func (a *App) LoadActivityFramesByID(id uint) ([]string, error) {
	frames, err := a.activityMeta.LoadFramesByID(id)
	if err != nil {
		pp.Warn(pp.App, "LoadActivityFramesByID: failed to load frames for activity ID %d: %v", id, err)
		return nil, err
	}
	return frames, nil
}

func (a *App) PerformActivity(acti activities.Activity) {
	pp.Info(pp.App, "PerformActivity : %s", acti.Name)

	a.petMeta.PerformActivity(acti.Name, acti.Experience, acti.Water, acti.Hunger, acti.Health, acti.Mood, acti.Energy, acti.Money, acti.DurationMinute)
}

func (a *App) StopActivity() {
	a.petMeta.StopActivity()
}
