## ADDED Requirements

### Requirement: gRPC Rate Limiting Interceptor

The system SHALL provide a gRPC interceptor that applies rate limiting to gRPC requests.

#### Scenario: Interceptor wraps gRPC handler
- **WHEN** a gRPC interceptor is applied to a gRPC service
- **THEN** the interceptor intercepts all RPC calls to that service

#### Scenario: Return ResourceExhausted when rate limited
- **WHEN** a gRPC request exceeds the rate limit
- **THEN** the interceptor returns gRPC status code ResourceExhausted

#### Scenario: Pass through allowed requests
- **WHEN** a gRPC request is within rate limits
- **THEN** the request is passed to the wrapped handler

### Requirement: Three-Tier Rate Limiting

The system SHALL provide three independent rate limiters: global (all requests), HTTP-only, and gRPC-only.

#### Scenario: Global limit applies to all requests
- **WHEN** a user makes HTTP or gRPC requests exceeding the global limit
- **THEN** requests beyond the limit are rejected regardless of protocol

#### Scenario: HTTP-only limit applies to HTTP requests only
- **WHEN** a user makes HTTP requests exceeding the HTTP-only limit
- **THEN** HTTP requests beyond the limit are rejected, but gRPC requests are unaffected

#### Scenario: gRPC-only limit applies to gRPC requests only
- **WHEN** a user makes gRPC requests exceeding the gRPC-only limit
- **THEN** gRPC requests beyond the limit are rejected, but HTTP requests are unaffected

#### Scenario: All three limits are enforced
- **WHEN** a user makes requests
- **THEN** the request must pass all applicable limits (global, protocol-specific, and method-specific) to be allowed

### Requirement: Per-Method Rate Limits

The system SHALL enforce per-method rate limits for both HTTP and gRPC requests, where each method can have its own rate limit.

#### Scenario: HTTP method has specific rate limit
- **WHEN** an HTTP method (e.g., GET /api/users) has a configured rate limit
- **THEN** requests to that method are limited according to the configured rate

#### Scenario: gRPC method has specific rate limit
- **WHEN** a gRPC method (e.g., /UserService/GetUser) has a configured rate limit
- **THEN** requests to that method are limited according to the configured rate

#### Scenario: Different methods have different limits
- **WHEN** multiple HTTP or gRPC methods have configured rate limits
- **THEN** each method is limited independently according to its configured rate

#### Scenario: Default rate for unconfigured methods
- **WHEN** a method does not have an explicit rate limit configured
- **THEN** the method uses the default rate limit for its protocol type

### Requirement: Configuration File Loading

The system SHALL load rate limiting configuration from a JSON or YAML file specified via environment variable.

#### Scenario: Load config from JSON file
- **WHEN** RATE_LIMIT_CONFIG_PATH environment variable points to a JSON file
- **THEN** the system loads configuration from that JSON file

#### Scenario: Load config from YAML file
- **WHEN** RATE_LIMIT_CONFIG_PATH environment variable points to a YAML file
- **THEN** the system loads configuration from that YAML file

#### Scenario: Validate configuration on load
- **WHEN** the configuration file is loaded
- **THEN** the system validates all values and returns an error if invalid

#### Scenario: Use environment variables as fallback
- **WHEN** RATE_LIMIT_CONFIG_PATH is not set
- **THEN** the system uses environment variables for configuration

### Requirement: User Identification via gRPC Metadata

The system SHALL identify users for gRPC rate limiting using a configurable metadata key.

#### Scenario: Extract user ID from metadata
- **WHEN** a gRPC request includes the configured user identification metadata
- **THEN** the metadata value is used as the user identifier for rate limiting

#### Scenario: Default to anonymous user
- **WHEN** a gRPC request does not include the configured metadata
- **THEN** the request is treated as from an anonymous user

#### Scenario: Configurable metadata key
- **WHEN** the configuration specifies a metadata key for user identification
- **THEN** the system uses that key to extract the user ID

## MODIFIED Requirements

### Requirement: Per-Endpoint Rate Limiting

The system SHALL provide per-endpoint rate limiting that enforces a separate rate limit for each unique HTTP endpoint (method + path combination) per user, where each endpoint can have its own configured rate limit.

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

#### Scenario: Use configured rate for specific endpoint
- **WHEN** an endpoint has a specific rate configured in the config file
- **THEN** that endpoint uses the configured rate instead of the default

#### Scenario: Use default rate for unconfigured endpoint
- **WHEN** an endpoint does not have a specific rate configured
- **THEN** that endpoint uses the default HTTP method rate

### Requirement: Configuration via Environment Variables

The system SHALL load rate limiting configuration from environment variables, with support for optional configuration file loading.

#### Scenario: Load user header from environment
- **WHEN** RATE_LIMIT_USER_HEADER environment variable is set
- **THEN** the middleware uses that header name for user identification

#### Scenario: Load config file path from environment
- **WHEN** RATE_LIMIT_CONFIG_PATH environment variable is set
- **THEN** the system loads configuration from the specified file

#### Scenario: Use environment variables when config file not set
- **WHEN** RATE_LIMIT_CONFIG_PATH is not set
- **THEN** the system uses RATE_LIMIT_USER_HEADER, RATE_LIMIT_PER_ENDPOINT, and RATE_LIMIT_GLOBAL for configuration

#### Scenario: Config file overrides environment variables
- **WHEN** both RATE_LIMIT_CONFIG_PATH and environment variables are set
- **THEN** the configuration file values take precedence over environment variables
