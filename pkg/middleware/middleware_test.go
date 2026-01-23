package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"rate_limiter_service/internal/config"
)

func TestNewMiddleware(t *testing.T) {
	cfg := config.Config{
		UserHeader:            "X-User-ID",
		PerEndpointRate:       10,
		GlobalRate:            100,
		GlobalBurstSize:       5,
		PerEndpointBurstSize:  5,
	}

	middleware := NewMiddleware(cfg)

	if middleware == nil {
		t.Fatal("NewMiddleware returned nil")
	}

	// Check that config is set (can't compare structs with maps directly)
	if middleware.config.UserHeader != cfg.UserHeader {
		t.Errorf("config UserHeader not set correctly")
	}

	if middleware.perEndpointLimiter == nil {
		t.Error("perEndpointLimiter not initialized")
	}

	if middleware.globalLimiter == nil {
		t.Error("globalLimiter not initialized")
	}

	if middleware.httpLimiter == nil {
		t.Error("httpLimiter not initialized")
	}
}

func TestMiddleware_Handler_Allow(t *testing.T) {
	cfg := config.Config{
		UserHeader:            "X-User-ID",
		PerEndpointRate:       10,
		GlobalRate:            100,
		GlobalBurstSize:       10,
		PerEndpointBurstSize:  10,
		HTTPRate:              50,
		HTTPBurstSize:         5,
		HTTPDefaultMethodRate: 10,
	}

	middleware := NewMiddleware(cfg)

	// Create a simple handler that returns 200
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	wrappedHandler := middleware.Handler(handler)

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("X-User-ID", "user123")

	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if body := w.Body.String(); body != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", body)
	}
}

func TestMiddleware_Handler_GlobalLimitExceeded(t *testing.T) {
	cfg := config.Config{
		UserHeader:            "X-User-ID",
		PerEndpointRate:       10,
		GlobalRate:            1, // Only 1 request per second globally
		GlobalBurstSize:       1,
		PerEndpointBurstSize:  10,
		HTTPRate:              100, // High HTTP rate so global limit is hit first
		HTTPBurstSize:         10,
		HTTPDefaultMethodRate: 10,
	}

	middleware := NewMiddleware(cfg)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := middleware.Handler(handler)

	// First request should be allowed
	req1 := httptest.NewRequest("GET", "/api/test", nil)
	req1.Header.Set("X-User-ID", "user123")

	w1 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("First request should be allowed, got status %d", w1.Code)
	}

	// Second request should be rate limited
	req2 := httptest.NewRequest("GET", "/api/test", nil)
	req2.Header.Set("X-User-ID", "user123")

	w2 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("Second request should be rate limited, got status %d", w2.Code)
	}

	if !strings.Contains(w2.Body.String(), "rate limit exceeded") {
		t.Errorf("Response should contain rate limit message, got: %s", w2.Body.String())
	}

	if !strings.Contains(w2.Body.String(), "global") {
		t.Errorf("Response should indicate global limit, got: %s", w2.Body.String())
	}
}

func TestMiddleware_Handler_PerEndpointLimitExceeded(t *testing.T) {
	cfg := config.Config{
		UserHeader:            "X-User-ID",
		PerEndpointRate:       1, // Only 1 request per second per endpoint
		GlobalRate:            100,
		GlobalBurstSize:       10,
		PerEndpointBurstSize:  1,
		HTTPRate:              100, // High HTTP rate so per-method limit is hit first
		HTTPBurstSize:         10,
		HTTPDefaultMethodRate: 1,
	}

	middleware := NewMiddleware(cfg)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := middleware.Handler(handler)

	// First request should be allowed
	req1 := httptest.NewRequest("GET", "/api/test", nil)
	req1.Header.Set("X-User-ID", "user123")

	w1 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("First request should be allowed, got status %d", w1.Code)
	}

	// Second request to same endpoint should be rate limited
	req2 := httptest.NewRequest("GET", "/api/test", nil)
	req2.Header.Set("X-User-ID", "user123")

	w2 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("Second request should be rate limited, got status %d", w2.Code)
	}

	if !strings.Contains(w2.Body.String(), "per-method") {
		t.Errorf("Response should indicate per-method limit, got: %s", w2.Body.String())
	}

	// Request to different endpoint should be allowed
	req3 := httptest.NewRequest("POST", "/api/test", nil)
	req3.Header.Set("X-User-ID", "user123")

	w3 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("Request to different endpoint should be allowed, got status %d", w3.Code)
	}
}

func TestMiddleware_ExtractUserID(t *testing.T) {
	cfg := config.Config{
		UserHeader:           "X-Custom-User",
		HTTPRate:             50,
		HTTPBurstSize:        5,
		HTTPDefaultMethodRate: 10,
	}

	middleware := NewMiddleware(cfg)

	tests := []struct {
		name     string
		headers  map[string]string
		expected string
	}{
		{
			name:     "header present",
			headers:  map[string]string{"X-Custom-User": "user123"},
			expected: "user123",
		},
		{
			name:     "header missing",
			headers:  map[string]string{},
			expected: "anonymous",
		},
		{
			name:     "header empty",
			headers:  map[string]string{"X-Custom-User": ""},
			expected: "anonymous",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			result := middleware.extractUserID(req)
			if result != tt.expected {
				t.Errorf("extractUserID() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestMiddleware_Handler_RateLimitHeaders(t *testing.T) {
	cfg := config.Config{
		UserHeader:            "X-User-ID",
		PerEndpointRate:       1,
		GlobalRate:            100,
		GlobalBurstSize:       10,
		PerEndpointBurstSize:  1,
		HTTPRate:              100,
		HTTPBurstSize:         10,
		HTTPDefaultMethodRate: 1,
	}

	middleware := NewMiddleware(cfg)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := middleware.Handler(handler)

	// Make two requests to trigger rate limit
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("X-User-ID", "user123")

		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)

		if i == 1 { // Second request should be rate limited
			if w.Code != http.StatusTooManyRequests {
				t.Errorf("Request %d should be rate limited", i+1)
			}

			// Check rate limit headers
			if limit := w.Header().Get("X-RateLimit-Limit"); limit != "1" {
				t.Errorf("X-RateLimit-Limit should be '1', got '%s'", limit)
			}

			if retryAfter := w.Header().Get("Retry-After"); retryAfter != "1" {
				t.Errorf("Retry-After should be '1', got '%s'", retryAfter)
			}
		}
	}
}