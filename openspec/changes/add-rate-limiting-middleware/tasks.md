## 1. Project Setup
- [ ] 1.1 Initialize Go module with go.mod
- [ ] 1.2 Create project directory structure (pkg/middleware, internal/config, etc.)
- [ ] 1.3 Add .golangci.yml configuration file
- [ ] 1.4 Set up Makefile with common targets (build, test, lint)

## 2. Configuration
- [ ] 2.1 Define configuration struct with rate limit settings
- [ ] 2.2 Implement environment variable loading for configuration
- [ ] 2.3 Add validation for configuration values
- [ ] 2.4 Write table-driven tests for configuration loading

## 3. Token Bucket Implementation
- [ ] 3.1 Define TokenBucket interface and struct
- [ ] 3.2 Implement token bucket logic (consume, refill)
- [ ] 3.3 Add thread-safe operations using sync.Mutex
- [ ] 3.4 Write table-driven tests for token bucket behavior

## 4. Per-Endpoint Rate Limiter
- [ ] 4.1 Define PerEndpointLimiter struct with endpoint-keyed buckets
- [ ] 4.2 Implement per-endpoint rate limiting logic
- [ ] 4.3 Add user identification from configured header
- [ ] 4.4 Write table-driven tests for per-endpoint limiting

## 5. Global Rate Limiter
- [ ] 5.1 Define GlobalLimiter struct with user-keyed buckets
- [ ] 5.2 Implement global rate limiting logic
- [ ] 5.3 Add user identification from configured header
- [ ] 5.4 Write table-driven tests for global limiting

## 6. HTTP Middleware
- [ ] 6.1 Define Middleware interface and implementation
- [ ] 6.2 Implement HTTP handler wrapper for rate limiting
- [ ] 6.3 Add HTTP 429 response with appropriate headers (Retry-After, X-RateLimit-*)
- [ ] 6.4 Write integration tests with http.ServeMux

## 7. Documentation
- [ ] 7.1 Update README.md with usage examples
- [ ] 7.2 Add godoc comments to exported functions and types
- [ ] 7.3 Document environment variables

## 8. Validation
- [ ] 8.1 Run golangci-lint and fix all issues
- [ ] 8.2 Ensure all tests pass with go test ./...
- [ ] 8.3 Verify test coverage meets requirements
