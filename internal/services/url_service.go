package services

import (
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

func Shorten(longUrl string) string {
	shortUrl := shortUrl(longUrl)

	return shortUrl
}

func Resolve(shortUrl string) string {
	longUrl := "testing"
	return longUrl
}
