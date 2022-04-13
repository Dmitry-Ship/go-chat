package postgres

import (
	"GitHub/go-chat/backend/pkg/domain"
	"GitHub/go-chat/backend/pkg/readModel"
)

func toUserDTO(user *User) *readModel.UserDTO {
	return &readModel.UserDTO{
		ID:     user.ID,
		Avatar: user.Avatar,
		Name:   user.Name,
	}
}

func ToUserPersistence(user *domain.User) *User {
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
