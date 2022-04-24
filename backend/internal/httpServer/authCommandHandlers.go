package httpServer

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (s *CommandController) handleLogin(w http.ResponseWriter, r *http.Request) {
	request := struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokens, err := s.commands.AuthService.Login(request.UserName, request.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Expires:  time.Now().Add(s.commands.AuthService.GetAccessTokenExpiration()),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(s.commands.AuthService.GetRefreshTokenExpiration()),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	expiration := struct {
		AccessTokenExpiration time.Duration `json:"access_token_expiration"`
	}{
		AccessTokenExpiration: s.commands.AuthService.GetAccessTokenExpiration(),
	}

	json.NewEncoder(w).Encode(expiration)
}

func (s *CommandController) handleLogout(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userId").(uuid.UUID)
	err := s.commands.AuthService.Logout(userID)

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

func (s *CommandController) handleSignUp(w http.ResponseWriter, r *http.Request) {
	request := struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokens, err := s.commands.AuthService.SignUp(request.UserName, request.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(s.commands.AuthService.GetAccessTokenExpiration()),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(s.commands.AuthService.GetRefreshTokenExpiration()),
		SameSite: http.SameSiteNoneMode,
	})

	expiration := struct {
		AccessTokenExpiration time.Duration `json:"access_token_expiration"`
	}{
		AccessTokenExpiration: s.commands.AuthService.GetAccessTokenExpiration(),
	}

	json.NewEncoder(w).Encode(expiration)
}

func (s *CommandController) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")

	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTokens, err := s.commands.AuthService.RotateTokens(refreshToken.Value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    newTokens.AccessToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(s.commands.AuthService.GetAccessTokenExpiration()),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newTokens.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(s.commands.AuthService.GetRefreshTokenExpiration()),
		SameSite: http.SameSiteNoneMode,
	})

	expiration := struct {
		AccessTokenExpiration time.Duration `json:"access_token_expiration"`
	}{
		AccessTokenExpiration: s.commands.AuthService.GetAccessTokenExpiration(),
	}

	json.NewEncoder(w).Encode(expiration)
}
