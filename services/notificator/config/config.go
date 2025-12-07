package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Rabbit RabbitConfig `yaml:"rabbit"`
}

type RabbitConfig struct {
	URL        string `yaml:"url" envconfig:"RAHMET_SURPO_RABBIT_URL"`
	Queue      string `yaml:"queue" envconfig:"RAHMET_SURPO_RABBIT_QUEUE"`
	Exchange   string `yaml:"exchange" envconfig:"RAHMET_SURPO_RABBIT_EXCHANGE"`
	RoutingKey string `yaml:"routing_key" envconfig:"RAHMET_SURPO_RABBIT_ROUTING_KEY"`
	Consumer   string `yaml:"consumer" envconfig:"RAHMET_SURPO_RABBIT_CONSUMER"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	if err = envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("env override: %w", err)
	}
	return cfg, nil
}
