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

func wireProductModule(deps ModuleDeps) {
	productUseCase := usecase.NewProductUseCase(&repository.ProductRepository{}, deps.DB)

	handler := consumer.NewProductConsumerHandler(deps.RdsPublisher, productUseCase, deps.S3Helper)

	productConsumer := rabbitmq.NewRabbitMqConsumer(
		configs.ProductImportConsumer(deps.Viper),
		handler.HandleMessage,
	)

	go func() {
		if err := productConsumer.Run(deps.Ctx); err != nil {
			log.Printf("product consumer stopped: %v", err)
		}
	}()
}
