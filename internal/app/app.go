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
	ctx context.Context
}

func NewApp(assets embed.FS) *App {
	embeddedPetAssets = assets
	loadDefaultSprite()
	return &App{}
}

func (a *App) setContextWithCancelAndSignalHandler(c context.Context) context.Context {
	if c == nil {
		c = context.Background()
	}
	ctx, _ := context.WithCancel(c)
	a.ctx = ctx
	return a.ctx
}

func (a *App) Startup(ctx context.Context) {
	/* app initialization */
	a.setContextWithCancelAndSignalHandler(ctx)
	sCfg := configs.LoadConfig("./configs/system.yaml", configLogics.System{})

	var myDB = database.DB{}
	myDB.InitDB(ctx, sCfg.DBFile)

	/* app services */
	go pet.Service(a.ctx, myDB.GetDB())
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
