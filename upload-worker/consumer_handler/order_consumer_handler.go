package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/S3_helper"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/excel"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/redis"
	"github.com/Andhika-GIT/go-message-broker-monorepo/usecase"
	"github.com/rabbitmq/amqp091-go"
)

type OrderConsumer struct {
	RdsPublisher *redis.Publisher
	UseCase      *usecase.OrderUseCase
	s3Helper     S3_helper.S3Helper
}

func NewOrderConsumerHandler(RdsPublisher *redis.Publisher, UseCase *usecase.OrderUseCase, s3Helper S3_helper.S3Helper) *OrderConsumer {
	return &OrderConsumer{
		RdsPublisher: RdsPublisher,
		UseCase:      UseCase,
		s3Helper:     s3Helper,
	}
}

func (w *OrderConsumer) HandleMessage(c context.Context, msg amqp091.Delivery) {

	var uploadMsg model.UploadMessage

	err := json.Unmarshal(msg.Body, &uploadMsg)

	if err != nil {
		log.Println("error when converting message json", err.Error())
		_ = msg.Nack(false, false)
		return
	}

	data, err := w.s3Helper.Download(c, uploadMsg.Key)

	if err != nil {
		log.Printf("error when downloading s3 file: %v", err)
		_ = msg.Nack(false, true)
		return
	}

	rows, err := excel.ReadExcel(bytes.NewReader(data))

	if err != nil {
		log.Printf("error when reading excel file: %v", err)
		_ = msg.Nack(false, true)
		return
	}

	newOrders := w.UseCase.ReadOrderExcel(rows)

	err = w.UseCase.CreateOrders(c, newOrders)

	if err != nil {
		log.Printf("error when creating orders from excel: %v", err)
		_ = msg.Nack(false, true)
		return
	}

	if err := w.s3Helper.Delete(c, uploadMsg.Key); err != nil {
		log.Printf("error when deleting s3 file: %v", err)
	}

	err = w.RdsPublisher.PublishMessage(c, "notifications", fmt.Sprintf("successfully uploaded %s", path.Base(uploadMsg.Key)))

	if err != nil {
		log.Printf("error when publishing message: %v", err)
	}

	_ = msg.Ack(false)

}
