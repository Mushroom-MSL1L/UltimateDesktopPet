package main

import (
	"context"
	"embed"

	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/app"
	_ "github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/app"
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/configs"
	windowservice "github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/window"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

const configPath = "./configs/system.yaml"

//go:embed all:frontend/dist
var frontendBuild embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Create an instance of the app structure
	windowService := windowservice.NewWindowService()
	myapp := app.NewApp(configPath)
	configService := configs.NewConfigService(configPath)

	err := wails.Run(&options.App{
		Title:         "Ultimate Desktop Pet",
		Width:         150,
		Height:        150,
		DisableResize: true,
		Frameless:     true,
		AlwaysOnTop:   true,

		AssetServer: &assetserver.Options{
			Assets: frontendBuild,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		CSSDragProperty:  "--wails-draggable",
		CSSDragValue:     "drag",
		OnStartup: func(ctx context.Context) {
			windowService.SetContext(ctx)
			myapp.Startup(ctx)
		},
		OnShutdown: myapp.Shutdown,
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               true,
			DisableFramelessWindowDecorations: true,
		},
		Mac: &mac.Options{
			About: &mac.AboutInfo{
				Title: "Ultimate Desktop Pet",
				Icon:  icon,
			},
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
		Linux: &linux.Options{
			Icon:        icon,
			ProgramName: "Ultimate Desktop Pet",
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},
		Bind: []interface{}{
			myapp,
			windowService,
			myapp.PetMeta,
			myapp.ItemsMeta,
			myapp.ActivityMeta,
			myapp.ChatMeta,
			configService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
