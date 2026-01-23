package memcache

import "time"

// ClientInterface defines the interface for Memcache operations
// This allows for mock implementations in tests
type ClientInterface interface {
	// Get retrieves a value from Memcache
	Get(key string) (uint64, error)
	// Set sets a value in Memcache with expiration
	Set(key string, value uint64, expiration time.Duration) error
	// IncrementWithExpiration atomically increments a counter and sets expiration if key doesn't exist
	IncrementWithExpiration(key string, delta uint64, expiration time.Duration) (uint64, error)
	// Delete removes a key from Memcache
	Delete(key string) error
	// HealthCheck checks if Memcache is accessible
	HealthCheck() error
	// Close closes the Memcache client
	Close() error
}

// Ensure Client implements ClientInterface
var _ ClientInterface = (*Client)(nil)
var _ ClientInterface = (*MockClient)(nil)
