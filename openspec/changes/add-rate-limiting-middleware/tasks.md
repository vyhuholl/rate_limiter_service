## 1. Project Setup
- [x] 1.1 Initialize Go module with go.mod
- [x] 1.2 Create project directory structure (pkg/middleware, internal/config, etc.)
- [x] 1.3 Add .golangci.yml configuration file
- [x] 1.4 Set up Makefile with common targets (build, test, lint)

## 2. Configuration
- [x] 2.1 Define configuration struct with rate limit settings
- [x] 2.2 Implement environment variable loading for configuration
- [x] 2.3 Add validation for configuration values
- [x] 2.4 Write table-driven tests for configuration loading

## 3. Token Bucket Implementation
- [x] 3.1 Define TokenBucket interface and struct
- [x] 3.2 Implement token bucket logic (consume, refill)
- [x] 3.3 Add thread-safe operations using sync.Mutex
- [x] 3.4 Write table-driven tests for token bucket behavior

## 4. Per-Endpoint Rate Limiter
- [x] 4.1 Define PerEndpointLimiter struct with endpoint-keyed buckets
- [x] 4.2 Implement per-endpoint rate limiting logic
- [x] 4.3 Add user identification from configured header
- [x] 4.4 Write table-driven tests for per-endpoint limiting

## 5. Global Rate Limiter
- [x] 5.1 Define GlobalLimiter struct with user-keyed buckets
- [x] 5.2 Implement global rate limiting logic
- [x] 5.3 Add user identification from configured header
- [x] 5.4 Write table-driven tests for global limiting

## 6. HTTP Middleware
- [x] 6.1 Define Middleware interface and implementation
- [x] 6.2 Implement HTTP handler wrapper for rate limiting
- [x] 6.3 Add HTTP 429 response with appropriate headers (Retry-After, X-RateLimit-*)
- [x] 6.4 Write integration tests with http.ServeMux

## 7. Documentation
- [x] 7.1 Update README.md with usage examples
- [x] 7.2 Add godoc comments to exported functions and types
- [x] 7.3 Document environment variables

## 8. Validation
- [x] 8.1 Run golangci-lint and fix all issues
- [x] 8.2 Ensure all tests pass with go test ./...
- [x] 8.3 Verify test coverage meets requirements
