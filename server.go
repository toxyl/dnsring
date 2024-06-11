package dnsring

type Server struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Timeout int    `yaml:"timeout"`
}

func NewServer(host string, port, timeout int) *Server {
	srv := &Server{
		Host:    host,
		Port:    port,
		Timeout: timeout,
	}
	return srv
}
