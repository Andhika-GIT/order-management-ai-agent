package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/excel"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/redis"
	"github.com/Andhika-GIT/go-message-broker-monorepo/usecase"
	"github.com/pkg/sftp"
	"github.com/rabbitmq/amqp091-go"
)

type UserConsumer struct {
	RdsPublisher *redis.Publisher
	UseCase      *usecase.UserUseCase
	sftpClient   *sftp.Client
}

func NewUserConsumerHandler(RdsPublisher *redis.Publisher, UseCase *usecase.UserUseCase, sftpClient *sftp.Client) *UserConsumer {
	return &UserConsumer{

		RdsPublisher: RdsPublisher,
		UseCase:      UseCase,
		sftpClient:   sftpClient,
	}
}

func (w *UserConsumer) HandleMessage(c context.Context, msg amqp091.Delivery) {

	var uploadMsg model.UploadMessage

	err := json.Unmarshal(msg.Body, &uploadMsg)

	if err != nil {
		log.Println("error when converting message json", err.Error())
		_ = msg.Nack(false, false)
		return
	}

	remoteFile, err := w.sftpClient.Open(uploadMsg.Filepath)

	if err != nil {
		log.Printf("error when reading sftp file: %v", err)
		_ = msg.Nack(false, true)
		return
	}

	rows, err := excel.ReadExcel(remoteFile)

	if err != nil {
		log.Printf("error when reading excel file: %v", err)
		_ = msg.Nack(false, true)
		return
	}

	newUsers := w.UseCase.ReadUsersExcel(rows)

	err = w.UseCase.CreateNewUsers(c, newUsers)

	if err != nil {
		log.Printf("error when creating users from excel: %v", err)
		_ = msg.Nack(false, true)
		return
	}

	err = w.RdsPublisher.PublishMessage(c, "notifications", fmt.Sprintf("successfully uploaded %s", uploadMsg.Filename))

	if err != nil {
		log.Printf("error when publishing message: %v", err)
	}

	_ = msg.Ack(false)
}
