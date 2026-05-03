package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Storage uploads files to AWS S3.
type S3Storage struct {
	client    *s3.Client
	bucket    string
	publicURL string // e.g. "https://<bucket>.s3.<region>.amazonaws.com"
}

func NewS3(client *s3.Client, bucket, publicURL string) *S3Storage {
	return &S3Storage{client: client, bucket: bucket, publicURL: publicURL}
}

// Upload puts an object at the given key and returns its public URL.
func (s *S3Storage) Upload(ctx context.Context, key, contentType string, body io.Reader) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("s3 upload %q: %w", key, err)
	}
	return s.publicURL + "/" + key, nil
}

// NoopStorage is used when S3 is not configured.
type NoopStorage struct{}

func (n *NoopStorage) Upload(_ context.Context, _, _ string, _ io.Reader) (string, error) {
	return "", fmt.Errorf("avatar storage is not configured")
}
