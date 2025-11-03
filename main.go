package main

import (
	"embed"

	"UltimateDesktopPet/internal/app"
	_ "UltimateDesktopPet/internal/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed assets/petImages/*/*
var petAssets embed.FS

func main() {
	// Create an instance of the app structure
	myapp := app.NewApp(petAssets)

	// Create application with options
	err := wails.Run(&options.App{
		Title:         "Ultimate Desktop Pet",
		Width:         280,
		Height:        280,
		DisableResize: true,
		Frameless:     true,
		AlwaysOnTop:   true,

		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		CSSDragProperty:  "--wails-draggable",
		CSSDragValue:     "drag",
		OnStartup:        myapp.Startup,
		OnShutdown:       myapp.Shutdown,
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               true,
			DisableFramelessWindowDecorations: true,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  true,
				HideTitleBar:               true,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
		},
		Linux: &linux.Options{},
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},
		Bind: []interface{}{
			myapp,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
