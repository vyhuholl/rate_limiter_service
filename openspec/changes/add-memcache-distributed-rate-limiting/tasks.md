## 1. Memcache Configuration

- [x] 1.1 Add Memcache configuration fields to Config struct
- [x] 1.2 Add Memcache section to FileConfig struct
- [x] 1.3 Implement environment variable loading for Memcache settings
- [x] 1.4 Implement config file loading for Memcache settings
- [x] 1.5 Add validation for Memcache configuration values
- [x] 1.6 Add failure mode enum (Allow, Deny) for Memcache failures
- [x] 1.7 Write table-driven tests for Memcache configuration loading

## 2. Memcache Client Wrapper

- [x] 2.1 Create pkg/memcache package
- [x] 2.2 Define MemcacheClient struct wrapping gomemcache.Client
- [x] 2.3 Implement NewMemcacheClient function with connection pooling
- [x] 2.4 Implement IncrementWithExpiration method for atomic increments
- [x] 2.5 Implement Get method for retrieving counter values
- [x] 2.6 Implement Set method for setting counter values with expiration
- [x] 2.7 Add health check method for Memcache connectivity
- [x] 2.8 Write table-driven tests for MemcacheClient wrapper
- [x] 2.9 Add mock implementation for testing

## 3. Distributed Limiter Interface

- [x] 3.1 Define Limiter interface that both in-memory and distributed limiters implement
- [x] 3.2 Define DistributedLimiter interface extending Limiter with Memcache-specific methods
- [x] 3.3 Add failure mode handling to DistributedLimiter interface
- [x] 3.4 Write interface tests for Limiter and DistributedLimiter

## 4. Distributed Global Limiter

- [x] 4.1 Create pkg/middleware/distributed package
- [x] 4.2 Define DistributedGlobalLimiter struct
- [x] 4.3 Implement NewDistributedGlobalLimiter function
- [x] 4.4 Implement Allow method with Memcache counter increment
- [x] 4.5 Implement GetRemainingTokens method
- [x] 4.6 Add Memcache failure handling (allow/deny modes)
- [x] 4.7 Add logging for Memcache failures
- [x] 4.8 Write table-driven tests for DistributedGlobalLimiter
- [x] 4.9 Write integration tests with mock Memcache

## 5. Distributed Per-Endpoint Limiter

- [x] 5.1 Define DistributedPerEndpointLimiter struct
- [x] 5.2 Implement NewDistributedPerEndpointLimiter function
- [x] 5.3 Implement Allow method with endpoint-specific keys
- [x] 5.4 Implement GetRemainingTokens method
- [x] 5.5 Add Memcache failure handling
- [x] 5.6 Write table-driven tests for DistributedPerEndpointLimiter
- [x] 5.7 Write integration tests with mock Memcache

## 6. Distributed HTTP-Only Limiter

- [x] 6.1 Define DistributedHTTPLimiter struct
- [x] 6.2 Implement NewDistributedHTTPLimiter function
- [x] 6.3 Implement Allow method with HTTP-specific keys
- [x] 6.4 Implement GetRemainingTokens method
- [x] 6.5 Add Memcache failure handling
- [x] 6.6 Write table-driven tests for DistributedHTTPLimiter
- [x] 6.7 Write integration tests with mock Memcache

## 7. Distributed gRPC-Only Limiter

- [x] 7.1 Define DistributedGRPCLimiter struct
- [x] 7.2 Implement NewDistributedGRPCLimiter function
- [x] 7.3 Implement Allow method with gRPC-specific keys
- [x] 7.4 Implement GetRemainingTokens method
- [x] 7.5 Add Memcache failure handling
- [x] 7.6 Write table-driven tests for DistributedGRPCLimiter
- [x] 7.7 Write integration tests with mock Memcache

## 8. Limiter Factory

- [x] 8.1 Create limiter factory function that returns in-memory or distributed limiters
- [x] 8.2 Add logic to detect if Memcache is configured
- [x] 8.3 Return in-memory limiters when Memcache is not configured
- [x] 8.4 Return distributed limiters when Memcache is configured
- [x] 8.5 Write tests for limiter factory function

## 9. HTTP Middleware Integration

- [x] 9.1 Update HTTP middleware to use limiter factory
- [x] 9.2 Ensure backward compatibility with in-memory limiters
- [x] 9.3 Add Memcache failure mode handling in middleware
- [x] 9.4 Update rate limit headers to work with distributed limiters
- [x] 9.5 Write integration tests for HTTP middleware with distributed limiters
- [x] 9.6 Write end-to-end tests with real Memcache (optional)

## 10. gRPC Interceptor Integration

- [x] 10.1 Update gRPC interceptor to use limiter factory
- [x] 10.2 Ensure backward compatibility with in-memory limiters
- [x] 10.3 Add Memcache failure mode handling in interceptor
- [x] 10.4 Write integration tests for gRPC interceptor with distributed limiters
- [x] 10.5 Write end-to-end tests with real Memcache (optional)

## 11. Dependency Management

- [x] 11.1 Add github.com/bradfitz/gomemcache to go.mod
- [x] 11.2 Run go mod tidy to update dependencies
- [x] 11.3 Verify no dependency conflicts

## 12. Documentation

- [x] 12.1 Update README.md with Memcache distributed rate limiting usage examples
- [x] 12.2 Document Memcache configuration options (environment variables and config file)
- [x] 12.3 Document failure modes and their implications
- [x] 12.4 Add godoc comments to all new exported functions and types
- [x] 12.5 Document key structure in Memcache
- [x] 12.6 Add example configuration files with Memcache settings
- [x] 12.7 Document performance considerations and best practices

## 13. Testing

- [x] 13.1 Add golangci-lint rules for new code
- [x] 13.2 Ensure all tests pass with go test ./...
- [x] 13.3 Verify test coverage meets requirements
- [x] 13.4 Write tests for Memcache failure scenarios
- [x] 13.5 Write tests for concurrent access to distributed limiters
- [x] 13.6 Write tests for key expiration behavior

## 14. Validation

- [x] 14.1 Run golangci-lint and fix all issues
- [x] 14.2 Ensure all tests pass
- [x] 14.3 Verify backward compatibility with existing in-memory limiters
- [x] 14.4 Test with real Memcache server (optional integration test)
- [x] 14.5 Validate configuration file format with Memcache settings
