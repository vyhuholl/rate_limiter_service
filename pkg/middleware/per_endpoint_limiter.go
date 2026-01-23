package middleware

import (
	"fmt"
	"sync"

	"rate_limiter_service/internal/config"
)

// PerEndpointLimiter enforces per-endpoint rate limits per user
type PerEndpointLimiter struct {
	// config holds the rate limiting configuration
	config config.Config

	// buckets stores token buckets keyed by "userID:endpointKey"
	// endpointKey is "method:path" (e.g., "GET:/api/users")
	buckets sync.Map // map[string]*TokenBucket
}

// NewPerEndpointLimiter creates a new per-endpoint rate limiter
func NewPerEndpointLimiter(cfg config.Config) *PerEndpointLimiter {
	return &PerEndpointLimiter{
		config: cfg,
	}
}

// Allow checks if the request for the given user and endpoint is allowed
// Returns true if allowed, false if rate limited
func (pel *PerEndpointLimiter) Allow(userID, method, path string) bool {
	endpointKey := fmt.Sprintf("%s:%s", method, path)
	bucketKey := fmt.Sprintf("%s:%s", userID, endpointKey)

	// Get or create bucket for this user-endpoint combination
	bucket := pel.getOrCreateBucket(bucketKey, endpointKey)

	return bucket.Allow()
}

// getOrCreateBucket retrieves or creates a token bucket for the given key
func (pel *PerEndpointLimiter) getOrCreateBucket(bucketKey, endpointKey string) *TokenBucket {
	// Try to load existing bucket
	if bucket, ok := pel.buckets.Load(bucketKey); ok {
		return bucket.(*TokenBucket)
	}

	// Determine rate for this endpoint
	rate := pel.getRateForEndpoint(endpointKey)
	burstSize := pel.config.PerEndpointBurstSize

	// Create new bucket
	bucket := NewTokenBucket(burstSize, rate)

	// Store it (may have been created by another goroutine in the meantime)
	actual, loaded := pel.buckets.LoadOrStore(bucketKey, bucket)
	if loaded {
		// Another goroutine created it, return that one instead
		return actual.(*TokenBucket)
	}

	return bucket
}

// getRateForEndpoint returns the rate limit for a specific endpoint, falling back to default
func (pel *PerEndpointLimiter) getRateForEndpoint(endpointKey string) int {
	if rate, ok := pel.config.HTTPMethods[endpointKey]; ok {
		return rate
	}
	return pel.config.HTTPDefaultMethodRate
}

// GetRemainingTokens returns the number of remaining tokens for a user-endpoint combination
func (pel *PerEndpointLimiter) GetRemainingTokens(userID, method, path string) int {
	endpointKey := fmt.Sprintf("%s:%s", method, path)
	bucketKey := fmt.Sprintf("%s:%s", userID, endpointKey)

	if bucket, ok := pel.buckets.Load(bucketKey); ok {
		return bucket.(*TokenBucket).GetTokens()
	}

	// If no bucket exists yet, return the full capacity
	return pel.config.PerEndpointBurstSize
}

// Reset clears all rate limiting state for testing purposes
func (pel *PerEndpointLimiter) Reset() {
	pel.buckets = sync.Map{}
}