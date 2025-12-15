package configs

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (Server) DefaultConfig() *Server {
	return &Server{
		Host: "localhost",
		Port: "8080",
	}
}

func (s *Server) Address() string {
	return s.Host + ":" + s.Port
}
