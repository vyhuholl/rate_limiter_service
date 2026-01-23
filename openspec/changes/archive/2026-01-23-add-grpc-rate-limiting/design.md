## Context

This change extends the existing HTTP rate limiting middleware to support gRPC services. The current implementation provides per-endpoint and global rate limiting for HTTP requests. The new implementation must:

1. Support gRPC interceptors for rate limiting
2. Provide separate rate limiters for HTTP-only, gRPC-only, and global (all) requests
3. Enable per-method rate limits for both HTTP and gRPC
4. Load method-specific rate limits from a JSON/YAML configuration file
5. Use default rates for methods without explicit configuration

### Constraints

- In-memory only (no distributed storage)
- Single-process enforcement
- Token bucket algorithm
- Go 1.22+ language features
- Must maintain backward compatibility with existing HTTP middleware

## Goals / Non-Goals

### Goals

- Provide gRPC interceptor for rate limiting
- Implement three-tier rate limiting (global, HTTP-only, gRPC-only)
- Support per-method rate limits for both HTTP and gRPC
- Load method-specific rate limits from JSON/YAML config file
- Use default rate for methods without explicit configuration
- Support user identification via gRPC metadata
- Return appropriate gRPC status code when rate limited

### Non-Goals

- Distributed rate limiting (single process only)
- Persistent storage of rate limit state
- Dynamic configuration reloading (requires restart)
- Metrics or monitoring integration
- Protocol-independent rate limiting abstraction

## Decisions

### Decision 1: Three-Tier Rate Limiting Architecture

**What**: Implement three independent rate limiters:
1. Global limiter - applies to all requests (HTTP + gRPC)
2. HTTP-only limiter - applies only to HTTP requests
3. gRPC-only limiter - applies only to gRPC requests

**Why**: This provides fine-grained control over traffic. Users can set different limits for each protocol type while maintaining an overall global limit. The global limiter acts as a safeguard, while protocol-specific limiters allow for different traffic patterns.

**Alternatives considered**:
- Single limiter with protocol tags: More complex, harder to reason about
- Only per-method limits: Would lose the ability to set protocol-wide limits

### Decision 2: Per-Method Rate Limits with Fallback

**What**: Each HTTP method (e.g., GET /api/users) and gRPC method (e.g., /UserService/GetUser) can have its own rate limit. Methods without explicit configuration use a default rate.

**Why**: Different methods have different costs and usage patterns. Per-method limits allow for granular control. Fallback to default ensures all methods are rate limited even without configuration.

**Alternatives considered**:
- Reject unconfigured methods: Would break existing services without configuration
- Single rate for all methods: Too restrictive for diverse workloads

### Decision 3: JSON/YAML Configuration File

**What**: Load method-specific rate limits from a JSON or YAML file specified via environment variable.

**Why**: Environment variables become unwieldy for many method-specific configurations. A config file is more maintainable and supports complex nested structures.

**Alternatives considered**:
- Environment variables with patterns: Too many variables, hard to manage
- Database lookup: Adds external dependency, complexity

### Decision 4: gRPC Interceptor Pattern

**What**: Implement gRPC rate limiting as a unary interceptor that wraps each RPC call.

**Why**: Interceptors are the standard way to add cross-cutting concerns in gRPC. They work similarly to HTTP middleware and allow for clean composition.

**Alternatives considered**:
- Custom gRPC server wrapper: Less flexible, breaks gRPC patterns
- Client-side limiting: Not enforceable, easy to bypass

### Decision 5: User Identification via Metadata

**What**: Extract user ID from gRPC metadata (similar to HTTP headers) using a configurable key.

**Why**: gRPC metadata is the equivalent of HTTP headers. Using a configurable key maintains flexibility across different services.

**Alternatives considered**:
- Hardcoded metadata key: Less flexible
- Peer IP: Doesn't work well with proxies/NAT

### Decision 6: Reuse Existing Token Bucket Implementation

**What**: Extend the existing token bucket implementation rather than creating a new one.

