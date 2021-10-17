package application

import (
	"GitHub/go-chat/backend/domain"
)

type UserService interface {
	CreateUser(user *domain.User) error
}

type userService struct {
	users domain.UserRepository
}

func NewUserService(users domain.UserRepository) *userService {
	return &userService{
		users: users,
	}
}

func (s *userService) CreateUser(user *domain.User) error {
	return s.users.Store(user)
}
