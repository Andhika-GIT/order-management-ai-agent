package worker

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/S3_helper"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/rabbitmq"
)

type UploadTask struct {
	File            io.Reader
	Key             string
	Bucket          string
	ContentType     string
	QueueRoutingKey string
}

var UploadQueue = make(chan UploadTask, 100)

type UploadWorker struct {
	publisher  *rabbitmq.Publisher
	s3Helper   S3_helper.S3Helper
	numWorkers int
}

func NewUploadWorker(s3Helper S3_helper.S3Helper, publisher *rabbitmq.Publisher, numWorkers int) *UploadWorker {
	return &UploadWorker{
		s3Helper:   s3Helper,
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
	err := w.s3Helper.Upload(ctx, task.Key, task.File, task.ContentType)

	if err != nil {
		log.Printf("failed to upload to s3: %v", err)
		return
	}

	uploadMsg := model.UploadMessage{
		Key:    task.Key,
		Bucket: task.Bucket,
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
