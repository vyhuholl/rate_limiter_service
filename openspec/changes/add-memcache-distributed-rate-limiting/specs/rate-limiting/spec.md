## ADDED Requirements

### Requirement: Memcache-based Distributed Rate Limiting

The system SHALL provide distributed rate limiting using Memcache as a shared storage backend to synchronize rate limit state across multiple application instances.

#### Scenario: Synchronize rate limits across multiple instances
- **WHEN** multiple application instances are running with Memcache configured
- **THEN** all instances share the same rate limit state for each user
- **AND** requests are rate limited consistently across all instances

#### Scenario: Use counter-based algorithm with atomic operations
- **WHEN** a rate limit check is performed
- **THEN** the system uses atomic increment/decrement operations in Memcache
- **AND** the counter value determines if the request is allowed

#### Scenario: Allow requests when counter is below limit
- **WHEN** the request counter in Memcache is below the configured rate limit
- **THEN** the request is allowed and the counter is incremented

#### Scenario: Reject requests when counter exceeds limit
- **WHEN** the request counter in Memcache equals or exceeds the configured rate limit
- **THEN** the request is rejected with appropriate error response

### Requirement: Memcache Configuration

The system SHALL load Memcache configuration from environment variables or configuration file, with support for connection settings and failure behavior.

#### Scenario: Load Memcache server addresses from environment
- **WHEN** RATE_LIMIT_MEMCACHE_SERVERS environment variable is set
- **THEN** the system connects to the specified Memcache servers

#### Scenario: Load Memcache timeout from environment
- **WHEN** RATE_LIMIT_MEMCACHE_TIMEOUT environment variable is set
- **THEN** the system uses the specified timeout for Memcache operations

#### Scenario: Load Memcache failure mode from environment
- **WHEN** RATE_LIMIT_MEMCACHE_FAILURE_MODE environment variable is set to "allow"
- **THEN** requests are allowed when Memcache is unavailable
- **WHEN** RATE_LIMIT_MEMCACHE_FAILURE_MODE environment variable is set to "deny"
- **THEN** requests are rejected when Memcache is unavailable

#### Scenario: Load Memcache configuration from config file
- **WHEN** the config file includes memcache section
- **THEN** the system uses the specified Memcache configuration

### Requirement: Memcache Key Structure

The system SHALL use a structured key format in Memcache to store rate limit counters for different scopes and users.

#### Scenario: Use structured keys for rate limit counters
- **WHEN** storing rate limit state in Memcache
- **THEN** keys follow the format "rate_limit:{scope}:{user_id}:{identifier}"
- **AND** scope is one of "global", "http", "grpc", or "endpoint"

#### Scenario: Set key expiration
- **WHEN** a rate limit counter is stored or updated in Memcache
- **THEN** the key has an expiration time based on the rate limit window

### Requirement: Distributed Global Rate Limiting

The system SHALL provide distributed global rate limiting that enforces a shared rate limit across all endpoints per user using Memcache.

#### Scenario: Synchronize global limits across instances
- **WHEN** a user makes requests to any endpoint across multiple instances
- **THEN** the total count across all instances is used for global rate limiting

#### Scenario: Reject requests when global limit exceeded
- **WHEN** a user makes requests exceeding the global limit across all instances
- **THEN** requests beyond the limit are rejected regardless of which instance receives them

### Requirement: Distributed Per-Endpoint Rate Limiting

The system SHALL provide distributed per-endpoint rate limiting that enforces a separate rate limit for each unique HTTP endpoint per user using Memcache.

#### Scenario: Synchronize per-endpoint limits across instances
- **WHEN** a user makes requests to the same endpoint across multiple instances
- **THEN** each endpoint has its own independent rate limit shared across instances

#### Scenario: Different endpoints have independent distributed limits
- **WHEN** a user makes requests to different endpoints across multiple instances
- **THEN** each endpoint maintains its own independent rate limit in Memcache

### Requirement: Distributed HTTP-Only Rate Limiting

The system SHALL provide distributed HTTP-only rate limiting that enforces a shared rate limit for HTTP requests per user using Memcache.

#### Scenario: Synchronize HTTP limits across instances
- **WHEN** a user makes HTTP requests across multiple instances
- **THEN** the HTTP-only rate limit is shared across all instances

#### Scenario: gRPC requests unaffected by HTTP limit
- **WHEN** a user makes gRPC requests across multiple instances
- **THEN** gRPC requests are not counted against the HTTP-only rate limit

### Requirement: Distributed gRPC-Only Rate Limiting

The system SHALL provide distributed gRPC-only rate limiting that enforces a shared rate limit for gRPC requests per user using Memcache.

