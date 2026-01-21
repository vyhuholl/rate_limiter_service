# Project Context

## Purpose

This project provides a Go-based rate limiter middleware service that can be integrated with any HTTP server. The service implements two rate limiting strategies:

1. **Per-endpoint rate limiting** - Limits the number of requests per second to a single endpoint
2. **Global rate limiting** - Limits the number of requests per second across all endpoints

Users are identified via a configurable HTTP header specified through environment variables. The middleware uses the token bucket algorithm for rate limiting and operates entirely in-memory.

## Tech Stack

- **Language**: Go 1.22+
- **Rate Limiting Algorithm**: Token Bucket
- **Storage**: In-memory only
- **Configuration**: Environment variables
- **Testing**: Table-driven tests
- **Linting**: golangci-lint
- **Formatting**: gofmt

## Project Conventions

### Code Style

- **Formatting**: Always use `gofmt` for code formatting
- **Linting**: Run `golangci-lint` before committing changes
- **Naming Conventions**:
  - Exported names (public) start with capital letters
  - Unexported names (private) use lowercase
  - Use camelCase for variable and function names
  - Keep names concise but clear - prefer clarity over brevity
- **Package Structure**: Follow standard Go project layout
- **Error Handling**: Use explicit error returns, avoid panic in production code

### Architecture Patterns

- **Middleware Pattern**: The rate limiter is implemented as HTTP middleware that wraps handlers
- **Interface-based Design**: Define interfaces for rate limiters to enable testing and future extensibility
- **Separation of Concerns**:
  - Rate limiting logic separated from HTTP handling
  - Configuration loading separated from core logic
- **Clean Architecture**: Keep business logic independent of external dependencies

### Testing Strategy

- **Test Style**: Table-driven tests for all functions with multiple scenarios
- **Coverage**: Aim for high test coverage, especially for rate limiting logic
- **Test Organization**: Place test files in the same package as the code being tested
- **Test Naming**: Use descriptive test names that explain what is being tested
- **Mocking**: Use interfaces to mock dependencies when needed
- **Integration Tests**: Test the middleware integration with HTTP handlers

### Git Workflow

- **Branching**: Simple workflow with `master` as the primary branch
- **Commit Messages**: Use conventional commits format:
  - `feat:` for new features
  - `fix:` for bug fixes
  - `refactor:` for code refactoring
  - `test:` for test changes
  - `docs:` for documentation updates
  - `chore:` for maintenance tasks
- **Pull Requests**: Required before merging to master
- **Code Review**: All changes must be reviewed before merging

## Domain Context

### Rate Limiting Concepts

- **Token Bucket Algorithm**: A rate limiting algorithm where tokens are added to a bucket at a fixed rate. Each request consumes a token. If the bucket is empty, requests are rejected.
- **Per-endpoint Limiting**: Each unique HTTP endpoint (method + path) has its own rate limit
- **Global Limiting**: All requests share a single rate limit regardless of endpoint
- **User Identification**: Users are identified by a configurable HTTP header (e.g., `X-User-ID`, `X-API-Key`)
- **Rate Limit Response**: When rate limited, return HTTP 429 Too Many Requests with appropriate headers

### Middleware Integration

The rate limiter is designed as middleware that can be added to any Go HTTP server:

```go
http.Handle("/", rateLimiterMiddleware(handler))
```

### Configuration Parameters

- `RATE_LIMIT_USER_HEADER`: HTTP header name used to identify users (required)
- `RATE_LIMIT_PER_ENDPOINT`: Requests per second per endpoint (default: 10)
- `RATE_LIMIT_GLOBAL`: Requests per second globally (default: 100)
- `RATE_LIMIT_BURST_SIZE`: Maximum burst size for token bucket (default: 10)

## Important Constraints

- **In-Memory Only**: Rate limits are stored in memory and reset on restart. Not suitable for distributed systems.
- **Single Process**: Rate limits are enforced per process instance. Multiple instances will have independent rate limits.
- **No Persistence**: Rate limit state is not persisted and is lost on application restart.
- **Header Required**: Requests must include the configured user identification header to be rate limited.
- **No Authentication**: The middleware does not authenticate users; it only identifies them via the configured header.
- **Go Version**: Requires Go 1.22+ for language features used in the implementation.

## External Dependencies

- **Standard Library**: Primary reliance on Go standard library (`net/http`, `time`, `sync`, etc.)
- **golangci-lint**: Development tool for linting (not a runtime dependency)
- **No External Services**: No external services or databases required for operation

### Potential Future Dependencies

- **Redis**: For distributed rate limiting (not currently implemented)
- **Prometheus**: For metrics and monitoring (not currently implemented)
