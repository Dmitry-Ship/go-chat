package domain

import (
	"errors"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("username is empty")
	}

	if len(username) > 100 {
		return errors.New("username too long")
	}

	return nil
}

func HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", errors.New("password must be at least 8 characters")
	}

	var found int
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			found |= 1 << 0
		case unicode.IsLower(c):
			found |= 1 << 1
		case unicode.IsDigit(c):
			found |= 1 << 2
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			found |= 1 << 3
		}
		if found == 0xF {
			break
		}
	}

	if found != 0xF {
		return "", errors.New("password must contain uppercase, lowercase, digit, and special character")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ComparePassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

type User struct {
	ID           uuid.UUID
	Avatar       string
	Name         string
	PasswordHash string
	RefreshToken string
}

func NewUser(userID uuid.UUID, username string, passwordHash string) *User {
	return &User{
		ID:           userID,
		Avatar:       string(username[0]),
		Name:         username,
		PasswordHash: passwordHash,
	}
}

func (u *User) SetRefreshToken(refreshToken string) {
	u.RefreshToken = refreshToken
}
