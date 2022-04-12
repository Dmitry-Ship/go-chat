package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserAggregate struct {
	ID           uuid.UUID
	Avatar       string
	Name         string
	Password     string
	CreatedAt    time.Time
	RefreshToken string
}

func NewUser(username string, password string) *UserAggregate {
	return &UserAggregate{
		ID:        uuid.New(),
		Avatar:    string(username[0]),
		CreatedAt: time.Now(),
		Name:      username,
		Password:  password,
	}
}

func (u *UserAggregate) SetRefreshToken(refreshToken string) {
	u.RefreshToken = refreshToken
}
