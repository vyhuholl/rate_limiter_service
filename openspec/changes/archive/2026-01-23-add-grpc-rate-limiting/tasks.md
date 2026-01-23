## 1. Configuration File Support

- [x] 1.1 Define configuration struct for JSON/YAML config file
- [x] 1.2 Implement JSON configuration file loading
- [x] 1.3 Implement YAML configuration file loading
- [x] 1.4 Add validation for configuration file contents
- [x] 1.5 Add support for RATE_LIMIT_CONFIG_PATH environment variable
- [x] 1.6 Maintain backward compatibility with existing environment variables
- [x] 1.7 Write table-driven tests for configuration file loading

## 2. Per-Method Rate Limiting

- [x] 2.1 Extend Config struct to include method-specific rate mappings
- [x] 2.2 Modify PerEndpointLimiter to use method-specific rates
- [x] 2.3 Add default rate fallback for unconfigured methods
- [x] 2.4 Update HTTP middleware to use per-method limits
- [x] 2.5 Write tests for per-method rate limiting

## 3. Three-Tier Rate Limiting

- [x] 3.1 Define HTTP-only rate limiter struct
- [x] 3.2 Define gRPC-only rate limiter struct
- [x] 3.3 Modify GlobalLimiter to track both HTTP and gRPC requests
- [x] 3.4 Update HTTP middleware to check HTTP-only limiter
- [x] 3.5 Write tests for three-tier rate limiting

## 4. gRPC Interceptor

- [x] 4.1 Create pkg/grpc package
- [x] 4.2 Define gRPC interceptor struct
- [x] 4.3 Implement unary interceptor for rate limiting
- [x] 4.4 Add user ID extraction from gRPC metadata
- [x] 4.5 Implement gRPC-only rate limiting check
- [x] 4.6 Implement per-method rate limiting for gRPC
- [x] 4.7 Return ResourceExhausted status when rate limited
- [x] 4.8 Write table-driven tests for gRPC interceptor

## 5. Integration and Testing

- [x] 5.1 Add golangci-lint rules for new code
- [x] 5.2 Write integration tests for HTTP middleware with config file
- [x] 5.3 Write integration tests for gRPC interceptor
- [x] 5.4 Write end-to-end tests for three-tier rate limiting
- [x] 5.5 Add example configuration files (JSON and YAML)

## 6. Documentation

- [x] 6.1 Update README.md with gRPC rate limiting usage examples
- [x] 6.2 Document configuration file format and options
- [x] 6.3 Add godoc comments to exported functions and types
- [x] 6.4 Document environment variables and config file precedence
- [x] 6.5 Add migration guide from environment-only configuration

## 7. Validation

- [x] 7.1 Run golangci-lint and fix all issues
- [x] 7.2 Ensure all tests pass with go test ./...
- [x] 7.3 Verify test coverage meets requirements
- [x] 7.4 Validate backward compatibility with existing HTTP middleware
