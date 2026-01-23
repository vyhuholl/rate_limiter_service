package memcache

import (
	"testing"
	"time"
)

const (
	testUserID = "user123"
)

func TestNewClient(t *testing.T) {
	client := NewClient([]string{"localhost:11211"}, 100*time.Millisecond, 10)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.client == nil {
		t.Error("client.client not initialized")
	}
}

func TestMockClient_Get(t *testing.T) {
	mock := NewMockClient()

	// Test getting non-existent key
	value, err := mock.Get("nonexistent")
	if err != nil {
		t.Errorf("Get() error = %v, want nil", err)
	}
	if value != 0 {
		t.Errorf("Get() = %d, want 0", value)
	}

	// Test setting and getting a key
	err = mock.Set("test_key", 42, time.Minute)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	value, err = mock.Get("test_key")
	if err != nil {
		t.Errorf("Get() error = %v, want nil", err)
	}
	if value != 42 {
		t.Errorf("Get() = %d, want 42", value)
	}
}

func TestMockClient_Set(t *testing.T) {
	mock := NewMockClient()

	// Test setting a value
	err := mock.Set("test_key", 123, time.Minute)
	if err != nil {
		t.Errorf("Set() error = %v, want nil", err)
	}

	// Test overwriting a value
	err = mock.Set("test_key", 456, time.Minute)
	if err != nil {
		t.Errorf("Set() error = %v, want nil", err)
	}

	value, err := mock.Get("test_key")
	if err != nil {
		t.Errorf("Get() error = %v, want nil", err)
	}
	if value != 456 {
		t.Errorf("Get() = %d, want 456", value)
	}
}

func TestMockClient_IncrementWithExpiration(t *testing.T) {
	mock := NewMockClient()

	// Test incrementing non-existent key (should initialize)
	value, err := mock.IncrementWithExpiration("test_key", 5, time.Minute)
	if err != nil {
		t.Errorf("IncrementWithExpiration() error = %v, want nil", err)
	}
	if value != 5 {
		t.Errorf("IncrementWithExpiration() = %d, want 5", value)
	}

	// Test incrementing existing key
	value, err = mock.IncrementWithExpiration("test_key", 3, time.Minute)
	if err != nil {
		t.Errorf("IncrementWithExpiration() error = %v, want nil", err)
	}
	if value != 8 {
		t.Errorf("IncrementWithExpiration() = %d, want 8", value)
	}
}

func TestMockClient_Delete(t *testing.T) {
	mock := NewMockClient()

	// Set a value
	err := mock.Set("test_key", 42, time.Minute)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// Verify it exists
	value, err := mock.Get("test_key")
	if err != nil {
		t.Errorf("Get() error = %v, want nil", err)
	}
	if value != 42 {
		t.Errorf("Get() = %d, want 42", value)
	}

	// Delete it
	err = mock.Delete("test_key")
	if err != nil {
		t.Errorf("Delete() error = %v, want nil", err)
	}

	// Verify it's gone
	value, err = mock.Get("test_key")
	if err != nil {
		t.Errorf("Get() error = %v, want nil", err)
	}
	if value != 0 {
		t.Errorf("Get() = %d, want 0", value)
	}
}

func TestMockClient_HealthCheck(t *testing.T) {
	mock := NewMockClient()

	// Test health check on open client
	err := mock.HealthCheck()
	if err != nil {
		t.Errorf("HealthCheck() error = %v, want nil", err)
	}

	// Test health check on closed client
	err = mock.Close()
	if err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}

	err = mock.HealthCheck()
	if err == nil {
		t.Error("HealthCheck() on closed client should return error")
	}
}

func TestMockClient_Expiration(t *testing.T) {
	mock := NewMockClient()

	// Set a value with short expiration
	err := mock.Set("expiring_key", 42, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// Wait for expiration
	time.Sleep(20 * time.Millisecond)

	// Verify it's expired
	value, err := mock.Get("expiring_key")
	if err != nil {
		t.Errorf("Get() error = %v, want nil", err)
	}
	if value != 0 {
		t.Errorf("Get() = %d, want 0 (expired)", value)
	}
}

func TestMockClient_Clear(t *testing.T) {
	mock := NewMockClient()

	// Set some values
	mock.Set("key1", 1, time.Minute)
	mock.Set("key2", 2, time.Minute)
	mock.Set("key3", 3, time.Minute)

	// Clear all
	mock.Clear()

	// Verify all are gone
	if value, _ := mock.Get("key1"); value != 0 {
		t.Errorf("key1 should be cleared, got %d", value)
	}
	if value, _ := mock.Get("key2"); value != 0 {
		t.Errorf("key2 should be cleared, got %d", value)
	}
	if value, _ := mock.Get("key3"); value != 0 {
		t.Errorf("key3 should be cleared, got %d", value)
	}
}
