package app

import (
	"log"

	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"github.com/Andhika-GIT/go-message-broker-monorepo/database"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/rabbitmq"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/sftpclient"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/worker"
	"github.com/go-chi/chi/v5"
)

func InitApp() *chi.Mux {
	r := NewRouter()

	v := configs.NewViper()
	cfg := configs.InitConfig(v)
	db := database.NewDatabase(&cfg.Database)
	rmq, err := rabbitmq.NewRabbitMqProducer(cfg.RabbitMQConnectURL)

	if err != nil {
		log.Fatalf("failed to initialize RabbitMQ connection: %v", err)
	}

	err = InitQueue(rmq, cfg)

	if err != nil {
		log.Fatalf("failed to bind RabbitMQ queues: %v", err)

	}

	sftpClient, err := sftpclient.NewSFTPClient(&cfg.SftpClient)

	if err != nil {
		log.Fatalf("failed to bind to sftp client %v", err)
	}

	uploadWorker := worker.NewUploadWorker(sftpClient, rmq, 3)
	orderUseCase := wireOrderModule(r, rmq, uploadWorker, db, cfg)
	userUseCase := wireUserModule(r, rmq, uploadWorker, db, cfg)
	wireDashboardModule(r, userUseCase, orderUseCase)

	go uploadWorker.Start()

	return r
}
