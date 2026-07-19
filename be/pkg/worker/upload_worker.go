package worker

import (
	"context"
	"encoding/json"
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
	publisher  *rabbitmq.Publisher
	sftpClient *sftp.Client
	numWorkers int
}

func NewUploadWorker(sftp *sftp.Client, publisher *rabbitmq.Publisher, numWorkers int) *UploadWorker {
	return &UploadWorker{
		sftpClient: sftp,
		publisher:  publisher,
		numWorkers: numWorkers,
	}
}

func (w *UploadWorker) Start(ctx context.Context) {
	for i := 0; i < w.numWorkers; i++ {
		go func() {
			for task := range UploadQueue {
				w.processUpload(ctx, task)
			}
		}()
	}
}

func (w *UploadWorker) processUpload(ctx context.Context, task UploadTask) {

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

	body, err := json.Marshal(uploadMsg)

	if err != nil {
		log.Printf("failed to marshal upload message: %v", err)
		return
	}

	if err := w.publisher.Publish(ctx, task.QueueRoutingKey, body); err != nil {
		log.Printf("failed to publish upload message: %v", err)
	}

}

func (w *UploadWorker) Queue(task UploadTask) {
	UploadQueue <- task
}
