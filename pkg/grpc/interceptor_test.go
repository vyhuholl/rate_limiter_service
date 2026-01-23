package grpc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"rate_limiter_service/internal/config"
)

const testSuccessResponse = "success"

// testRateLimitHelper is a helper for testing rate limiting scenarios
func testRateLimitHelper(
	t *testing.T,
	cfg config.Config,
	expectedMessage string,
) {
	interceptor := NewInterceptor(cfg)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return testSuccessResponse, nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/TestService/TestMethod",
	}

	md := metadata.New(map[string]string{"user-id": "user123"})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	// First request should be allowed
	_, err := interceptor.UnaryInterceptor()(ctx, "request", info, handler)
	if err != nil {
		t.Errorf("First request should be allowed, got error: %v", err)
	}

	// Second request should be rate limited
	_, err = interceptor.UnaryInterceptor()(ctx, "request", info, handler)
	if err == nil {
		t.Error("Second request should be rate limited")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Errorf("Expected gRPC status error, got %T", err)
	}

	if st.Code() != codes.ResourceExhausted {
		t.Errorf("Expected ResourceExhausted, got %v", st.Code())
	}

	if st.Message() != expectedMessage {
		t.Errorf("Expected %q, got %q", expectedMessage, st.Message())
	}
}

func TestNewInterceptor(t *testing.T) {
	cfg := config.Config{
		GrpcMetadataKey:     "user-id",
		GlobalRate:          100,
		GlobalBurstSize:     10,
		GRPCRate:            50,
		GRPCBurstSize:       5,
		GRPCDefaultMethodRate: 10,
	}

	interceptor := NewInterceptor(cfg)

	if interceptor == nil {
		t.Fatal("NewInterceptor returned nil")
	}

	if interceptor.config.GrpcMetadataKey != cfg.GrpcMetadataKey {
		t.Errorf("config not set correctly")
	}

	if interceptor.globalLimiter == nil {
		t.Error("globalLimiter not initialized")
	}

	if interceptor.grpcLimiter == nil {
		t.Error("grpcLimiter not initialized")
	}

	if interceptor.perMethodLimiter == nil {
		t.Error("perMethodLimiter not initialized")
	}
}

func TestInterceptor_ExtractUserID(t *testing.T) {
	cfg := config.Config{
		GrpcMetadataKey: "user-id",
	}

	interceptor := NewInterceptor(cfg)

	tests := []struct {
		name     string
		metadata metadata.MD
		expected string
	}{
		{
			name:     "metadata present",
			metadata: metadata.New(map[string]string{"user-id": "user123"}),
			expected: "user123",
		},
		{
			name:     "metadata missing",
			metadata: metadata.New(map[string]string{}),
			expected: "anonymous",
		},
		{
			name:     "metadata empty",
			metadata: metadata.New(map[string]string{"user-id": ""}),
			expected: "anonymous",
		},
		{
			name:     "metadata with whitespace",
			metadata: metadata.New(map[string]string{"user-id": "  user123  "}),
			expected: "user123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := metadata.NewIncomingContext(context.Background(), tt.metadata)
			result := interceptor.extractUserID(ctx)
			if result != tt.expected {
				t.Errorf("extractUserID() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestInMemoryGRPCMethodLimiter_Allow(t *testing.T) {
	cfg := config.Config{
		GRPCBurstSize:         2,
		GRPCDefaultMethodRate: 2,
		GRPCMethods:           map[string]int{"/TestService/TestMethod": 5},
	}

	limiter := NewInMemoryGRPCMethodLimiter(cfg)

	userID := "user123"

	// Test default method rate
	if !limiter.Allow(userID, "/TestService/DefaultMethod") {
		t.Error("First request to default method should be allowed")
	}
	if !limiter.Allow(userID, "/TestService/DefaultMethod") {
		t.Error("Second request to default method should be allowed")
	}
	if limiter.Allow(userID, "/TestService/DefaultMethod") {
		t.Error("Third request to default method should be denied")
	}

	// Test configured method rate
	cfg.GRPCBurstSize = 3
	limiter2 := NewInMemoryGRPCMethodLimiter(cfg)

	for i := 0; i < 3; i++ {
		if !limiter2.Allow(userID, "/TestService/TestMethod") {
			t.Errorf("Request %d to configured method should be allowed", i+1)
		}
	}
	if limiter2.Allow(userID, "/TestService/TestMethod") {
		t.Error("Fourth request to configured method should be denied")
	}
}

func TestInterceptor_UnaryInterceptor_Allow(t *testing.T) {
	cfg := config.Config{
		GrpcMetadataKey:       "user-id",
		GlobalRate:            10,
		GlobalBurstSize:       10,
		GRPCRate:              10,
		GRPCBurstSize:         10,
		GRPCDefaultMethodRate: 10,
	}

	interceptor := NewInterceptor(cfg)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/TestService/TestMethod",
	}

	// Create context with user ID
	md := metadata.New(map[string]string{"user-id": "user123"})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := interceptor.UnaryInterceptor()(ctx, "request", info, handler)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp != testSuccessResponse {
		t.Errorf("Expected %q, got %v", testSuccessResponse, resp)
	}
}

func TestInterceptor_UnaryInterceptor_GlobalLimitExceeded(t *testing.T) {
	cfg := config.Config{
		GrpcMetadataKey:       "user-id",
		GlobalRate:            1, // Only 1 request per second globally
		GlobalBurstSize:       1,
		GRPCRate:              10,
		GRPCBurstSize:         10,
		GRPCDefaultMethodRate: 10,
	}

	testRateLimitHelper(t, cfg, "rate limit exceeded: global")
}

func TestInterceptor_UnaryInterceptor_GRPCLimitExceeded(t *testing.T) {
	cfg := config.Config{
		GrpcMetadataKey:       "user-id",
		GlobalRate:            10,
		GlobalBurstSize:       10,
		GRPCRate:              1, // Only 1 gRPC request per second
		GRPCBurstSize:         1,
		GRPCDefaultMethodRate: 10,
	}

	testRateLimitHelper(t, cfg, "rate limit exceeded: grpc")
}