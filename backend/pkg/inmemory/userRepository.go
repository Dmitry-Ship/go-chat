package inmemory

import (
	"GitHub/go-chat/backend/domain"
	"errors"
)

type userRepository struct {
	users map[int32]*domain.User
}

func NewUserRepository() domain.UserRepository {
	return &userRepository{
		users: make(map[int32]*domain.User),
	}
}

func (r *userRepository) FindByID(id int32) (*domain.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepository) FindByName(name string) (*domain.User, error) {
	for _, user := range r.users {
		if user.Name == name {
			return user, nil
		}
	}
	return nil, errors.New("not found")
}

func (r *userRepository) FindAll() ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) Create(user *domain.User) (*domain.User, error) {
	r.users[user.Id] = user
	return user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	_, ok := r.users[user.Id]
	if !ok {
		return errors.New("not found")
	}
	r.users[user.Id] = user
	return nil
}

func (r *userRepository) Delete(id int32) error {
	_, ok := r.users[id]
	if !ok {
		return errors.New("not found")
	}
	delete(r.users, id)
	return nil
}
