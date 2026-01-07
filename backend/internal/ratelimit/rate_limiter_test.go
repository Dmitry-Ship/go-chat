package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSlidingWindowRateLimiter_CheckLimit_AllowsWithinLimit(t *testing.T) {
	config := Config{
		MaxConnections: 5,
		WindowDuration: 1 * time.Minute,
	}
	rateLimiter := NewSlidingWindowRateLimiter(config)

	for i := 0; i < 5; i++ {
		allowed, _ := rateLimiter.CheckLimit("test-key")
		assert.True(t, allowed)
		rateLimiter.RecordAttempt("test-key")
	}
}

func TestSlidingWindowRateLimiter_CheckLimit_ExceedsLimit(t *testing.T) {
	config := Config{
		MaxConnections: 3,
		WindowDuration: 1 * time.Minute,
	}
	rateLimiter := NewSlidingWindowRateLimiter(config)

	for i := 0; i < 3; i++ {
		allowed, _ := rateLimiter.CheckLimit("test-key")
		assert.True(t, allowed)
		rateLimiter.RecordAttempt("test-key")
	}

	allowed, retryAfter := rateLimiter.CheckLimit("test-key")
	assert.False(t, allowed)
	assert.Greater(t, retryAfter, 0)
}

func TestSlidingWindowRateLimiter_CheckLimit_SlidingWindow(t *testing.T) {
	config := Config{
		MaxConnections: 3,
		WindowDuration: 100 * time.Millisecond,
	}
	rateLimiter := NewSlidingWindowRateLimiter(config)

	for i := 0; i < 3; i++ {
		rateLimiter.RecordAttempt("test-key")
	}

	allowed, _ := rateLimiter.CheckLimit("test-key")
	assert.False(t, allowed)

	time.Sleep(110 * time.Millisecond)

	allowed, _ = rateLimiter.CheckLimit("test-key")
	assert.True(t, allowed)
}

func TestSlidingWindowRateLimiter_CheckLimit_DifferentKeys(t *testing.T) {
	config := Config{
		MaxConnections: 2,
		WindowDuration: 1 * time.Minute,
	}
	rateLimiter := NewSlidingWindowRateLimiter(config)

	for i := 0; i < 2; i++ {
		rateLimiter.RecordAttempt("key1")
	}

	allowed, _ := rateLimiter.CheckLimit("key1")
	assert.False(t, allowed)

	allowed, _ = rateLimiter.CheckLimit("key2")
	assert.True(t, allowed)
}

func TestSlidingWindowRateLimiter_ConcurrentAccess(t *testing.T) {
	config := Config{
		MaxConnections: 10,
		WindowDuration: 1 * time.Minute,
	}
	rateLimiter := NewSlidingWindowRateLimiter(config)

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
