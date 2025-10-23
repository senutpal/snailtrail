package main

import (
	"encoding/json"
	"net"
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

func SimpleLimiter(limit int,window time.Duration) *SimpleRateLimiter {
	return &SimpleRateLimiter{
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

func perClientRateLimiter(next http.Handler) http.Handler {
	type client struct {
		limiter *SimpleRateLimiter
		lastSeen time.Time
	}

	clients := make(map[string]*client)
	var mu sync.Mutex

	go func () {
		for{
			time.Sleep(time.Minute)
			mu.Lock()
			for ip,c := range clients {
				if time.Since(c.lastSeen) > 2*time.Minute {
					delete(clients,ip)
				}
			}
			mu.Unlock()
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	ip,_,err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w,"Unable to parse IP address",http.StatusInternalServerError)
		return 
	}
	mu.Lock()
	c,exists := clients[ip]
if !exists {
	c= &client{
		limiter: SimpleLimiter(2,time.Second),
		lastSeen: time.Now(),
	}
	clients[ip] = c
}
c.lastSeen = time.Now()
mu.Unlock()

if !c.limiter.Allow() {
	message := Message {
		Status: "Request Failed",
		Body : "Rate limit exceeded for your IP",
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	json.NewEncoder(w).Encode(&message)
	return
}
next.ServeHTTP(w,r)
})
}