package configLogics

type System struct {
	UDPDBDir              string `yaml:"udpDBDir"`
	ImageDBDir            string `yaml:"imageDBDir"`
	PetImageFolder        string `yaml:"petImageFolder"`
	ItemsImageFolder      string `yaml:"itemsImageFolder"`
	ActivitiesImageFolder string `yaml:"activitiesImageFolder"`
}

func (System) DefaultConfig() *System {
	return &System{
		UDPDBDir:              "./assets/db/udp.db",
		ImageDBDir:            "./assets/db/images.db",
		PetImageFolder:        "./assets/petImages/default/",
		ItemsImageFolder:      "./assets/itemImages/default/",
		ActivitiesImageFolder: "./assets/activityImages/default/",
	}
}
