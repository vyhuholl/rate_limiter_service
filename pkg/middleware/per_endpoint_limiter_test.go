package middleware

import (
	"testing"

	"rate_limiter_service/internal/config"
)

func TestNewPerEndpointLimiter(t *testing.T) {
	cfg := config.Config{
		UserHeader:            "X-User-ID",
		PerEndpointRate:       10,
		GlobalRate:            100,
		GlobalBurstSize:       5,
		PerEndpointBurstSize:  5,
	}

	limiter := NewPerEndpointLimiter(cfg)

	if limiter == nil {
		t.Fatal("NewPerEndpointLimiter returned nil")
	}

	if limiter.config != cfg {
		t.Errorf("config not set correctly, got %+v, want %+v", limiter.config, cfg)
	}
}

func TestPerEndpointLimiter_Allow(t *testing.T) {
	cfg := config.Config{
		UserHeader:            "X-User-ID",
		PerEndpointRate:       2, // 2 requests per second
		GlobalRate:            100,
		GlobalBurstSize:       2,
		PerEndpointBurstSize:  2,
	}

	limiter := NewPerEndpointLimiter(cfg)

	const userID = "user123"
	const method = "GET"
	path := "/api/users"

	// Should allow first 2 requests
	if !limiter.Allow(userID, method, path) {
		t.Error("First request should be allowed")
	}
	if !limiter.Allow(userID, method, path) {
		t.Error("Second request should be allowed")
	}

	// Third request should be denied
	if limiter.Allow(userID, method, path) {
		t.Error("Third request should be denied")
	}

	// Different endpoint should be allowed
	if !limiter.Allow(userID, "POST", path) {
		t.Error("Different endpoint should be allowed")
	}

	// Different user should be allowed
	if !limiter.Allow("user456", method, path) {
		t.Error("Different user should be allowed")
	}
}

func TestPerEndpointLimiter_GetRemainingTokens(t *testing.T) {
	cfg := config.Config{
		UserHeader:            "X-User-ID",
		PerEndpointRate:       5,
		GlobalRate:            100,
		GlobalBurstSize:       5,
		PerEndpointBurstSize:  5,
	}

	limiter := NewPerEndpointLimiter(cfg)

	userID := "user123"
	method := "GET"
	path := "/api/data"

	// Initially should have full capacity
	if remaining := limiter.GetRemainingTokens(userID, method, path); remaining != 5 {
		t.Errorf("Initial remaining tokens = %d, want 5", remaining)
	}

	// After consuming some tokens
	limiter.Allow(userID, method, path)
	limiter.Allow(userID, method, path)

	if remaining := limiter.GetRemainingTokens(userID, method, path); remaining != 3 {
		t.Errorf("Remaining tokens after 2 requests = %d, want 3", remaining)
	}
}

func TestPerEndpointLimiter_SeparateEndpoints(t *testing.T) {
	cfg := config.Config{
		UserHeader:            "X-User-ID",
		PerEndpointRate:       1,
		GlobalRate:            100,
		GlobalBurstSize:       1,
		PerEndpointBurstSize:  1,
	}

	limiter := NewPerEndpointLimiter(cfg)

	userID := "user123"

	// Each endpoint should have independent limits
	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/users"},
		{"POST", "/api/users"},
		{"GET", "/api/orders"},
	}

	for _, endpoint := range endpoints {
		// First request to this endpoint should be allowed
		if !limiter.Allow(userID, endpoint.method, endpoint.path) {
			t.Errorf("First request to %s %s should be allowed", endpoint.method, endpoint.path)
		}

		// Second request should be denied
		if limiter.Allow(userID, endpoint.method, endpoint.path) {
			t.Errorf("Second request to %s %s should be denied", endpoint.method, endpoint.path)
		}
	}
}

func TestPerEndpointLimiter_SeparateUsers(t *testing.T) {
	cfg := config.Config{
		UserHeader:            "X-User-ID",
		PerEndpointRate:       1,
		GlobalRate:            100,
		GlobalBurstSize:       1,
		PerEndpointBurstSize:  1,
	}

	limiter := NewPerEndpointLimiter(cfg)

	method := "GET"
	path := "/api/data"

	users := []string{"user1", "user2", "user3"}

	for _, userID := range users {
		// First request for this user should be allowed
		if !limiter.Allow(userID, method, path) {
			t.Errorf("First request for user %s should be allowed", userID)
		}

		// Second request should be denied
		if limiter.Allow(userID, method, path) {
			t.Errorf("Second request for user %s should be denied", userID)
		}
	}
}