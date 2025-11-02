package app

import (
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"UltimateDesktopPet/internal/configLogics"
	"UltimateDesktopPet/internal/database"
	"UltimateDesktopPet/internal/pet"
	_ "UltimateDesktopPet/internal/pet"
	_ "UltimateDesktopPet/internal/system"

	"UltimateDesktopPet/pkg/configs"
)

var embeddedPetAssets embed.FS

var petSpriteDataURI string

func loadDefaultSprite() {
	defPath := filepath.ToSlash("assets/petImages/default/cat_move.gif")
	if data, err := embeddedPetAssets.ReadFile(defPath); err == nil && len(data) > 0 {
		petSpriteDataURI = "data:image/gif;base64," + base64.StdEncoding.EncodeToString(data)
		return
	}

	relPath := filepath.Join("..", "..", "assets", "petImages", "default", "cat_move.gif")
	data, err := os.ReadFile(relPath)
	if err != nil {
		if exe, eerr := os.Executable(); eerr == nil {
			exeDir := filepath.Dir(exe)
			tryPath := filepath.Join(exeDir, "..", "..", "assets", "petImages", "default", "cat_move.gif")
			if d2, e2 := os.ReadFile(tryPath); e2 == nil {
				data = d2
				err = nil
			}
		}
	}

	if err == nil && len(data) > 0 {
		petSpriteDataURI = "data:image/gif;base64," + base64.StdEncoding.EncodeToString(data)
	} else {
		petSpriteDataURI = ""
	}
}

// helper: load a sprite by logical name, e.g. "default/cat_move.gif"
func loadSpriteByName(name string) (string, error) {
	// normalize name
	name = filepath.ToSlash(strings.TrimPrefix(name, "/"))

	// try embedded first; embedded paths use the same relative path as in main.go embed
	embeddedPath := filepath.ToSlash("assets/petImages/" + name)
	if data, err := embeddedPetAssets.ReadFile(embeddedPath); err == nil && len(data) > 0 {
		return "data:image/gif;base64," + base64.StdEncoding.EncodeToString(data), nil
	}

	// fallback filesystem (relative to repo)
	fsPath := filepath.Join("..", "..", "assets", "petImages", filepath.FromSlash(name))
	if data, err := os.ReadFile(fsPath); err == nil && len(data) > 0 {
		return "data:image/gif;base64," + base64.StdEncoding.EncodeToString(data), nil
	}

	// fallback relative to executable
	if exe, eerr := os.Executable(); eerr == nil {
		exeDir := filepath.Dir(exe)
		tryPath := filepath.Join(exeDir, "..", "..", "assets", "petImages", filepath.FromSlash(name))
		if data, err := os.ReadFile(tryPath); err == nil && len(data) > 0 {
			return "data:image/gif;base64," + base64.StdEncoding.EncodeToString(data), nil
		}
	}

	return "", fmt.Errorf("sprite not found: %s", name)
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
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
	// go func() {
	// 	defer cancelFunc()
	// 	signalChannel := make(chan os.Signal, 1)
	// 	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	// 	pp.Assert(pp.System, "signal handler registed")
	// 	select {
	// 	case s := <-signalChannel:
	// 		pp.Assert(pp.System, "trap signal %v", s)
	// 	case <-a.ctx.Done():
	// 		pp.Assert(pp.System, "context done")
	// 	}
	// }()
	return a.ctx
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	/* app initialization */
	a.setContextWithCancelAndSignalHandler(ctx)
	sCfg := configs.LoadConfig("./configs/system.yaml", configLogics.System{})

	var myDB = database.DB{}
	myDB.InitDB(ctx, sCfg.DBFile)

	/* app services */
	go pet.Service(a.ctx, myDB.GetDB())
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// PetSprite returns the embedded pet sprite as a data URI.
func (a *App) PetSprite() string {
	if petSpriteDataURI == "" {
		log.Println("PetSprite: petSpriteDataURI is empty")
	} else {
		log.Println("PetSprite: length =", len(petSpriteDataURI))
	}
	return petSpriteDataURI
}

// PetSpriteBy returns a specific sprite by logical path under assets/petImages.
func (a *App) PetSpriteBy(name string) (string, error) {
	data, err := loadSpriteByName(name)
	if err != nil {
		log.Printf("PetSpriteBy(%s) error: %v\n", name, err)
		return "", err
	}
	log.Printf("PetSpriteBy(%s) length=%d\n", name, len(data))
	return data, nil
}
