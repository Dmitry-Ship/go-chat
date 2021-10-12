package application

import (
	"GitHub/go-chat/backend/domain"

	"github.com/google/uuid"
)

type UserService interface {
	GetUser(id uuid.UUID) (*domain.User, error)
	CreateUser(user *domain.User) (*domain.User, error)
}

type userService struct {
	users domain.UserRepository
}

func NewUserService(users domain.UserRepository) *userService {
	return &userService{
		users: users,
	}
}

func (s *userService) GetUser(id uuid.UUID) (*domain.User, error) {
	return s.users.FindByID(id)
}

func (s *userService) CreateUser(user *domain.User) (*domain.User, error) {
	return s.users.Create(user)
}
