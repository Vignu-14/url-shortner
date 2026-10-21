package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next(w, r)

		fmt.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Since(start))
	}
}
