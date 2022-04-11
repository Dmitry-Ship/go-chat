package postgres

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/readModel"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid"`
	Avatar       string
	Name         string
	Password     string
	CreatedAt    time.Time
	RefreshToken string `gorm:"column:refresh_token"`
}

func (User) TableName() string {
	return "users"
}

func ToUserDTO(user *User) *readModel.UserDTO {
	return &readModel.UserDTO{
		ID:     user.ID,
		Avatar: user.Avatar,
		Name:   user.Name,
	}
}

func ToUserDAO(user *domain.User) *User {
	return &User{
		ID:           user.ID,
		Avatar:       user.Avatar,
		Name:         user.Name,
		Password:     user.Password,
		CreatedAt:    user.CreatedAt,
		RefreshToken: user.RefreshToken,
	}
}

func ToUserDomain(user *User) *domain.User {
	return &domain.User{
		ID:           user.ID,
		Avatar:       user.Avatar,
		Name:         user.Name,
		Password:     user.Password,
		CreatedAt:    user.CreatedAt,
		RefreshToken: user.RefreshToken,
	}
}
