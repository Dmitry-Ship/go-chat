package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
)

func toUserPersistence(user *domain.User) *User {
	return &User{
		ID:           user.ID,
		Avatar:       user.Avatar,
		Name:         user.Name,
		Password:     user.Password,
		RefreshToken: user.RefreshToken,
	}
}

func toUserDomain(user *User) *domain.User {
	return &domain.User{
		ID:           user.ID,
		Avatar:       user.Avatar,
		Name:         user.Name,
		Password:     user.Password,
		RefreshToken: user.RefreshToken,
	}
}
