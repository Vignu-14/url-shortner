package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"url-shortner/internal/handlers"
	"url-shortner/internal/services"
	"url-shortner/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {

	pool, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	rdb, err := redisClient()
	if err != nil {
		log.Print(err)
	}

	redisDb := storage.NewRedisStore(rdb)
	db := storage.NewDBStorage(pool)

	store := storage.NewCachedStore(db, redisDb)

	ser := services.NewURLService(store)
	// mem := storage.NewMemoryStore()
	hand := handlers.NewHandler(*ser)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /", hand.GenerateUrl)
	mux.HandleFunc("GET /{shortUrl}", hand.GetUrl)

	fmt.Println("server started on port 3000")
	http.ListenAndServe(":3000", mux)

}

func ConnectDB() (*pgxpool.Pool, error) {
	
	connStr := os.Getenv("DBURL")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	Pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	if err := Pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}
	return Pool, nil
}

func redisClient() (*redis.Client, error) {
	addr := os.Getenv("REDISADDR")
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     "",
		PoolSize:     200,
		MinIdleConns: 50,
		DB:           0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	return rdb, nil
}
