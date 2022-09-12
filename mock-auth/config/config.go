package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Server struct {
	URL  string `yaml:"url"`
	Port uint16 `yaml:"port"`
}

type Config struct {
	Server *Server `yaml:"server"`
}

func (c *Config) Init() error {
	bytes, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bytes, c)
}
