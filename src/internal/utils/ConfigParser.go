package utils

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Volume      string             `yaml:"volume"`
	WorkerImage string             `yaml:"worker"`
	Projects    map[string]Project `yaml:"projects"`
	Remotes     map[string]Remote  `yaml:"remotes"`
}

type Project struct {
	Containers []string `yaml:"containers"`
	Passphrase string   `yaml:"passphrase"`
}

type Remote struct {
	Type string `yaml:"type"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Path string `yaml:"path"`
}

func ParseConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH not set")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	content := os.ExpandEnv(string(data))

	var config Config
	err = yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
