package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type userIDKeyType string

const userIDKey userIDKeyType = "userId"

func (s *Server) private(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie("access_token")

		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID, err := s.authCommands.ParseAccessToken(accessToken.Value)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type paginationKeyType string

const paginationKey paginationKeyType = "pagination"

type pagination struct {
	page     int
	pageSize int
}

func (p pagination) GetPage() int {
	return p.page
}

func (p pagination) GetPageSize() int {
	return p.pageSize
}

func returnError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	log.Printf("Internal error: %v", err)
	errorResponse := struct {
		Error string `json:"error"`
	}{
		Error: "An internal error occurred",
	}

	if err = json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Error encoding error response: %v", err)
		return
	}
}

func (s *Server) securityHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.Itoa(HSTSMaxAgeSeconds)+"; includeSubDomains")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next(w, r)
	}
}

func (s *Server) limitRequestBodySize(maxBytes int64, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		next(w, r)
	}
}

func (s *Server) corsHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", s.config.ClientOrigin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.WriteHeader(http.StatusOK)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", s.config.ClientOrigin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	}
}

func withPagination(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		page, _ := strconv.Atoi(query.Get("page"))
		pageSize, _ := strconv.Atoi(query.Get("page_size"))

		p := pagination{
			page:     page,
			pageSize: pageSize,
		}

		ctx := context.WithValue(r.Context(), paginationKey, p)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
