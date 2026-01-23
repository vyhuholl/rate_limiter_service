package memcache

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// Client wraps the gomemcache client with additional functionality
type Client struct {
	client *memcache.Client
}

// NewClient creates a new Memcache client wrapper
func NewClient(servers []string, timeout time.Duration, maxIdleConns int) *Client {
	client := memcache.New(servers...)
	client.Timeout = timeout
	client.MaxIdleConns = maxIdleConns
	return &Client{
		client: client,
	}
}

// Get retrieves a value from Memcache
func (c *Client) Get(key string) (uint64, error) {
	item, err := c.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return 0, nil // Not found, return 0
		}
		return 0, fmt.Errorf("failed to get key %q: %w", key, err)
	}

	value, err := strconv.ParseUint(string(item.Value), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse value for key %q: %w", key, err)
	}

	return value, nil
}

// Set sets a value in Memcache with expiration
func (c *Client) Set(key string, value uint64, expiration time.Duration) error {
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(strconv.FormatUint(value, 10)),
		Expiration: int32(expiration.Seconds()),
	}
	return c.client.Set(item)
}

// IncrementWithExpiration atomically increments a counter and sets expiration if key doesn't exist
// Returns the new value after increment
func (c *Client) IncrementWithExpiration(key string, delta uint64, expiration time.Duration) (uint64, error) {
	// Try to increment first
	newValue, err := c.client.Increment(key, delta)
	if err == nil {
		return newValue, nil
	}

	// If key doesn't exist, initialize it
	if err == memcache.ErrCacheMiss {
		// Set initial value with expiration
		if err := c.Set(key, delta, expiration); err != nil {
			return 0, fmt.Errorf("failed to initialize key %q: %w", key, err)
		}
		return delta, nil
	}

	return 0, fmt.Errorf("failed to increment key %q: %w", key, err)
}

// Delete removes a key from Memcache
func (c *Client) Delete(key string) error {
	return c.client.Delete(key)
}

// HealthCheck checks if Memcache is accessible
func (c *Client) HealthCheck() error {
	// Try to set and get a test key
	testKey := "health_check"
	if err := c.Set(testKey, 1, time.Second); err != nil {
		return fmt.Errorf("memcache health check failed: %w", err)
	}

	if _, err := c.Get(testKey); err != nil {
		return fmt.Errorf("memcache health check failed: %w", err)
	}

	// Clean up
	_ = c.Delete(testKey)

	return nil
}

// Close closes the Memcache client connections
func (c *Client) Close() error {
	// gomemcache doesn't have an explicit Close method
	// The connections will be garbage collected
	return nil
}
