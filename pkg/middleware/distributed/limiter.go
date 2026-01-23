package distributed

import (
	"time"

	"rate_limiter_service/internal/config"
	"rate_limiter_service/pkg/memcache"
)

// Limiter defines the interface for rate limiting
// Both in-memory and distributed limiters implement this interface
type Limiter interface {
	// Allow checks if a request is allowed
	Allow(userID string) bool
	// GetRemainingTokens returns the number of remaining tokens
	GetRemainingTokens(userID string) int
}

// DistributedLimiter extends Limiter with Memcache-specific functionality
type DistributedLimiter interface {
	Limiter

	// GetMemcacheClient returns the underlying Memcache client
	GetMemcacheClient() memcache.ClientInterface
	// GetConfig returns the configuration
	GetConfig() config.Config
	// GetScope returns the scope of the limiter (global, http, grpc, endpoint)
	GetScope() string
	// GetWindowDuration returns the duration of the rate limit window
	GetWindowDuration() time.Duration
}

// BaseLimiter provides common functionality for all distributed limiters
type BaseLimiter struct {
	client memcache.ClientInterface
	config config.Config
	scope  string
	rate   int
}

// NewBaseLimiter creates a new base limiter
func NewBaseLimiter(client memcache.ClientInterface, cfg config.Config, scope string, rate int) *BaseLimiter {
	return &BaseLimiter{
		client: client,
		config: cfg,
		scope:  scope,
		rate:   rate,
	}
}

// GetMemcacheClient returns the underlying Memcache client
func (bl *BaseLimiter) GetMemcacheClient() memcache.ClientInterface {
	return bl.client
}

// GetConfig returns the configuration
func (bl *BaseLimiter) GetConfig() config.Config {
	return bl.config
}

// GetScope returns the scope of the limiter
func (bl *BaseLimiter) GetScope() string {
	return bl.scope
}

// GetWindowDuration returns the duration of the rate limit window
// For distributed limiters, we use a 1-second sliding window
func (bl *BaseLimiter) GetWindowDuration() time.Duration {
	return time.Second
}

// checkRateLimit checks if the count is within the rate limit
func (bl *BaseLimiter) checkRateLimit(count uint64) bool {
	return count <= uint64(bl.rate)
}

// getExpiration returns the expiration time for Memcache keys
// We add a small buffer to handle edge cases at window boundaries
func (bl *BaseLimiter) getExpiration() time.Duration {
	// Use 2 seconds to allow for some buffer at window boundaries
	return 2 * time.Second
}
