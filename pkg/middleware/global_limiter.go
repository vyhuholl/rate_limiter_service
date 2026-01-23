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

// HTTPLimiter enforces HTTP-only rate limits per user
type HTTPLimiter struct {
	// config holds the rate limiting configuration
	config config.Config

	// buckets stores token buckets keyed by userID
	buckets sync.Map // map[string]*TokenBucket
}

// NewHTTPLimiter creates a new HTTP-only rate limiter
func NewHTTPLimiter(cfg config.Config) *HTTPLimiter {
	return &HTTPLimiter{
		config: cfg,
	}
}

// Allow checks if the HTTP request for the given user is allowed
// Returns true if allowed, false if rate limited
func (hl *HTTPLimiter) Allow(userID string) bool {
	bucket := hl.getOrCreateBucket(userID)
	return bucket.Allow()
}

// getOrCreateBucket retrieves or creates a token bucket for the given user
func (hl *HTTPLimiter) getOrCreateBucket(userID string) *TokenBucket {
	// Try to load existing bucket
	if bucket, ok := hl.buckets.Load(userID); ok {
		return bucket.(*TokenBucket)
	}

	// Create new bucket
	bucket := NewTokenBucket(hl.config.HTTPBurstSize, hl.config.HTTPRate)

	// Store it (may have been created by another goroutine in the meantime)
	actual, loaded := hl.buckets.LoadOrStore(userID, bucket)
	if loaded {
		// Another goroutine created it, return that one instead
		return actual.(*TokenBucket)
	}

	return bucket
}

// GetRemainingTokens returns the number of remaining tokens for a user for HTTP requests
func (hl *HTTPLimiter) GetRemainingTokens(userID string) int {
	if bucket, ok := hl.buckets.Load(userID); ok {
		return bucket.(*TokenBucket).GetTokens()
	}

	// If no bucket exists yet, return the full capacity
	return hl.config.HTTPBurstSize
}

// Reset clears all rate limiting state for testing purposes
func (hl *HTTPLimiter) Reset() {
	hl.buckets = sync.Map{}
}

// GRPCLimiter enforces gRPC-only rate limits per user
type GRPCLimiter struct {
	// config holds the rate limiting configuration
	config config.Config

	// buckets stores token buckets keyed by userID
	buckets sync.Map // map[string]*TokenBucket
}

// NewGRPCLimiter creates a new gRPC-only rate limiter
func NewGRPCLimiter(cfg config.Config) *GRPCLimiter {
	return &GRPCLimiter{
		config: cfg,
	}
}

// Allow checks if the gRPC request for the given user is allowed
// Returns true if allowed, false if rate limited
func (gl *GRPCLimiter) Allow(userID string) bool {
	bucket := gl.getOrCreateBucket(userID)
	return bucket.Allow()
}

// getOrCreateBucket retrieves or creates a token bucket for the given user
func (gl *GRPCLimiter) getOrCreateBucket(userID string) *TokenBucket {
	// Try to load existing bucket
	if bucket, ok := gl.buckets.Load(userID); ok {
		return bucket.(*TokenBucket)
	}

	// Create new bucket
	bucket := NewTokenBucket(gl.config.GRPCBurstSize, gl.config.GRPCRate)

	// Store it (may have been created by another goroutine in the meantime)
	actual, loaded := gl.buckets.LoadOrStore(userID, bucket)
	if loaded {
		// Another goroutine created it, return that one instead
		return actual.(*TokenBucket)
	}

	return bucket
}

// GetRemainingTokens returns the number of remaining tokens for a user for gRPC requests
func (gl *GRPCLimiter) GetRemainingTokens(userID string) int {
	if bucket, ok := gl.buckets.Load(userID); ok {
		return bucket.(*TokenBucket).GetTokens()
	}

	// If no bucket exists yet, return the full capacity
	return gl.config.GRPCBurstSize
}

// Reset clears all rate limiting state for testing purposes
func (gl *GRPCLimiter) Reset() {
	gl.buckets = sync.Map{}
}