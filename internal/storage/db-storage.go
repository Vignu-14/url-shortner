package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)


type DBStorage struct {
	Pool *pgxpool.Pool
}

func NewDBStorage (pool *pgxpool.Pool) *DBStorage {
	return &DBStorage{Pool: pool}
}

func (db *DBStorage) Save(ctx context.Context, shortUrl, originalURL string) error {
	query := `INSERT INTO short_urls(short_url, long_url)
			  VALUES($1, $2)
			  ON CONFLICT (short_url) DO UPDATE SET long_url = EXCLUDED.long_url`

	cmgTag, err := db.Pool.Exec(ctx, query, shortUrl, originalURL)
	if err != nil {
		return fmt.Errorf("cannot store the url in database: %w", err)
	}
	if cmgTag.RowsAffected() == 0 {
		return fmt.Errorf("url storage failed in database: no rows affected")
	}
	return nil
}

func (db *DBStorage) Get(ctx context.Context, shortUrl string) (string, error) {
	query := `SELECT long_url FROM short_urls WHERE short_url = $1`
	var originalURL string
	err := db.Pool.QueryRow(ctx, query, shortUrl).Scan(&originalURL)
	if err != nil {
		return "", fmt.Errorf("cannot retreive original url or does not exist: %w", err)
	}
	return originalURL, nil
}

// "postgres://postgres:admin@localhost:5432/school?sslmode=disable"