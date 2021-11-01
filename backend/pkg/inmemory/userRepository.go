package inmemory

import (
	"GitHub/go-chat/backend/domain"
	"errors"

	"github.com/google/uuid"
)

type userRepository struct {
	users         map[uuid.UUID]*domain.User
	refreshTokens map[uuid.UUID]string
}

func NewUserRepository() *userRepository {
	return &userRepository{
		users:         make(map[uuid.UUID]*domain.User),
		refreshTokens: make(map[uuid.UUID]string),
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

func (r *userRepository) StoreRefreshToken(userID uuid.UUID, refreshToken string) error {
	r.refreshTokens[userID] = refreshToken
	return nil
}

func (r *userRepository) GetRefreshTokenByUserId(userID uuid.UUID) (string, error) {
	refreshToken, ok := r.refreshTokens[userID]
	if !ok {
		return "", errors.New("refresh token not found")
	}
	return refreshToken, nil
}

func (r *userRepository) DeleteRefreshToken(userID uuid.UUID) error {
	delete(r.refreshTokens, userID)
	return nil
}
