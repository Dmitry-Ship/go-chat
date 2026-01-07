package domain

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Store(user *User) error
	Update(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	FindByUsername(username string) (*User, error)
}

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
		return "", errors.New("password too short")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ComparePassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

type User struct {
	aggregate
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
