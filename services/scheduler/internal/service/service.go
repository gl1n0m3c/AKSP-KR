package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"

	"github.com/gl1n0m3c/AKSP-KR/services/scheduler/config"
)

type App struct {
	repo      Repository
	publisher Publisher
	period    time.Duration
	window    time.Duration
}

func Run(ctx context.Context, cfg *config.Config) error {
	period, err := time.ParseDuration(cfg.Schedule.Period)
	if err != nil {
		return fmt.Errorf("parse period: %w", err)
	}
	window, err := time.ParseDuration(cfg.Schedule.Window)
	if err != nil {
		return fmt.Errorf("parse window: %w", err)
	}

	db, err := sql.Open("postgres", dsn(cfg.DB))
	if err != nil {
		return fmt.Errorf("sql open: %w", err)
	}
	defer db.Close()

	repository := NewRepository(db)

	rmqConn, err := amqp.Dial(cfg.Rabbit.URL)
	if err != nil {
		return fmt.Errorf("rabbit dial: %w", err)
	}
	defer rmqConn.Close()

	ch, err := rmqConn.Channel()
	if err != nil {
		return fmt.Errorf("rabbit channel: %w", err)
	}
	defer ch.Close()

	pub, err := NewPublisher(ch, cfg.Rabbit)
	if err != nil {
		return fmt.Errorf("publisher: %w", err)
	}

	app := &App{
		repo:      repository,
		publisher: pub,
		period:    period,
		window:    window,
	}

	return app.loop(ctx)
}

func (a *App) loop(ctx context.Context) error {
	ticker := time.NewTicker(a.period)
	defer ticker.Stop()

	// first tick immediately
	if err := a.tick(ctx); err != nil {
		slog.Error("tick failed: " + err.Error())
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := a.tick(ctx); err != nil {
				slog.Error("tick failed: " + err.Error())
			}
		}
	}
}

func (a *App) tick(ctx context.Context) error {
	events, err := a.repo.FetchMeetings(ctx, a.window, a.period)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}
	fmt.Println(events)
	if len(events) == 0 {
		return nil
	}
	if err = a.publisher.Publish(ctx, events); err != nil {
		return fmt.Errorf("publish: %w", err)
	}
	return nil
}

func dsn(db config.DBConfig) string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		db.Host, db.Port, db.Name, db.User, db.Pass)
}
