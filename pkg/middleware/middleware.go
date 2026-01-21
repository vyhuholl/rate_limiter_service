package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"rate_limiter_service/internal/config"
)

// Middleware wraps an HTTP handler with rate limiting
type Middleware struct {
	config         config.Config
	perEndpointLimiter *PerEndpointLimiter
	globalLimiter   *GlobalLimiter
}

// NewMiddleware creates a new rate limiting middleware
func NewMiddleware(cfg config.Config) *Middleware {
	return &Middleware{
		config:         cfg,
		perEndpointLimiter: NewPerEndpointLimiter(cfg),
		globalLimiter:   NewGlobalLimiter(cfg),
	}
}

// Handler wraps an HTTP handler with rate limiting
func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := m.extractUserID(r)

		// Check global limit first
		if !m.globalLimiter.Allow(userID) {
			m.writeRateLimitResponse(w, "global")
			return
		}

		// Check per-endpoint limit
		if !m.perEndpointLimiter.Allow(userID, r.Method, r.URL.Path) {
			m.writeRateLimitResponse(w, "per-endpoint")
			return
		}

		// Request allowed, call next handler
		next.ServeHTTP(w, r)
	})
}

// extractUserID extracts the user ID from the configured header
func (m *Middleware) extractUserID(r *http.Request) string {
	userID := r.Header.Get(m.config.UserHeader)
	if userID == "" {
		// If no user header is present, treat as anonymous user
		// In a real implementation, you might want to handle this differently
		return "anonymous"
	}
	return userID
}

// writeRateLimitResponse writes an HTTP 429 response with appropriate headers
func (m *Middleware) writeRateLimitResponse(w http.ResponseWriter, limitType string) {
	w.Header().Set("Content-Type", "application/json")

	// Set rate limit headers
	// Note: These are simplified; in production you might want more detailed headers
	w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", m.getRateLimitForType(limitType)))
	w.Header().Set("Retry-After", strconv.Itoa(m.getRetryAfterSeconds()))

	w.WriteHeader(http.StatusTooManyRequests)

	// Write a simple JSON response
	response := fmt.Sprintf(`{"error": "rate limit exceeded", "type": "%s"}`, limitType)
	_, _ = w.Write([]byte(response))
}

// getRateLimitForType returns the rate limit for the given type
func (m *Middleware) getRateLimitForType(limitType string) int {
	switch limitType {
	case "global":
		return m.config.GlobalRate
	case "per-endpoint":
		return m.config.PerEndpointRate
	default:
		return 0
	}
}

// getRetryAfterSeconds returns a reasonable retry-after time in seconds
func (m *Middleware) getRetryAfterSeconds() int {
	// Calculate based on token refill interval
	// Use the more restrictive of per-endpoint or global rate
	refillInterval := m.config.GetRefillInterval(m.config.PerEndpointRate)
	if globalInterval := m.config.GetRefillInterval(m.config.GlobalRate); globalInterval > refillInterval {
		refillInterval = globalInterval
	}

	// Return approximately how long until next token is available
	return int(refillInterval.Seconds())
}

// Reset clears all rate limiting state for testing purposes
func (m *Middleware) Reset() {
	m.perEndpointLimiter.Reset()
	m.globalLimiter.Reset()
}