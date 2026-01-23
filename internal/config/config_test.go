package config

import (
	"os"
	"reflect"
	"testing"
	"time"
)

// configsEqual compares two Config structs for equality
func configsEqual(a, b Config) bool {
	return a.UserHeader == b.UserHeader &&
		a.GrpcMetadataKey == b.GrpcMetadataKey &&
		a.PerEndpointRate == b.PerEndpointRate &&
		a.GlobalRate == b.GlobalRate &&
		a.GlobalBurstSize == b.GlobalBurstSize &&
		a.PerEndpointBurstSize == b.PerEndpointBurstSize &&
		a.HTTPRate == b.HTTPRate &&
		a.HTTPBurstSize == b.HTTPBurstSize &&
		a.GRPCRate == b.GRPCRate &&
		a.GRPCBurstSize == b.GRPCBurstSize &&
		a.HTTPDefaultMethodRate == b.HTTPDefaultMethodRate &&
		a.GRPCDefaultMethodRate == b.GRPCDefaultMethodRate &&
		reflect.DeepEqual(a.HTTPMethods, b.HTTPMethods) &&
		reflect.DeepEqual(a.GRPCMethods, b.GRPCMethods)
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	expected := Config{
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
		HTTPMethods:           nil,
		GRPCMethods:           nil,
	}

	if config.UserHeader != expected.UserHeader ||
		config.GrpcMetadataKey != expected.GrpcMetadataKey ||
		config.PerEndpointRate != expected.PerEndpointRate ||
		config.GlobalRate != expected.GlobalRate ||
		config.GlobalBurstSize != expected.GlobalBurstSize ||
		config.PerEndpointBurstSize != expected.PerEndpointBurstSize ||
		config.HTTPRate != expected.HTTPRate ||
		config.HTTPBurstSize != expected.HTTPBurstSize ||
		config.GRPCRate != expected.GRPCRate ||
		config.GRPCBurstSize != expected.GRPCBurstSize ||
		config.HTTPDefaultMethodRate != expected.HTTPDefaultMethodRate ||
		config.GRPCDefaultMethodRate != expected.GRPCDefaultMethodRate ||
		len(config.HTTPMethods) != 0 ||
		len(config.GRPCMethods) != 0 {
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
				GrpcMetadataKey:       "user-id",
				PerEndpointRate:       20,
				GlobalRate:            100,
				GlobalBurstSize:       10,
				PerEndpointBurstSize:  10,
				HTTPRate:              50,
				HTTPBurstSize:         5,
				GRPCRate:              50,
				GRPCBurstSize:         5,
				HTTPDefaultMethodRate: 10,
				GRPCDefaultMethodRate: 10,
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
				GrpcMetadataKey:       "user-id",
				PerEndpointRate:       10,
				GlobalRate:            200,
				GlobalBurstSize:       10,
				PerEndpointBurstSize:  10,
				HTTPRate:              50,
				HTTPBurstSize:         5,
				GRPCRate:              50,
				GRPCBurstSize:         5,
				HTTPDefaultMethodRate: 10,
				GRPCDefaultMethodRate: 10,
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
				GrpcMetadataKey:       "user-id",
				PerEndpointRate:       10,
				GlobalRate:            100,
				GlobalBurstSize:       20,
				PerEndpointBurstSize:  20,
				HTTPRate:              50,
				HTTPBurstSize:         5,
				GRPCRate:              50,
				GRPCBurstSize:         5,
				HTTPDefaultMethodRate: 10,
				GRPCDefaultMethodRate: 10,
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
				GrpcMetadataKey:       "user-id",
				PerEndpointRate:       5,
				GlobalRate:            50,
				GlobalBurstSize:       5,
				PerEndpointBurstSize:  5,
				HTTPRate:              50,
				HTTPBurstSize:         5,
				GRPCRate:              50,
				GRPCBurstSize:         5,
				HTTPDefaultMethodRate: 10,
				GRPCDefaultMethodRate: 10,
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
				GrpcMetadataKey:       "user-id",
				PerEndpointRate:       150,
				GlobalRate:            100,
				GlobalBurstSize:       10,
				PerEndpointBurstSize:  10,
				HTTPRate:              50,
				HTTPBurstSize:         5,
				GRPCRate:              50,
				GRPCBurstSize:         5,
				HTTPDefaultMethodRate: 10,
				GRPCDefaultMethodRate: 10,
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
			_ = os.Unsetenv("RATE_LIMIT_HTTP_RATE")
			_ = os.Unsetenv("RATE_LIMIT_HTTP_BURST_SIZE")
			_ = os.Unsetenv("RATE_LIMIT_GRPC_RATE")
			_ = os.Unsetenv("RATE_LIMIT_GRPC_BURST_SIZE")
			_ = os.Unsetenv("RATE_LIMIT_GRPC_METADATA_KEY")
			_ = os.Unsetenv("RATE_LIMIT_CONFIG_PATH")

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
				if !configsEqual(config, tt.expected) {
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

func TestLoadFromFile(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		content  string
		expected Config
		hasError bool
	}{
		{
			name:     "valid JSON config",
			filePath: "/tmp/test_config.json",
			content: `{
				"rate_limits": {
					"global": {"rate": 100, "burst": 10},
					"http": {"rate": 50, "burst": 5, "default_method_rate": 10, "methods": {"GET /api/users": 20}},
					"grpc": {"rate": 30, "burst": 3, "default_method_rate": 5, "methods": {"/UserService/GetUser": 15}}
				},
				"user_identification": {
					"http_header": "X-API-Key",
					"grpc_metadata_key": "user"
				}
			}`,
			expected: Config{
				UserHeader:            "X-API-Key",
				GrpcMetadataKey:       "user",
				GlobalRate:            100,
				GlobalBurstSize:       10,
				HTTPRate:              50,
				HTTPBurstSize:         5,
				HTTPDefaultMethodRate: 10,
				HTTPMethods:           map[string]int{"GET /api/users": 20},
				GRPCRate:              30,
				GRPCBurstSize:         3,
				GRPCDefaultMethodRate: 5,
				GRPCMethods:           map[string]int{"/UserService/GetUser": 15},
				PerEndpointRate:       10,
				PerEndpointBurstSize:  5,
			},
			hasError: false,
		},
		{
			name:     "unsupported file extension",
			filePath: "/tmp/test_config.txt",
			content:  "invalid",
			expected: Config{},
			hasError: true,
		},
		{
			name:     "invalid JSON",
			filePath: "/tmp/test_config.json",
			content:  `{"invalid": json}`,
			expected: Config{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write test file
			err := os.WriteFile(tt.filePath, []byte(tt.content), 0600)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}
			defer func() {
				_ = os.Remove(tt.filePath) // Clean up
			}()

			config, err := LoadFromFile(tt.filePath)

			if tt.hasError {
				if err == nil {
					t.Errorf("LoadFromFile() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("LoadFromFile() unexpected error: %v", err)
				}
				if !configsEqual(config, tt.expected) {
					t.Errorf("LoadFromFile() = %v, want %v", config, tt.expected)
				}
			}
		})
	}
}