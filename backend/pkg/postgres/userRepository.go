package postgres

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/readModel"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Store(user *domain.User) error {
	err := r.db.Create(ToUserPersistence(user)).Error

	return err
}

func (r *userRepository) GetUserByID(id uuid.UUID) (*readModel.UserDTO, error) {
	user := User{}
	err := r.db.Where("id = ?", id).First(&user).Error

	return toUserDTO(&user), err
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	user := User{}

	err := r.db.Where("name = ?", username).First(&user).Error

	return ToUserDomain(&user), err
}

func (r *userRepository) StoreRefreshToken(userID uuid.UUID, refreshToken string) error {
	err := r.db.Where("id = ?", userID).First(&User{}).Update("refresh_token", refreshToken).Error

	return err
}

func (r *userRepository) GetRefreshTokenByUserID(userID uuid.UUID) (string, error) {
	user := User{}
	err := r.db.Where("id = ?", userID).First(&user).Error

	return user.RefreshToken, err
}

func (r *userRepository) DeleteRefreshToken(userID uuid.UUID) error {
	user := User{}
	err := r.db.Where("id = ?", userID).First(&user).Update("refresh_token", nil).Error

	return err
}

func (r *userRepository) FindAll() ([]*readModel.UserDTO, error) {
	users := []*User{}
	err := r.db.Limit(50).Find(&users).Error

	dtoUsers := make([]*readModel.UserDTO, len(users))

	for i, user := range users {
		dtoUsers[i] = toUserDTO(user)
	}

	return dtoUsers, err
}
