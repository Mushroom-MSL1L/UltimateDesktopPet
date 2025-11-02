package configLogics

type System struct {
	DBFile string `yaml:"dbfile"`
}

func (System) DefaultConfig() *System {
	return &System{
		DBFile: "./assets/db/udp.db",
	}
}
