package order

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Andhika-GIT/go-message-broker-monorepo/internal/shared"
	redispubsub "github.com/Andhika-GIT/go-message-broker-monorepo/internal/shared/redis"
	"github.com/pkg/sftp"
)

type OrderDirectWorker struct {
	Rmq          *shared.RabbitMqConsumer
	RdsPublisher *redispubsub.Publisher
	UseCase      *OrderUseCase
	QueueCfg     *shared.RabbitMQQueue
	sftpClient   *sftp.Client
}

func NewOrderDirectWorker(Rmq *shared.RabbitMqConsumer, RdsPublisher *redispubsub.Publisher, UseCase *OrderUseCase, cfg *shared.RabbitMQQueue, sftpClient *sftp.Client) *OrderDirectWorker {
	return &OrderDirectWorker{
		Rmq:          Rmq,
		RdsPublisher: RdsPublisher,
		UseCase:      UseCase,
		QueueCfg:     cfg,
		sftpClient:   sftpClient,
	}
}

func (w *OrderDirectWorker) Start() {

	c := context.Background()

	msgs, err := w.Rmq.Consume(w.QueueCfg.OrderDirectImport)

	if err != nil {
		log.Println(err.Error())
	}

	var uploadMsg shared.UploadMessage
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

		rows, err := shared.ReadExcel(remoteFile)

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
