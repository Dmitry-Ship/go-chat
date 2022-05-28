package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
)

func toUserPersistence(user *domain.User) *User {
	return &User{
		ID:           user.ID,
		Avatar:       user.Avatar,
		Name:         user.Name.String(),
		Password:     user.Password.String(),
		RefreshToken: user.RefreshToken,
	}
}

func toUserDomain(user *User) *domain.User {
	name, err := domain.NewUserName(user.Name)

	if err != nil {
		panic(err)
	}

	password, err := domain.NewUserPassword(user.Password, func(p []byte) ([]byte, error) {
		return []byte(p), nil
	})

	if err != nil {
		panic(err)
	}

	return &domain.User{
		ID:           user.ID,
		Avatar:       user.Avatar,
		Name:         name,
		Password:     password,
		RefreshToken: user.RefreshToken,
	}
}
