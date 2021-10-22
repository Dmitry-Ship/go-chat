package inmemory

import (
	"GitHub/go-chat/backend/domain"
	"errors"

	"github.com/google/uuid"
)

type userRepository struct {
	users map[uuid.UUID]*domain.User
}

func NewUserRepository() *userRepository {
	return &userRepository{
		users: make(map[uuid.UUID]*domain.User),
	}
}

func (r *userRepository) Store(user *domain.User) error {
	r.users[user.Id] = user
	return nil
}

func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	for _, user := range r.users {
		if user.Name == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}
