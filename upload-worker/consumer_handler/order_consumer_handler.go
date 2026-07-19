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

type OrderConsumer struct {
	RdsPublisher *redis.Publisher
	UseCase      *usecase.OrderUseCase
	sftpClient   *sftp.Client
}

func NewOrderConsumerHandler(RdsPublisher *redis.Publisher, UseCase *usecase.OrderUseCase, sftpClient *sftp.Client) *OrderConsumer {
	return &OrderConsumer{
		RdsPublisher: RdsPublisher,
		UseCase:      UseCase,
		sftpClient:   sftpClient,
	}
}

func (w *OrderConsumer) HandleMessage(c context.Context, msg amqp091.Delivery) {

	var uploadMsg model.UploadMessage

	err := json.Unmarshal(msg.Body, &uploadMsg)

	if err != nil {
		log.Print(err.Error())
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

	orders := w.UseCase.ReadOrderExcel(rows)

	err = w.UseCase.CreateOrders(c, orders)

	if err != nil {
		log.Print(err.Error())
		return
	}

	err = w.RdsPublisher.PublishMessage(c, "notifications", fmt.Sprintf("successfully uploaded %s", uploadMsg.Filename))

	if err != nil {
		log.Print(err.Error())
	}

}
