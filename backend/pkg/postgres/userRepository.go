package postgres

import (
	"GitHub/go-chat/backend/domain"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	users         *gorm.DB
	refreshTokens map[uuid.UUID]string
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		users:         db,
		refreshTokens: make(map[uuid.UUID]string),
	}
}

func (r *userRepository) Store(user *domain.User) error {
	err := r.users.Create(&user).Error

	return err
}

func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	user := domain.User{}
	err := r.users.Where("id = ?", id).First(&user).Error

	return &user, err
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	user := domain.User{}

	err := r.users.Where("name = ?", username).First(&user).Error

	return &user, err
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
