package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gl1n0m3c/AKSP-KR/services/notificator/config"
	"github.com/gl1n0m3c/AKSP-KR/services/notificator/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config/config.yaml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		logger.Error("load config", "error", err)
		os.Exit(1)
	}

	if err = service.Run(ctx, cfg); err != nil {
		logger.Error("run", "error", err)
		os.Exit(1)
	}
}
