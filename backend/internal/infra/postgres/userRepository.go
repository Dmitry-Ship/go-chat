package postgres

import (
	"GitHub/go-chat/backend/internal/domain"

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
	return r.db.Create(ToUserPersistence(user)).Error
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(ToUserPersistence(user)).Error
}

func (r *userRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	user := User{}
	err := r.db.Where("id = ?", id).First(&user).Error

	return ToUserDomain(&user), err
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	user := User{}

	err := r.db.Where("name = ?", username).First(&user).Error

	return ToUserDomain(&user), err
}
