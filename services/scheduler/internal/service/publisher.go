package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/streadway/amqp"

	"github.com/gl1n0m3c/AKSP-KR/services/scheduler/config"
)

type Publisher interface {
	Publish(ctx context.Context, events []MeetingNotification) error
}

type publisher struct {
	ch        *amqp.Channel
	exchange  string
	routing   string
	queueName string
}

func NewPublisher(ch *amqp.Channel, cfg config.RabbitConfig) (Publisher, error) {
	q, err := ch.QueueDeclare(
		cfg.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("queue declare: %w", err)
	}

	if cfg.Exchange != "" {
		if err = ch.ExchangeDeclare(cfg.Exchange, "direct", true, false, false, false, nil); err != nil {
			return nil, fmt.Errorf("exchange declare: %w", err)
		}
		if err = ch.QueueBind(q.Name, cfg.RoutingKey, cfg.Exchange, false, nil); err != nil {
			return nil, fmt.Errorf("queue bind: %w", err)
		}
	}

	return &publisher{
		ch:        ch,
		exchange:  cfg.Exchange,
		routing:   cfg.RoutingKey,
		queueName: q.Name,
	}, nil
}

func (p *publisher) Publish(ctx context.Context, events []MeetingNotification) error {
	for _, ev := range events {
		body, err := marshalEvent(ev)
		if err != nil {
			return fmt.Errorf("marshal: %w", err)
		}
		if err = p.ch.Publish(
			p.exchange,
			p.routing,
			false,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				DeliveryMode: amqp.Persistent,
				Timestamp:    time.Now(),
			},
		); err != nil {
			return fmt.Errorf("publish: %w", err)
		}
		slog.Info("publish event", ev.ID)
	}
	return nil
}

func (p *publisher) Close() error {
	return nil
}
