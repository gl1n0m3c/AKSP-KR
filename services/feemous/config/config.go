package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
}

type ServerConfig struct {
	Addr string `yaml:"addr" envconfig:"FEEMOUS_ADDR"`
}

type DBConfig struct {
	Host string `yaml:"host" envconfig:"FEEMOUS_HOST"`
	Port string `yaml:"port" envconfig:"FEEMOUS_PORT"`
	Name string `yaml:"name" envconfig:"FEEMOUS_DB"`
	User string `yaml:"user" envconfig:"FEEMOUS_USER"`
	Pass string `yaml:"pass" envconfig:"FEEMOUS_PASS"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	if err = envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("env override: %w", err)
	}

	return cfg, nil
}
