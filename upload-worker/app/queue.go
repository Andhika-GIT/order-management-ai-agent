package app

import (
	"log"

	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/rabbitmq"
)

func InitQueue(rmq *rabbitmq.RabbitMqConsumer, cfg *configs.Config) error {
	log.Println("Initializing RabbitMQ...")

	// Declare exchange
	if err := rmq.DeclareExchange(cfg.RabbitMQExchange, "direct"); err != nil {
		return err
	}

	// --- User Queues ---
	if err := rmq.QueueBind(cfg.RabbitMQQueue.UserDirectImport, cfg.RabbitMQExchange, cfg.RabbitMQRoutingKey.UserDirectImport); err != nil {
		return err
	}
	// --- Order Queues ---
	if err := rmq.QueueBind(cfg.RabbitMQQueue.OrderDirectImport, cfg.RabbitMQExchange, cfg.RabbitMQRoutingKey.OrderDirectImport); err != nil {
		return err
	}

	return nil
}
