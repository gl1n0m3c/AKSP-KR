package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gl1n0m3c/AKSP-KR/services/migrator/config"
	"github.com/gl1n0m3c/AKSP-KR/services/migrator/internal/service"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config/config.yaml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		slog.Error("load config", err)
		os.Exit(1)
	}

	if err = service.Run(ctx, cfg); err != nil {
		slog.Error("migrate failed", err)
		os.Exit(1)
	}
}
