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

type UserDAO struct {
	ID           uuid.UUID `gorm:"type:uuid"`
	Avatar       string
	Name         string
	Password     string
	CreatedAt    time.Time
	RefreshToken string `gorm:"column:refresh_token"`
}

func (UserDAO) TableName() string {
	return "users"
}

func ToUserDTO(user *UserDAO) *UserDTO {
	return &UserDTO{
		ID:     user.ID,
		Avatar: user.Avatar,
		Name:   user.Name,
	}
}

func ToUserDAO(user *User) *UserDAO {
	return &UserDAO{
		ID:           user.ID,
		Avatar:       user.Avatar,
		Name:         user.Name,
		Password:     user.Password,
		CreatedAt:    user.CreatedAt,
		RefreshToken: user.RefreshToken,
	}
}

func ToUserDomain(user *UserDAO) *User {
	return &User{
		ID:           user.ID,
		Avatar:       user.Avatar,
		Name:         user.Name,
		Password:     user.Password,
		CreatedAt:    user.CreatedAt,
		RefreshToken: user.RefreshToken,
	}
}
