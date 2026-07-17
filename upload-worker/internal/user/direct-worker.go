package user

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Andhika-GIT/go-message-broker-monorepo/internal/shared"
	redispubsub "github.com/Andhika-GIT/go-message-broker-monorepo/internal/shared/redis"
	"github.com/pkg/sftp"
)

type UserDirectUploadWorker struct {
	Rmq          *shared.RabbitMqConsumer
	RdsPublisher *redispubsub.Publisher
	UseCase      *UserUseCase
	QueueCfg     *shared.RabbitMQQueue
	sftpClient   *sftp.Client
}

func NewUserDirectUploadWorker(Rmq *shared.RabbitMqConsumer, RdsPublisher *redispubsub.Publisher, UseCase *UserUseCase, cfg *shared.RabbitMQQueue, sftpClient *sftp.Client) *UserDirectUploadWorker {
	return &UserDirectUploadWorker{
		Rmq:          Rmq,
		RdsPublisher: RdsPublisher,
		UseCase:      UseCase,
		QueueCfg:     cfg,
		sftpClient:   sftpClient,
	}
}

func (w *UserDirectUploadWorker) Start() {

	c := context.Background()

	msgs, err := w.Rmq.Consume(w.QueueCfg.UserDirectImport)

	if err != nil {
		log.Println(err)
	}

	var uploadMsg shared.UploadMessage
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

		rows, err := shared.ReadExcel(remoteFile)

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
