package migrate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Config struct {
	Host string `yaml:"host" envconfig:"FEEMOUS_HOST"`
	Port string `yaml:"port" envconfig:"FEEMOUS_PORT"`
	DB   string `yaml:"db" envconfig:"FEEMOUS_DB"`
	User string `yaml:"user" envconfig:"FEEMOUS_USER"`
	Pass string `yaml:"pass" envconfig:"FEEMOUS_PASS"`
}

func Run(ctx context.Context, cfg Config) error {
	conn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.DB, cfg.User, cfg.Pass,
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return fmt.Errorf("sql open: %w", err)
	}
	defer db.Close()

	if err = db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	goose.SetBaseFS(migrations)

	if err = goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err = goose.UpContext(ctx, db, "migrations"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
