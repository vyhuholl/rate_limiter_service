# Change: Add gRPC Rate Limiting

## Why

The current rate limiting implementation only supports HTTP requests. To support gRPC services, we need to extend the rate limiter to handle both HTTP and gRPC traffic with separate limits per protocol type, while maintaining global rate limiting across all requests.

## What Changes

- Add gRPC interceptor for rate limiting
- Implement three-tier rate limiting:
  - Global rate limiter (all requests - both HTTP and gRPC)
  - HTTP-only rate limiter
  - gRPC-only rate limiter
- Add per-method rate limits for both HTTP and gRPC (each method can have its own rate)
- Add JSON/YAML configuration file support for method-to-rate mappings
- Use default rate for methods without explicit configuration
- Support user identification via gRPC metadata
- Return gRPC status code `ResourceExhausted` when rate limited

## Impact

- Affected specs: `rate-limiting` (MODIFIED to add gRPC support)
- Affected code:
  - New `pkg/grpc/` package for gRPC interceptor
  - Extended `internal/config/` for file-based configuration
  - Modified `pkg/middleware/` for per-method limits
