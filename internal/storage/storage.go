package storage

import (
	"context"
	"log"

	"github.com/mitchellh/mapstructure"
)

type Backend interface {
	Save(ctx context.Context, path string, data []byte) error
}

func New(services map[string]any) Backend {
	var storage Backend

	for k, v := range services {
		switch k {
		case "s3":
			var attrs S3Attributes
			if err := mapstructure.Decode(v, &attrs); err != nil {
				log.Fatal(err)
			}
			storage = NewS3(attrs)
		}
	}

	return storage
}
