package main

import (
	"flag"
	"os"

	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/sync_server/internal/configs"
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/sync_server/internal/routes"

	"github.com/Mushroom-MSL1L/UltimateDesktopPet/pkg/configLoader"
	pp "github.com/Mushroom-MSL1L/UltimateDesktopPet/pkg/print"
)

func main() {
	cfgPath := getConfigPath()
	cfg := configLoader.LoadConfig(cfgPath, &configs.Server{})
	pp.Info(pp.System, "using config : %s", cfgPath)

	r := routes.SetupRouter(cfg.Address())
	if err := r.Run(cfg.Address()); err != nil {
		pp.Fatal(pp.System, "server stopped: %v", err)
	}
}

func getConfigPath() string {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	if configPath != "" {
		return configPath
	}

	if env := os.Getenv("SYNC_SERVER_CONFIG_PATH"); env != "" {
		return env
	}

	return "./configs/server.yaml"
}
