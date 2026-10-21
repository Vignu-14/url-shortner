package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"url-shortner/internal/handlers"
	"url-shortner/internal/services"
	"url-shortner/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main(){

	pool, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	Pool := storage.NewDBStorage(pool)

	mux := http.NewServeMux()

	// mem := storage.NewMemoryStore()
	ser := services.NewURLService(Pool)
	hand := handlers.NewHandler(*ser)

	mux.HandleFunc("POST /", hand.GenerateUrl)
	mux.HandleFunc("GET /{shortUrl}", hand.GetUrl)
	

	fmt.Println("server started on port 3000")
	http.ListenAndServe(":3000", mux)

}

func ConnectDB () (*pgxpool.Pool, error) {
	connStr := "postgres://postgres:admin@localhost:5432/urlshortner?sslmode=disable"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	Pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return Pool, nil
}
