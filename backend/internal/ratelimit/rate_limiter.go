package ratelimit

import (
	"math"
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

type tokenBucketRateLimiter struct {
	config     Config
	mu         sync.Mutex
	store      map[string]*bucketState
	capacity   float64
	refillRate float64
}

type bucketState struct {
	tokens     float64
	lastRefill time.Time
}

func NewTokenBucketRateLimiter(config Config) RateLimiter {
	capacity := float64(config.MaxConnections)
	if capacity < 0 {
		capacity = 0
	}

	refillRate := 0.0
	if capacity > 0 && config.WindowDuration > 0 {
		refillRate = capacity / config.WindowDuration.Seconds()
	}

	return &tokenBucketRateLimiter{
		config:     config,
		store:      make(map[string]*bucketState),
		capacity:   capacity,
		refillRate: refillRate,
	}
}

func NewSlidingWindowRateLimiter(config Config) RateLimiter {
	return NewTokenBucketRateLimiter(config)
}

func (r *tokenBucketRateLimiter) CheckLimit(key string) (bool, int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.capacity <= 0 || r.refillRate <= 0 {
		return false, 0
	}

	bucket := r.getBucket(key, time.Now())
	if bucket.tokens >= 1 {
		return true, 0
	}

	retryAfter := int(math.Ceil((1 - bucket.tokens) / r.refillRate))
	if retryAfter < 0 {
		retryAfter = 0
	}

	return false, retryAfter
}

func (r *tokenBucketRateLimiter) RecordAttempt(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.capacity <= 0 || r.refillRate <= 0 {
		return
	}

	bucket := r.getBucket(key, time.Now())
	if bucket.tokens >= 1 {
		bucket.tokens -= 1
		return
	}

	bucket.tokens = 0
}

func (r *tokenBucketRateLimiter) getBucket(key string, now time.Time) *bucketState {
	bucket, exists := r.store[key]
	if !exists {
		bucket = &bucketState{
			tokens:     r.capacity,
			lastRefill: now,
		}
		r.store[key] = bucket
		return bucket
	}

	elapsed := now.Sub(bucket.lastRefill).Seconds()
	if elapsed <= 0 {
		return bucket
	}

	bucket.tokens = math.Min(r.capacity, bucket.tokens+elapsed*r.refillRate)
	bucket.lastRefill = now
	return bucket
}
