## ADDED Requirements

### Requirement: Token Bucket Rate Limiter

The system SHALL implement a token bucket rate limiter that allows requests up to a configured rate and burst size.

#### Scenario: Consume token when available
- **WHEN** a token is available in the bucket
- **THEN** the request is allowed and one token is consumed

#### Scenario: Reject request when bucket empty
- **WHEN** no tokens are available in the bucket
- **THEN** the request is rejected

#### Scenario: Refill tokens over time
- **WHEN** time passes
- **THEN** tokens are added to the bucket at the configured rate up to the burst size

#### Scenario: Allow burst up to burst size
- **WHEN** the bucket is full
- **THEN** up to burst size tokens can be consumed immediately

### Requirement: Per-Endpoint Rate Limiting

The system SHALL provide per-endpoint rate limiting that enforces a separate rate limit for each unique HTTP endpoint (method + path combination) per user.

#### Scenario: Limit requests to single endpoint
- **WHEN** a user makes requests to a single endpoint exceeding the per-endpoint limit
- **THEN** requests beyond the limit are rejected

#### Scenario: Allow requests to different endpoints
- **WHEN** a user makes requests to different endpoints
- **THEN** each endpoint has its own independent rate limit

#### Scenario: Different users have independent limits
- **WHEN** multiple users access the same endpoint
- **THEN** each user has their own independent rate limit

#### Scenario: Identify user from configured header
- **WHEN** a request includes the configured user identification header
- **THEN** the header value is used as the user identifier for rate limiting

### Requirement: Global Rate Limiting

The system SHALL provide global rate limiting that enforces a shared rate limit across all endpoints per user.

#### Scenario: Limit requests across all endpoints
- **WHEN** a user makes requests to any endpoint exceeding the global limit
- **THEN** requests beyond the limit are rejected regardless of endpoint

#### Scenario: Different users have independent global limits
- **WHEN** multiple users make requests
- **THEN** each user has their own independent global rate limit

#### Scenario: Global limit applies to all endpoints
- **WHEN** a user makes requests to multiple endpoints
- **THEN** the total count across all endpoints is used for global rate limiting

### Requirement: HTTP Middleware Integration

The system SHALL provide HTTP middleware that wraps existing handlers to apply rate limiting.

#### Scenario: Middleware wraps handler
- **WHEN** middleware is applied to an HTTP handler
- **THEN** the middleware intercepts all requests to that handler

#### Scenario: Return HTTP 429 when rate limited
- **WHEN** a request exceeds the rate limit
- **THEN** the middleware returns HTTP 429 Too Many Requests status

#### Scenario: Include rate limit headers in response
- **WHEN** a request is rate limited
- **THEN** the response includes X-RateLimit-Limit, X-RateLimit-Remaining, and Retry-After headers

#### Scenario: Pass through allowed requests
- **WHEN** a request is within rate limits
- **THEN** the request is passed to the wrapped handler

### Requirement: Configuration via Environment Variables

The system SHALL load rate limiting configuration from environment variables.

#### Scenario: Load user header from environment
- **WHEN** RATE_LIMIT_USER_HEADER environment variable is set
- **THEN** the middleware uses that header name for user identification

#### Scenario: Load per-endpoint limit from environment
- **WHEN** RATE_LIMIT_PER_ENDPOINT environment variable is set
- **THEN** the middleware uses that value as requests per second per endpoint

#### Scenario: Load global limit from environment
- **WHEN** RATE_LIMIT_GLOBAL environment variable is set
- **THEN** the middleware uses that value as requests per second globally

#### Scenario: Load burst size from environment
- **WHEN** RATE_LIMIT_BURST_SIZE environment variable is set
- **THEN** the middleware uses that value as the maximum burst size

#### Scenario: Use default values when not configured
- **WHEN** environment variables are not set
- **THEN** the middleware uses default values (10 per-endpoint, 100 global, 10 burst)

### Requirement: Thread-Safe Operations

The system SHALL ensure all rate limiting operations are thread-safe for concurrent requests.

#### Scenario: Handle concurrent requests to same endpoint
- **WHEN** multiple concurrent requests are made to the same endpoint by the same user
- **THEN** rate limiting is enforced correctly without race conditions

#### Scenario: Handle concurrent requests by different users
- **WHEN** multiple concurrent requests are made by different users
- **THEN** each user's rate limit is enforced independently without interference
