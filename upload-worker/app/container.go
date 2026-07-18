package app

import (
	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"github.com/Andhika-GIT/go-message-broker-monorepo/consumer"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/rabbitmq"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/redis"
	"github.com/Andhika-GIT/go-message-broker-monorepo/repository"
	"github.com/Andhika-GIT/go-message-broker-monorepo/usecase"
	"github.com/pkg/sftp"
	"gorm.io/gorm"
)

func wireUserModule(rmq *rabbitmq.RabbitMqConsumer, rdsPublisher *redis.Publisher, db *gorm.DB, queueCfg *configs.RabbitMQQueue, sftpClient *sftp.Client) *usecase.UserUseCase {
	userUseCase := usecase.NewUserUseCase(&repository.UserRepository{}, db)

	c := consumer.NewUserConsumer(rmq, rdsPublisher, userUseCase, queueCfg, sftpClient)

	go c.Start()

	return userUseCase
}

func wireOrderModule(rmq *rabbitmq.RabbitMqConsumer, rdsPublisher *redis.Publisher, db *gorm.DB, userUseCase *usecase.UserUseCase, queueCfg *configs.RabbitMQQueue, sftpClient *sftp.Client) {
	orderUseCase := usecase.NewOrderUseCase(&repository.OrderRepository{}, db, userUseCase)

	c := consumer.NewOrderConsumer(rmq, rdsPublisher, orderUseCase, queueCfg, sftpClient)

	go c.Start()
}
