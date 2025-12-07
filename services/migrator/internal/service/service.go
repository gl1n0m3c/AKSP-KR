package service

import (
	"context"

	"github.com/gl1n0m3c/AKSP-KR/services/migrator/config"
	"github.com/gl1n0m3c/AKSP-KR/services/migrator/internal/migrate"
)

func Run(ctx context.Context, cfg *config.Config) error {
	return migrate.Run(ctx, cfg.Migrate)
}
