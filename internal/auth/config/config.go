package config

import (
	"errors"
	"log"

	"github.com/profclems/go-dotenv"
)

type Server struct {
	Port uint16 `yaml:"port"`
}

type Config struct {
	Server *Server `yaml:"server"`
	DBURI  string
}

func Init() (*Config, error) {
	err := dotenv.LoadConfig()
	if err != nil {
		return nil, err
	}

	port := dotenv.GetInt("PORT")
	log.Printf("Port: %d\n", port)
	if port == 0 {
		return nil, errors.New("cannot find PORT variable")
	}

	dbURI := dotenv.GetString("DB_URI")
	log.Printf("DB URI: %s\n", dbURI)
	if dbURI == "" {
		return nil, errors.New("cannot find MONGODB_URI variable")
	}

	return &Config{
		Server: &Server{
			Port: uint16(port),
		},
		DBURI: dbURI,
	}, nil
}
