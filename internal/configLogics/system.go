package configLogics

type System struct {
	DBFile      string `yaml:"dbfile"`
	HostAddress string `yaml:"hostaddress"`
	Port        string `yaml:"port"`
}
