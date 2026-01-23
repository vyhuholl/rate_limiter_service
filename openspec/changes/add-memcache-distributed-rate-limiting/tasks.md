## 1. Memcache Configuration

- [ ] 1.1 Add Memcache configuration fields to Config struct
- [ ] 1.2 Add Memcache section to FileConfig struct
- [ ] 1.3 Implement environment variable loading for Memcache settings
- [ ] 1.4 Implement config file loading for Memcache settings
- [ ] 1.5 Add validation for Memcache configuration values
- [ ] 1.6 Add failure mode enum (Allow, Deny) for Memcache failures
- [ ] 1.7 Write table-driven tests for Memcache configuration loading

## 2. Memcache Client Wrapper

- [ ] 2.1 Create pkg/memcache package
- [ ] 2.2 Define MemcacheClient struct wrapping gomemcache.Client
- [ ] 2.3 Implement NewMemcacheClient function with connection pooling
- [ ] 2.4 Implement IncrementWithExpiration method for atomic increments
- [ ] 2.5 Implement Get method for retrieving counter values
- [ ] 2.6 Implement Set method for setting counter values with expiration
- [ ] 2.7 Add health check method for Memcache connectivity
- [ ] 2.8 Write table-driven tests for MemcacheClient wrapper
- [ ] 2.9 Add mock implementation for testing

## 3. Distributed Limiter Interface

- [ ] 3.1 Define Limiter interface that both in-memory and distributed limiters implement
- [ ] 3.2 Define DistributedLimiter interface extending Limiter with Memcache-specific methods
- [ ] 3.3 Add failure mode handling to DistributedLimiter interface
- [ ] 3.4 Write interface tests for Limiter and DistributedLimiter

## 4. Distributed Global Limiter

- [ ] 4.1 Create pkg/middleware/distributed package
- [ ] 4.2 Define DistributedGlobalLimiter struct
- [ ] 4.3 Implement NewDistributedGlobalLimiter function
- [ ] 4.4 Implement Allow method with Memcache counter increment
- [ ] 4.5 Implement GetRemainingTokens method
- [ ] 4.6 Add Memcache failure handling (allow/deny modes)
- [ ] 4.7 Add logging for Memcache failures
- [ ] 4.8 Write table-driven tests for DistributedGlobalLimiter
- [ ] 4.9 Write integration tests with mock Memcache

## 5. Distributed Per-Endpoint Limiter

- [ ] 5.1 Define DistributedPerEndpointLimiter struct
- [ ] 5.2 Implement NewDistributedPerEndpointLimiter function
- [ ] 5.3 Implement Allow method with endpoint-specific keys
- [ ] 5.4 Implement GetRemainingTokens method
- [ ] 5.5 Add Memcache failure handling
- [ ] 5.6 Write table-driven tests for DistributedPerEndpointLimiter
- [ ] 5.7 Write integration tests with mock Memcache

## 6. Distributed HTTP-Only Limiter

- [ ] 6.1 Define DistributedHTTPLimiter struct
- [ ] 6.2 Implement NewDistributedHTTPLimiter function
- [ ] 6.3 Implement Allow method with HTTP-specific keys
- [ ] 6.4 Implement GetRemainingTokens method
- [ ] 6.5 Add Memcache failure handling
- [ ] 6.6 Write table-driven tests for DistributedHTTPLimiter
- [ ] 6.7 Write integration tests with mock Memcache

## 7. Distributed gRPC-Only Limiter

- [ ] 7.1 Define DistributedGRPCLimiter struct
- [ ] 7.2 Implement NewDistributedGRPCLimiter function
- [ ] 7.3 Implement Allow method with gRPC-specific keys
- [ ] 7.4 Implement GetRemainingTokens method
- [ ] 7.5 Add Memcache failure handling
- [ ] 7.6 Write table-driven tests for DistributedGRPCLimiter
- [ ] 7.7 Write integration tests with mock Memcache

## 8. Limiter Factory

- [ ] 8.1 Create limiter factory function that returns in-memory or distributed limiters
- [ ] 8.2 Add logic to detect if Memcache is configured
- [ ] 8.3 Return in-memory limiters when Memcache is not configured
- [ ] 8.4 Return distributed limiters when Memcache is configured
- [ ] 8.5 Write tests for limiter factory function

## 9. HTTP Middleware Integration

- [ ] 9.1 Update HTTP middleware to use limiter factory
- [ ] 9.2 Ensure backward compatibility with in-memory limiters
- [ ] 9.3 Add Memcache failure mode handling in middleware
- [ ] 9.4 Update rate limit headers to work with distributed limiters
- [ ] 9.5 Write integration tests for HTTP middleware with distributed limiters
- [ ] 9.6 Write end-to-end tests with real Memcache (optional)

## 10. gRPC Interceptor Integration

- [ ] 10.1 Update gRPC interceptor to use limiter factory
- [ ] 10.2 Ensure backward compatibility with in-memory limiters
- [ ] 10.3 Add Memcache failure mode handling in interceptor
- [ ] 10.4 Write integration tests for gRPC interceptor with distributed limiters
- [ ] 10.5 Write end-to-end tests with real Memcache (optional)

## 11. Dependency Management

- [ ] 11.1 Add github.com/bradfitz/gomemcache to go.mod
- [ ] 11.2 Run go mod tidy to update dependencies
- [ ] 11.3 Verify no dependency conflicts

## 12. Documentation

- [ ] 12.1 Update README.md with Memcache distributed rate limiting usage examples
- [ ] 12.2 Document Memcache configuration options (environment variables and config file)
- [ ] 12.3 Document failure modes and their implications
- [ ] 12.4 Add godoc comments to all new exported functions and types
- [ ] 12.5 Document key structure in Memcache
- [ ] 12.6 Add example configuration files with Memcache settings
- [ ] 12.7 Document performance considerations and best practices

## 13. Testing

- [ ] 13.1 Add golangci-lint rules for new code
- [ ] 13.2 Ensure all tests pass with go test ./...
- [ ] 13.3 Verify test coverage meets requirements
- [ ] 13.4 Write tests for Memcache failure scenarios
- [ ] 13.5 Write tests for concurrent access to distributed limiters
- [ ] 13.6 Write tests for key expiration behavior

## 14. Validation

- [ ] 14.1 Run golangci-lint and fix all issues
- [ ] 14.2 Ensure all tests pass
- [ ] 14.3 Verify backward compatibility with existing in-memory limiters
- [ ] 14.4 Test with real Memcache server (optional integration test)
- [ ] 14.5 Validate configuration file format with Memcache settings
