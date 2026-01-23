package middleware

import (
	"rate_limiter_service/internal/config"
	"rate_limiter_service/pkg/memcache"
	"rate_limiter_service/pkg/middleware/distributed"
)

// LimiterFactory creates limiters based on configuration
// Returns in-memory limiters by default, distributed limiters when Memcache is configured
type LimiterFactory struct {
	config config.Config
}

// NewLimiterFactory creates a new limiter factory
func NewLimiterFactory(cfg config.Config) *LimiterFactory {
	return &LimiterFactory{
		config: cfg,
	}
}

// CreateGlobalLimiter creates a global limiter (in-memory or distributed)
func (lf *LimiterFactory) CreateGlobalLimiter() GlobalLimiterInterface {
	if lf.config.IsDistributedEnabled() {
		client := memcache.NewClient(
			lf.config.MemcacheServers,
			lf.config.MemcacheTimeout,
			lf.config.MemcacheMaxIdleConns,
		)
		return distributed.NewGlobalLimiter(client, lf.config)
	}
	return NewGlobalLimiter(lf.config)
}

// CreatePerEndpointLimiter creates a per-endpoint limiter (in-memory or distributed)
func (lf *LimiterFactory) CreatePerEndpointLimiter() PerEndpointLimiterInterface {
	if lf.config.IsDistributedEnabled() {
		client := memcache.NewClient(
			lf.config.MemcacheServers,
			lf.config.MemcacheTimeout,
			lf.config.MemcacheMaxIdleConns,
		)
		return distributed.NewPerEndpointLimiter(client, lf.config)
	}
	return NewPerEndpointLimiter(lf.config)
}

// CreateHTTPLimiter creates an HTTP-only limiter (in-memory or distributed)
func (lf *LimiterFactory) CreateHTTPLimiter() HTTPLimiterInterface {
	if lf.config.IsDistributedEnabled() {
		client := memcache.NewClient(
			lf.config.MemcacheServers,
			lf.config.MemcacheTimeout,
			lf.config.MemcacheMaxIdleConns,
		)
		return distributed.NewHTTPLimiter(client, lf.config)
	}
	return NewHTTPLimiter(lf.config)
}

// CreateGRPCLimiter creates a gRPC-only limiter (in-memory or distributed)
func (lf *LimiterFactory) CreateGRPCLimiter() GRPCLimiterInterface {
	if lf.config.IsDistributedEnabled() {
		client := memcache.NewClient(
			lf.config.MemcacheServers,
			lf.config.MemcacheTimeout,
			lf.config.MemcacheMaxIdleConns,
		)
		return distributed.NewGRPCLimiter(client, lf.config)
	}
	return NewGRPCLimiter(lf.config)
}

// GlobalLimiterInterface defines the interface for global limiters
type GlobalLimiterInterface interface {
	Allow(userID string) bool
	GetRemainingTokens(userID string) int
	Reset()
}

// PerEndpointLimiterInterface defines the interface for per-endpoint limiters
type PerEndpointLimiterInterface interface {
	Allow(userID, method, path string) bool
	GetRemainingTokens(userID, method, path string) int
	Reset()
}

// HTTPLimiterInterface defines the interface for HTTP-only limiters
type HTTPLimiterInterface interface {
	Allow(userID string) bool
	GetRemainingTokens(userID string) int
	Reset()
}

// GRPCLimiterInterface defines the interface for gRPC-only limiters
type GRPCLimiterInterface interface {
	Allow(userID string) bool
	GetRemainingTokens(userID string) int
	Reset()
}
