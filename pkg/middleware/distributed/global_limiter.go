package distributed

import (
	"rate_limiter_service/internal/config"
	"rate_limiter_service/pkg/memcache"
)

const (
	scopeGlobal = "global"
)

// GlobalLimiter enforces global rate limits per user across all endpoints using Memcache
type GlobalLimiter struct {
	*CommonLimiter
}

// NewGlobalLimiter creates a new distributed global rate limiter
func NewGlobalLimiter(client memcache.ClientInterface, cfg config.Config) *GlobalLimiter {
	return &GlobalLimiter{
		CommonLimiter: NewCommonLimiter(client, cfg, scopeGlobal, cfg.GlobalRate),
	}
}

// Allow checks if the request for the given user is allowed globally
// Returns true if allowed, false if rate limited
func (gl *GlobalLimiter) Allow(userID string) bool {
	key := gl.config.GetMemcacheKey(gl.scope, userID, "")

	// Increment counter with expiration
	newCount, err := gl.client.IncrementWithExpiration(key, 1, gl.GetExpiration())
	if err != nil {
		// Handle Memcache failure based on failure mode
		gl.CommonLimiter.LogError(userID, err)
		return gl.CommonLimiter.HandleFailure()
	}

	// Check if within rate limit
	return gl.CommonLimiter.CheckRateLimit(newCount)
}

// GetRemainingTokens returns the number of remaining tokens for a user globally
func (gl *GlobalLimiter) GetRemainingTokens(userID string) int {
	key := gl.config.GetMemcacheKey(gl.scope, userID, "")

	count, err := gl.client.Get(key)
	if err != nil {
		// Handle Memcache failure
		gl.CommonLimiter.LogError(userID, err)
		// On failure, return full capacity (conservative approach)
		return gl.CommonLimiter.GetRate()
	}

	remaining := gl.CommonLimiter.GetRate() - int(count)
	if remaining < 0 {
		remaining = 0
	}
	return remaining
}

// Reset clears all rate limiting state for testing purposes
// For distributed limiters, this is a no-op since state is in Memcache
func (gl *GlobalLimiter) Reset() {
	// No-op: distributed state is managed by Memcache
	// Tests should use mock Memcache client for state management
}
