package httpServer

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (s *HTTPServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	request := struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokens, err := s.app.Commands.AuthService.Login(request.UserName, request.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Expires:  time.Now().Add(s.app.Commands.AuthService.GetAccessTokenExpiration()),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(s.app.Commands.AuthService.GetRefreshTokenExpiration()),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	expiration := struct {
		AccessTokenExpiration time.Duration `json:"access_token_expiration"`
	}{
		AccessTokenExpiration: s.app.Commands.AuthService.GetAccessTokenExpiration(),
	}

	json.NewEncoder(w).Encode(expiration)
}

func (s *HTTPServer) handleLogout(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userId").(uuid.UUID)
	err := s.app.Commands.AuthService.Logout(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	json.NewEncoder(w).Encode("OK")
}

func (s *HTTPServer) handleSignUp(w http.ResponseWriter, r *http.Request) {
	request := struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokens, err := s.app.Commands.AuthService.SignUp(request.UserName, request.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(s.app.Commands.AuthService.GetAccessTokenExpiration()),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(s.app.Commands.AuthService.GetRefreshTokenExpiration()),
		SameSite: http.SameSiteNoneMode,
	})

	expiration := struct {
		AccessTokenExpiration time.Duration `json:"access_token_expiration"`
	}{
		AccessTokenExpiration: s.app.Commands.AuthService.GetAccessTokenExpiration(),
	}

	json.NewEncoder(w).Encode(expiration)
}

func (s *HTTPServer) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")

	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTokens, err := s.app.Commands.AuthService.RotateTokens(refreshToken.Value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    newTokens.AccessToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(s.app.Commands.AuthService.GetAccessTokenExpiration()),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newTokens.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(s.app.Commands.AuthService.GetRefreshTokenExpiration()),
		SameSite: http.SameSiteNoneMode,
	})

	expiration := struct {
		AccessTokenExpiration time.Duration `json:"access_token_expiration"`
	}{
		AccessTokenExpiration: s.app.Commands.AuthService.GetAccessTokenExpiration(),
	}

	json.NewEncoder(w).Encode(expiration)
}
