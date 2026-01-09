package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func isSecureRequest(r *http.Request) bool {
	if r.URL.Scheme != "" {
		return r.URL.Scheme == "https"
	}
	return r.TLS != nil
}

func setAuthCookie(w http.ResponseWriter, r *http.Request, name, value string, expires time.Time) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expires,
		HttpOnly: true,
		Path:     "/",
	}

	if isSecureRequest(r) {
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	} else {
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	}

	http.SetCookie(w, cookie)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	request := struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	tokens, err := s.authCommands.Login(r.Context(), request.UserName, request.Password)

	if err != nil {
		returnError(w, http.StatusUnauthorized, err)
		return
	}

	setAuthCookie(w, r, "access_token", tokens.AccessToken, time.Now().Add(tokens.AccessTokenExpiration))
	setAuthCookie(w, r, "refresh_token", tokens.RefreshToken, time.Now().Add(tokens.RefreshTokenExpiration))

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

	if err := s.authCommands.Logout(r.Context(), userID); err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
	}
	if !isSecureRequest(r) {
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	} else {
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	}
	http.SetCookie(w, cookie)

	cookie = &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
	}
	if !isSecureRequest(r) {
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	} else {
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	}
	http.SetCookie(w, cookie)

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

	tokens, err := s.authCommands.SignUp(r.Context(), request.UserName, request.Password)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	setAuthCookie(w, r, "access_token", tokens.AccessToken, time.Now().Add(tokens.AccessTokenExpiration))
	setAuthCookie(w, r, "refresh_token", tokens.RefreshToken, time.Now().Add(tokens.RefreshTokenExpiration))

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

	newTokens, err := s.authCommands.RotateTokens(r.Context(), refreshToken.Value)

	if err != nil {
		returnError(w, http.StatusUnauthorized, err)
		return
	}

	setAuthCookie(w, r, "access_token", newTokens.AccessToken, time.Now().Add(newTokens.AccessTokenExpiration))
	setAuthCookie(w, r, "refresh_token", newTokens.RefreshToken, time.Now().Add(newTokens.RefreshTokenExpiration))

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
