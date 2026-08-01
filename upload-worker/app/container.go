package app

import (
	"context"
	"log"

	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	consumer "github.com/Andhika-GIT/go-message-broker-monorepo/consumer_handler"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/S3_helper"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/rabbitmq"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/redis"
	"github.com/Andhika-GIT/go-message-broker-monorepo/repository"
	"github.com/Andhika-GIT/go-message-broker-monorepo/usecase"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ModuleDeps struct {
	DB           *gorm.DB
	RdsPublisher *redis.Publisher
	S3Helper     S3_helper.S3Helper
	Viper        *viper.Viper
	Ctx          context.Context
}

func wireUserModule(deps ModuleDeps) *usecase.UserUseCase {
	userUseCase := usecase.NewUserUseCase(&repository.UserRepository{}, deps.DB)

	handler := consumer.NewUserConsumerHandler(deps.RdsPublisher, userUseCase, deps.S3Helper)

	userConsumer := rabbitmq.NewRabbitMqConsumer(
		configs.UserImportConsumer(deps.Viper),
		handler.HandleMessage,
	)

	go func() {
		if err := userConsumer.Run(deps.Ctx); err != nil {
			log.Printf("user consumer stopped: %v", err)
		}
	}()

	return userUseCase
}

func wireOrderModule(deps ModuleDeps, userUseCase *usecase.UserUseCase) {
	orderUseCase := usecase.NewOrderUseCase(&repository.OrderRepository{}, deps.DB, userUseCase)

	handler := consumer.NewOrderConsumerHandler(deps.RdsPublisher, orderUseCase, deps.S3Helper)

	orderConsumer := rabbitmq.NewRabbitMqConsumer(
		configs.OrderImportConsumer(deps.Viper),
		handler.HandleMessage,
	)

	go func() {
		if err := orderConsumer.Run(deps.Ctx); err != nil {
			log.Printf("order consumer stopped: %v", err)
		}
	}()
}
