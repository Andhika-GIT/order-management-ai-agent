package rabbitmq

import (
	"context"
	"fmt"
	"sync"

	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	config configs.RabbitMQConfig
	conn   *amqp091.Connection
	mu     sync.Mutex
}

func NewPublisher(config configs.RabbitMQConfig) (*Publisher, error) {
	p := &Publisher{config: config}

	if err := p.connect(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Publisher) connect() error {
	conn, err := amqp091.Dial(p.config.URL)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(
		p.config.ExchangeName,
		p.config.ExchangeType,
		true, false, false, false, nil,
	); err != nil {
		return fmt.Errorf("exchange declare: %w", err)
	}

	p.conn = conn
	return nil
}

func (p *Publisher) Publish(ctx context.Context, routingKey string, body []byte) error {
	p.mu.Lock()
	if p.conn == nil || p.conn.IsClosed() {
		if err := p.connect(); err != nil {
			p.mu.Unlock()
			return fmt.Errorf("reconnect: %w", err)
		}
	}
	conn := p.conn
	p.mu.Unlock()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}
	defer ch.Close()

	return ch.PublishWithContext(ctx,
		p.config.ExchangeName,
		routingKey,
		false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
