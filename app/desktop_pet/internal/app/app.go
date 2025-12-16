package app

import (
	"context"

	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/activities"
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/chat"
	_ "github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/chat"
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/configs"
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/database"
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/items"
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/pet"
	_ "github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/pet"

	"github.com/Mushroom-MSL1L/UltimateDesktopPet/pkg/configLoader"
	pp "github.com/Mushroom-MSL1L/UltimateDesktopPet/pkg/print"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx          context.Context
	PetMeta      *pet.PetMeta
	ItemsMeta    *items.ItemsMeta
	ActivityMeta *activities.ActivityMeta
	configs      *configs.System
	ChatMeta     *chat.ChatMeta
}

func NewApp(configPath string) *App {
	app := &App{
		PetMeta:      &pet.PetMeta{},
		ItemsMeta:    items.NewItemMeta(),
		ActivityMeta: activities.NewActivityMeta(),
	}
	app.ctx = context.Background()

	app.configs = configLoader.LoadConfig(configPath, &configs.System{})
	sCfg := app.configs

	app.ItemsMeta.DB.InitDB(app.ctx, sCfg.StaticAssetsDBDir, database.StaticAssets)
	app.ItemsMeta.ST.SpecifiedImageFolder = sCfg.ItemsImageFolder
	app.ActivityMeta.DB.InitDB(app.ctx, sCfg.StaticAssetsDBDir, database.StaticAssets)
	app.ActivityMeta.ST.SpecifiedImageFolder = sCfg.ActivitiesImageFolder

	app.ItemsMeta.DB.LoadSQLFileIfEmpty(sCfg.StaticAssetsSQLDir)

	app.PetMeta.DB.InitDB(app.ctx, sCfg.UDPDBDir, database.Pets)
	app.PetMeta = pet.NewPetMeta(app.PetMeta.DB, app.ItemsMeta, app.ActivityMeta)
	app.PetMeta.ST.SpecifiedImageFolder = sCfg.PetImageFolder

	app.ChatMeta = chat.NewChatMeta(app.ctx, sCfg.GeminiAPIKey)
	app.ChatMeta.DB.InitDB(app.ctx, sCfg.UDPDBDir, database.Pets)
	app.ChatMeta.RolePlayContext = sCfg.ChatRolePlayContext

	return app
}

func (a *App) Startup(parentCtx context.Context) {
	/* app initialization */
	a.ctx = parentCtx
	pp.Assert(pp.App, "startup")

	/* app services */
	go a.PetMeta.Service(a.ctx)
}

func (a *App) Shutdown(parentCtx context.Context) {
	pp.Info(pp.App, "app shutdowning")
	a.PetMeta.Shutdown()
	a.ItemsMeta.Shutdown()
	a.ActivityMeta.Shutdown()
	a.ChatMeta.Shutdown()
	pp.Assert(pp.App, "app shutdown complete")
}

func (a *App) Quit() {
	pp.Info(pp.App, "app quitting")
	runtime.Quit(a.ctx)
}
