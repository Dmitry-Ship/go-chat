package postgres

import (
	"GitHub/go-chat/backend/domain"

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

	err := r.users.Create(domain.ToUserPersistence(user)).Error

	return err
}

func (r *userRepository) FindByID(id uuid.UUID) (*domain.UserDTO, error) {
	user := domain.UserPersistence{}
	err := r.users.Where("id = ?", id).First(&user).Error

	return domain.ToUserDTO(&user), err
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	user := domain.UserPersistence{}

	err := r.users.Where("name = ?", username).First(&user).Error

	return domain.ToUserDomain(&user), err
}

func (r *userRepository) StoreRefreshToken(userID uuid.UUID, refreshToken string) error {
	user := domain.UserPersistence{}
	err := r.users.Where("id = ?", userID).First(&user).Update("refresh_token", refreshToken).Error

	return err
}

func (r *userRepository) GetRefreshTokenByUserID(userID uuid.UUID) (string, error) {
	user := domain.UserPersistence{}
	err := r.users.Where("id = ?", userID).First(&user).Error

	return user.RefreshToken, err
}

func (r *userRepository) DeleteRefreshToken(userID uuid.UUID) error {
	user := domain.UserPersistence{}
	err := r.users.Where("id = ?", userID).First(&user).Update("refresh_token", nil).Error

	return err

}

func (r *userRepository) FindAll() ([]*domain.UserDTO, error) {
	users := []*domain.UserPersistence{}
	err := r.users.Limit(50).Find(&users).Error

	dtoUsers := make([]*domain.UserDTO, len(users))

	for i, user := range users {
		dtoUsers[i] = domain.ToUserDTO(user)
	}

	return dtoUsers, err
}
