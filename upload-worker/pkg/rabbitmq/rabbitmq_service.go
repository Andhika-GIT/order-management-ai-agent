package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMqConsumer struct {
	config  configs.RabbitMQConfig
	handler func(ctx context.Context, msg amqp091.Delivery)
}

func NewRabbitMqConsumer(config configs.RabbitMQConfig, handler func(ctx context.Context, msg amqp091.Delivery)) *RabbitMqConsumer {
	return &RabbitMqConsumer{
		config:  config,
		handler: handler,
	}
}

func (c *RabbitMqConsumer) Start(ctx context.Context) error {
	for {
		err := c.Run(ctx)

		if err == nil {
			return nil
		}

		log.Printf("consumer err, will reconnect in 5 seconds ...")

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(5 * time.Second):

		}

	}
}

func (c *RabbitMqConsumer) Run(ctx context.Context) error {
	conn, err := amqp091.Dial(c.config.URL)

	if err != nil {
		return fmt.Errorf("dial: %w", err)

	}

	ch, err := conn.Channel()

	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		c.config.ExchangeName,
		c.config.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("exchange declare: %w", err)
	}

	_, err = ch.QueueDeclare(c.config.Queue, true, false, false, false, nil)

	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	msgs, err := ch.Consume(c.config.Queue, "", false, false, false, false, nil)

	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	connClose := conn.NotifyClose(make(chan *amqp091.Error, 1))

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-connClose:
			return fmt.Errorf("connection closed: %w", err)
		case msg := <-msgs:
			c.handler(ctx, msg)
		}
	}

}
