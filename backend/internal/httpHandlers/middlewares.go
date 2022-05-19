package httpHandlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

func (s *HTTPHandlers) withHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientURL := os.Getenv("CLIENT_ORIGIN")

		w.Header().Set("Access-Control-Allow-Origin", clientURL)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

type userIDKeyType string

const userIDKey userIDKeyType = "userId"

func (s *HTTPHandlers) private(next http.HandlerFunc) http.HandlerFunc {
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

		userId, err := s.commandController.commands.AuthService.ParseAccessToken(accessToken.Value)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userId)

		s.withHeaders(next).ServeHTTP(w, r.WithContext(ctx))
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

func (s *HTTPHandlers) paginate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		page, _ := strconv.Atoi(query.Get("page"))
		pageSize, _ := strconv.Atoi(query.Get("page_size"))

		p := pagination{
			page:     page,
			pageSize: pageSize,
		}

		ctx := context.WithValue(r.Context(), paginationKey, p)

		s.withHeaders(next).ServeHTTP(w, r.WithContext(ctx))
	})
}

func returnError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	errorResponse := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	err = json.NewEncoder(w).Encode(errorResponse)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
