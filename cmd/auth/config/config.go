package config

import (
	"errors"
	"github.com/profclems/go-dotenv"
	"log"
)

type Server struct {
	Port uint16 `yaml:"port"`
}

type Config struct {
	Server *Server `yaml:"server"`
}

func (c *Config) Init() error {
	err := dotenv.LoadConfig()
	if err != nil {
		return err
	}

	port := dotenv.GetInt("PORT")
	if port == 0 {
		return errors.New("cannot find PORT variable")
	}

	log.Printf("Port: %d\n", port)

	c.Server = &Server{
		Port: uint16(port),
	}
	return nil
}
