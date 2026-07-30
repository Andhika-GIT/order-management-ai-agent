package S3_helper

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Helper interface {
	GetPresignedInsertURL(ctx context.Context, key, contentType string, expiry time.Duration) (string, error) // generate URL for FE to upload directly to S3
}

type s3Helper struct {
	client *s3.Client
	bucket string
}

func NewS3Helper(client *s3.Client, bucket string) S3Helper {
	return &s3Helper{client: client, bucket: bucket}
}

// GetPresignedInsertURL implements [S3Helper].
func (h *s3Helper) GetPresignedInsertURL(ctx context.Context, key string, contentType string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(h.client)
	req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      &h.bucket,
		Key:         &key,
		ContentType: &contentType,
	}, s3.WithPresignExpires(expiry))

	if err != nil {
		return "", err
	}

	return req.URL, nil
}
