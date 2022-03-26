package domain

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid"`
	Avatar       string    `json:"avatar"`
	Name         string    `json:"name"`
	Password     string    `json:"-"`
	RefreshToken string    `json:"-" gorm:"column:refresh_token" `
}

func NewUser(username string, password string) *User {
	return &User{
		ID:       uuid.New(),
		Avatar:   string(username[0]),
		Name:     username,
		Password: password,
	}
}

func (u *User) SetRefreshToken(refreshToken string) {
	u.RefreshToken = refreshToken
}
