package app

import (
	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"github.com/Andhika-GIT/go-message-broker-monorepo/controller"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/rabbitmq"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/worker"
	"github.com/Andhika-GIT/go-message-broker-monorepo/repository"
	"github.com/Andhika-GIT/go-message-broker-monorepo/usecase"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func wireOrderModule(r chi.Router, rmq *rabbitmq.RabbitMqProducer, uploadWorker *worker.UploadWorker, db *gorm.DB, cfg *configs.Config) *usecase.OrderUseCase {
	repo := repository.NewOrderRepository(db)
	uc := usecase.NewOrderUseCase(repo, rmq)
	ctrl := controller.NewOrderController(uc, uploadWorker, &cfg.RabbitMQRoutingKey, cfg.SftpClient.Path)
	registerOrderRoutes(r, ctrl)

	return uc
}

func wireUserModule(r chi.Router, rmq *rabbitmq.RabbitMqProducer, uploadWorker *worker.UploadWorker, db *gorm.DB, cfg *configs.Config) *usecase.UserUseCase {
	repo := repository.NewUserRepository(db)
	uc := usecase.NewUserUseCase(repo, rmq, db)
	ctrl := controller.NewUserController(uc, uploadWorker, &cfg.RabbitMQRoutingKey, cfg.SftpClient.Path)
	registerUserRoutes(r, ctrl)

	return uc
}

func wireDashboardModule(r chi.Router, userUseCase *usecase.UserUseCase, orderUseCase *usecase.OrderUseCase) {
	ctrl := controller.NewDashboardController(userUseCase, orderUseCase)
	registerDashboardRoutes(r, ctrl)
}
