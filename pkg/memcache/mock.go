package memcache

import (
	"fmt"
	"sync"
	"time"
)

// MockClient is a mock implementation of the Memcache client for testing
type MockClient struct {
	mu     sync.RWMutex
	data   map[string]mockItem
	closed bool
}

type mockItem struct {
	value     uint64
	expiresAt time.Time
}

// NewMockClient creates a new mock Memcache client
func NewMockClient() *MockClient {
	return &MockClient{
		data: make(map[string]mockItem),
	}
}

// Get retrieves a value from the mock Memcache
func (m *MockClient) Get(key string) (uint64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return 0, fmt.Errorf("client is closed")
	}

	item, exists := m.data[key]
	if !exists {
		return 0, nil // Not found, return 0
	}

	// Check expiration
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		return 0, nil // Expired
	}

	return item.value, nil
}

// Set sets a value in the mock Memcache with expiration
func (m *MockClient) Set(key string, value uint64, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return fmt.Errorf("client is closed")
	}

	var expiresAt time.Time
	if expiration > 0 {
		expiresAt = time.Now().Add(expiration)
	}

	m.data[key] = mockItem{
		value:     value,
		expiresAt: expiresAt,
	}

	return nil
}

// IncrementWithExpiration atomically increments a counter and sets expiration if key doesn't exist
func (m *MockClient) IncrementWithExpiration(key string, delta uint64, expiration time.Duration) (uint64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return 0, fmt.Errorf("client is closed")
	}

	item, exists := m.data[key]
	isExpired := exists && !item.expiresAt.IsZero() && time.Now().After(item.expiresAt)

	if !exists || isExpired {
		// Key doesn't exist or is expired, initialize it
		var expiresAt time.Time
		if expiration > 0 {
			expiresAt = time.Now().Add(expiration)
		}
		m.data[key] = mockItem{
			value:     delta,
			expiresAt: expiresAt,
		}
		return delta, nil
	}

	// Increment existing value
	item.value += delta
	m.data[key] = item
	return item.value, nil
}

// Delete removes a key from the mock Memcache
func (m *MockClient) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return fmt.Errorf("client is closed")
	}

	delete(m.data, key)
	return nil
}

// HealthCheck checks if the mock Memcache is accessible
func (m *MockClient) HealthCheck() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return fmt.Errorf("client is closed")
	}
	return nil
}

// Close closes the mock Memcache client
func (m *MockClient) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.closed = true
	return nil
}

// Clear clears all data in the mock Memcache
func (m *MockClient) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string]mockItem)
}
