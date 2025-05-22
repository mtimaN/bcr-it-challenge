package main

import (
	"sync"
	"time"
)

// Rate limiter
type RateLimiter struct {
	clients map[string]*clientInfo
	mutex   sync.RWMutex
	limit   int
	window  time.Duration
}

type clientInfo struct {
	requests []time.Time
	lastSeen time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*clientInfo),
		limit:   limit,
		window:  window,
	}

	// Cleanup goroutine
	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *RateLimiter) Allow(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	client, exists := rl.clients[clientIP]

	if !exists {
		client = &clientInfo{
			requests: make([]time.Time, 0),
			lastSeen: now,
		}
		rl.clients[clientIP] = client
	}

	client.lastSeen = now

	// Remove old requests outside the window
	cutoff := now.Add(-rl.window)
	validRequests := make([]time.Time, 0)
	for _, reqTime := range client.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	client.requests = validRequests

	// Check if limit exceeded
	if len(client.requests) >= rl.limit {
		return false
	}

	// Add current request
	client.requests = append(client.requests, now)
	return true
}

func (rl *RateLimiter) cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	cutoff := time.Now().Add(-rl.window * 2)
	for ip, client := range rl.clients {
		if client.lastSeen.Before(cutoff) {
			delete(rl.clients, ip)
		}
	}
}
