package application

import (
	"GitHub/go-chat/backend/domain"
)

type UserService interface {
	SendToAllUserWSClients(userID int32, message Notification)
	NewNotification(notificationType string, data interface{}) Notification
	GetUser(id int32) (*domain.User, error)
	AddWSClient(userID int32) chan Notification
	CreateUser(user *domain.User) (*domain.User, error)
}

type userService struct {
	users         domain.UserRepository
	userWSClients map[int32][]chan Notification
}

type Notification struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewUserService(users domain.UserRepository) *userService {
	return &userService{users: users, userWSClients: make(map[int32][]chan Notification)}
}

func (s *userService) GetUser(id int32) (*domain.User, error) {
	return s.users.FindByID(id)
}

func (s *userService) CreateUser(user *domain.User) (*domain.User, error) {
	return s.users.Create(user)
}

func (s *userService) AddWSClient(userID int32) chan Notification {
	channel := make(chan Notification, 1024)

	s.userWSClients[userID] = append(s.userWSClients[userID], channel)
	return channel
}

func (s *userService) SendToAllUserWSClients(userID int32, message Notification) {
	clients := s.userWSClients[userID]

	for _, client := range clients {
		client <- message
	}
}

func (c *userService) NewNotification(notificationType string, data interface{}) Notification {
	return Notification{
		Type: notificationType,
		Data: data,
	}
}
