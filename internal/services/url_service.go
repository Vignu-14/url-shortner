package services

import (
	"context"
	"fmt"
	"url-shortner/internal/storage"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type URLService struct {
	store storage.Store
}

func NewURLService(store storage.Store) *URLService {
	return &URLService{
		store: store,
	}
}

func (s URLService) Shorten(ctx context.Context, longUrl string) (string, error) {
	shortUrl := shortUrl(longUrl)
	err := s.store.Save(ctx, shortUrl, longUrl)
	if err != nil {
		return "", fmt.Errorf("can not generate short url: %w", err)
	}
	s.store.Get(ctx, "ksdf")
	return shortUrl, nil
}

func (s URLService) Resolve(ctx context.Context, shortUrl string) (string, error) {
	originalUrl, err := s.store.Get(ctx, shortUrl)
	if err != nil {
		return "", fmt.Errorf("cannot get original request or does not exist: %w", err)
	}
	
	return originalUrl, nil
}

func (s URLService) Test(
	ctx context.Context,
	code string,
	originalURL string,
	dbStore *storage.DBStorage,
	redisStore *storage.RedisStore,
) error {
	if err := dbStore.Save(ctx, code, originalURL); err != nil {
		return fmt.Errorf("save to database: %w", err)
	}

	dbURL, err := dbStore.Get(ctx, code)
	if err != nil {
		return fmt.Errorf("get from database: %w", err)
	}
	fmt.Printf("database URL: %s\n", dbURL)

	if err := redisStore.Save(ctx, code, originalURL); err != nil {
		return fmt.Errorf("save to Redis: %w", err)
	}

	redisURL, err := redisStore.Get(ctx, code)
	if err != nil {
		return fmt.Errorf("get from Redis: %w", err)
	}
	fmt.Printf("Redis URL: %s\n", redisURL)

	return nil
}
