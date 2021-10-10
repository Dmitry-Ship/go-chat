package application

import (
	"GitHub/go-chat/backend/domain"
)

type NotificationService interface {
	Run()
}

type notificationService struct {
	participants domain.ParticipantRepository
	Broadcast    chan *MessageFull
	userService  UserService
}

func NewNotificationService(participants domain.ParticipantRepository, userService UserService) *notificationService {
	return &notificationService{participants: participants, userService: userService, Broadcast: make(chan *MessageFull, 1024)}
}

func (s *notificationService) Run() {
	for {
		select {
		case message := <-s.Broadcast:
			notification := s.userService.NewNotification("message", message)
			participants, err := s.participants.FindByRoomID(message.RoomId)

			if err != nil {
				continue
			}

			for _, participant := range participants {
				s.userService.SendToAllUserWSClients(participant.UserId, notification)
			}

		}
	}
}
