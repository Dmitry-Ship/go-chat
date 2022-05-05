package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	Store(user *User) error
	Update(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	FindByUsername(username string) (*User, error)
}

type User struct {
	ID           uuid.UUID
	Avatar       string
	Name         string
	Password     string
	CreatedAt    time.Time
	RefreshToken string
}

func NewUser(username string, password string) *User {
	return &User{
		ID:        uuid.New(),
		Avatar:    string(username[0]),
		CreatedAt: time.Now(),
		Name:      username,
		Password:  password,
	}
}

func (u *User) SetRefreshToken(refreshToken string) {
	u.RefreshToken = refreshToken
}
