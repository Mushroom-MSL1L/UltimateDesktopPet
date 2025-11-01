package main

import (
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed assets/petImages/default/cat_move.gif
var petSpriteBytes []byte

var petSpriteDataURI = "data:image/gif;base64," + base64.StdEncoding.EncodeToString(petSpriteBytes)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.WindowSetAlwaysOnTop(ctx, true)
	runtime.WindowSetBackgroundColour(ctx, 0, 0, 0, 0)
	runtime.WindowSetMinSize(ctx, 200, 200)
	runtime.WindowSetSize(ctx, 280, 280)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// PetSprite returns the embedded pet sprite as a data URI.
func (a *App) PetSprite() string {
	return petSpriteDataURI
}
