package postgres

import (
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	repository
}

func NewUserRepository(db *gorm.DB, eventPublisher *infra.EventBus) *userRepository {
	return &userRepository{
		repository: *newRepository(db, eventPublisher),
	}
}

func (r *userRepository) Store(user *domain.User) error {
	return r.store(user, &User{
		ID:           user.ID,
		Avatar:       user.Avatar,
		Name:         user.Name,
		Password:     user.PasswordHash,
		RefreshToken: user.RefreshToken,
	})
}

func (r *userRepository) Update(user *domain.User) error {
	return r.update(user, &User{
		ID:           user.ID,
		Avatar:       user.Avatar,
		Name:         user.Name,
		Password:     user.PasswordHash,
		RefreshToken: user.RefreshToken,
	})
}

func (r *userRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	user := User{}
	err := r.db.Where(&User{ID: id}).First(&user).Error

	if err != nil {
		return nil, fmt.Errorf("get user by id error: %w", err)
	}

	return &domain.User{
		ID:           user.ID,
		Avatar:       user.Avatar,
		Name:         user.Name,
		PasswordHash: user.Password,
		RefreshToken: user.RefreshToken,
	}, nil
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	user := User{}

	err := r.db.Where(&User{Name: username}).First(&user).Error

	if err != nil {
		return nil, fmt.Errorf("find user by username error: %w", err)
	}

	return &domain.User{
		ID:           user.ID,
		Avatar:       user.Avatar,
		Name:         user.Name,
		PasswordHash: user.Password,
		RefreshToken: user.RefreshToken,
	}, nil
}
