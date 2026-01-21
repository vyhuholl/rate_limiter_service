package middleware

import (
	"sync"
	"time"
)

// TokenBucket represents a token bucket rate limiter
type TokenBucket struct {
	// mu protects concurrent access to the bucket
	mu sync.Mutex

	// capacity is the maximum number of tokens the bucket can hold
	capacity int

	// tokens is the current number of tokens in the bucket
	tokens int

	// refillRate is the number of tokens added per second
	refillRate int

	// lastRefill is the timestamp of the last token refill
	lastRefill time.Time

	// refillInterval is the time between adding one token
	refillInterval time.Duration
}

// NewTokenBucket creates a new token bucket with the specified capacity and refill rate
func NewTokenBucket(capacity int, refillRate int) *TokenBucket {
	if capacity <= 0 {
		capacity = 1
	}
	if refillRate <= 0 {
		refillRate = 1
	}

	return &TokenBucket{
		capacity:       capacity,
		tokens:         capacity, // Start with full bucket
		refillRate:     refillRate,
		lastRefill:     time.Now(),
		refillInterval: time.Second / time.Duration(refillRate),
	}
}

// Allow attempts to consume one token from the bucket.
// Returns true if the token was consumed, false if the bucket was empty.
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// refill adds tokens to the bucket based on elapsed time since last refill
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	// Calculate how many tokens should be added
	tokensToAdd := int(elapsed / tb.refillInterval)

	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		// Update lastRefill to the time when the last token was added
		tb.lastRefill = tb.lastRefill.Add(time.Duration(tokensToAdd) * tb.refillInterval)
	}
}

// GetTokens returns the current number of tokens in the bucket
func (tb *TokenBucket) GetTokens() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	return tb.tokens
}

// GetCapacity returns the maximum capacity of the bucket
func (tb *TokenBucket) GetCapacity() int {
	return tb.capacity
}

// Reset fills the bucket to capacity and resets the last refill time
func (tb *TokenBucket) Reset() {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.tokens = tb.capacity
	tb.lastRefill = time.Now()
}