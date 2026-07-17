package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

const ExchangeGoApp = "go-app-exchange"

type RabbitMqProducer struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
}

func NewRabbitMqProducer(connectionUrl string) (*RabbitMqProducer, error) {
	conn, err := amqp091.Dial(connectionUrl)

	if err != nil {
		return nil, err

	}

	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	return &RabbitMqProducer{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (c *RabbitMqProducer) Close() {
	if c.Channel != nil {
		_ = c.Channel.Close()
	}

	if c.Conn != nil {
		_ = c.Conn.Close()
	}
}

func (r *RabbitMqProducer) DeclareExchange(exchangeName, exchangeType string) error {
	return r.Channel.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMqProducer) QueueBind(queueName, exchangeName, routingKey string) error {
	return r.Channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
}

func (c *RabbitMqProducer) Publish(routingKey string, payload any) error {
	body, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	err = c.Channel.Publish(
		ExchangeGoApp, routingKey, false, false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}
