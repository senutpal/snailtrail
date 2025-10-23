package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type SimpleRateLimiter struct {
	mu sync.Mutex
	lastChecked time.Time
	count int
	limit int 
	window time.Duration
} 

func SimpleLimiter (limit int, window time.Duration) *SimpleRateLimiter {
	return &SimpleRateLimiter {
		lastChecked: time.Now(),
		count: 0,
		limit: limit,
		window: window,
	}
}

func (l *SimpleRateLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	if now.Sub(l.lastChecked) > l.window {
		l.count = 0
		l.lastChecked = now
	}

	if l.count < l.limit {
		l.count++
		return true
	}

	return false
}

func rateLimiter(next http.Handler) http.Handler {
	limiter := SimpleLimiter(2, time.Second) 

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			message := Message{
				Status: "Request Failed",
				Body:   "The API is at capacity",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(&message)
			return
		}

		next.ServeHTTP(w, r) 
	})
}
