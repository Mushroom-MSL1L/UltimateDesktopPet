package app

import (
	"context"
	"embed"
	"fmt"
	"time"

	items "UltimateDesktopPet/internal/Items"
	"UltimateDesktopPet/internal/activities"
	"UltimateDesktopPet/internal/configLogics"
	"UltimateDesktopPet/internal/database"
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
}

func NewApp(assets embed.FS) *App {
	embeddedPetAssets = assets
	loadDefaultSprite()
	return &App{
		petMeta:      &pet.PetMeta{},
		itemsMeta:    &items.ItemsMeta{},
		activityMeta: &activities.ActivityMeta{},
	}
}

func (a *App) Startup(parentCtx context.Context) {
	/* app initialization */
	a.ctx = parentCtx
	a.useConfigurations(a.ctx, "./configs/system.yaml")
	pp.Assert(pp.App, "startup")

	/* app services */
	go a.petMeta.Service(a.ctx)

}

func (a *App) useConfigurations(c context.Context, configPath string) {
	sCfg := configs.LoadConfig(configPath, configLogics.System{})

	a.petMeta.DB.InitDB(c, sCfg.DBFile, database.Pets)
	a.petMeta.ImagePath = sCfg.PetImageDir

	a.itemsMeta.DB.InitDB(c, sCfg.ItemsDBFile, database.Items)
	a.itemsMeta.ImagePath = sCfg.ItemsDir

	a.activityMeta.DB.InitDB(c, sCfg.ActivitiesDBFile, database.Activities)
	a.activityMeta.ImagePath = sCfg.ActivitiesDir
}

func (a *App) Shutdown(parentCtx context.Context) {
	pp.Info(pp.App, "app shutdowning")
	a.petMeta.Shutdown()
	a.itemsMeta.Shutdown()
	a.activityMeta.Shutdown()
	pp.Assert(pp.App, "app shutdown complete")
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) PetSprite() string {
	if petSpriteDataURI == "" {
		pp.Fatal(pp.App, "PetSprite: no default sprite loaded")
	}
	return petSpriteDataURI
}

func (a *App) PetSpriteBy(path string) (string, error) {
	data, err := loadSpriteByName(path)
	if err != nil {
		pp.Fatal(pp.App, "PetSpriteBy: failed to load sprite %s: %v", path, err)
		return "", err
	}
	return data, nil
}
