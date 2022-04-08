package postgres

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/mappers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	users *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		users: db,
	}
}

func (r *userRepository) Store(user *domain.User) error {

	err := r.users.Create(mappers.ToUserPersistence(user)).Error

	return err
}

func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	user := mappers.UserPersistence{}
	err := r.users.Where("id = ?", id).First(&user).Error

	return mappers.ToUserDomain(&user), err
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	user := mappers.UserPersistence{}

	err := r.users.Where("name = ?", username).First(&user).Error

	return mappers.ToUserDomain(&user), err
}

func (r *userRepository) StoreRefreshToken(userID uuid.UUID, refreshToken string) error {
	user := mappers.UserPersistence{}
	err := r.users.Where("id = ?", userID).First(&user).Update("refresh_token", refreshToken).Error

	return err
}

func (r *userRepository) GetRefreshTokenByUserID(userID uuid.UUID) (string, error) {
	user := mappers.UserPersistence{}
	err := r.users.Where("id = ?", userID).First(&user).Error

	return user.RefreshToken, err
}

func (r *userRepository) DeleteRefreshToken(userID uuid.UUID) error {
	user := mappers.UserPersistence{}
	err := r.users.Where("id = ?", userID).First(&user).Update("refresh_token", nil).Error

	return err

}

func (r *userRepository) FindAll() ([]*domain.User, error) {
	users := []*mappers.UserPersistence{}
	err := r.users.Limit(50).Find(&users).Error

	domainUsers := make([]*domain.User, len(users))

	for i, user := range users {
		domainUsers[i] = mappers.ToUserDomain(user)
	}

	return domainUsers, err
}
