package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"

	"github.com/gl1n0m3c/AKSP-KR/services/migrator/internal/migrate"
)

type Config struct {
	Migrate migrate.Config `yaml:"migrate"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("unmarshall: %w", err)
	}
	if err = envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("process: %w", err)
	}

	return cfg, nil
}
