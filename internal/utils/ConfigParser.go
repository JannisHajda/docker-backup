package utils

import (
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Projects map[string]Project `yaml:"projects"`
	Remotes  map[string]Remote  `yaml:"remotes"`
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
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	data, err := os.ReadFile("config.yaml")
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
