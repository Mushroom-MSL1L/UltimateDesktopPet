package configLogics

type System struct {
	DBFile      string `yaml:"dbfile"`
	HostAddress string `yaml:"hostaddress"`
	Port        string `yaml:"port"`
}

func (System) DefaultConfig() *System {
	return &System{
		DBFile:      "./assets/db/udp.db",
		HostAddress: "localhost",
		Port:        "8080",
	}
}
