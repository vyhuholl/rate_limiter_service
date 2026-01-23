## Context

This change introduces rate limiting middleware to protect HTTP endpoints from abuse. The middleware must be reusable across any Go HTTP server and support two rate limiting strategies: per-endpoint and global. Users are identified via a configurable HTTP header.

### Constraints

- In-memory only (no distributed storage)
- Single-process enforcement
- Token bucket algorithm
- Go 1.22+ language features

## Goals / Non-Goals

### Goals

- Provide reusable HTTP middleware for rate limiting
- Implement per-endpoint rate limiting (separate limit per endpoint)
- Implement global rate limiting (shared limit across all endpoints)
- Support configurable rate limits via environment variables
- Return HTTP 429 with appropriate headers when rate limited

### Non-Goals

- Distributed rate limiting (single process only)
- Persistent storage of rate limit state
- Authentication or authorization (user ID comes from header)
- Metrics or monitoring integration
- Dynamic configuration changes (requires restart)

## Decisions

### Decision 1: Token Bucket Algorithm

**What**: Use token bucket algorithm for rate limiting.

**Why**: Token bucket provides smooth rate limiting with burst capacity, is well-understood, and works well for HTTP request limiting. It allows short bursts while maintaining average rate.

**Alternatives considered**:
- Sliding window: More accurate but more complex and memory-intensive
- Fixed window: Simpler but can allow bursts at window boundaries
- Leaky bucket: Good for traffic shaping but less suitable for request limiting

### Decision 2: Separate Limiters for Per-Endpoint and Global

**What**: Implement two independent limiter types that can be composed.

**Why**: Separation allows flexibility - users can use either or both limiters independently. Composition via middleware chaining provides clean separation of concerns.

**Alternatives considered**:
- Single limiter with both modes: More complex, harder to test independently
- Combined limiter: Would couple the two concerns, reducing flexibility

### Decision 3: User Identification via Configurable Header

**What**: Identify users by reading a configurable HTTP header (e.g., X-User-ID).

**Why**: Keeps the middleware simple and reusable. Different services can use different headers for user identification without code changes.

**Alternatives considered**:
- Hardcoded header name: Less flexible, requires code changes for different services
- IP-based identification: Doesn't work well with NAT, proxies, or shared IPs
- Cookie-based: Too specific to web applications

### Decision 4: In-Memory Storage with sync.Map

**What**: Use sync.Map for storing per-user/per-endpoint rate limit buckets.

**Why**: sync.Map provides thread-safe access without explicit locking, optimized for read-heavy workloads typical of rate limiting scenarios.

**Alternatives considered**:
- map with sync.RWMutex: More boilerplate, similar performance
- map with sync.Mutex: Simpler but may have contention under high load
- Custom sharding: More complex, not needed for initial implementation

### Decision 5: Middleware Interface

**What**: Define a Middleware function type that wraps http.Handler.

**Why**: Follows standard Go middleware patterns, works with any http.Handler, and allows easy composition with other middleware.

**Alternatives considered**:
- Custom middleware interface: Would break compatibility with standard library
- Middleware struct with ServeHTTP method: Less flexible, harder to compose

## Risks / Trade-offs

### Risk 1: Memory Usage Under High Cardinality

**Risk**: With many unique users and endpoints, memory usage could grow unbounded.

**Mitigation**: Add periodic cleanup of idle buckets in a future iteration. For now, document this limitation.

### Risk 2: No Distributed Coordination

**Risk**: Multiple service instances will have independent rate limits, allowing higher aggregate throughput.

**Mitigation**: Document this limitation clearly. For distributed scenarios, recommend using a Redis-based implementation (future work).

### Trade-off: Simplicity vs Features

**Trade-off**: Initial implementation prioritizes simplicity and clarity over advanced features like dynamic configuration, metrics, or distributed support.

**Rationale**: Simple implementation is easier to understand, test, and maintain. Features can be added incrementally as needed.

## Migration Plan

No migration needed - this is new functionality added to an empty project.

## Open Questions

- Should we add a cleanup mechanism for idle buckets to prevent unbounded memory growth?
- Should we expose metrics for rate limit hits and rejections?
- Should we support custom key generation functions for more complex user identification?