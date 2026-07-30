package app

import (
	"context"
	"log"

	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"github.com/Andhika-GIT/go-message-broker-monorepo/database"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/S3_helper"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/rabbitmq"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/sftpclient"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/worker"
	"github.com/go-chi/chi/v5"
)

func InitApp() *chi.Mux {
	r := NewRouter()
	ctx := context.Background()

	v := configs.NewViper()
	cfg := configs.InitConfig(v)
	db := database.NewDatabase(&cfg.Database)

	sftpClient, err := sftpclient.NewSFTPClient(&cfg.SftpClient)

	if err != nil {
		log.Fatalf("failed to bind to sftp client %v", err)
	}

	publisher, err := rabbitmq.NewPublisher(configs.Producer(v))

	if err != nil {
		log.Fatalf("failed to initialize RabbitMQ publisher: %v", err)
	}

	uploadWorker := worker.NewUploadWorker(sftpClient, publisher, 3)
	uploadWorker.Start(ctx)

	s3Client, err := database.NewS3Client(cfg.S3Config, ctx)

	if err != nil {
		log.Fatalf("failed to connect to s3: %v", err)
	}

	s3Helper := S3_helper.NewS3Helper(s3Client, cfg.S3Config.Bucket)

	deps := ModuleDeps{
		Router:       r,
		DB:           db,
		UploadWorker: uploadWorker,
		RoutingKey:   &cfg.RabbitMQRoutingKey,
		SftpPath:     cfg.SftpClient.Path,
		S3Helper:     s3Helper,
	}

	orderUseCase := wireOrderModule(deps)
	userUseCase := wireUserModule(deps)
	wireDashboardModule(deps, userUseCase, orderUseCase)
	wireUploadModule(deps, orderUseCase)

	return r
}
