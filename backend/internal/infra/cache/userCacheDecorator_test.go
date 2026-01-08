package cache

import (
	"context"
	"testing"
	"time"

	"GitHub/go-chat/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCacheClient struct {
	mock.Mock
}

func (m *MockCacheClient) Get(ctx context.Context, key string) ([]byte, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCacheClient) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCacheClient) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheClient) DeletePattern(ctx context.Context, pattern string) error {
	args := m.Called(ctx, pattern)
	return args.Error(0)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Store(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestUserCacheDecorator_GetByID_CacheHit(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCacheClient)

	userID := uuid.New()
	userName := "testuser"
	userPassword, _ := domain.HashPassword("password123")
	user := &domain.User{
		ID:           userID,
		Name:         userName,
		PasswordHash: userPassword,
		Avatar:       "T",
	}

	userData, _ := SerializeUser(user)

	mockCache.On("Get", mock.Anything, UserKey(userID.String())).Return(userData, nil)

	decorator := NewUserCacheDecorator(mockRepo, mockCache)
	result, err := decorator.GetByID(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.ID)

	mockCache.AssertExpectations(t)
	assert.Equal(t, 0, len(mockRepo.Calls))
}

func TestUserCacheDecorator_GetByID_CacheMiss(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCacheClient)

	userID := uuid.New()
	userName := "testuser"
	userPassword, _ := domain.HashPassword("password123")
	user := &domain.User{
		ID:           userID,
		Name:         userName,
		PasswordHash: userPassword,
		Avatar:       "T",
	}

	mockCache.On("Get", mock.Anything, UserKey(userID.String())).Return(nil, nil)
	mockRepo.On("GetByID", mock.Anything, userID).Return(user, nil)
	mockCache.On("Set", mock.Anything, UserKey(userID.String()), mock.Anything, TTLUser).Return(nil)

	decorator := NewUserCacheDecorator(mockRepo, mockCache)
	result, err := decorator.GetByID(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.ID)

	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestUserCacheDecorator_Store_InvalidatesCache(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCacheClient)

	userID := uuid.New()
	userName := "testuser"
	userPassword, _ := domain.HashPassword("password123")
	user := &domain.User{
		ID:           userID,
		Name:         userName,
		PasswordHash: userPassword,
		Avatar:       "T",
	}

	mockRepo.On("Store", mock.Anything, user).Return(nil)
	mockCache.On("Delete", mock.Anything, UserKey(userID.String())).Return(nil)
	mockCache.On("Delete", mock.Anything, UsernameKey("testuser")).Return(nil)

	decorator := NewUserCacheDecorator(mockRepo, mockCache)
	err := decorator.Store(context.Background(), user)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUserCacheDecorator_Update_InvalidatesCache(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCacheClient)

	userID := uuid.New()
	userName := "testuser"
	userPassword, _ := domain.HashPassword("password123")
	user := &domain.User{
		ID:           userID,
		Name:         userName,
		PasswordHash: userPassword,
		Avatar:       "T",
	}

	mockRepo.On("Update", mock.Anything, user).Return(nil)
	mockCache.On("Delete", mock.Anything, UserKey(userID.String())).Return(nil)
	mockCache.On("Delete", mock.Anything, UsernameKey("testuser")).Return(nil)

	decorator := NewUserCacheDecorator(mockRepo, mockCache)
	err := decorator.Update(context.Background(), user)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
