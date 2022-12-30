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
	DBURL  string
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

	dbURL := dotenv.GetString("DB_URL")
	log.Printf("DB URL: %d\n", dbURL)
	if dbURL == "" {
		return nil, errors.New("cannot find MONGODB_URL variable")
	}

	return &Config{
		Server: &Server{
			Port: uint16(port),
		},
		DBURL: dbURL,
	}, nil
}