**Why**: The token bucket algorithm is protocol-agnostic. Reusing existing code reduces duplication and maintains consistency.

**Alternatives considered**:
- Separate gRPC token bucket: Unnecessary duplication
- Different algorithm: No clear benefit, adds complexity

## Architecture

### Rate Limiting Flow

```mermaid
sequenceDiagram
    participant Client
    participant HTTPMiddleware
    participant GRPCInterceptor
    participant GlobalLimiter
    participant HTTPLimiter
    participantGRPCLimiter
    participant PerMethodLimiter
    participant Handler

    Client->>HTTPMiddleware: HTTP Request
    HTTPMiddleware->>GlobalLimiter: Check(userID)
    alt Global limit exceeded
        HTTPMiddleware-->>Client: 429 Too Many Requests
    else Global limit OK
        HTTPMiddleware->>HTTPLimiter: Check(userID)
        alt HTTP limit exceeded
            HTTPMiddleware-->>Client: 429 Too Many Requests
        else HTTP limit OK
            HTTPMiddleware->>PerMethodLimiter: Check(userID, method)
            alt Method limit exceeded
                HTTPMiddleware-->>Client: 429 Too Many Requests
            else All limits OK
                HTTPMiddleware->>Handler: Forward request
            end
        end
    end

    Client->>GRPCInterceptor: gRPC Request
    GRPCInterceptor->>GlobalLimiter: Check(userID)
    alt Global limit exceeded
        GRPCInterceptor-->>Client: ResourceExhausted
    else Global limit OK
        GRPCInterceptor->>GRPCLimiter: Check(userID)
        alt gRPC limit exceeded
            GRPCInterceptor-->>Client: ResourceExhausted
        else gRPC limit OK
            GRPCInterceptor->>PerMethodLimiter: Check(userID, method)
            alt Method limit exceeded
                GRPCInterceptor-->>Client: ResourceExhausted
            else All limits OK
                GRPCInterceptor->>Handler: Forward request
            end
        end
    end
```

### Configuration Structure

```yaml
# Example config.yaml
rate_limits:
  global:
    rate: 100  # requests per second
    burst: 10

  http:
    rate: 50   # requests per second
    burst: 5
    default_method_rate: 10
    methods:
      "GET /api/users": 20
      "POST /api/users": 5
      "DELETE /api/users": 2

  grpc:
    rate: 50   # requests per second
    burst: 5
    default_method_rate: 10
    methods:
      "/UserService/GetUser": 20
      "/UserService/CreateUser": 5
      "/UserService/DeleteUser": 2

user_identification:
  http_header: "X-User-ID"
  grpc_metadata_key: "user-id"
```

## Risks / Trade-offs

### Risk 1: Increased Memory Usage

**Risk**: Per-method rate limits require storing buckets for each (user, method) combination, potentially increasing memory usage.

**Mitigation**: Document this limitation. Consider adding periodic cleanup of idle buckets in future iterations.

### Risk 2: Configuration Complexity

**Risk**: The config file adds complexity and a new failure mode (invalid configuration).

**Mitigation**: Provide clear validation errors and example configurations. Document all options thoroughly.

### Risk 3: Backward Compatibility

**Risk**: Changes to configuration loading may break existing deployments using environment variables.

**Mitigation**: Maintain support for existing environment variables. The config file is optional; environment variables can still be used for basic configuration.

### Trade-off: Flexibility vs Simplicity

**Trade-off**: Per-method limits add significant flexibility but also complexity in configuration and implementation.

**Rationale**: The user explicitly requested per-method limits. The complexity is justified by the use case.

## Migration Plan

1. Add new configuration loading code alongside existing environment variable loading
2. Maintain backward compatibility - environment variables still work
3. Add gRPC interceptor as a new feature (HTTP middleware unchanged)
4. Users can gradually adopt the new features

## Open Questions

- Should the config file support hot reloading (without restart)?
- Should we add Prometheus metrics for rate limit hits?
- Should we support streaming gRPC interceptors in addition to unary?
