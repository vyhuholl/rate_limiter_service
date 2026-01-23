# Rate Limiter Service

A reusable rate limiting service for Go applications supporting both HTTP and gRPC, implementing token bucket algorithm with three-tier rate limiting (global, protocol-specific, and per-method).

## Features

- **Three-Tier Rate Limiting**: Global, HTTP-only, and gRPC-only limits
- **Per-Method Rate Limiting**: Separate rate limits for each HTTP endpoint or gRPC method per user
- **Global Rate Limiting**: Overall rate limits across all requests per user
- **Token Bucket Algorithm**: Smooth rate limiting with burst capacity
- **Configurable User Identification**: Identify users via HTTP headers or gRPC metadata
- **HTTP 429 Responses**: Proper rate limit exceeded responses with headers
- **gRPC ResourceExhausted Status**: Proper gRPC error responses
- **Configuration Files**: JSON/YAML configuration support
- **Thread-Safe**: Concurrent request handling

## Configuration

Configure rate limits using environment variables or JSON/YAML configuration files.

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `RATE_LIMIT_CONFIG_PATH` | Path to JSON/YAML config file | - |
| `RATE_LIMIT_USER_HEADER` | HTTP header for user identification | `X-User-ID` |
| `RATE_LIMIT_GRPC_METADATA_KEY` | gRPC metadata key for user identification | `user-id` |
| `RATE_LIMIT_GLOBAL` | Global requests per second per user | `100` |
| `RATE_LIMIT_GLOBAL_BURST_SIZE` | Global burst capacity | `10` |
| `RATE_LIMIT_HTTP_RATE` | HTTP requests per second per user | `50` |
| `RATE_LIMIT_HTTP_BURST_SIZE` | HTTP burst capacity | `5` |
| `RATE_LIMIT_GRPC_RATE` | gRPC requests per second per user | `50` |
| `RATE_LIMIT_GRPC_BURST_SIZE` | gRPC burst capacity | `5` |
| `RATE_LIMIT_BURST_SIZE` | Legacy burst capacity for both limiters | `10` |

### Configuration File

Use `RATE_LIMIT_CONFIG_PATH` to specify a JSON or YAML configuration file. See `examples/config.json` and `examples/config.yaml` for format examples.

## Usage

```go
package main

import (
    "net/http"
    "rate_limiter_service/internal/config"
    "rate_limiter_service/pkg/middleware"
)

func main() {
    // Load configuration from environment or config file
    cfg, err := config.Load()
    if err != nil {
        panic(err)
    }

    // HTTP server with rate limiting
    rateLimiter := middleware.NewMiddleware(cfg)
    mux := http.NewServeMux()
    mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("User data"))
    })
    wrappedMux := rateLimiter.Handler(mux)

    go http.ListenAndServe(":8080", wrappedMux)

    // gRPC server with rate limiting
    grpcInterceptor := grpc.NewInterceptor(cfg)
    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(grpcInterceptor.UnaryInterceptor()),
    )
    // Register your gRPC services...

    grpcServer.Serve(listener)
}
```

## Rate Limiting Behavior

- **Global**: All requests from a user count toward the global limit
- **Protocol-Specific**: HTTP and gRPC requests have separate limits per user
- **Per-Method**: Each HTTP endpoint or gRPC method has configurable rate limits per user
- **HTTP Responses**: Rate limited HTTP requests return 429 with `X-RateLimit-Limit` and `Retry-After` headers
- **gRPC Responses**: Rate limited gRPC requests return `ResourceExhausted` status
- **User Identification**: HTTP uses headers, gRPC uses metadata, falls back to "anonymous" if missing

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
- **GlobalLimiter**: Manages global rate limits across all requests per user
- **HTTPLimiter**: Manages HTTP-specific rate limits per user
- **GRPCLimiter**: Manages gRPC-specific rate limits per user
- **PerEndpointLimiter**: Manages per-method rate limits for HTTP requests
- **GRPCMethodLimiter**: Manages per-method rate limits for gRPC requests
- **Middleware**: HTTP handler wrapper with three-tier rate limiting
- **Interceptor**: gRPC unary interceptor with three-tier rate limiting
- **Config**: Supports both environment variables and JSON/YAML files