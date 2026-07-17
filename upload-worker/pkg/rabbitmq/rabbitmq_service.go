package rabbitmq

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMqConsumer struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
}

func NewRabbitMqConsumer(connectionUrl string) (*RabbitMqConsumer, error) {
	conn, err := amqp091.Dial(connectionUrl)

	if err != nil {
		return nil, err

	}

	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	return &RabbitMqConsumer{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (c *RabbitMqConsumer) Close() {
	if c.Channel != nil {
		_ = c.Channel.Close()
	}

	if c.Conn != nil {
		_ = c.Conn.Close()
	}
}

func (r *RabbitMqConsumer) DeclareExchange(exchangeName, exchangeType string) error {
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

func (r *RabbitMqConsumer) QueueBind(queueName, exchangeName, routingKey string) error {
	return r.Channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
}

func (c *RabbitMqConsumer) Consume(queueName string) (<-chan amqp091.Delivery, error) {

	msgs, err := c.Channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	if err != nil {
		return nil, fmt.Errorf("failed to consume message : %v", err.Error())
	}

	return msgs, nil
}
