# Change: Add Memcache-based Distributed Rate Limiting

## Why

The current rate limiting implementation is in-memory only, which means each application instance maintains independent rate limit state. In a distributed environment with multiple application instances, this results in:
- Users can bypass rate limits by hitting different instances
- No coordination between instances leads to inconsistent enforcement
- Rate limit state is lost on application restart

To enable true distributed rate limiting across multiple applications, we need to synchronize rate limit state using an external shared storage like Memcache.

## What Changes

- Add Memcache client integration using `github.com/bradfitz/gomemcache`
- Implement counter-based distributed rate limiting using Memcache atomic operations
- Add configurable Memcache connection settings (server addresses, timeout, max idle connections)
- Add configurable failure behavior (fail-open vs fail-closed) when Memcache is unavailable
- Extend configuration to support Memcache settings via environment variables and config file
- Create distributed limiter implementations for Global, PerEndpoint, HTTP-only, and gRPC-only limiters
- Add sliding window counter algorithm for distributed rate limiting
- Maintain backward compatibility with in-memory limiters (Memcache is opt-in)

## Impact

- Affected specs: `rate-limiting` (MODIFIED to add distributed rate limiting)
- Affected code:
  - New `pkg/memcache/` package for Memcache client wrapper
  - New `pkg/middleware/distributed/` package for distributed limiters
  - Extended `internal/config/` for Memcache configuration
  - Modified `pkg/middleware/` to support both in-memory and distributed limiters
- New external dependency: `github.com/bradfitz/gomemcache`
