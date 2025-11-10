package activities

import (
	"UltimateDesktopPet/internal/database"
	pp "UltimateDesktopPet/pkg/print"
)

type ActivityMeta struct {
	Controller *database.BaseController[Activity]
	ImagePath  string
	DB         database.DB
	Activity   *Activity
}

func init() {
	a := newActivityController(nil)
	database.RegisterSchema(database.Images, a)
	pp.Assert(pp.Activities, "Activities init complete")
}

func newActivityController(model **Activity) *database.BaseController[Activity] {
	return &database.BaseController[Activity]{Model: model}
}

func (p *ActivityMeta) Shutdown() {
	p.DB.CloseDB()
	pp.Assert(pp.Activities, "activities service stopped")
}
