package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DB       DBConfig       `yaml:"db"`
	Rabbit   RabbitConfig   `yaml:"rabbit"`
	Schedule ScheduleConfig `yaml:"schedule"`
}

type DBConfig struct {
	Host string `yaml:"host" envconfig:"FEEMOUS_HOST"`
	Port string `yaml:"port" envconfig:"FEEMOUS_PORT"`
	Name string `yaml:"name" envconfig:"FEEMOUS_DB"`
	User string `yaml:"user" envconfig:"FEEMOUS_USER"`
	Pass string `yaml:"pass" envconfig:"FEEMOUS_PASS"`
}

type RabbitConfig struct {
	URL        string `yaml:"url" envconfig:"RAHMET_SURPO_RABBIT_URL"`
	Queue      string `yaml:"queue" envconfig:"RAHMET_SURPO_RABBIT_QUEUE"`
	Exchange   string `yaml:"exchange" envconfig:"RAHMET_SURPO_RABBIT_EXCHANGE"`
	RoutingKey string `yaml:"routing_key" envconfig:"RAHMET_SURPO_RABBIT_ROUTING_KEY"`
}

type ScheduleConfig struct {
	Period string `yaml:"period" envconfig:"SCHED_PERIOD"` // e.g. "1m"
	Window string `yaml:"window" envconfig:"SCHED_WINDOW"` // e.g. "15m"
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
