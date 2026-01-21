# Rate Limiter Service

A reusable HTTP rate limiting middleware for Go applications, implementing token bucket algorithm with both per-endpoint and global rate limits.

## Features

- **Per-Endpoint Rate Limiting**: Separate rate limits for each API endpoint per user
- **Global Rate Limiting**: Overall rate limits across all endpoints per user
- **Token Bucket Algorithm**: Smooth rate limiting with burst capacity
- **Configurable User Identification**: Identify users via HTTP headers
- **HTTP 429 Responses**: Proper rate limit exceeded responses with headers
- **Thread-Safe**: Concurrent request handling

## Configuration

Configure rate limits using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `RATE_LIMIT_USER_HEADER` | HTTP header for user identification | `X-User-ID` |
| `RATE_LIMIT_PER_ENDPOINT` | Requests per second per endpoint per user | `10` |
| `RATE_LIMIT_GLOBAL` | Requests per second globally per user | `100` |
| `RATE_LIMIT_BURST_SIZE` | Burst capacity for both limiters | `10` |
| `RATE_LIMIT_GLOBAL_BURST_SIZE` | Burst capacity for global limiter | `10` |
| `RATE_LIMIT_PER_ENDPOINT_BURST_SIZE` | Burst capacity for per-endpoint limiter | `10` |

## Usage

```go
package main

import (
    "net/http"
    "rate_limiter_service/internal/config"
    "rate_limiter_service/pkg/middleware"
)

func main() {
    // Load configuration from environment
    cfg, err := config.LoadFromEnv()
    if err != nil {
        panic(err)
    }

    // Create rate limiting middleware
    rateLimiter := middleware.NewMiddleware(cfg)

    // Create your HTTP handler
    mux := http.NewServeMux()
    mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("User data"))
    })
    mux.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Order data"))
    })

    // Wrap with rate limiting middleware
    wrappedMux := rateLimiter.Handler(mux)

    http.ListenAndServe(":8080", wrappedMux)
}
```

## Rate Limiting Behavior

- **Per-Endpoint**: Each endpoint (method + path) has independent rate limits per user
- **Global**: All requests from a user count toward the global limit
- **Headers**: Rate limited requests return HTTP 429 with `X-RateLimit-Limit` and `Retry-After` headers
- **User Identification**: Uses configured header (default: `X-User-ID`), falls back to "anonymous" if missing

## Example API Usage

```bash
# First request - allowed
curl -H "X-User-ID: user123" http://localhost:8080/api/users

# Exceed per-endpoint limit - 429 response
curl -H "X-User-ID: user123" http://localhost:8080/api/users
# Response: 429 Too Many Requests
# {"error": "rate limit exceeded", "type": "per-endpoint"}
```

## Testing

```bash
# Run tests
make test

# Run linter
make lint

# Build
make build
```

## Architecture

- **TokenBucket**: Implements token bucket algorithm with thread-safe operations
- **PerEndpointLimiter**: Manages per-endpoint rate limits using endpoint-keyed buckets
- **GlobalLimiter**: Manages global rate limits using user-keyed buckets
- **Middleware**: HTTP handler wrapper that applies both limiters