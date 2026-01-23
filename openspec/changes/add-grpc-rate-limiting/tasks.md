## 1. Configuration File Support

- [ ] 1.1 Define configuration struct for JSON/YAML config file
- [ ] 1.2 Implement JSON configuration file loading
- [ ] 1.3 Implement YAML configuration file loading
- [ ] 1.4 Add validation for configuration file contents
- [ ] 1.5 Add support for RATE_LIMIT_CONFIG_PATH environment variable
- [ ] 1.6 Maintain backward compatibility with existing environment variables
- [ ] 1.7 Write table-driven tests for configuration file loading

## 2. Per-Method Rate Limiting

- [ ] 2.1 Extend Config struct to include method-specific rate mappings
- [ ] 2.2 Modify PerEndpointLimiter to use method-specific rates
- [ ] 2.3 Add default rate fallback for unconfigured methods
- [ ] 2.4 Update HTTP middleware to use per-method limits
- [ ] 2.5 Write tests for per-method rate limiting

## 3. Three-Tier Rate Limiting

- [ ] 3.1 Define HTTP-only rate limiter struct
- [ ] 3.2 Define gRPC-only rate limiter struct
- [ ] 3.3 Modify GlobalLimiter to track both HTTP and gRPC requests
- [ ] 3.4 Update HTTP middleware to check HTTP-only limiter
- [ ] 3.5 Write tests for three-tier rate limiting

## 4. gRPC Interceptor

- [ ] 4.1 Create pkg/grpc package
- [ ] 4.2 Define gRPC interceptor struct
- [ ] 4.3 Implement unary interceptor for rate limiting
- [ ] 4.4 Add user ID extraction from gRPC metadata
- [ ] 4.5 Implement gRPC-only rate limiting check
- [ ] 4.6 Implement per-method rate limiting for gRPC
- [ ] 4.7 Return ResourceExhausted status when rate limited
- [ ] 4.8 Write table-driven tests for gRPC interceptor

## 5. Integration and Testing

- [ ] 5.1 Add golangci-lint rules for new code
- [ ] 5.2 Write integration tests for HTTP middleware with config file
- [ ] 5.3 Write integration tests for gRPC interceptor
- [ ] 5.4 Write end-to-end tests for three-tier rate limiting
- [ ] 5.5 Add example configuration files (JSON and YAML)

## 6. Documentation

- [ ] 6.1 Update README.md with gRPC rate limiting usage examples
- [ ] 6.2 Document configuration file format and options
- [ ] 6.3 Add godoc comments to exported functions and types
- [ ] 6.4 Document environment variables and config file precedence
- [ ] 6.5 Add migration guide from environment-only configuration

## 7. Validation

- [ ] 7.1 Run golangci-lint and fix all issues
- [ ] 7.2 Ensure all tests pass with go test ./...
- [ ] 7.3 Verify test coverage meets requirements
- [ ] 7.4 Validate backward compatibility with existing HTTP middleware
