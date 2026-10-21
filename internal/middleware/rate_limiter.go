package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type ratelimiter struct {
	mu       sync.Mutex
	visitors map[string]int
	limit    int
	duration time.Duration
}

func NewRateLimiter(limit int, duration time.Duration) *ratelimiter {
	rl := &ratelimiter{
		visitors: make(map[string]int),
		limit:    limit,
		duration: duration,
	}
	go rl.resetVisitorCount()
	return rl
}

func (rl *ratelimiter) resetVisitorCount() {
	for {
		time.Sleep(rl.duration)
		rl.mu.Lock()
		rl.visitors = make(map[string]int)
		rl.mu.Unlock()
	}
}

func (rl *ratelimiter) RateLimiter(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}
		
		rl.mu.Lock()
		rl.visitors[host]++
		rl.mu.Unlock()

		if rl.visitors[host] > rl.limit {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}
