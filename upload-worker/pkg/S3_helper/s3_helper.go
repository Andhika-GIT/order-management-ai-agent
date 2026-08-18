package S3_helper

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Helper interface {
	Upload(ctx context.Context, key string, body io.Reader, contentType string) error
	Delete(ctx context.Context, key string) error
	Download(ctx context.Context, key string) ([]byte, error)
	GetPresignedReadURL(ctx context.Context, key string, expiry time.Duration) (string, error)                // for reading s3 file
	GetPresignedInsertURL(ctx context.Context, key, contentType string, expiry time.Duration) (string, error) // generate URL to post / insert image
}

type s3Helper struct {
	client *s3.Client
	bucket string
}

func NewS3Helper(client *s3.Client, bucket string) S3Helper {
	return &s3Helper{client: client, bucket: bucket}
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

// Delete implements [S3Helper].
func (h *s3Helper) Delete(ctx context.Context, key string) error {
	_, err := h.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &h.bucket,
		Key:    &key,
	})

	return err
}

func (h *s3Helper) Download(ctx context.Context, key string) ([]byte, error) {
	out, err := h.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &h.bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(out.Body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GetPresignedReadURL implements [S3Helper].
func (h *s3Helper) GetPresignedReadURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	presignedClient := s3.NewPresignClient(h.client)
	req, err := presignedClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &h.bucket,
		Key:    &key,
	}, s3.WithPresignExpires(expiry))

	if err != nil {
		return "", err
	}

	return req.URL, nil
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
