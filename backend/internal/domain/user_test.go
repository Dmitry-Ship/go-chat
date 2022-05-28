package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	name, _ := NewUserName("John")
	password, _ := NewUserPassword("12345678", func(p []byte) ([]byte, error) {
		return []byte(p), nil
	})
	user := NewUser(name, password)

	assert.NotNil(t, user.ID)
	assert.NotNil(t, user.Avatar)
	assert.NotNil(t, user.Password, password)
	assert.Equal(t, user.Name, name)
}

func TestNewUsername(t *testing.T) {
	name, err := NewUserName("John")

	assert.Nil(t, err)
	assert.Equal(t, name.String(), "John")
}

func TestNewUsernameEmpty(t *testing.T) {
	_, err := NewUserName("")

	assert.Equal(t, err.Error(), "username is empty")
}

func TestNewUsernameEmptyLong(t *testing.T) {
	name := ""

	for i := 0; i < 101; i++ {
		name += "a"
	}

	_, err := NewUserName(name)

	assert.Equal(t, err.Error(), "username is too long")
}

func TestNewUserPassword(t *testing.T) {
	password, err := NewUserPassword("asdasdasdasdasdasd", func(p []byte) ([]byte, error) {
		return []byte(p), nil
	})

	assert.Nil(t, err)
	assert.Equal(t, password.String(), "asdasdasdasdasdasd")
}

func TestNewUserPasswordEmpty(t *testing.T) {
	_, err := NewUserPassword("", func(p []byte) ([]byte, error) {
		return []byte(p), nil
	})

	assert.Equal(t, err.Error(), "password is empty")
}

func TestNewUserPasswordShort(t *testing.T) {
	_, err := NewUserPassword("123", func(p []byte) ([]byte, error) {
		return []byte(p), nil
	})

	assert.Equal(t, err.Error(), "password is too short")
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
	user := NewUser(name, password)

	user.SetRefreshToken("test")

	assert.Equal(t, user.RefreshToken, "test")
}
