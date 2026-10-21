package storage

import (
	"context"
	"errors"
	"log"
)

var ErrNotFound = errors.New("short URL not found")

type Store interface {
	Save(ctx context.Context, code string, originalURL string) error
	Get(ctx context.Context, code string) (string, error)
}

type CachedStore struct {
	db    Store
	cache Store
}

func NewCachedStore(db, cache Store) *CachedStore {
	return &CachedStore{
		db:    db,
		cache: cache,
	}
}

var _ Store = (*CachedStore)(nil)

func (s *CachedStore) Save(ctx context.Context, code, originalURL string) error {
	if err := s.db.Save(ctx, code, originalURL); err != nil {
		return err
	}

	if err := s.cache.Save(ctx, code, originalURL); err != nil {
		log.Printf("Redis cache write failed: %v", err)
	}

	return nil
}

func (s *CachedStore) Get(ctx context.Context, code string) (string, error) {
	originalURL, err := s.cache.Get(ctx, code)
	if err == nil {
		return originalURL, nil
	}
	if !errors.Is(err, ErrNotFound) {
		log.Printf("Redis cache read failed: %v", err)
	}

	originalURL, err = s.db.Get(ctx, code)
	if err != nil {
		return "", err
	}

	if err := s.cache.Save(ctx, code, originalURL); err != nil {
		log.Printf("Redis cache fill failed: %v", err)
	}

	return originalURL, nil
}
