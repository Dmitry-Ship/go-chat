package ratelimit

import (
	"sync"
	"time"
)

type Config struct {
	MaxConnections int
	WindowDuration time.Duration
}

type RateLimiter interface {
	CheckLimit(key string) (bool, int)
	RecordAttempt(key string)
}

type slidingWindowRateLimiter struct {
	config Config
	mu     sync.RWMutex
	store  map[string][]time.Time
}

func NewSlidingWindowRateLimiter(config Config) RateLimiter {
	return &slidingWindowRateLimiter{
		config: config,
		store:  make(map[string][]time.Time),
	}
}

func (r *slidingWindowRateLimiter) CheckLimit(key string) (bool, int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-r.config.WindowDuration)

	timestamps, exists := r.store[key]
	if !exists {
		return true, 0
	}

	validTimestamps := make([]time.Time, 0, len(timestamps))
	for _, ts := range timestamps {
		if ts.After(cutoff) {
			validTimestamps = append(validTimestamps, ts)
		}
	}

	r.store[key] = validTimestamps

	if len(validTimestamps) >= r.config.MaxConnections {
		oldest := validTimestamps[0]
		retryAfter := int(oldest.Add(r.config.WindowDuration).Sub(now).Seconds())
		if retryAfter < 0 {
			retryAfter = 0
		}
		return false, retryAfter
	}

	return true, 0
}

func (r *slidingWindowRateLimiter) RecordAttempt(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-r.config.WindowDuration)

	timestamps := r.store[key]
	validTimestamps := make([]time.Time, 0, len(timestamps)+1)
	for _, ts := range timestamps {
		if ts.After(cutoff) {
			validTimestamps = append(validTimestamps, ts)
		}
	}

	validTimestamps = append(validTimestamps, now)
	r.store[key] = validTimestamps
}
