package domain

import (
	"errors"

	"github.com/google/uuid"
)

type UserRepository interface {
	Store(user *User) error
	Update(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	FindByUsername(username string) (*User, error)
}

type User struct {
	aggregate
	ID           uuid.UUID
	Avatar       string
	Name         string
	Password     string
	RefreshToken string
}

func NewUser(username string, password string) (*User, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}

	if len(username) > 100 {
		return nil, errors.New("username is too long")
	}

	if password == "" {
		return nil, errors.New("password is empty")
	}

	if len(password) < 8 {
		return nil, errors.New("password is too short")
	}

	return &User{
		ID:       uuid.New(),
		Avatar:   string(username[0]),
		Name:     username,
		Password: password,
	}, nil
}

func (u *User) SetRefreshToken(refreshToken string) {
	u.RefreshToken = refreshToken
}
