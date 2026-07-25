package utilities

import (
	"sync"
	"time"
)

// rate limiting mechanism
type TokenBucket struct {
	capacity   int
	tokens     int
	rate       time.Duration
	lastFilled time.Time
	mu         sync.Mutex
}

// Will create a new token bucket.
func NewTokenBucket(capacity int, rate time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		rate:       rate,
		lastFilled: time.Now(),
	}
}

// will check if a request is allowed based on remaining tokens
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastFilled)
	addTokens := int(elapsed / tb.rate)

	tb.tokens += addTokens

	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	if addTokens > 0 {
		tb.lastFilled = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}
