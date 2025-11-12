package app

import (
	"context"
	"fmt"
	"strings"

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
		itemsMeta:    &items.ItemsMeta{},
		activityMeta: &activities.ActivityMeta{},
	}
	app.configs = configs.LoadConfig(configPath, &configLogics.System{})

	return app
}

func (a *App) Startup(parentCtx context.Context) {
	/* app initialization */
	a.ctx = parentCtx
	a.useConfigurations()
	pp.Assert(pp.App, "startup")

	/* app services */
	go a.petMeta.Service(a.ctx)
}

func (a *App) useConfigurations() {
	sCfg := a.configs

	a.petMeta.DB.InitDB(a.ctx, sCfg.UDPDBDir, database.Pets)
	a.petMeta.ImagePath = sCfg.PetImageFolder

	a.itemsMeta.DB.InitDB(a.ctx, sCfg.StaticAssetsDBDir, database.StaticAssets)
	a.itemsMeta.ImagePath = sCfg.ItemsImageFolder
	a.activityMeta.ImagePath = sCfg.ActivitiesImageFolder
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

func (a *App) PetFrames() ([]string, error) {
	frames, err := loadDefaultFrames(a.configs.PetImageFolder)
	if err != nil {
		pp.Fatal(pp.App, "PetFrames: failed to load default frames: %v", err)
		return nil, err
	}
	return frames, nil
}

func (a *App) PetFramesBy(animationType string) ([]string, error) {
	frames, err := loadFramesFromDir(a.configs.PetImageFolder, animationType)
	if err != nil {
		pp.Fatal(pp.App, "PetFramesBy: failed to load frames for animation type %s: %v", animationType, err)
		return nil, err
	}
	return frames, nil
}

// ChatWithPet is a stub that will later become the real conversation handler.
func (a *App) ChatWithPet(userInput string) string {
	trimmed := strings.TrimSpace(userInput)
	pp.Info(pp.App, "ChatWithPet: received message %q", trimmed)
	if trimmed == "" {
		return "I didn't catch that. Try saying hi!"
	}
	return fmt.Sprintf("I'm still learning, but I heard: %s", trimmed)
}
