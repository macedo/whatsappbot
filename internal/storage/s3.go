package storage

import (
	"bytes"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	client *s3.Client
	Bucket string
}

type S3Attributes struct {
	Bucket           string   `mapstructure:"bucket"`
	ConfigFiles      []string `mapstructure:"config_files"`
	CredentialsFiles []string `mapstructure:"credentials_files"`
}

func NewS3(attrs S3Attributes) *S3 {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithSharedConfigFiles(attrs.ConfigFiles),
		config.WithSharedCredentialsFiles(attrs.CredentialsFiles),
	)
	if err != nil {
		log.Fatal(err)
	}

	return &S3{
		client: s3.NewFromConfig(cfg),
		Bucket: attrs.Bucket,
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
