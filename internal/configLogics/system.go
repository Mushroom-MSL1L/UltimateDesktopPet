package configLogics

type System struct {
	DBFile           string `yaml:"dbFile"`
	PetImageDir      string `yaml:"petImageDir"`
	ItemsDBFile      string `yaml:"itemsDBFile"`
	ItemsDir         string `yaml:"itemsDir"`
	ActivitiesDBFile string `yaml:"activitiesDBFile"`
	ActivitiesDir    string `yaml:"activitiesDir"`
}

func (System) DefaultConfig() *System {
	return &System{
		DBFile:           "./assets/db/udp.db",
		PetImageDir:      "./assets/petImages/default/",
		ItemsDBFile:      "./assets/db/items.db",
		ItemsDir:         "./assets/items/default/",
		ActivitiesDBFile: "./assets/db/activities.db",
		ActivitiesDir:    "./assets/activities/default/",
	}
}
