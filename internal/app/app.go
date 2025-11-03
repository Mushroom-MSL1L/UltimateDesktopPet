package app

import (
	"context"
	"embed"
	"fmt"

	"UltimateDesktopPet/internal/configLogics"
	"UltimateDesktopPet/internal/database"
	"UltimateDesktopPet/internal/pet"
	_ "UltimateDesktopPet/internal/pet"
	_ "UltimateDesktopPet/internal/system"

	"UltimateDesktopPet/pkg/configs"
	pp "UltimateDesktopPet/pkg/print"
)

type App struct {
	ctx     context.Context
	myDB    database.DB
}

func NewApp(assets embed.FS) *App {
	embeddedPetAssets = assets
	loadDefaultSprite()
	return &App{}
}

func (a *App) Startup(parentCtx context.Context) {
	/* app initialization */
	a.ctx = parentCtx
	sCfg := configs.LoadConfig("./configs/system.yaml", configLogics.System{})
	a.petMeta = &pet.PetMeta{}
	a.myDB.InitDB(a.ctx, sCfg.DBFile)

	/* app services */
	go pet.Service(a.ctx, myDB.GetDB())
}

func (a *App) Shutdown(parentCtx context.Context) {
	pp.Assert(pp.App, "shutdown")
	a.myDB.CloseDB()
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
