package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/excel"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/rabbitmq"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/redis"
	"github.com/Andhika-GIT/go-message-broker-monorepo/usecase"
	"github.com/pkg/sftp"
)

type UserConsumer struct {
	Rmq          *rabbitmq.RabbitMqConsumer
	RdsPublisher *redis.Publisher
	UseCase      *usecase.UserUseCase
	QueueCfg     *configs.RabbitMQQueue
	sftpClient   *sftp.Client
}

func NewUserConsumer(Rmq *rabbitmq.RabbitMqConsumer, RdsPublisher *redis.Publisher, UseCase *usecase.UserUseCase, cfg *configs.RabbitMQQueue, sftpClient *sftp.Client) *UserConsumer {
	return &UserConsumer{
		Rmq:          Rmq,
		RdsPublisher: RdsPublisher,
		UseCase:      UseCase,
		QueueCfg:     cfg,
		sftpClient:   sftpClient,
	}
}

func (w *UserConsumer) Start() {

	c := context.Background()

	msgs, err := w.Rmq.Consume(w.QueueCfg.UserDirectImport)

	if err != nil {
		log.Println(err)
	}

	var uploadMsg model.UploadMessage
	for msg := range msgs {
		err := json.Unmarshal(msg.Body, &uploadMsg)

		if err != nil {
			log.Println("error when converting", err.Error())
			continue
		}

		remoteFile, err := w.sftpClient.Open(uploadMsg.Filepath)

		if err != nil {
			log.Printf("error when reading sftp file: %v", err)
			continue
		}

		rows, err := excel.ReadExcel(remoteFile)

		if err != nil {
			log.Print(err.Error())
		}

		newUsers := w.UseCase.ReadUsersExcel(rows)

		err = w.UseCase.CreateNewUsers(c, newUsers)

		if err != nil {
			log.Print(err.Error())
		}

		err = w.RdsPublisher.PublishMessage(c, "notifications", fmt.Sprintf("successfully uploaded %s", uploadMsg.Filename))

		if err != nil {
			log.Print(err.Error())
		}

		log.Printf("filepath is %s", uploadMsg.Filepath)

	}

}
