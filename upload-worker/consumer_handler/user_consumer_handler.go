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
		log.Println("error when converting", err.Error())
		return
	}

	remoteFile, err := w.sftpClient.Open(uploadMsg.Filepath)

	if err != nil {
		log.Printf("error when reading sftp file: %v", err)
		return
	}

	rows, err := excel.ReadExcel(remoteFile)

	if err != nil {
		log.Print(err.Error())
		return
	}

	newUsers := w.UseCase.ReadUsersExcel(rows)

	err = w.UseCase.CreateNewUsers(c, newUsers)

	if err != nil {
		log.Print(err.Error())
		return
	}

	err = w.RdsPublisher.PublishMessage(c, "notifications", fmt.Sprintf("successfully uploaded %s", uploadMsg.Filename))

	if err != nil {
		log.Print(err.Error())
	}

}
