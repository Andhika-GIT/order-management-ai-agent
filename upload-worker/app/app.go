package app

import (
	"context"
	"log"

	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"github.com/Andhika-GIT/go-message-broker-monorepo/database"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/S3_helper"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/redis"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/sftpclient"
	"github.com/go-chi/chi/v5"
)

func InitApp() *chi.Mux {
	r := chi.NewRouter()
	ctx := context.Background()

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

	redisClient, err := database.NewRedisClient(&cfg.RedisClient)
	if err != nil {
		log.Printf("failed to connect to redis: %v", err)

	}

	redisPublisher := redis.NewPublisher(redisClient)

	s3Client, err := database.NewS3Client(cfg.S3Config, ctx)

	if err != nil {
		log.Printf("failed to connect to s3: %v", err)
	}

	S3_helper.NewS3Helper(s3Client, cfg.S3Config.Bucket)

	deps := ModuleDeps{
		DB:           db,
		RdsPublisher: redisPublisher,
		SftpClient:   sftpClient,
		Viper:        v,
		Ctx:          ctx,
	}

	userUseCase := wireUserModule(deps)
	wireOrderModule(deps, userUseCase)

	return r
}
