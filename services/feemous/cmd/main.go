package main

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/gl1n0m3c/AKSP-KR/services/feemous/config"
	"github.com/gl1n0m3c/AKSP-KR/services/feemous/internal/server"
	"github.com/gl1n0m3c/AKSP-KR/services/feemous/internal/store"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config/config.yaml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		logger.Error("load config", "error", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", dsn(cfg))
	if err != nil {
		logger.Error("db open", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		logger.Error("db ping", "error", err)
		os.Exit(1)
	}

	repo := store.New(db)
	s := server.New(repo)

	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      s.Routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info("feemous listening", "addr", cfg.Server.Addr)
	if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server", "error", err)
		os.Exit(1)
	}
}

func dsn(cfg *config.Config) string {
	return "host=" + cfg.DB.Host +
		" port=" + cfg.DB.Port +
		" dbname=" + cfg.DB.Name +
		" user=" + cfg.DB.User +
		" password=" + cfg.DB.Pass +
		" sslmode=disable"
}
