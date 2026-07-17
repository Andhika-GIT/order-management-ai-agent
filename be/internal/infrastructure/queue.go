package infrastructure

import (
	"log"

	"github.com/Andhika-GIT/go-message-broker-monorepo/internal/shared"
)

func InitQueue(rmq *shared.RabbitMqProducer, cfg *shared.Config) error {
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
