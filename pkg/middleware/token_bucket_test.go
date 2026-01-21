package middleware

import (
	"testing"
	"time"
)

func TestNewTokenBucket(t *testing.T) {
	tests := []struct {
		name        string
		capacity    int
		refillRate  int
		expectedCap int
	}{
		{
			name:        "normal capacity and rate",
			capacity:    10,
			refillRate:  5,
			expectedCap: 10,
		},
		{
			name:        "zero capacity defaults to 1",
			capacity:    0,
			refillRate:  5,
			expectedCap: 1,
		},
		{
			name:        "negative capacity defaults to 1",
			capacity:    -5,
			refillRate:  5,
			expectedCap: 1,
		},
		{
			name:        "zero refill rate defaults to 1",
			capacity:    10,
			refillRate:  0,
			expectedCap: 10,
		},
		{
			name:        "negative refill rate defaults to 1",
			capacity:    10,
			refillRate:  -5,
			expectedCap: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := NewTokenBucket(tt.capacity, tt.refillRate)
			if tb.capacity != tt.expectedCap {
				t.Errorf("capacity = %d, want %d", tb.capacity, tt.expectedCap)
			}
			if tb.tokens != tt.expectedCap {
				t.Errorf("initial tokens = %d, want %d", tb.tokens, tt.expectedCap)
			}
		})
	}
}

func TestTokenBucket_Allow(t *testing.T) {
	tb := NewTokenBucket(3, 10) // 3 capacity, 10 tokens/second

	// Should allow 3 requests initially
	for i := 0; i < 3; i++ {
		if !tb.Allow() {
			t.Errorf("Allow() should return true for request %d", i+1)
		}
	}

	// Should reject the 4th request
	if tb.Allow() {
		t.Error("Allow() should return false when bucket is empty")
	}

	// Wait for tokens to refill (200ms for 2 tokens at 10/sec)
	time.Sleep(200 * time.Millisecond)

	// Should allow requests again
	if !tb.Allow() {
		t.Error("Allow() should return true after refill")
	}
}

func TestTokenBucket_Refill(t *testing.T) {
	tb := NewTokenBucket(10, 2) // 10 capacity, 2 tokens/second

	// Empty the bucket
	for i := 0; i < 10; i++ {
		if !tb.Allow() {
			t.Errorf("Failed to consume token %d", i+1)
		}
	}

	// Bucket should be empty
	if tb.GetTokens() != 0 {
		t.Errorf("Expected 0 tokens, got %d", tb.GetTokens())
	}

	// Wait for 1 second (should add 2 tokens)
	time.Sleep(550 * time.Millisecond) // A bit more than 500ms for 1 token

	tokens := tb.GetTokens()
	if tokens < 1 {
		t.Errorf("Expected at least 1 token after refill, got %d", tokens)
	}

	// Wait another second (should add more tokens, up to capacity)
	time.Sleep(550 * time.Millisecond)

	tokens = tb.GetTokens()
	if tokens < 2 {
		t.Errorf("Expected at least 2 tokens after second refill, got %d", tokens)
	}
}

func TestTokenBucket_GetTokens(t *testing.T) {
	tb := NewTokenBucket(5, 10)

	// Initially full
	if tokens := tb.GetTokens(); tokens != 5 {
		t.Errorf("Initial tokens = %d, want 5", tokens)
	}

	// Consume some tokens
	tb.Allow()
	tb.Allow()

	if tokens := tb.GetTokens(); tokens != 3 {
		t.Errorf("Tokens after consuming 2 = %d, want 3", tokens)
	}
}

func TestTokenBucket_GetCapacity(t *testing.T) {
	tb := NewTokenBucket(15, 5)
	if capacity := tb.GetCapacity(); capacity != 15 {
		t.Errorf("GetCapacity() = %d, want 15", capacity)
	}
}

func TestTokenBucket_Reset(t *testing.T) {
	tb := NewTokenBucket(5, 10)

	// Consume all tokens
	for i := 0; i < 5; i++ {
		tb.Allow()
	}

	// Bucket should be empty
	if tokens := tb.GetTokens(); tokens != 0 {
		t.Errorf("Tokens before reset = %d, want 0", tokens)
	}

	// Reset
	tb.Reset()

	// Bucket should be full again
	if tokens := tb.GetTokens(); tokens != 5 {
		t.Errorf("Tokens after reset = %d, want 5", tokens)
	}
}

func TestTokenBucket_ConcurrentAccess(t *testing.T) {
	tb := NewTokenBucket(100, 1000)

	// Test concurrent access
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				tb.Allow()
			}
			done <- true
		}()
	}

	// Wait for all goroutines to finish
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should have consumed exactly 100 tokens
	if tokens := tb.GetTokens(); tokens != 0 {
		t.Errorf("Tokens after concurrent access = %d, want 0", tokens)
	}
}

func TestTokenBucket_RefillInterval(t *testing.T) {
	tb := NewTokenBucket(10, 2) // 2 tokens per second = 500ms per token

	expectedInterval := 500 * time.Millisecond
	if tb.refillInterval != expectedInterval {
		t.Errorf("refillInterval = %v, want %v", tb.refillInterval, expectedInterval)
	}

	// Test with different rates
	tb2 := NewTokenBucket(10, 10) // 10 tokens per second = 100ms per token
	expectedInterval2 := 100 * time.Millisecond
	if tb2.refillInterval != expectedInterval2 {
		t.Errorf("refillInterval = %v, want %v", tb2.refillInterval, expectedInterval2)
	}
}