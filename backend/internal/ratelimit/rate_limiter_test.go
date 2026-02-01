package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucketRateLimiter_CheckLimit_AllowsWithinCapacity(t *testing.T) {
	config := Config{
		MaxConnections: 5,
		WindowDuration: 1 * time.Minute,
	}
	rateLimiter := NewTokenBucketRateLimiter(config)

	for i := 0; i < 5; i++ {
		allowed, _ := rateLimiter.CheckLimit("test-key")
		assert.True(t, allowed)
		rateLimiter.RecordAttempt("test-key")
	}
}

func TestTokenBucketRateLimiter_CheckLimit_ExceedsCapacity(t *testing.T) {
	config := Config{
		MaxConnections: 3,
		WindowDuration: 1 * time.Minute,
	}
	rateLimiter := NewTokenBucketRateLimiter(config)

	for i := 0; i < 3; i++ {
		allowed, _ := rateLimiter.CheckLimit("test-key")
		assert.True(t, allowed)
		rateLimiter.RecordAttempt("test-key")
	}

	allowed, retryAfter := rateLimiter.CheckLimit("test-key")
	assert.False(t, allowed)
	assert.Greater(t, retryAfter, 0)
}

func TestTokenBucketRateLimiter_CheckLimit_RefillsAfterWindow(t *testing.T) {
	config := Config{
		MaxConnections: 3,
		WindowDuration: 100 * time.Millisecond,
	}
	rateLimiter := NewTokenBucketRateLimiter(config)

	for i := 0; i < 3; i++ {
		rateLimiter.RecordAttempt("test-key")
	}

	allowed, _ := rateLimiter.CheckLimit("test-key")
	assert.False(t, allowed)

	time.Sleep(110 * time.Millisecond)

	allowed, _ = rateLimiter.CheckLimit("test-key")
	assert.True(t, allowed)
}

func TestTokenBucketRateLimiter_CheckLimit_DifferentKeys(t *testing.T) {
	config := Config{
		MaxConnections: 2,
		WindowDuration: 1 * time.Minute,
	}
	rateLimiter := NewTokenBucketRateLimiter(config)

	for i := 0; i < 2; i++ {
		rateLimiter.RecordAttempt("key1")
	}

	allowed, _ := rateLimiter.CheckLimit("key1")
	assert.False(t, allowed)

	allowed, _ = rateLimiter.CheckLimit("key2")
	assert.True(t, allowed)
}

func TestTokenBucketRateLimiter_ConcurrentAccess(t *testing.T) {
	config := Config{
		MaxConnections: 10,
		WindowDuration: 1 * time.Minute,
	}
	rateLimiter := NewTokenBucketRateLimiter(config)

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			allowed, _ := rateLimiter.CheckLimit("concurrent-key")
			assert.True(t, allowed)
			rateLimiter.RecordAttempt("concurrent-key")
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	allowed, _ := rateLimiter.CheckLimit("concurrent-key")
	assert.False(t, allowed)
}
