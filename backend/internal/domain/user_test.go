package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	name := "John"
	password, _ := HashPassword("Password123!")
	userID := uuid.New()
	user := NewUser(userID, name, password)

	assert.Equal(t, user.ID, userID)
	assert.Equal(t, user.Avatar, string(name[0]))
	assert.NotNil(t, user.PasswordHash, password)
	assert.Equal(t, user.Name, name)
}

func TestValidateUsername(t *testing.T) {
	name := "John"
	err := ValidateUsername(name)

	assert.Nil(t, err)
	assert.Equal(t, name, "John")
}

func TestValidateUsernameErrors(t *testing.T) {
	type testCase struct {
		name        string
		expectedErr error
	}

	longName := ""

	for i := 0; i < 101; i++ {
		longName += "a"
	}

	testCases := []testCase{
		{
			name:        "",
			expectedErr: errors.New("username is empty"),
		}, {
			name:        longName,
			expectedErr: errors.New("username too long"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateUsername(tc.name)

			assert.Equal(t, err, tc.expectedErr)
		})
	}
}

func TestHashPassword(t *testing.T) {
	password, err := HashPassword("Password123!")

	assert.Nil(t, err)
	assert.NotEqual(t, password, "Password123!")
}

func TestHashPasswordErrors(t *testing.T) {
	type testCase struct {
		password    string
		expectedErr error
	}

	testCases := []testCase{
		{
			password:    "",
			expectedErr: errors.New("password must be at least 8 characters"),
		}, {
			password:    "123",
			expectedErr: errors.New("password must be at least 8 characters"),
		}, {
			password:    "Password123",
			expectedErr: errors.New("password must contain uppercase, lowercase, digit, and special character"),
		}, {
			password:    "password!",
			expectedErr: errors.New("password must contain uppercase, lowercase, digit, and special character"),
		}, {
			password:    "PASSWORD1!",
			expectedErr: errors.New("password must contain uppercase, lowercase, digit, and special character"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.password, func(t *testing.T) {
			_, err := HashPassword(tc.password)

			assert.Equal(t, err, tc.expectedErr)
		})
	}
}

func TestComparePassword(t *testing.T) {
	hashed, _ := HashPassword("Password123!")

	err := ComparePassword(hashed, "Password123!")

	assert.Nil(t, err)
}

func TestComparePasswordFail(t *testing.T) {
	hashed, _ := HashPassword("Password123!")

	err := ComparePassword(hashed, "wrongpassword")

	assert.NotNil(t, err)
}

func TestSetRefreshToken(t *testing.T) {
	password, _ := HashPassword("Password123!")
	name := "John"
	userID := uuid.New()
	user := NewUser(userID, name, password)

	user.SetRefreshToken("test")

	assert.Equal(t, user.RefreshToken, "test")
}
