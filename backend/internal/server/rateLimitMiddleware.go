package server

import (
	"net"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) wsRateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		if allowed, retryAfter := s.ipRateLimiter.CheckLimit(ip); !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		accessToken, err := r.Cookie("access_token")
		var userID string

		if err == nil {
			parsedUserID, parseErr := s.authCommands.ParseAccessToken(accessToken.Value)
			if parseErr == nil {
				userID = parsedUserID.String()
			}
		}

		if userID != "" {
			if allowed, retryAfter := s.userRateLimiter.CheckLimit(userID); !allowed {
				w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
		}

		s.ipRateLimiter.RecordAttempt(ip)
		if userID != "" {
			s.userRateLimiter.RecordAttempt(userID)
		}

		next.ServeHTTP(w, r)
	}
}

func (s *Server) httpRateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		if allowed, retryAfter := s.ipRateLimiter.CheckLimit(ip); !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		s.ipRateLimiter.RecordAttempt(ip)

		next.ServeHTTP(w, r)
	}
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}
