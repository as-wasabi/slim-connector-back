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

// LoadConfig 設定ファイルを読み込む
func LoadConfig(path string) (*Config, error) {
	// デフォルトパス
	if path == "" {
		path = "./config/config.yml"
	}

	// ファイル読み込み
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config

	// YAML解析
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// 環境変数による上書き
	overrideConfigFromEnv(&config)

	return &config, nil
}

// overrideConfigFromEnv 環境変数で設定を上書き
func overrideConfigFromEnv(config *Config) {
	// サーバー設定
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		// 文字列からintへの変換が必要
		if p, err := strconv.Atoi(port); err == nil {
			config.Server.Port = p
		}
	}

	// MongoDB設定
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		config.MongoDB.URI = uri
	}
}
