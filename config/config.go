package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`

	MongoDB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		URI      string `yaml:"uri"`
	} `yaml:"mongodb"`
}

func LoadConfig(path string) (*Config, error) {

	if path == "" {
		path = "./config/config.yml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	overrideConfigFromEnv(&config)

	return &config, nil
}

func overrideConfigFromEnv(config *Config) {

	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {

		if p, err := strconv.Atoi(port); err == nil {
			config.Server.Port = p
		}
	}

	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		config.MongoDB.URI = uri
	}
}
