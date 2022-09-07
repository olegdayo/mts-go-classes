package main

type Server struct {
	URL  string
	Port uint16
}

type Config struct {
	Server *Server
}

func InitConfig() (*Config, error) {
	return &Config{
			Server: &Server{
				URL:  "localhost",
				Port: 8080,
			},
		},
		nil
}
