package S3_helper

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Helper interface {
	Upload(ctx context.Context, key string, body io.Reader, contentType string) error
	GetPresignedInsertURL(ctx context.Context, key, contentType string, expiry time.Duration) (string, error) // generate URL for FE to upload directly to S3
	Delete(ctx context.Context, key string) error
}

type s3Helper struct {
	client         *s3.Client
	bucket         string
	publicEndpoint string
}

func NewS3Helper(client *s3.Client, bucket string, publicEndpoint string) S3Helper {
	return &s3Helper{client: client, bucket: bucket, publicEndpoint: publicEndpoint}
}

// Upload implements [S3Helper].
func (h *s3Helper) Upload(ctx context.Context, key string, body io.Reader, contentType string) error {
	_, err := h.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &h.bucket,
		Key:         &key,
		Body:        body,
		ContentType: &contentType,
	})

	return err
}

// GetPresignedInsertURL implements [S3Helper].
func (h *s3Helper) GetPresignedInsertURL(ctx context.Context, key string, contentType string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(h.client, func(po *s3.PresignOptions) {
		if h.publicEndpoint == "" {
			return
		}

		po.ClientOptions = append(po.ClientOptions, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(h.publicEndpoint)
		})
	})

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

// Delete implements [S3Helper].
func (h *s3Helper) Delete(ctx context.Context, key string) error {
	_, err := h.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &h.bucket,
		Key:    &key,
	})

	return err
}
