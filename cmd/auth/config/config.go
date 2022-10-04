package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	PathToConfig = "cmd/auth/config/config.yaml"
)

type Server struct {
	URL  string `yaml:"url"`
	Port uint16 `yaml:"port"`
}

type Config struct {
	Server *Server `yaml:"server"`
}

func (c *Config) Init() error {
	bytes, err := ioutil.ReadFile(PathToConfig)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bytes, c)
}
