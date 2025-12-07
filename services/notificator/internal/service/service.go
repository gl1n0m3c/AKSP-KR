package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/streadway/amqp"

	"github.com/gl1n0m3c/AKSP-KR/services/notificator/config"
)

type App struct {
	consumer *Consumer
}

func Run(ctx context.Context, cfg *config.Config) error {
	conn, err := amqp.Dial(cfg.Rabbit.URL)
	if err != nil {
		return fmt.Errorf("rabbit dial: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("rabbit channel: %w", err)
	}
	defer ch.Close()

	consumer, err := NewConsumer(ch, cfg.Rabbit)
	if err != nil {
		return fmt.Errorf("consumer: %w", err)
	}

	app := &App{consumer: consumer}
	return app.loop(ctx)
}

func (a *App) loop(ctx context.Context) error {
	msgs, err := a.consumer.Consume()
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}
			slog.Info("received notification", "body", string(msg.Body))
			_ = msg.Ack(false)
		}
	}
}
