package utils

import (
	"sync"
	"time"
)

// RateLimiter struct to hold rate limit configuration
type RateLimiter struct {
	maxRequests int
	window      time.Duration
	mu          sync.Mutex
	requests    []time.Time
}

// NewRateLimiter creates a new RateLimiter instance
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests: maxRequests,
		window:      window,
		requests:    make([]time.Time, 0, maxRequests),
	}
}

// Allow checks if a request is allowed based on the rate limit configuration
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	// Remove outdated requests
	cutoff := now.Add(-rl.window)
	for len(rl.requests) > 0 && rl.requests[0].Before(cutoff) {
		rl.requests = rl.requests[1:]
	}

	if len(rl.requests) < rl.maxRequests {
		// Allow the request
		rl.requests = append(rl.requests, now)
		return true
	}

	// Deny the request
	return false
}
