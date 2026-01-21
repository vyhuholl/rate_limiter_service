package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds the rate limiter configuration
type Config struct {
	// UserHeader is the HTTP header name used to identify users
	UserHeader string
	// PerEndpointRate is the rate limit per endpoint per user (requests per second)
	PerEndpointRate int
	// GlobalRate is the global rate limit per user across all endpoints (requests per second)
	GlobalRate int
	// GlobalBurstSize is the maximum burst size for global token bucket
	GlobalBurstSize int
	// PerEndpointBurstSize is the maximum burst size for per-endpoint token buckets
	PerEndpointBurstSize int
}

// DefaultConfig returns the default configuration values
func DefaultConfig() Config {
	return Config{
		UserHeader:            "X-User-ID",
		PerEndpointRate:       10,
		GlobalRate:            100,
		GlobalBurstSize:       10,
		PerEndpointBurstSize:  10,
	}
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (Config, error) {
	config := DefaultConfig()

	if userHeader := os.Getenv("RATE_LIMIT_USER_HEADER"); userHeader != "" {
		config.UserHeader = userHeader
	}

	if perEndpointRate := os.Getenv("RATE_LIMIT_PER_ENDPOINT"); perEndpointRate != "" {
		rate, err := strconv.Atoi(perEndpointRate)
		if err != nil {
			return config, fmt.Errorf("invalid RATE_LIMIT_PER_ENDPOINT value %q: %w", perEndpointRate, err)
		}
		if rate <= 0 {
			return config, fmt.Errorf("RATE_LIMIT_PER_ENDPOINT must be positive, got %d", rate)
		}
		config.PerEndpointRate = rate
	}

	if globalRate := os.Getenv("RATE_LIMIT_GLOBAL"); globalRate != "" {
		rate, err := strconv.Atoi(globalRate)
		if err != nil {
			return config, fmt.Errorf("invalid RATE_LIMIT_GLOBAL value %q: %w", globalRate, err)
		}
		if rate <= 0 {
			return config, fmt.Errorf("RATE_LIMIT_GLOBAL must be positive, got %d", rate)
		}
		config.GlobalRate = rate
	}

	if globalBurstSize := os.Getenv("RATE_LIMIT_GLOBAL_BURST_SIZE"); globalBurstSize != "" {
		size, err := strconv.Atoi(globalBurstSize)
		if err != nil {
			return config, fmt.Errorf("invalid RATE_LIMIT_GLOBAL_BURST_SIZE value %q: %w", globalBurstSize, err)
		}
		if size <= 0 {
			return config, fmt.Errorf("RATE_LIMIT_GLOBAL_BURST_SIZE must be positive, got %d", size)
		}
		config.GlobalBurstSize = size
	}

	if burstSize := os.Getenv("RATE_LIMIT_BURST_SIZE"); burstSize != "" {
		size, err := strconv.Atoi(burstSize)
		if err != nil {
			return config, fmt.Errorf("invalid RATE_LIMIT_BURST_SIZE value %q: %w", burstSize, err)
		}
		if size <= 0 {
			return config, fmt.Errorf("RATE_LIMIT_BURST_SIZE must be positive, got %d", size)
		}
		config.GlobalBurstSize = size
		config.PerEndpointBurstSize = size
	}

	if globalBurstSize := os.Getenv("RATE_LIMIT_GLOBAL_BURST_SIZE"); globalBurstSize != "" {
		size, err := strconv.Atoi(globalBurstSize)
		if err != nil {
			return config, fmt.Errorf("invalid RATE_LIMIT_GLOBAL_BURST_SIZE value %q: %w", globalBurstSize, err)
		}
		if size <= 0 {
			return config, fmt.Errorf("RATE_LIMIT_GLOBAL_BURST_SIZE must be positive, got %d", size)
		}
		config.GlobalBurstSize = size
	}

	if perEndpointBurstSize := os.Getenv("RATE_LIMIT_PER_ENDPOINT_BURST_SIZE"); perEndpointBurstSize != "" {
		size, err := strconv.Atoi(perEndpointBurstSize)
		if err != nil {
			return config, fmt.Errorf("invalid RATE_LIMIT_PER_ENDPOINT_BURST_SIZE value %q: %w", perEndpointBurstSize, err)
		}
		if size <= 0 {
			return config, fmt.Errorf("RATE_LIMIT_PER_ENDPOINT_BURST_SIZE must be positive, got %d", size)
		}
		config.PerEndpointBurstSize = size
	}

	// Validate that per-endpoint rate doesn't exceed global rate
	if config.PerEndpointRate > config.GlobalRate {
		return config, fmt.Errorf("per-endpoint rate (%d) cannot exceed global rate (%d)",
			config.PerEndpointRate, config.GlobalRate)
	}

	return config, nil
}

// GetRefillInterval returns the time interval for refilling tokens based on the rate
func (c Config) GetRefillInterval(rate int) time.Duration {
	// Calculate how often to add one token (in nanoseconds)
	// rate tokens per second = rate tokens / 1000000000 nanoseconds
	refillInterval := int64(time.Second) / int64(rate)
	return time.Duration(refillInterval)
}