package storage

import (
	"context"
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
	return s.client.Set(ctx, code, originalURL, time.Second * 60 ).Err()
}
