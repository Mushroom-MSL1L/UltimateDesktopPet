package configLogics

type System struct {
	UDPDBDir              string `yaml:"udpDBDir"`
	PetImageFolder        string `yaml:"petImageFolder"`
	ItemsDBDir            string `yaml:"itemsDBDir"`
	ItemsImageFolder      string `yaml:"itemsImageFolder"`
	ActivitiesDBDir       string `yaml:"activitiesDBDir"`
	ActivitiesImageFolder string `yaml:"activitiesImageFolder"`
}

func (System) DefaultConfig() *System {
	return &System{
		UDPDBDir:              "./assets/db/udp.db",
		PetImageFolder:        "./assets/petImages/default/",
		ItemsDBDir:            "./assets/db/items.db",
		ItemsImageFolder:      "./assets/itemImages/default/",
		ActivitiesDBDir:       "./assets/db/activities.db",
		ActivitiesImageFolder: "./assets/activityImages/default/",
	}
}
