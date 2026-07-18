package worker

import (
	"fmt"
	"io"
	"log"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/rabbitmq"
	"github.com/pkg/sftp"
)

type UploadTask struct {
	File            io.Reader
	Filename        string
	Filepath        string
	QueueRoutingKey string
}

var UploadQueue = make(chan UploadTask, 100)

type UploadWorker struct {
	rmq        *rabbitmq.RabbitMqProducer
	sftpClient *sftp.Client
	numWorkers int
}

func NewUploadWorker(sftp *sftp.Client, rmq *rabbitmq.RabbitMqProducer, numWorkers int) *UploadWorker {
	return &UploadWorker{
		sftpClient: sftp,
		rmq:        rmq,
		numWorkers: numWorkers,
	}
}

func (w *UploadWorker) Start() {
	for i := 0; i < w.numWorkers; i++ {
		go func() {
			for task := range UploadQueue {
				w.processUpload(task)
			}
		}()
	}
}

func (w *UploadWorker) processUpload(task UploadTask) {

	path := fmt.Sprintf("%v/%s", task.Filepath, task.Filename)

	log.Print(path)

	dstFile, err := w.sftpClient.Create(path)

	if err != nil {
		log.Fatalf("error when connecting to sftp : %s", err.Error())
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, task.File)

	if err != nil {
		log.Fatalf("error when insert file to sftp %s", err.Error())
	}

	uploadMsg := model.UploadMessage{
		Filename: task.Filename,
		Filepath: path,
	}

	w.rmq.Publish(task.QueueRoutingKey, uploadMsg)

}

func (w *UploadWorker) Queue(task UploadTask) {
	UploadQueue <- task
}
