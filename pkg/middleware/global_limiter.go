package middleware

import (
	"sync"

	"rate_limiter_service/internal/config"
)

// GlobalLimiter enforces global rate limits per user across all endpoints
type GlobalLimiter struct {
	// config holds the rate limiting configuration
	config config.Config

	// buckets stores token buckets keyed by userID
	buckets sync.Map // map[string]*TokenBucket
}

// NewGlobalLimiter creates a new global rate limiter
func NewGlobalLimiter(cfg config.Config) *GlobalLimiter {
	return &GlobalLimiter{
		config: cfg,
	}
}

// Allow checks if the request for the given user is allowed globally
// Returns true if allowed, false if rate limited
func (gl *GlobalLimiter) Allow(userID string) bool {
	bucket := gl.getOrCreateBucket(userID)
	return bucket.Allow()
}

// getOrCreateBucket retrieves or creates a token bucket for the given user
func (gl *GlobalLimiter) getOrCreateBucket(userID string) *TokenBucket {
	// Try to load existing bucket
	if bucket, ok := gl.buckets.Load(userID); ok {
		return bucket.(*TokenBucket)
	}

	// Create new bucket
	bucket := NewTokenBucket(gl.config.GlobalBurstSize, gl.config.GlobalRate)

	// Store it (may have been created by another goroutine in the meantime)
	actual, loaded := gl.buckets.LoadOrStore(userID, bucket)
	if loaded {
		// Another goroutine created it, return that one instead
		return actual.(*TokenBucket)
	}

	return bucket
}

// GetRemainingTokens returns the number of remaining tokens for a user globally
func (gl *GlobalLimiter) GetRemainingTokens(userID string) int {
	if bucket, ok := gl.buckets.Load(userID); ok {
		return bucket.(*TokenBucket).GetTokens()
	}

	// If no bucket exists yet, return the full capacity
	return gl.config.GlobalBurstSize
}

// Reset clears all rate limiting state for testing purposes
func (gl *GlobalLimiter) Reset() {
	gl.buckets = sync.Map{}
}