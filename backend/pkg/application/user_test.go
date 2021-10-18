package application

import (
	"GitHub/go-chat/backend/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct{}

func (mr *mockUserRepo) Store(user *domain.User) error {
	return nil
}

func (mr *mockUserRepo) FindByID(id uuid.UUID) (*domain.User, error) {
	return nil, nil
}

var mockedUserRepo = &mockUserRepo{}

func TestCreateUser(t *testing.T) {
	// Given
	userService := NewUserService(mockedUserRepo)
	user := domain.NewUser()
	// When
	err := userService.CreateUser(user)

	// Then
	assert.Nil(t, err)
}
