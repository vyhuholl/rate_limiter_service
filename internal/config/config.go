package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the rate limiter configuration
type Config struct {
	// UserHeader is the HTTP header name used to identify users
	UserHeader string
	// GrpcMetadataKey is the gRPC metadata key used to identify users
	GrpcMetadataKey string
	// PerEndpointRate is the rate limit per endpoint per user (requests per second)
	PerEndpointRate int
	// GlobalRate is the global rate limit per user across all endpoints (requests per second)
	GlobalRate int
	// GlobalBurstSize is the maximum burst size for global token bucket
	GlobalBurstSize int
	// PerEndpointBurstSize is the maximum burst size for per-endpoint token buckets
	PerEndpointBurstSize int
	// HTTPRate is the rate limit for HTTP requests only (requests per second)
	HTTPRate int
	// HTTPBurstSize is the maximum burst size for HTTP-only token bucket
	HTTPBurstSize int
	// GRPCRate is the rate limit for gRPC requests only (requests per second)
	GRPCRate int
	// GRPCBurstSize is the maximum burst size for gRPC-only token bucket
	GRPCBurstSize int
	// HTTPMethods is a map of HTTP method+path to rate limit
	HTTPMethods map[string]int
	// HTTPDefaultMethodRate is the default rate for HTTP methods not explicitly configured
	HTTPDefaultMethodRate int
	// GRPCMethods is a map of gRPC method to rate limit
	GRPCMethods map[string]int
	// GRPCDefaultMethodRate is the default rate for gRPC methods not explicitly configured
	GRPCDefaultMethodRate int
}

// FileConfig represents the structure of the configuration file
type FileConfig struct {
	RateLimits struct {
		Global struct {
			Rate  int `json:"rate" yaml:"rate"`
			Burst int `json:"burst" yaml:"burst"`
		} `json:"global" yaml:"global"`
		HTTP struct {
			Rate              int            `json:"rate" yaml:"rate"`
			Burst             int            `json:"burst" yaml:"burst"`
			DefaultMethodRate int            `json:"default_method_rate" yaml:"default_method_rate"`
			Methods           map[string]int `json:"methods" yaml:"methods"`
		} `json:"http" yaml:"http"`
		GRPC struct {
			Rate              int            `json:"rate" yaml:"rate"`
			Burst             int            `json:"burst" yaml:"burst"`
			DefaultMethodRate int            `json:"default_method_rate" yaml:"default_method_rate"`
			Methods           map[string]int `json:"methods" yaml:"methods"`
		} `json:"grpc" yaml:"grpc"`
	} `json:"rate_limits" yaml:"rate_limits"`
	UserIdentification struct {
		HTTPHeader     string `json:"http_header" yaml:"http_header"`
		GRPCMetadataKey string `json:"grpc_metadata_key" yaml:"grpc_metadata_key"`
	} `json:"user_identification" yaml:"user_identification"`
}

// DefaultConfig returns the default configuration values
func DefaultConfig() Config {
	return Config{
		UserHeader:            "X-User-ID",
		GrpcMetadataKey:       "user-id",
		PerEndpointRate:       10,
		GlobalRate:            100,
		GlobalBurstSize:       10,
		PerEndpointBurstSize:  10,
		HTTPRate:              50,
		HTTPBurstSize:         5,
		GRPCRate:              50,
		GRPCBurstSize:         5,
		HTTPDefaultMethodRate: 10,
		GRPCDefaultMethodRate: 10,
	}
}

// loadEnvInt loads a positive integer value from environment variable
func loadEnvInt(envVar string, defaultValue int) (int, error) {
	if value := os.Getenv(envVar); value != "" {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue, fmt.Errorf("invalid %s value %q: %w", envVar, value, err)
		}
		if intValue <= 0 {
			return defaultValue, fmt.Errorf("%s must be positive, got %d", envVar, intValue)
		}
		return intValue, nil
	}
	return defaultValue, nil
}

// loadBasicEnvConfig loads basic rate limiting configuration from environment variables
func loadBasicEnvConfig(config *Config) error {
	var err error

	config.UserHeader = loadEnvString("RATE_LIMIT_USER_HEADER", config.UserHeader)

	if config.PerEndpointRate, err = loadEnvInt("RATE_LIMIT_PER_ENDPOINT", config.PerEndpointRate); err != nil {
		return err
	}

	if config.GlobalRate, err = loadEnvInt("RATE_LIMIT_GLOBAL", config.GlobalRate); err != nil {
		return err
	}

	if config.GlobalBurstSize, err = loadEnvInt("RATE_LIMIT_GLOBAL_BURST_SIZE", config.GlobalBurstSize); err != nil {
		return err
	}

	if burstSize, err := loadEnvInt("RATE_LIMIT_BURST_SIZE", 0); err != nil {
		return err
	} else if burstSize > 0 {
		config.GlobalBurstSize = burstSize
		config.PerEndpointBurstSize = burstSize
	}

	if config.PerEndpointBurstSize, err = loadEnvInt(
		"RATE_LIMIT_PER_ENDPOINT_BURST_SIZE",
		config.PerEndpointBurstSize,
	); err != nil {
		return err
	}

	return nil
}

