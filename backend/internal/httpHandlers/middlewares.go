package httpHandlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

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

		userID, err := s.commandHandlers.commands.AuthService.ParseAccessToken(accessToken.Value)

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

func (s *HTTPHandlers) get(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		clientURL := os.Getenv("CLIENT_ORIGIN")
		w.Header().Set("Access-Control-Allow-Origin", clientURL)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func (s *HTTPHandlers) post(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		clientURL := os.Getenv("CLIENT_ORIGIN")

		w.Header().Set("Access-Control-Allow-Origin", clientURL)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Origin")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *HTTPHandlers) getPaginated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		page, _ := strconv.Atoi(query.Get("page"))
		pageSize, _ := strconv.Atoi(query.Get("page_size"))

		p := pagination{
			page:     page,
			pageSize: pageSize,
		}

		ctx := context.WithValue(r.Context(), paginationKey, p)

		s.get(next).ServeHTTP(w, r.WithContext(ctx))
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
