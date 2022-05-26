package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	repository
}

func NewUserRepository(db *gorm.DB, eventPublisher infra.EventPublisher) *userRepository {
	return &userRepository{
		repository: *newRepository(db, eventPublisher),
	}
}

func (r *userRepository) Store(user *domain.User) error {
	err := r.db.Create(ToUserPersistence(user)).Error

	if err != nil {
		return err
	}

	r.dispatchEvents(user)

	return nil
}

func (r *userRepository) Update(user *domain.User) error {
	err := r.db.Save(ToUserPersistence(user)).Error

	if err != nil {
		return err
	}

	r.dispatchEvents(user)

	return nil
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
