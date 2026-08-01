package app

import (
	"context"
	"log"

	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"github.com/Andhika-GIT/go-message-broker-monorepo/database"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/S3_helper"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/rabbitmq"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/worker"
	"github.com/go-chi/chi/v5"
)

func InitApp() *chi.Mux {
	r := NewRouter()
	ctx := context.Background()

	v := configs.NewViper()
	cfg := configs.InitConfig(v)
	db := database.NewDatabase(&cfg.Database)

	publisher, err := rabbitmq.NewPublisher(configs.Producer(v))

	if err != nil {
		log.Fatalf("failed to initialize RabbitMQ publisher: %v", err)
	}

	s3Client, err := database.NewS3Client(cfg.S3Config, ctx)

	if err != nil {
		log.Fatalf("failed to connect to s3: %v", err)
	}

	s3Helper := S3_helper.NewS3Helper(s3Client, cfg.S3Config.Bucket)

	uploadWorker := worker.NewUploadWorker(s3Helper, publisher, 3)
	uploadWorker.Start(ctx)

	deps := ModuleDeps{
		Router:       r,
		DB:           db,
		UploadWorker: uploadWorker,
		config:       cfg,
		S3Helper:     s3Helper,
	}

	orderUseCase := wireOrderModule(deps)
	userUseCase := wireUserModule(deps)
	wireDashboardModule(deps, userUseCase, orderUseCase)
	wireUploadModule(deps, orderUseCase)

	return r
}
