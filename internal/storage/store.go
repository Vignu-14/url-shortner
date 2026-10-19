package storage

import (
	"context"
	"errors"
)

var ErrNotFounc = errors.New("short URL not found")

type Store interface {
	Save(ctx context.Context, code string, originalURL string) error
	Get(ctx context.Context, code string) (string, error)
}

