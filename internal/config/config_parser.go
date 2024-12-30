package config

import (
	"os"

	"errors"

	"gopkg.in/yaml.v2"
)

var ErrParsingYaml = errors.New("Error parsing YAML")

type Config struct {
	Port    int            `yaml:"entry_point"`
	Servers []ServerConfig `yaml:"servers"`
}

type ServerConfig struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func GetConfigFromPath(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, ErrParsingYaml
	}
	return &config, nil
}
