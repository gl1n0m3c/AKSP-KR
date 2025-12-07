package service

import (
	"fmt"

	"github.com/streadway/amqp"

	"github.com/gl1n0m3c/AKSP-KR/services/notificator/config"
)

type Consumer struct {
	ch      *amqp.Channel
	cfg     config.RabbitConfig
	queue   string
	tag     string
	closing chan struct{}
}

func NewConsumer(ch *amqp.Channel, cfg config.RabbitConfig) (*Consumer, error) {
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

	return &Consumer{
		ch:      ch,
		cfg:     cfg,
		queue:   q.Name,
		tag:     cfg.Consumer,
		closing: make(chan struct{}),
	}, nil
}

func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := c.ch.Consume(
		c.queue,
		c.tag,
		false, // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	return msgs, nil
}
