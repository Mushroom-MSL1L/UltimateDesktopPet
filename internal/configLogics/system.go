package configLogics

type System struct {
	UDPDBDir              string `yaml:"udpDBDir"`
	StaticAssetsDBDir     string `yaml:"staticAssetsDBDir"`
	PetImageFolder        string `yaml:"petImageFolder"`
	ItemsImageFolder      string `yaml:"itemsImageFolder"`
	ActivitiesImageFolder string `yaml:"activitiesImageFolder"`
}

func (System) DefaultConfig() *System {
	return &System{
		UDPDBDir:              "./assets/db/udp.db",
		StaticAssetsDBDir:     "./assets/db/static_assets.db",
		PetImageFolder:        "./assets/petImages/default/",
		ItemsImageFolder:      "./assets/itemImages/default/",
		ActivitiesImageFolder: "./assets/activityImages/default/",
	}
}
