package distributed

import (
	"log"
	"time"

	"rate_limiter_service/internal/config"
	"rate_limiter_service/pkg/memcache"
)

// CommonLimiter provides common functionality for all distributed limiters
type CommonLimiter struct {
	client memcache.ClientInterface
	config config.Config
	scope  string
	rate   int
}

// NewCommonLimiter creates a new common limiter
func NewCommonLimiter(client memcache.ClientInterface, cfg config.Config, scope string, rate int) *CommonLimiter {
	return &CommonLimiter{
		client: client,
		config: cfg,
		scope:  scope,
		rate:   rate,
	}
}

// GetClient returns the underlying Memcache client
func (cl *CommonLimiter) GetClient() memcache.ClientInterface {
	return cl.client
}

// GetConfig returns the configuration
func (cl *CommonLimiter) GetConfig() config.Config {
	return cl.config
}

// GetScope returns the scope of the limiter
func (cl *CommonLimiter) GetScope() string {
	return cl.scope
}

// GetRate returns the rate limit
func (cl *CommonLimiter) GetRate() int {
	return cl.rate
}

// GetExpiration returns the expiration time for Memcache keys
// We add a small buffer to handle edge cases at window boundaries
func (cl *CommonLimiter) GetExpiration() time.Duration {
	// Use 2 seconds to allow for some buffer at window boundaries
	return 2 * time.Second
}

// CheckRateLimit checks if count is within the rate limit
func (cl *CommonLimiter) CheckRateLimit(count uint64) bool {
	return count <= uint64(cl.rate)
}

// HandleFailure handles Memcache failures based on configured failure mode
func (cl *CommonLimiter) HandleFailure() bool {
	switch cl.config.MemcacheFailureMode {
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

// LogError logs a Memcache error with context
func (cl *CommonLimiter) LogError(userID string, err error) {
	log.Printf("memcache error incrementing %s counter for user %s: %v", cl.scope, userID, err)
}
