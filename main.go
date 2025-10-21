package main

import (
	"UltimateDesktopPet/internal/configLogics"
	"UltimateDesktopPet/internal/database"
	"UltimateDesktopPet/internal/synchronization"
	"UltimateDesktopPet/pkg/configs"
)

//	@title			Ultimate Desktop Pet API
//	@version		1.0
//	@description	Ultimate Desktop Pet API

//	@contact.name	Mushroom_MSL1L
//	@contact.url	https://github.com/Mushroom-MSL1L

//	@license.name	no license yet

//	@host	localhost:8080
//	@BasePath  /

//	@tag.name	temp
//	@tag.name	docs

func main() {
	var myDB = database.DB{}
	sCfg := configs.LoadConfig[configLogics.System]("./configs/system.yaml")

	myDB.InitDB(sCfg.DBFile)
	defer myDB.CloseDB()

	synchronization.Server(sCfg.HostAddress, sCfg.Port)
}
