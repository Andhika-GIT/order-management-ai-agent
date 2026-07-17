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

type OrderConsumer struct {
	Rmq          *rabbitmq.RabbitMqConsumer
	RdsPublisher *redis.Publisher
	UseCase      *usecase.OrderUseCase
	QueueCfg     *configs.RabbitMQQueue
	sftpClient   *sftp.Client
}

func NewOrderConsumer(Rmq *rabbitmq.RabbitMqConsumer, RdsPublisher *redis.Publisher, UseCase *usecase.OrderUseCase, cfg *configs.RabbitMQQueue, sftpClient *sftp.Client) *OrderConsumer {
	return &OrderConsumer{
		Rmq:          Rmq,
		RdsPublisher: RdsPublisher,
		UseCase:      UseCase,
		QueueCfg:     cfg,
		sftpClient:   sftpClient,
	}
}

func (w *OrderConsumer) Start() {
	defer w.Rmq.Close()

	c := context.Background()

	msgs, err := w.Rmq.Consume(w.QueueCfg.OrderDirectImport)

	if err != nil {
		log.Println(err.Error())
	}

	var uploadMsg model.UploadMessage
	for msg := range msgs {
		err := json.Unmarshal(msg.Body, &uploadMsg)

		if err != nil {
			log.Panicln(err.Error())
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

		orders := w.UseCase.ReadOrderExcel(rows)

		err = w.UseCase.CreateOrders(c, orders)

		if err != nil {
			log.Print(err.Error())
		}

		err = w.RdsPublisher.PublishMessage(c, "notifications", fmt.Sprintf("successfully uploaded %s", uploadMsg.Filename))

		if err != nil {
			log.Print(err.Error())
		}

		log.Printf("rows are %v", rows)

	}
}
