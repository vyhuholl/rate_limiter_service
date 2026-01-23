package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"rate_limiter_service/internal/config"
	"rate_limiter_service/pkg/middleware"
)

// Interceptor provides gRPC rate limiting functionality
type Interceptor struct {
	config             config.Config
	globalLimiter      middleware.GlobalLimiterInterface
	grpcLimiter        middleware.GRPCLimiterInterface
	perMethodLimiter   GRPCMethodLimiterInterface
}

// GRPCMethodLimiterInterface defines the interface for gRPC per-method limiters
type GRPCMethodLimiterInterface interface {
	Allow(userID, method string) bool
	Reset()
}

// InMemoryGRPCMethodLimiter enforces per-method rate limits for gRPC using in-memory storage
type InMemoryGRPCMethodLimiter struct {
	config config.Config
	// buckets stores token buckets keyed by "userID:method"
	buckets map[string]*middleware.TokenBucket
}

// NewInMemoryGRPCMethodLimiter creates a new in-memory gRPC per-method rate limiter
func NewInMemoryGRPCMethodLimiter(cfg config.Config) *InMemoryGRPCMethodLimiter {
	return &InMemoryGRPCMethodLimiter{
		config: cfg,
		buckets: make(map[string]*middleware.TokenBucket),
	}
}

// Allow checks if the gRPC request for the given user and method is allowed
func (gml *InMemoryGRPCMethodLimiter) Allow(userID, method string) bool {
	key := userID + ":" + method
	bucket, exists := gml.buckets[key]
	if !exists {
		// Determine rate for this method
		rate := gml.getRateForMethod(method)
		burstSize := gml.config.GRPCBurstSize
		bucket = middleware.NewTokenBucket(burstSize, rate)
		gml.buckets[key] = bucket
	}
	return bucket.Allow()
}

// getRateForMethod returns the rate limit for a specific gRPC method
func (gml *InMemoryGRPCMethodLimiter) getRateForMethod(method string) int {
	if rate, ok := gml.config.GRPCMethods[method]; ok {
		return rate
	}
	return gml.config.GRPCDefaultMethodRate
}

// Reset clears all rate limiting state for testing
func (gml *InMemoryGRPCMethodLimiter) Reset() {
	gml.buckets = make(map[string]*middleware.TokenBucket)
}

// NewInterceptor creates a new gRPC rate limiting interceptor
func NewInterceptor(cfg config.Config) *Interceptor {
	factory := middleware.NewLimiterFactory(cfg)
	return &Interceptor{
		config:           cfg,
		globalLimiter:    factory.CreateGlobalLimiter(),
		grpcLimiter:      factory.CreateGRPCLimiter(),
		perMethodLimiter: NewInMemoryGRPCMethodLimiter(cfg),
	}
}

// UnaryInterceptor returns a gRPC unary interceptor for rate limiting
func (i *Interceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		userID := i.extractUserID(ctx)

		// Check global limit first
		if !i.globalLimiter.Allow(userID) {
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded: global")
		}

		// Check gRPC-only limit
		if !i.grpcLimiter.Allow(userID) {
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded: grpc")
		}

		// Check per-method limit
		if !i.perMethodLimiter.Allow(userID, info.FullMethod) {
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded: per-method")
		}

		// Request allowed, call handler
		return handler(ctx, req)
	}
}

// extractUserID extracts the user ID from gRPC metadata
func (i *Interceptor) extractUserID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "anonymous"
	}

	values := md.Get(i.config.GrpcMetadataKey)
	if len(values) == 0 || strings.TrimSpace(values[0]) == "" {
		return "anonymous"
	}

	return strings.TrimSpace(values[0])
}

// Reset clears all rate limiting state for testing
func (i *Interceptor) Reset() {
	i.globalLimiter.Reset()
	i.grpcLimiter.Reset()
	i.perMethodLimiter.Reset()
}