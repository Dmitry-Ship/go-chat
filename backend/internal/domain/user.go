package domain

import (
	"errors"

	"github.com/google/uuid"
)

type UserRepository interface {
	Store(user *User) error
	Update(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	FindByUsername(username string) (*User, error)
}

type userName struct {
	name string
}

func (n *userName) String() string {
	return n.name
}

func NewUserName(name string) (*userName, error) {
	if name == "" {
		return nil, errors.New("username is empty")
	}

	if len(name) > 100 {
		return nil, errors.New("username is too long")
	}

	return &userName{
		name: name,
	}, nil
}

type userPassword struct {
	password string
}

func (n *userPassword) String() string {
	return n.password
}

func (n *userPassword) Compare(password *userPassword, compare func(p1 []byte, p2 []byte) error) error {
	err := compare([]byte(n.password), []byte(password.String()))

	if err != nil {
		return errors.New("password is incorrect")
	}

	return nil
}

func NewUserPassword(password string, hash func(p []byte) ([]byte, error)) (*userPassword, error) {
	if password == "" {
		return nil, errors.New("password is empty")
	}

	if len(password) < 8 {
		return nil, errors.New("password is too short")
	}

	bytes, err := hash([]byte(password))

	if err != nil {
		return nil, err
	}

	hashedPassword := string(bytes)

	return &userPassword{
		password: hashedPassword,
	}, nil
}

type User struct {
	aggregate
	ID           uuid.UUID
	Avatar       string
	Name         *userName
	Password     *userPassword
	RefreshToken string
}

func NewUser(username *userName, password *userPassword) *User {
	return &User{
		ID:       uuid.New(),
		Avatar:   string(username.String()[0]),
		Name:     username,
		Password: password,
	}
}

func (u *User) SetRefreshToken(refreshToken string) {
	u.RefreshToken = refreshToken
}
