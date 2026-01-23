package distributed

import (
	"testing"

	"rate_limiter_service/internal/config"
	"rate_limiter_service/pkg/memcache"
)

func TestNewGlobalLimiter(t *testing.T) {
	mock := memcache.NewMockClient()
	cfg := config.DefaultConfig()
	cfg.GlobalRate = 100

	limiter := NewGlobalLimiter(mock, cfg)

	if limiter == nil {
		t.Fatal("NewGlobalLimiter returned nil")
	}

	if limiter.client != mock {
		t.Error("client not set correctly")
	}

	if limiter.config.GlobalRate != 100 {
		t.Error("config not set correctly")
	}

	if limiter.scope != "global" {
		t.Error("scope not set correctly")
	}
}

func TestGlobalLimiter_Allow(t *testing.T) {
	mock := memcache.NewMockClient()
	cfg := config.DefaultConfig()
	cfg.GlobalRate = 2 // Only 2 requests per second
	cfg.MemcacheFailureMode = config.FailureModeAllow

	limiter := NewGlobalLimiter(mock, cfg)

	userID := "user123"

	// First request should be allowed
	if !limiter.Allow(userID) {
		t.Error("First request should be allowed")
	}

	// Second request should be allowed
	if !limiter.Allow(userID) {
		t.Error("Second request should be allowed")
	}

	// Third request should be denied
	if limiter.Allow(userID) {
		t.Error("Third request should be denied")
	}
}

func TestGlobalLimiter_GetRemainingTokens(t *testing.T) {
	mock := memcache.NewMockClient()
	cfg := config.DefaultConfig()
	cfg.GlobalRate = 10

	limiter := NewGlobalLimiter(mock, cfg)

	userID := "user123"

	// Initially should have full capacity
	remaining := limiter.GetRemainingTokens(userID)
	if remaining != 10 {
		t.Errorf("Initial remaining tokens = %d, want 10", remaining)
	}

	// Use some tokens
	for i := 0; i < 3; i++ {
		limiter.Allow(userID)
	}

	remaining = limiter.GetRemainingTokens(userID)
	if remaining != 7 {
		t.Errorf("Remaining tokens after 3 requests = %d, want 7", remaining)
	}
}

func TestGlobalLimiter_FailureModeAllow(t *testing.T) {
	mock := memcache.NewMockClient()
	cfg := config.DefaultConfig()
	cfg.GlobalRate = 10
	cfg.MemcacheFailureMode = config.FailureModeAllow

	limiter := NewGlobalLimiter(mock, cfg)

	// Close the mock to simulate failure
	mock.Close()

	// With fail-open, requests should be allowed even on failure
	if !limiter.Allow("user123") {
		t.Error("Request should be allowed in fail-open mode")
	}
}

func TestGlobalLimiter_FailureModeDeny(t *testing.T) {
	mock := memcache.NewMockClient()
	cfg := config.DefaultConfig()
	cfg.GlobalRate = 10
	cfg.MemcacheFailureMode = config.FailureModeDeny

	limiter := NewGlobalLimiter(mock, cfg)

	// Close the mock to simulate failure
	mock.Close()

	// With fail-closed, requests should be denied on failure
	if limiter.Allow("user123") {
		t.Error("Request should be denied in fail-closed mode")
	}
}

func TestGlobalLimiter_SeparateUsers(t *testing.T) {
	mock := memcache.NewMockClient()
	cfg := config.DefaultConfig()
	cfg.GlobalRate = 5

	limiter := NewGlobalLimiter(mock, cfg)

	user1 := "user1"
	user2 := "user2"

	// Each user should have their own limit
	for i := 0; i < 5; i++ {
		if !limiter.Allow(user1) {
			t.Errorf("Request %d for user1 should be allowed", i+1)
		}
	}

	for i := 0; i < 5; i++ {
		if !limiter.Allow(user2) {
			t.Errorf("Request %d for user2 should be allowed", i+1)
		}
	}

	// Both should be denied on 6th request
	if limiter.Allow(user1) {
		t.Error("6th request for user1 should be denied")
	}
	if limiter.Allow(user2) {
		t.Error("6th request for user2 should be denied")
	}
}

func TestGlobalLimiter_Reset(t *testing.T) {
	mock := memcache.NewMockClient()
	cfg := config.DefaultConfig()
	cfg.GlobalRate = 10

	limiter := NewGlobalLimiter(mock, cfg)

	userID := "user123"

	// Use some tokens
	for i := 0; i < 3; i++ {
		limiter.Allow(userID)
	}

	// Reset should be a no-op for distributed limiters
	limiter.Reset()

	// State is managed by Memcache, so we need to clear the mock
	mock.Clear()

	// After clearing mock, should have full capacity again
	remaining := limiter.GetRemainingTokens(userID)
	if remaining != 10 {
		t.Errorf("After reset and clear, remaining tokens = %d, want 10", remaining)
	}
}
