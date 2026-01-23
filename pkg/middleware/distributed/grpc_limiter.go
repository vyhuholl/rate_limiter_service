package distributed

import (
	"rate_limiter_service/internal/config"
	"rate_limiter_service/pkg/memcache"
)

const (
	scopeGRPC = "grpc"
)

// GRPCLimiter enforces gRPC-only rate limits per user using Memcache
type GRPCLimiter struct {
	*CommonLimiter
}

// NewGRPCLimiter creates a new distributed gRPC-only rate limiter
func NewGRPCLimiter(client memcache.ClientInterface, cfg config.Config) *GRPCLimiter {
	return &GRPCLimiter{
		CommonLimiter: NewCommonLimiter(client, cfg, scopeGRPC, cfg.GRPCRate),
	}
}

// Allow checks if the gRPC request for the given user is allowed
// Returns true if allowed, false if rate limited
func (gl *GRPCLimiter) Allow(userID string) bool {
	key := gl.config.GetMemcacheKey(gl.scope, userID, "")

	// Increment counter with expiration
	newCount, err := gl.client.IncrementWithExpiration(key, 1, gl.GetExpiration())
	if err != nil {
		// Handle Memcache failure based on failure mode
		gl.LogError(userID, err)
		return gl.HandleFailure()
	}

	// Check if within rate limit
	return gl.CheckRateLimit(newCount)
}

// GetRemainingTokens returns the number of remaining tokens for a user for gRPC requests
func (gl *GRPCLimiter) GetRemainingTokens(userID string) int {
	key := gl.config.GetMemcacheKey(gl.scope, userID, "")

	count, err := gl.client.Get(key)
	if err != nil {
		// Handle Memcache failure
		gl.LogError(userID, err)
		// On failure, return full capacity (conservative approach)
		return gl.GetRate()
	}

	remaining := gl.GetRate() - int(count)
	if remaining < 0 {
		remaining = 0
	}
	return remaining
}

// Reset clears all rate limiting state for testing purposes
// For distributed limiters, this is a no-op since state is in Memcache
func (gl *GRPCLimiter) Reset() {
	// No-op: distributed state is managed by Memcache
	// Tests should use mock Memcache client for state management
}
