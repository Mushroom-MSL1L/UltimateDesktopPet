package activities

import (
	pp "UltimateDesktopPet/pkg/print"
)

func (a *ActivityMeta) LoadAllActivities() ([]Activity, error) {
	activities, err := a.LoadAll()
	if err != nil {
		pp.Warn(pp.Activities, "LoadAllActivities: failed to load all activities: %v", err)
		return nil, err
	}
	return activities, err
}

func (a *ActivityMeta) LoadActivityFramesByID(id uint) ([]string, error) {
	frames, err := a.LoadFramesByID(id)
	if err != nil {
		pp.Warn(pp.Activities, "LoadActivityFramesByID: failed to load frames for activity ID %d: %v", id, err)
		return nil, err
	}
	return frames, nil
}
