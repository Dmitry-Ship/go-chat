package application

import (
	"GitHub/go-chat/backend/domain"
	ws "GitHub/go-chat/backend/pkg/websocket"

	"fmt"

	"github.com/google/uuid"
)

type RoomCommandService interface {
	CreatePublicRoom(id uuid.UUID, name string, userId uuid.UUID) error
	JoinPublicRoom(roomId uuid.UUID, userId uuid.UUID) error
	LeavePublicRoom(roomId uuid.UUID, userId uuid.UUID) error
	DeleteRoom(id uuid.UUID) error
	SendMessage(messageText string, messageType string, roomId uuid.UUID, userId uuid.UUID) error
}

type roomCommandService struct {
	rooms        domain.RoomRepository
	participants domain.ParticipantRepository
	users        domain.UserRepository
	messages     domain.ChatMessageRepository
	hub          ws.Hub
}

func NewRoomCommandService(rooms domain.RoomRepository, participants domain.ParticipantRepository, users domain.UserRepository, messages domain.ChatMessageRepository, hub ws.Hub) *roomCommandService {
	return &roomCommandService{
		rooms:        rooms,
		users:        users,
		participants: participants,
		messages:     messages,
		hub:          hub,
	}
}

func (s *roomCommandService) CreatePublicRoom(id uuid.UUID, name string, userId uuid.UUID) error {
	room := domain.NewRoom(id, name, false)
	err := s.rooms.Store(room)

	if err != nil {
		return err
	}

	err = s.JoinPublicRoom(room.ID, userId)

	if err != nil {
		return err
	}

	return nil
}

func (s *roomCommandService) JoinPublicRoom(roomID uuid.UUID, userId uuid.UUID) error {
	err := s.participants.Store(domain.NewParticipant(roomID, userId))

	if err != nil {
		return err
	}

	user, err := s.users.FindByID(userId)

	if err != nil {
		return err
	}

	err = s.SendMessage(fmt.Sprintf(" %s joined", user.Name), "system", roomID, user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *roomCommandService) LeavePublicRoom(roomID uuid.UUID, userId uuid.UUID) error {
	err := s.participants.DeleteByRoomIDAndUserID(roomID, userId)

	if err != nil {
		return err
	}

	user, err := s.users.FindByID(userId)

	if err != nil {
		return err
	}

	err = s.SendMessage(fmt.Sprintf("%s left", user.Name), "system", roomID, user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *roomCommandService) DeleteRoom(id uuid.UUID) error {
	notification := ws.OutgoingNotification{
		Type: "room_deleted",
		Payload: struct {
			RoomId uuid.UUID `json:"room_id"`
		}{
			RoomId: id,
		},
	}

	err := s.notifyAllParticipants(id, notification)

	if err != nil {
		return err
	}

	err = s.rooms.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *roomCommandService) notifyAllParticipants(roomID uuid.UUID, notification ws.OutgoingNotification) error {
	participants, err := s.participants.FindAllByRoomID(roomID)

	if err != nil {
		return err
	}

	for _, participant := range participants {
		if notification.Type == "message" {
			message := notification.Payload.(MessageFull)
			message.IsInbound = participant.UserID != message.User.ID
			notification.Payload = message

		}

		s.hub.BroadcastToClients(notification, participant.UserID)
	}

	return nil
}

func (s *roomCommandService) SendMessage(messageText string, messageType string, roomId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewChatMessage(messageText, messageType, roomId, userId)

	err := s.messages.Store(message)

	if err != nil {
		return err
	}

	user, err := s.users.FindByID(message.UserID)

	if err != nil {
		return err
	}

	notification := ws.OutgoingNotification{
		Type: "message",
		Payload: MessageFull{
			User:        user,
			ChatMessage: message,
		},
	}

	err = s.notifyAllParticipants(roomId, notification)

	if err != nil {
		return err
	}

	return nil
}
