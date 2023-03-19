package config

import "github.com/dlshle/authnz/pkg/yaml"

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
	GRPC string `yaml:"grpc"`
}

type DatabaseConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	DBName string `yaml:"db_name"`
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
}

func Load(path string) (Config, error) {
	var cfg Config
	err := yaml.LoadConfig(path, &cfg)
	return cfg, err
}
