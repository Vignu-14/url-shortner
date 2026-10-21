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
	return shortUrl, nil
}

func (s URLService) Resolve(ctx context.Context, shortUrl string) (string, error) {
	originalUrl, err := s.store.Get(ctx, shortUrl)
	if err != nil {
		return "", fmt.Errorf("cannot get original request or does not exist: %w", err)
	}
	
	return originalUrl, nil
}

