package main

import (
	"fmt"
	"net/http"
	"url-shortner/internal/handlers"
	"url-shortner/internal/services"
	"url-shortner/internal/storage"
)

func main(){

	mux := http.NewServeMux()

	mem := storage.NewMemoryStore()
	ser := services.NewURLService(mem)
	hand := handlers.NewHandler(*ser)

	mux.HandleFunc("POST /", hand.GenerateUrl)
	mux.HandleFunc("GET /{shortUrl}", hand.GetUrl)
	

	fmt.Println("server started on port 3000")
	http.ListenAndServe(":3000", mux)

}

http://localhost:3000/fUMQplJwuNh

http://localhost:3000/qHQ1tq2Nvi