// loadAdvancedEnvConfig loads advanced rate limiting configuration from environment variables
func loadAdvancedEnvConfig(config *Config) error {
	var err error

	if config.HTTPRate, err = loadEnvInt("RATE_LIMIT_HTTP_RATE", config.HTTPRate); err != nil {
		return err
	}

	if config.HTTPBurstSize, err = loadEnvInt("RATE_LIMIT_HTTP_BURST_SIZE", config.HTTPBurstSize); err != nil {
		return err
	}

	if config.GRPCRate, err = loadEnvInt("RATE_LIMIT_GRPC_RATE", config.GRPCRate); err != nil {
		return err
	}

	if config.GRPCBurstSize, err = loadEnvInt("RATE_LIMIT_GRPC_BURST_SIZE", config.GRPCBurstSize); err != nil {
		return err
	}

	config.GrpcMetadataKey = loadEnvString("RATE_LIMIT_GRPC_METADATA_KEY", config.GrpcMetadataKey)

	return nil
}

// loadEnvString loads a string value from environment variable
func loadEnvString(envVar, defaultValue string) string {
	if value := os.Getenv(envVar); value != "" {
		return value
	}
	return defaultValue
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (Config, error) {
	config := DefaultConfig()

	if err := loadBasicEnvConfig(&config); err != nil {
		return config, err
	}

	if err := loadAdvancedEnvConfig(&config); err != nil {
		return config, err
	}

	// Load new three-tier rate limiting fields
	if httpRate := os.Getenv("RATE_LIMIT_HTTP_RATE"); httpRate != "" {
		rate, err := strconv.Atoi(httpRate)
		if err != nil {
			return config, fmt.Errorf("invalid RATE_LIMIT_HTTP_RATE value %q: %w", httpRate, err)
		}
		if rate <= 0 {
			return config, fmt.Errorf("RATE_LIMIT_HTTP_RATE must be positive, got %d", rate)
		}
		config.HTTPRate = rate
	}

	if httpBurstSize := os.Getenv("RATE_LIMIT_HTTP_BURST_SIZE"); httpBurstSize != "" {
		size, err := strconv.Atoi(httpBurstSize)
		if err != nil {
			return config, fmt.Errorf("invalid RATE_LIMIT_HTTP_BURST_SIZE value %q: %w", httpBurstSize, err)
		}
		if size <= 0 {
			return config, fmt.Errorf("RATE_LIMIT_HTTP_BURST_SIZE must be positive, got %d", size)
		}
		config.HTTPBurstSize = size
	}

	if grpcRate := os.Getenv("RATE_LIMIT_GRPC_RATE"); grpcRate != "" {
		rate, err := strconv.Atoi(grpcRate)
		if err != nil {
			return config, fmt.Errorf("invalid RATE_LIMIT_GRPC_RATE value %q: %w", grpcRate, err)
		}
		if rate <= 0 {
			return config, fmt.Errorf("RATE_LIMIT_GRPC_RATE must be positive, got %d", rate)
		}
		config.GRPCRate = rate
	}

	if grpcBurstSize := os.Getenv("RATE_LIMIT_GRPC_BURST_SIZE"); grpcBurstSize != "" {
		size, err := strconv.Atoi(grpcBurstSize)
		if err != nil {
			return config, fmt.Errorf("invalid RATE_LIMIT_GRPC_BURST_SIZE value %q: %w", grpcBurstSize, err)
		}
		if size <= 0 {
			return config, fmt.Errorf("RATE_LIMIT_GRPC_BURST_SIZE must be positive, got %d", size)
		}
		config.GRPCBurstSize = size
	}

	if grpcMetadataKey := os.Getenv("RATE_LIMIT_GRPC_METADATA_KEY"); grpcMetadataKey != "" {
		config.GrpcMetadataKey = grpcMetadataKey
	}

	// Validate that per-endpoint rate doesn't exceed global rate
	if config.PerEndpointRate > config.GlobalRate {
		return config, fmt.Errorf("per-endpoint rate (%d) cannot exceed global rate (%d)",
			config.PerEndpointRate, config.GlobalRate)
	}

	return config, nil
}

// LoadFromFile loads configuration from a JSON or YAML file
func LoadFromFile(filePath string) (Config, error) {
	// Start with defaults
	config := DefaultConfig()

	// Read the file
	// #nosec G304 - filePath is provided by user/config, not attacker controlled
	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file %q: %w", filePath, err)
	}

	// Determine file type and parse
	var fileConfig FileConfig
	ext := filepath.Ext(filePath)
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &fileConfig); err != nil {
			return config, fmt.Errorf("failed to parse JSON config file %q: %w", filePath, err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &fileConfig); err != nil {
			return config, fmt.Errorf("failed to parse YAML config file %q: %w", filePath, err)
		}
	default:
		return config, fmt.Errorf("unsupported config file extension %q, supported: .json, .yaml, .yml", ext)
	}

	// Validate and convert to Config
	if err := validateAndConvertFileConfig(&config, &fileConfig); err != nil {
		return config, fmt.Errorf("invalid config file %q: %w", filePath, err)
	}

	return config, nil
}

