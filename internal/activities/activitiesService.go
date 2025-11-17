package activities

import (
	"UltimateDesktopPet/internal/database"
	"UltimateDesktopPet/pkg/file"
	pp "UltimateDesktopPet/pkg/print"
	"context"
	"fmt"
)

const activitiesStaticAssetPath = "assets/activityImages"
const activitiesDefaultImageFolder = "default"

type ActivityMeta struct {
	Controller *database.BaseController[Activity]
	DB         database.DB
	Activity   *Activity
	ST         file.SpriteTool
}

func init() {
	a := newActivityController(nil)
	database.RegisterSchema(database.StaticAssets, a)
	pp.Assert(pp.Activities, "Activities init complete")
}

func newActivityController(model **Activity) *database.BaseController[Activity] {
	return &database.BaseController[Activity]{Model: model}
}

func (a *ActivityMeta) Service(c context.Context) {
	a.Activity = &Activity{}
	a.Controller = newActivityController(&a.Activity)
	a.ST = file.NewSpriteTool(a.ST)
	a.ST.StaticAssetPath = activitiesStaticAssetPath
	a.ST.DefaultImageFolder = activitiesDefaultImageFolder
}

func (p *ActivityMeta) Shutdown() {
	p.DB.CloseDB()
	pp.Assert(pp.Activities, "activities service stopped")
}

func (p *ActivityMeta) LoadAll() ([]ActivityWithFrames, error) {
	activities, err := p.Controller.ReadAll(p.DB.GetDB())
	if activities == nil {
		return nil, fmt.Errorf("LoadAll return nil")
	}

	validActivities := make([]ActivityWithFrames, 0, len(*activities))
	for _, acti := range *activities {
		frames, err := p.ST.LoadFramesFromDir(acti.Path)
		if err != nil {
			pp.Warn(pp.Activities, "LoadAll: activity %s cannot load frame: %v", acti.Name, err)
			continue
		}
		var temp ActivityWithFrames
		temp.Frames = frames
		temp.Activity = acti
		validActivities = append(validActivities, temp)
	}
	return validActivities, err
}
