package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type UserRepository interface {
	GenericRepository[*User]
	GetByID(id uuid.UUID) (*User, error)
	FindByUsername(username string) (*User, error)
}

type userName struct {
	name string
}

func (n *userName) String() string {
	return n.name
}

func NewUserName(name string) (userName, error) {
	if name == "" {
		return userName{}, fmt.Errorf("username is empty")
	}

	if len(name) > 100 {
		return userName{}, fmt.Errorf("username is too long")
	}

	return userName{
		name: name,
	}, nil
}

type userPassword struct {
	password string
}

func (n userPassword) String() string {
	return n.password
}

func (n *userPassword) Compare(password userPassword, compare func(p1 []byte, p2 []byte) error) error {
	err := compare([]byte(n.password), []byte(password.String()))

	if err != nil {
		return fmt.Errorf("password is incorrect")
	}

	return nil
}

func NewUserPassword(password string, hash func(p []byte) ([]byte, error)) (userPassword, error) {
	if password == "" {
		return userPassword{}, fmt.Errorf("password is empty")
	}

	if len(password) < 8 {
		return userPassword{}, fmt.Errorf("password is too short")
	}

	bytes, err := hash([]byte(password))

	if err != nil {
		return userPassword{}, fmt.Errorf("hash password error: %w", err)
	}

	hashedPassword := string(bytes)

	return userPassword{
		password: hashedPassword,
	}, nil
}

type User struct {
	aggregate
	ID           uuid.UUID
	Avatar       string
	Name         userName
	Password     userPassword
	RefreshToken string
}

func NewUser(userID uuid.UUID, username userName, password userPassword) *User {
	return &User{
		ID:       userID,
		Avatar:   string(username.String()[0]),
		Name:     username,
		Password: password,
	}
}

func (u *User) SetRefreshToken(refreshToken string) {
	u.RefreshToken = refreshToken
}
