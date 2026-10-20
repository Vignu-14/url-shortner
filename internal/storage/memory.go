package storage

import (
	"context"
	"sync"
)

type MemoryStore struct {
	mu sync.RWMutex
	urls map[string]string
}

func NewMemoryStore() *MemoryStore{
	return &MemoryStore{
		urls: make(map[string]string),
	}
}

var _ Store = (*MemoryStore)(nil)

func (m *MemoryStore) Save(ctx context.Context, code string, originalURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.urls[code] = originalURL
	return nil
}

func (m *MemoryStore) Get(ctx context.Context, code string) (string, error) {
	m.mu.RLocker()
	defer m.mu.RUnlock()

	originalURL, ok := m.urls[code]
	if !ok {
		return "", ErrNotFound
	}
	return originalURL, nil
}