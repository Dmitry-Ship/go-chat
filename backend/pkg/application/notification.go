package application

import (
	"GitHub/go-chat/backend/domain"
)

type Notification struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type NotificationService interface {
	NewNotification(notificationType string, data interface{}) Notification
	AddWSClient(userID int32) chan Notification
	SendToAllUserWSClients(userID int32, message Notification)
	Run()
}

type notificationService struct {
	participants  domain.ParticipantRepository
	Broadcast     chan *MessageFull
	userWSClients map[int32][]chan Notification
}

func NewNotificationService(participants domain.ParticipantRepository) *notificationService {
	return &notificationService{
		participants:  participants,
		Broadcast:     make(chan *MessageFull, 1024),
		userWSClients: make(map[int32][]chan Notification),
	}
}

func (s *notificationService) Run() {
	for {
		select {
		case message := <-s.Broadcast:
			notification := s.NewNotification("message", message)
			participants, err := s.participants.FindByRoomID(message.RoomId)

			if err != nil {
				continue
			}

			for _, participant := range participants {
				s.SendToAllUserWSClients(participant.UserId, notification)
			}

		}
	}
}

func (s *notificationService) AddWSClient(userID int32) chan Notification {
	channel := make(chan Notification, 1024)

	s.userWSClients[userID] = append(s.userWSClients[userID], channel)
	return channel
}

func (s *notificationService) SendToAllUserWSClients(userID int32, message Notification) {
	clients := s.userWSClients[userID]

	for _, client := range clients {
		client <- message
	}
}

func (c *notificationService) NewNotification(notificationType string, data interface{}) Notification {
	return Notification{
		Type: notificationType,
		Data: data,
	}
}
