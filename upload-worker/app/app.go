package app

import (
	"log"

	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"github.com/Andhika-GIT/go-message-broker-monorepo/database"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/rabbitmq"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/redis"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/sftpclient"
	"github.com/go-chi/chi/v5"
)

func InitApp() *chi.Mux {
	r := chi.NewRouter()

	v, err := configs.NewViper()

	if err != nil {
		log.Print(err.Error())
	}

	cfg := configs.InitConfig(v)

	db, err := database.NewDatabase(&cfg.Database)

	if err != nil {
		log.Print(err.Error())
	}

	sftpClient, err := sftpclient.NewSFTPClient(&cfg.SftpClient)

	if err != nil {
		log.Printf("failed to initialize sftp: %v", err)
	}

	rmq, err := rabbitmq.NewRabbitMqConsumer(cfg.RabbitMQConnectURL)

	if err != nil {
		log.Printf("failed to initialize RabbitMQ connection: %v", err)
	}

	err = InitQueue(rmq, cfg)

	if err != nil {
		log.Printf("failed to bind RabbitMQ queues: %v", err)

	}

	redisClient, err := redis.NewRedisClient(&cfg.RedisClient)
	if err != nil {
		log.Printf("failed to connect to redis: %v", err)

	}

	redisPublisher := redis.NewPublisher(redisClient)

	userUseCase := wireUserModule(rmq, redisPublisher, db, &cfg.RabbitMQQueue, sftpClient)
	wireOrderModule(rmq, redisPublisher, db, userUseCase, &cfg.RabbitMQQueue, sftpClient)

	return r
}
