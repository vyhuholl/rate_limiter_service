package distributed

import (
	"fmt"
	"log"

	"rate_limiter_service/internal/config"
	"rate_limiter_service/pkg/memcache"
)

const (
	scopeEndpoint = "endpoint"
)

// PerEndpointLimiter enforces per-endpoint rate limits per user using Memcache
type PerEndpointLimiter struct {
	*BaseLimiter
}

// NewPerEndpointLimiter creates a new distributed per-endpoint rate limiter
func NewPerEndpointLimiter(client memcache.ClientInterface, cfg config.Config) *PerEndpointLimiter {
	return &PerEndpointLimiter{
		BaseLimiter: NewBaseLimiter(client, cfg, scopeEndpoint, cfg.HTTPDefaultMethodRate),
	}
}

// Allow checks if the request for the given user and endpoint is allowed
// Returns true if allowed, false if rate limited
func (pel *PerEndpointLimiter) Allow(userID, method, path string) bool {
	endpointKey := fmt.Sprintf("%s:%s", method, path)
	key := pel.config.GetMemcacheKey(pel.scope, userID, endpointKey)

	// Get rate for this specific endpoint
	rate := pel.getRateForEndpoint(endpointKey)

	// Increment counter with expiration
	newCount, err := pel.client.IncrementWithExpiration(key, 1, pel.getExpiration())
	if err != nil {
		// Handle Memcache failure based on failure mode
		log.Printf("memcache error incrementing per-endpoint counter for user %s, endpoint %s: %v", userID, endpointKey, err)
		return pel.handleFailure()
	}

	// Check if within rate limit
	return newCount <= uint64(rate)
}

// GetRemainingTokens returns the number of remaining tokens for a user-endpoint combination
func (pel *PerEndpointLimiter) GetRemainingTokens(userID, method, path string) int {
	endpointKey := fmt.Sprintf("%s:%s", method, path)
	key := pel.config.GetMemcacheKey(pel.scope, userID, endpointKey)

	rate := pel.getRateForEndpoint(endpointKey)

	count, err := pel.client.Get(key)
	if err != nil {
		// Handle Memcache failure
		log.Printf("memcache error getting per-endpoint counter for user %s, endpoint %s: %v", userID, endpointKey, err)
		// On failure, return full capacity (conservative approach)
		return rate
	}

	remaining := rate - int(count)
	if remaining < 0 {
		remaining = 0
	}
	return remaining
}

// getRateForEndpoint returns the rate limit for a specific endpoint
func (pel *PerEndpointLimiter) getRateForEndpoint(endpointKey string) int {
	if rate, ok := pel.config.HTTPMethods[endpointKey]; ok {
		return rate
	}
	return pel.config.HTTPDefaultMethodRate
}

// handleFailure handles Memcache failures based on configured failure mode
func (pel *PerEndpointLimiter) handleFailure() bool {
	switch pel.config.MemcacheFailureMode {
	case config.FailureModeAllow:
		// Fail-open: allow requests when Memcache is unavailable
		return true
	case config.FailureModeDeny:
		// Fail-closed: deny requests when Memcache is unavailable
		return false
	default:
		// Default to allow for safety
		return true
	}
}

// Reset clears all rate limiting state for testing purposes
// For distributed limiters, this is a no-op since state is in Memcache
func (pel *PerEndpointLimiter) Reset() {
	// No-op: distributed state is managed by Memcache
	// Tests should use mock Memcache client for state management
}
