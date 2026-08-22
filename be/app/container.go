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
	config       *configs.Config
	S3Helper     S3_helper.S3Helper
}

func wireProductModule(deps ModuleDeps) *usecase.ProductUseCase {
	repo := repository.NewProductRepository(deps.DB)
	uc := usecase.NewProductUseCase(repo)
	ctrl := controller.NewProductController(uc, deps.UploadWorker, deps.config)
	registerProductRoutes(deps.Router, ctrl)

	return uc
}

func wireDashboardModule(deps ModuleDeps, productUseCase *usecase.ProductUseCase) {
	ctrl := controller.NewDashboardController(productUseCase)
	registerDashboardRoutes(deps.Router, ctrl)
}

func wireUploadModule(deps ModuleDeps, productUsecase *usecase.ProductUseCase) {
	ctrl := controller.NewUploadController(deps.S3Helper, *productUsecase)
	registerUploadRoutes(deps.Router, ctrl)
}
