package activities

import (
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/desktop_pet/internal/database"
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/desktop_pet/pkg/file"
	pp "github.com/Mushroom-MSL1L/UltimateDesktopPet/desktop_pet/pkg/print"
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

func NewActivityMeta() *ActivityMeta {
	a := &ActivityMeta{}
	a.Activity = &Activity{}
	a.Controller = newActivityController(&a.Activity)
	a.ST = file.NewSpriteTool(a.ST)
	a.ST.StaticAssetPath = activitiesStaticAssetPath
	a.ST.DefaultImageFolder = activitiesDefaultImageFolder
	return a
}

func (a *ActivityMeta) Shutdown() {
	a.DB.CloseDB()
	pp.Assert(pp.Activities, "activities service stopped")
}

func (a *ActivityMeta) LoadAll() ([]Activity, error) {
	activities, err := a.Controller.ReadAll(a.DB.GetDB())
	if err != nil {
		return nil, err
	}

	validActivities := make([]Activity, 0, len(*activities))
	for _, acti := range *activities {
		validActivities = append(validActivities, acti)
	}
	return validActivities, err
}

func (a *ActivityMeta) LoadFramesByID(id uint) ([]string, error) {
	activity, err := a.Controller.Read(a.DB.GetDB(), id)
	if err != nil {
		return nil, err
	}
	return a.ST.LoadFramesFromDir(activity.Path)
}