// validateAndConvertFileConfig validates the file config and converts it to Config
func validateAndConvertFileConfig(config *Config, fileConfig *FileConfig) error {
	// User identification
	config.UserHeader = fileConfig.UserIdentification.HTTPHeader
	config.GrpcMetadataKey = fileConfig.UserIdentification.GRPCMetadataKey

	// Global rate limits
	if fileConfig.RateLimits.Global.Rate <= 0 {
		return fmt.Errorf("global rate must be positive, got %d", fileConfig.RateLimits.Global.Rate)
	}
	if fileConfig.RateLimits.Global.Burst <= 0 {
		return fmt.Errorf("global burst must be positive, got %d", fileConfig.RateLimits.Global.Burst)
	}
	config.GlobalRate = fileConfig.RateLimits.Global.Rate
	config.GlobalBurstSize = fileConfig.RateLimits.Global.Burst

	// HTTP rate limits
	if fileConfig.RateLimits.HTTP.Rate <= 0 {
		return fmt.Errorf("HTTP rate must be positive, got %d", fileConfig.RateLimits.HTTP.Rate)
	}
	if fileConfig.RateLimits.HTTP.Burst <= 0 {
		return fmt.Errorf("HTTP burst must be positive, got %d", fileConfig.RateLimits.HTTP.Burst)
	}
	if fileConfig.RateLimits.HTTP.DefaultMethodRate <= 0 {
		return fmt.Errorf("HTTP default method rate must be positive, got %d", fileConfig.RateLimits.HTTP.DefaultMethodRate)
	}
	config.HTTPRate = fileConfig.RateLimits.HTTP.Rate
	config.HTTPBurstSize = fileConfig.RateLimits.HTTP.Burst
	config.HTTPDefaultMethodRate = fileConfig.RateLimits.HTTP.DefaultMethodRate
	config.HTTPMethods = fileConfig.RateLimits.HTTP.Methods

	// gRPC rate limits
	if fileConfig.RateLimits.GRPC.Rate <= 0 {
		return fmt.Errorf("gRPC rate must be positive, got %d", fileConfig.RateLimits.GRPC.Rate)
	}
	if fileConfig.RateLimits.GRPC.Burst <= 0 {
		return fmt.Errorf("gRPC burst must be positive, got %d", fileConfig.RateLimits.GRPC.Burst)
	}
	if fileConfig.RateLimits.GRPC.DefaultMethodRate <= 0 {
		return fmt.Errorf("gRPC default method rate must be positive, got %d", fileConfig.RateLimits.GRPC.DefaultMethodRate)
	}
	config.GRPCRate = fileConfig.RateLimits.GRPC.Rate
	config.GRPCBurstSize = fileConfig.RateLimits.GRPC.Burst
	config.GRPCDefaultMethodRate = fileConfig.RateLimits.GRPC.DefaultMethodRate
	config.GRPCMethods = fileConfig.RateLimits.GRPC.Methods

	// Set per-endpoint to HTTP default for backward compatibility
	config.PerEndpointRate = config.HTTPDefaultMethodRate
	config.PerEndpointBurstSize = config.HTTPBurstSize

	return nil
}

// Load loads configuration from environment variables or config file
// Priority: config file > environment variables
func Load() (Config, error) {
	// Check for config file
	if configPath := os.Getenv("RATE_LIMIT_CONFIG_PATH"); configPath != "" {
		return LoadFromFile(configPath)
	}

	// Fall back to environment variables
	return LoadFromEnv()
}

// GetRefillInterval returns the time interval for refilling tokens based on the rate
func (c Config) GetRefillInterval(rate int) time.Duration {
	// Calculate how often to add one token (in nanoseconds)
	// rate tokens per second = rate tokens / 1000000000 nanoseconds
	refillInterval := int64(time.Second) / int64(rate)
	return time.Duration(refillInterval)
}