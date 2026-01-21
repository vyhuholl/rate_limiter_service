package config

import (
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	expected := Config{
		UserHeader:            "X-User-ID",
		PerEndpointRate:       10,
		GlobalRate:            100,
		GlobalBurstSize:       10,
		PerEndpointBurstSize:  10,
	}

	if config != expected {
		t.Errorf("DefaultConfig() = %v, want %v", config, expected)
	}
}

func TestLoadFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected Config
		hasError bool
	}{
		{
			name:     "no environment variables",
			env:      map[string]string{},
			expected: DefaultConfig(),
			hasError: false,
		},
		{
			name: "custom user header",
			env: map[string]string{
				"RATE_LIMIT_USER_HEADER": "X-API-Key",
			},
			expected: Config{
				UserHeader:            "X-API-Key",
				PerEndpointRate:       10,
				GlobalRate:            100,
				GlobalBurstSize:       10,
				PerEndpointBurstSize:  10,
			},
			hasError: false,
		},
		{
			name: "custom per-endpoint rate",
			env: map[string]string{
				"RATE_LIMIT_PER_ENDPOINT": "20",
			},
			expected: Config{
				UserHeader:            "X-User-ID",
				PerEndpointRate:       20,
				GlobalRate:            100,
				GlobalBurstSize:       10,
				PerEndpointBurstSize:  10,
			},
			hasError: false,
		},
		{
			name: "custom global rate",
			env: map[string]string{
				"RATE_LIMIT_GLOBAL": "200",
			},
			expected: Config{
				UserHeader:            "X-User-ID",
				PerEndpointRate:       10,
				GlobalRate:            200,
				GlobalBurstSize:       10,
				PerEndpointBurstSize:  10,
			},
			hasError: false,
		},
		{
			name: "custom burst size",
			env: map[string]string{
				"RATE_LIMIT_BURST_SIZE": "20",
			},
			expected: Config{
				UserHeader:            "X-User-ID",
				PerEndpointRate:       10,
				GlobalRate:            100,
				GlobalBurstSize:       20,
				PerEndpointBurstSize:  20,
			},
			hasError: false,
		},
		{
			name: "all custom values",
			env: map[string]string{
				"RATE_LIMIT_USER_HEADER":    "Authorization",
				"RATE_LIMIT_PER_ENDPOINT":   "5",
				"RATE_LIMIT_GLOBAL":         "50",
				"RATE_LIMIT_BURST_SIZE":     "5",
			},
			expected: Config{
				UserHeader:            "Authorization",
				PerEndpointRate:       5,
				GlobalRate:            50,
				GlobalBurstSize:       5,
				PerEndpointBurstSize:  5,
			},
			hasError: false,
		},
		{
			name: "invalid per-endpoint rate",
			env: map[string]string{
				"RATE_LIMIT_PER_ENDPOINT": "invalid",
			},
			expected: DefaultConfig(),
			hasError: true,
		},
		{
			name: "zero per-endpoint rate",
			env: map[string]string{
				"RATE_LIMIT_PER_ENDPOINT": "0",
			},
			expected: DefaultConfig(),
			hasError: true,
		},
		{
			name: "negative per-endpoint rate",
			env: map[string]string{
				"RATE_LIMIT_PER_ENDPOINT": "-5",
			},
			expected: DefaultConfig(),
			hasError: true,
		},
		{
			name: "invalid global rate",
			env: map[string]string{
				"RATE_LIMIT_GLOBAL": "not-a-number",
			},
			expected: DefaultConfig(),
			hasError: true,
		},
		{
			name: "zero global rate",
			env: map[string]string{
				"RATE_LIMIT_GLOBAL": "0",
			},
			expected: DefaultConfig(),
			hasError: true,
		},
		{
			name: "invalid burst size",
			env: map[string]string{
				"RATE_LIMIT_BURST_SIZE": "abc",
			},
			expected: DefaultConfig(),
			hasError: true,
		},
		{
			name: "zero burst size",
			env: map[string]string{
				"RATE_LIMIT_BURST_SIZE": "0",
			},
			expected: DefaultConfig(),
			hasError: true,
		},
		{
			name: "per-endpoint rate exceeds global rate",
			env: map[string]string{
				"RATE_LIMIT_PER_ENDPOINT": "150",
				"RATE_LIMIT_GLOBAL":       "100",
			},
			expected: Config{
				UserHeader:            "X-User-ID",
				PerEndpointRate:       150,
				GlobalRate:            100,
				GlobalBurstSize:       10,
				PerEndpointBurstSize:  10,
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all environment variables first
			_ = os.Unsetenv("RATE_LIMIT_USER_HEADER")
			_ = os.Unsetenv("RATE_LIMIT_PER_ENDPOINT")
			_ = os.Unsetenv("RATE_LIMIT_GLOBAL")
			_ = os.Unsetenv("RATE_LIMIT_BURST_SIZE")
			_ = os.Unsetenv("RATE_LIMIT_GLOBAL_BURST_SIZE")
			_ = os.Unsetenv("RATE_LIMIT_PER_ENDPOINT_BURST_SIZE")

			// Set test environment variables
			for key, value := range tt.env {
				_ = os.Setenv(key, value)
			}

			config, err := LoadFromEnv()

			if tt.hasError {
				if err == nil {
					t.Errorf("LoadFromEnv() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("LoadFromEnv() unexpected error: %v", err)
				}
				if config != tt.expected {
					t.Errorf("LoadFromEnv() = %v, want %v", config, tt.expected)
				}
			}
		})
	}
}

func TestGetRefillInterval(t *testing.T) {
	tests := []struct {
		name     string
		rate     int
		expected time.Duration
	}{
		{
			name:     "rate 1 per second",
			rate:     1,
			expected: time.Second,
		},
		{
			name:     "rate 2 per second",
			rate:     2,
			expected: 500 * time.Millisecond,
		},
		{
			name:     "rate 10 per second",
			rate:     10,
			expected: 100 * time.Millisecond,
		},
		{
			name:     "rate 100 per second",
			rate:     100,
			expected: 10 * time.Millisecond,
		},
	}

	config := DefaultConfig()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.GetRefillInterval(tt.rate)
			if result != tt.expected {
				t.Errorf("GetRefillInterval(%d) = %v, want %v", tt.rate, result, tt.expected)
			}
		})
	}
}