package storage

import "context"

type Backend interface {
	Save(ctx context.Context, path string, data []byte) error
}
