package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	request := struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	tokens, err := s.authCommands.Login(request.UserName, request.Password)

	if err != nil {
		returnError(w, http.StatusUnauthorized, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Expires:  time.Now().Add(tokens.AccessTokenExpiration),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(tokens.RefreshTokenExpiration),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	expiration := struct {
		AccessTokenExpiration time.Duration `json:"access_token_expiration"`
	}{
		AccessTokenExpiration: tokens.AccessTokenExpiration,
	}

	if err = json.NewEncoder(w).Encode(expiration); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
		return
	}

	if err := s.authCommands.Logout(userID); err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	if err := json.NewEncoder(w).Encode("OK"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleSignUp(w http.ResponseWriter, r *http.Request) {
	request := struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	tokens, err := s.authCommands.SignUp(request.UserName, request.Password)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(tokens.AccessTokenExpiration),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(tokens.RefreshTokenExpiration),
		SameSite: http.SameSiteNoneMode,
	})

	expiration := struct {
		AccessTokenExpiration time.Duration `json:"access_token_expiration"`
	}{
		AccessTokenExpiration: tokens.AccessTokenExpiration,
	}

	if err = json.NewEncoder(w).Encode(expiration); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")

	if err != nil {
		if err == http.ErrNoCookie {
			returnError(w, http.StatusUnauthorized, err)
			return
		}

		returnError(w, http.StatusBadRequest, err)
		return
	}

	newTokens, err := s.authCommands.RotateTokens(refreshToken.Value)

	if err != nil {
		returnError(w, http.StatusUnauthorized, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    newTokens.AccessToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(newTokens.AccessTokenExpiration),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newTokens.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(newTokens.RefreshTokenExpiration),
		SameSite: http.SameSiteNoneMode,
	})

	expiration := struct {
		AccessTokenExpiration time.Duration `json:"access_token_expiration"`
	}{
		AccessTokenExpiration: newTokens.AccessTokenExpiration,
	}

	if err := json.NewEncoder(w).Encode(expiration); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
