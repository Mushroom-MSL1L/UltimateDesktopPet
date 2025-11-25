package app

import (
	"context"

	"UltimateDesktopPet/internal/activities"
	"UltimateDesktopPet/internal/configLogics"
	"UltimateDesktopPet/internal/database"
	"UltimateDesktopPet/internal/items"
	"UltimateDesktopPet/internal/pet"
	_ "UltimateDesktopPet/internal/pet"

	"UltimateDesktopPet/pkg/configs"
	pp "UltimateDesktopPet/pkg/print"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx          context.Context
	petMeta      *pet.PetMeta
	itemsMeta    *items.ItemsMeta
	activityMeta *activities.ActivityMeta
	configs      *configLogics.System
}

func NewApp(configPath string) *App {
	app := &App{
		petMeta:      &pet.PetMeta{},
		itemsMeta:    items.NewItemMeta(),
		activityMeta: activities.NewActivityMeta(),
	}
	app.ctx = context.Background()

	app.configs = configs.LoadConfig(configPath, &configLogics.System{})
	sCfg := app.configs

	app.itemsMeta.DB.InitDB(app.ctx, sCfg.StaticAssetsDBDir, database.StaticAssets)
	app.itemsMeta.ST.SpecifiedImageFolder = sCfg.ItemsImageFolder
	app.activityMeta.DB.InitDB(app.ctx, sCfg.StaticAssetsDBDir, database.StaticAssets)
	app.activityMeta.ST.SpecifiedImageFolder = sCfg.ActivitiesImageFolder

	app.petMeta.DB.InitDB(app.ctx, sCfg.UDPDBDir, database.Pets)
	app.petMeta = pet.NewPetMeta(app.petMeta.DB, app.itemsMeta, app.activityMeta)
	app.petMeta.ST.SpecifiedImageFolder = sCfg.PetImageFolder

	return app
}

func (a *App) Startup(parentCtx context.Context) {
	/* app initialization */
	a.ctx = parentCtx
	pp.Assert(pp.App, "startup")

	/* app services */
	go a.petMeta.Service(a.ctx)
}

func (a *App) Shutdown(parentCtx context.Context) {
	pp.Info(pp.App, "app shutdowning")
	a.petMeta.Shutdown()
	a.itemsMeta.Shutdown()
	a.activityMeta.Shutdown()
	pp.Assert(pp.App, "app shutdown complete")
}

func (a *App) Quit() {
	pp.Info(pp.App, "app quitting")
	runtime.Quit(a.ctx)
}
