package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"GitHub/go-chat/backend/internal/ratelimit"
	"github.com/stretchr/testify/assert"
)

func TestWSRateLimitMiddleware_AllowsWithinLimit(t *testing.T) {
	config := ratelimit.Config{
		MaxConnections: 5,
		WindowDuration: 1 * time.Minute,
	}

	server := &Server{
		ipRateLimiter:   ratelimit.NewSlidingWindowRateLimiter(config),
		userRateLimiter: ratelimit.NewSlidingWindowRateLimiter(config),
	}

	handler := server.wsRateLimit(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/ws", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestWSRateLimitMiddleware_BlocksExceededIPLimit(t *testing.T) {
	config := ratelimit.Config{
		MaxConnections: 2,
		WindowDuration: 1 * time.Minute,
	}

	server := &Server{
		ipRateLimiter:   ratelimit.NewSlidingWindowRateLimiter(config),
		userRateLimiter: ratelimit.NewSlidingWindowRateLimiter(config),
	}

	handler := server.wsRateLimit(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/ws", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	req := httptest.NewRequest("GET", "/ws", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.NotEmpty(t, w.Header().Get("Retry-After"))
}

func TestWSRateLimitMiddleware_DifferentIPs(t *testing.T) {
	config := ratelimit.Config{
		MaxConnections: 2,
		WindowDuration: 1 * time.Minute,
	}

	server := &Server{
		ipRateLimiter:   ratelimit.NewSlidingWindowRateLimiter(config),
		userRateLimiter: ratelimit.NewSlidingWindowRateLimiter(config),
	}

	handler := server.wsRateLimit(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/ws", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		handler(w, req)
	}

	req := httptest.NewRequest("GET", "/ws", nil)
	req.RemoteAddr = "192.168.1.2:12345"
	w := httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWSRateLimitMiddleware_SlidingWindow(t *testing.T) {
	config := ratelimit.Config{
		MaxConnections: 2,
		WindowDuration: 100 * time.Millisecond,
	}

	server := &Server{
		ipRateLimiter:   ratelimit.NewSlidingWindowRateLimiter(config),
		userRateLimiter: ratelimit.NewSlidingWindowRateLimiter(config),
	}

	handler := server.wsRateLimit(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/ws", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		handler(w, req)
	}

	req := httptest.NewRequest("GET", "/ws", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)

	time.Sleep(110 * time.Millisecond)

	req = httptest.NewRequest("GET", "/ws", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w = httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWSRateLimitMiddleware_XForwardedForHeader(t *testing.T) {
	config := ratelimit.Config{
		MaxConnections: 1,
		WindowDuration: 1 * time.Minute,
	}

	server := &Server{
		ipRateLimiter:   ratelimit.NewSlidingWindowRateLimiter(config),
		userRateLimiter: ratelimit.NewSlidingWindowRateLimiter(config),
	}

	handler := server.wsRateLimit(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/ws", nil)
	req.Header.Set("X-Forwarded-For", "10.0.0.1, 10.0.0.2")
	w := httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req = httptest.NewRequest("GET", "/ws", nil)
	req.Header.Set("X-Forwarded-For", "10.0.0.1, 10.0.0.2")
	w = httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestWSRateLimitMiddleware_XRealIPHeader(t *testing.T) {
	config := ratelimit.Config{
		MaxConnections: 1,
		WindowDuration: 1 * time.Minute,
	}

	server := &Server{
		ipRateLimiter:   ratelimit.NewSlidingWindowRateLimiter(config),
		userRateLimiter: ratelimit.NewSlidingWindowRateLimiter(config),
	}

	handler := server.wsRateLimit(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/ws", nil)
	req.Header.Set("X-Real-IP", "10.0.0.5")
	w := httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req = httptest.NewRequest("GET", "/ws", nil)
	req.Header.Set("X-Real-IP", "10.0.0.5")
	w = httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}
