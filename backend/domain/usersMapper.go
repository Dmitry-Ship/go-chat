package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID     uuid.UUID `json:"id"`
	Avatar string    `json:"avatar"`
	Name   string    `json:"name"`
}

type UserPersistence struct {
	ID           uuid.UUID `gorm:"type:uuid"`
	Avatar       string
	Name         string
	Password     string
	CreatedAt    time.Time
	RefreshToken string `gorm:"column:refresh_token"`
}

func (UserPersistence) TableName() string {
	return "users"
}

func ToUserDTO(user *UserPersistence) *UserDTO {
	return &UserDTO{
		ID:     user.ID,
		Avatar: user.Avatar,
		Name:   user.Name,
	}
}

func ToUserPersistence(user *User) *UserPersistence {
	return &UserPersistence{
		ID:           user.ID,
		Avatar:       user.Avatar,
		Name:         user.Name,
		Password:     user.Password,
		CreatedAt:    user.CreatedAt,
		RefreshToken: user.RefreshToken,
	}
}

func ToUserDomain(user *UserPersistence) *User {
	return &User{
		ID:           user.ID,
		Avatar:       user.Avatar,
		Name:         user.Name,
		Password:     user.Password,
		CreatedAt:    user.CreatedAt,
		RefreshToken: user.RefreshToken,
	}
}
