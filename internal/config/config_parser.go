package config

import (
	"os"

	"gopkg.in/yaml.v2"
	"sigs.k8s.io/kind/pkg/errors"
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
