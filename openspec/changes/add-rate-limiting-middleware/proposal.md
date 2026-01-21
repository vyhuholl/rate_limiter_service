# Change: Add Rate Limiting Middleware

## Why

The project needs a reusable rate limiting middleware for HTTP servers that can protect endpoints from abuse. This middleware will provide two rate limiting strategies: per-endpoint limits (to protect individual endpoints) and global limits (to protect the overall service).

## What Changes

- Add rate limiting middleware using token bucket algorithm
- Implement per-endpoint rate limiting (requests per second to a single endpoint)
- Implement global rate limiting (requests per second across all endpoints)
- Configure user identification via environment variables
- Return HTTP 429 Too Many Requests when rate limited
- Support configurable rate limits and burst size

## Impact

- Affected specs: New capability `rate-limiting`
- Affected code: New middleware package, configuration loading
