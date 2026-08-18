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

type UserConsumer struct {
	RdsPublisher *redis.Publisher
	UseCase      *usecase.UserUseCase
	s3Helper     S3_helper.S3Helper
}

func NewUserConsumerHandler(RdsPublisher *redis.Publisher, UseCase *usecase.UserUseCase, s3Helper S3_helper.S3Helper) *UserConsumer {
	return &UserConsumer{

		RdsPublisher: RdsPublisher,
		UseCase:      UseCase,
		s3Helper:     s3Helper,
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

	newUsers := w.UseCase.ReadUsersExcel(rows)

	err = w.UseCase.CreateNewUsers(c, newUsers)

	if err != nil {
		log.Printf("error when creating users from excel: %v", err)
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
