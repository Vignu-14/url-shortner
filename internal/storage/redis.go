package storage

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func (s *RedisStore) Save(ctx context.Context, code string, originalURL string) error {
	return s.client.Set(ctx, code, originalURL, time.Hour * 12).Err()
}

func (s *RedisStore) Get(ctx context.Context, code string) (string, error) {
	originalUrl, err := s.client.Get(ctx, code).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrNotFound
		}
		return "", err
	}
	return originalUrl, nil
}
