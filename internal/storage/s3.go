package storage

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	client *s3.Client
	Bucket string
}

func NewS3(cli *s3.Client, bucket string) *S3 {
	return &S3{
		client: cli,
		Bucket: bucket,
	}
}

func (s *S3) Save(ctx context.Context, path string, data []byte) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	})

	return err
}