#### Scenario: Synchronize gRPC limits across instances
- **WHEN** a user makes gRPC requests across multiple instances
- **THEN** the gRPC-only rate limit is shared across all instances

#### Scenario: HTTP requests unaffected by gRPC limit
- **WHEN** a user makes HTTP requests across multiple instances
- **THEN** HTTP requests are not counted against the gRPC-only rate limit

### Requirement: Memcache Failure Handling

The system SHALL handle Memcache failures gracefully according to the configured failure mode.

#### Scenario: Fail-open when Memcache unavailable
- **WHEN** Memcache is unavailable and failure mode is set to "allow"
- **THEN** requests are allowed and a warning is logged
- **AND** rate limiting is in degraded mode

#### Scenario: Fail-closed when Memcache unavailable
- **WHEN** Memcache is unavailable and failure mode is set to "deny"
- **THEN** requests are rejected with rate limit error
- **AND** an error is logged

#### Scenario: Handle Memcache timeout
- **WHEN** a Memcache operation times out
- **THEN** the system handles the timeout according to the configured failure mode

### Requirement: Backward Compatibility with In-Memory Limiters

The system SHALL maintain backward compatibility with in-memory rate limiters when Memcache is not configured.

#### Scenario: Use in-memory limiters by default
- **WHEN** Memcache is not configured
- **THEN** the system uses in-memory rate limiters
- **AND** behavior is unchanged from previous versions

#### Scenario: Switch between in-memory and distributed limiters
- **WHEN** configuration changes between in-memory and Memcache
- **THEN** the system uses the appropriate limiter based on current configuration

### Requirement: Memcache Connection Pooling

The system SHALL use connection pooling for Memcache connections to improve performance.

#### Scenario: Reuse connections to Memcache
- **WHEN** multiple rate limit checks are performed
- **THEN** existing Memcache connections are reused from the pool

#### Scenario: Configure maximum idle connections
- **WHEN** RATE_LIMIT_MEMCACHE_MAX_IDLE_CONNECTIONS is set
- **THEN** the connection pool maintains at most the specified number of idle connections

## MODIFIED Requirements

### Requirement: Configuration via Environment Variables

The system SHALL load rate limiting configuration from environment variables, with support for optional configuration file loading and Memcache settings.

#### Scenario: Load Memcache settings from environment
- **WHEN** RATE_LIMIT_MEMCACHE_SERVERS environment variable is set
- **THEN** the system uses Memcache for distributed rate limiting
- **WHEN** RATE_LIMIT_MEMCACHE_TIMEOUT environment variable is set
- **THEN** the system uses the specified timeout for Memcache operations
- **WHEN** RATE_LIMIT_MEMCACHE_FAILURE_MODE environment variable is set
- **THEN** the system uses the specified failure mode when Memcache is unavailable
- **WHEN** RATE_LIMIT_MEMCACHE_MAX_IDLE_CONNECTIONS environment variable is set
- **THEN** the system uses the specified maximum number of idle connections

#### Scenario: Load config file path from environment
- **WHEN** RATE_LIMIT_CONFIG_PATH environment variable is set
- **THEN** the system loads configuration from the specified file including Memcache settings

#### Scenario: Use environment variables when config file not set
- **WHEN** RATE_LIMIT_CONFIG_PATH is not set
- **THEN** the system uses RATE_LIMIT_USER_HEADER, RATE_LIMIT_PER_ENDPOINT, RATE_LIMIT_GLOBAL, and Memcache environment variables for configuration

#### Scenario: Config file overrides environment variables
- **WHEN** both RATE_LIMIT_CONFIG_PATH and environment variables are set
- **THEN** the configuration file values take precedence over environment variables

### Requirement: Configuration File Loading

The system SHALL load rate limiting configuration from a JSON or YAML file specified via environment variable, including Memcache settings.

#### Scenario: Load Memcache config from JSON file
- **WHEN** RATE_LIMIT_CONFIG_PATH environment variable points to a JSON file with memcache section
- **THEN** the system loads Memcache configuration from that JSON file

#### Scenario: Load Memcache config from YAML file
- **WHEN** RATE_LIMIT_CONFIG_PATH environment variable points to a YAML file with memcache section
- **THEN** the system loads Memcache configuration from that YAML file

#### Scenario: Validate configuration on load
- **WHEN** the configuration file is loaded including Memcache settings
- **THEN** the system validates all values and returns an error if invalid

#### Scenario: Use environment variables as fallback
- **WHEN** RATE_LIMIT_CONFIG_PATH is not set
- **THEN** the system uses environment variables for configuration including Memcache settings
