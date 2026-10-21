package storage

import ( "github.com/redis/go-redis/v9" )

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore (client *redis.Client) *RedisStore{
	return &RedisStore{client: client}
}

