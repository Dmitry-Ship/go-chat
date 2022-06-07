package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	name, _ := NewUserName("John")
	password, _ := NewUserPassword("12345678", func(p []byte) ([]byte, error) {
		return []byte(p), nil
	})
	userID := uuid.New()
	user := NewUser(userID, name, password)

	assert.Equal(t, user.ID, userID)
	assert.Equal(t, user.Avatar, string(name.String()[0]))
	assert.NotNil(t, user.Password, password)
	assert.Equal(t, user.Name, name)
}

func TestNewUsername(t *testing.T) {
	name, err := NewUserName("John")

	assert.Nil(t, err)
	assert.Equal(t, name.String(), "John")
}

func TestNewUsernameErrors(t *testing.T) {
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
			expectedErr: errors.New("username is too long"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewUserName(tc.name)

			assert.Equal(t, err, tc.expectedErr)
		})
	}
}

func TestNewUserPassword(t *testing.T) {
	password, err := NewUserPassword("asdasdasdasdasdasd", func(p []byte) ([]byte, error) {
		return []byte(p), nil
	})

	assert.Nil(t, err)
	assert.Equal(t, password.String(), "asdasdasdasdasdasd")
}

func TestNewUserPasswordErrors(t *testing.T) {
	type testCase struct {
		password    string
		expectedErr error
	}

	testCases := []testCase{
		{
			password:    "",
			expectedErr: errors.New("password is empty"),
		}, {
			password:    "123",
			expectedErr: errors.New("password is too short"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.password, func(t *testing.T) {
			_, err := NewUserPassword(tc.password, func(p []byte) ([]byte, error) {
				return []byte(p), nil
			})

			assert.Equal(t, err, tc.expectedErr)
		})
	}
}

func TestNewUserPasswordCompare(t *testing.T) {
	password, _ := NewUserPassword("12345678", func(p []byte) ([]byte, error) {
		return []byte(p), nil
	})

	anotherPassword, _ := NewUserPassword("12345678", func(p []byte) ([]byte, error) {
		return []byte(p), nil
	})

	err := password.Compare(anotherPassword, func(p1 []byte, p2 []byte) error {
		return nil
	})

	assert.Nil(t, err)
}

func TestNewUserPasswordCompareFail(t *testing.T) {
	password, _ := NewUserPassword("12345678", func(p []byte) ([]byte, error) {
		return []byte(p), nil
	})

	anotherPassword, _ := NewUserPassword("12345678123", func(p []byte) ([]byte, error) {
		return []byte(p), nil
	})

	err := password.Compare(anotherPassword, func(p1 []byte, p2 []byte) error {
		return errors.New("")
	})

	assert.Equal(t, err.Error(), "password is incorrect")
}

func TestSetRefreshToken(t *testing.T) {
	password, _ := NewUserPassword("12345678", func(p []byte) ([]byte, error) {
		return []byte(p), nil
	})
	name, _ := NewUserName("John")
	userID := uuid.New()
	user := NewUser(userID, name, password)

	user.SetRefreshToken("test")

	assert.Equal(t, user.RefreshToken, "test")
}
