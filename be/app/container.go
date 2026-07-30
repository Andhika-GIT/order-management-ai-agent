package app

import (
	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"github.com/Andhika-GIT/go-message-broker-monorepo/controller"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/S3_helper"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/worker"
	"github.com/Andhika-GIT/go-message-broker-monorepo/repository"
	"github.com/Andhika-GIT/go-message-broker-monorepo/usecase"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type ModuleDeps struct {
	Router       chi.Router
	DB           *gorm.DB
	UploadWorker *worker.UploadWorker
	RoutingKey   *configs.RabbitMQRoutingKey
	SftpPath     string
	S3Helper     S3_helper.S3Helper
}

func wireOrderModule(deps ModuleDeps) *usecase.OrderUseCase {
	repo := repository.NewOrderRepository(deps.DB)
	uc := usecase.NewOrderUseCase(repo)
	ctrl := controller.NewOrderController(uc, deps.UploadWorker, deps.RoutingKey, deps.SftpPath)
	registerOrderRoutes(deps.Router, ctrl)

	return uc
}

func wireUserModule(deps ModuleDeps) *usecase.UserUseCase {
	repo := repository.NewUserRepository(deps.DB)
	uc := usecase.NewUserUseCase(repo, deps.DB)
	ctrl := controller.NewUserController(uc, deps.UploadWorker, deps.RoutingKey, deps.SftpPath)
	registerUserRoutes(deps.Router, ctrl)

	return uc
}

func wireDashboardModule(deps ModuleDeps, userUseCase *usecase.UserUseCase, orderUseCase *usecase.OrderUseCase) {
	ctrl := controller.NewDashboardController(userUseCase, orderUseCase)
	registerDashboardRoutes(deps.Router, ctrl)
}

func wireUploadModule(deps ModuleDeps) {
	ctrl := controller.NewUploadController(deps.S3Helper)
	registerUploadRoutes(deps.Router, ctrl)
}
