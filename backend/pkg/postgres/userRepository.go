package postgres

import (
	"GitHub/go-chat/backend/pkg/domain"
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

func (r *userRepository) Update(user *domain.User) error {
	err := r.db.Save(ToUserPersistence(user)).Error

	return err
}

func (r *userRepository) GetUserByID(id uuid.UUID) (*readModel.UserDTO, error) {
	user := User{}
	err := r.db.Where("id = ?", id).First(&user).Error

	return toUserDTO(&user), err
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

func (r *userRepository) FindContacts(userID uuid.UUID) ([]*readModel.ContactDTO, error) {
	users := []*User{}
	err := r.db.Limit(50).Where("id <> ?", userID).Find(&users).Error

	dtoContacts := make([]*readModel.ContactDTO, len(users))

	for i, user := range users {
		dtoContacts[i] = toContactDTO(user)
	}

	return dtoContacts, err
}
