package distributed

import (
	"rate_limiter_service/internal/config"
	"rate_limiter_service/pkg/memcache"
)

const (
	scopeHTTP = "http"
)

// HTTPLimiter enforces HTTP-only rate limits per user using Memcache
type HTTPLimiter struct {
	*CommonLimiter
}

// NewHTTPLimiter creates a new distributed HTTP-only rate limiter
func NewHTTPLimiter(client memcache.ClientInterface, cfg config.Config) *HTTPLimiter {
	return &HTTPLimiter{
		CommonLimiter: NewCommonLimiter(client, cfg, scopeHTTP, cfg.HTTPRate),
	}
}

// Allow checks if the HTTP request for the given user is allowed
// Returns true if allowed, false if rate limited
func (hl *HTTPLimiter) Allow(userID string) bool {
	key := hl.config.GetMemcacheKey(hl.scope, userID, "")

	// Increment counter with expiration
	newCount, err := hl.client.IncrementWithExpiration(key, 1, hl.GetExpiration())
	if err != nil {
		// Handle Memcache failure based on failure mode
		hl.LogError(userID, err)
		return hl.HandleFailure()
	}

	// Check if within rate limit
	return hl.CheckRateLimit(newCount)
}

// GetRemainingTokens returns the number of remaining tokens for a user for HTTP requests
func (hl *HTTPLimiter) GetRemainingTokens(userID string) int {
	key := hl.config.GetMemcacheKey(hl.scope, userID, "")

	count, err := hl.client.Get(key)
	if err != nil {
		// Handle Memcache failure
		hl.LogError(userID, err)
		// On failure, return full capacity (conservative approach)
		return hl.GetRate()
	}

	remaining := hl.GetRate() - int(count)
	if remaining < 0 {
		remaining = 0
	}
	return remaining
}

// Reset clears all rate limiting state for testing purposes
// For distributed limiters, this is a no-op since state is in Memcache
func (hl *HTTPLimiter) Reset() {
	// No-op: distributed state is managed by Memcache
	// Tests should use mock Memcache client for state management
}
