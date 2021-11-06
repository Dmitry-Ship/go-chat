package application

import (
	"GitHub/go-chat/backend/domain"
	ws "GitHub/go-chat/backend/pkg/websocket"

	"fmt"

	"github.com/google/uuid"
)

type RoomCommandService interface {
	CreateRoom(id uuid.UUID, name string, userId uuid.UUID) error
	JoinRoom(roomId uuid.UUID, userId uuid.UUID) error
	LeaveRoom(roomId uuid.UUID, userId uuid.UUID) error
	DeleteRoom(id uuid.UUID) error
	SendMessage(messageText string, messageType string, roomId uuid.UUID, userId uuid.UUID) error
}

type roomCommandService struct {
	rooms        domain.RoomRepository
	participants domain.ParticipantRepository
	users        domain.UserRepository
	messages     domain.ChatMessageRepository
	hub          ws.HubBroadcaster
}

func NewRoomCommandService(rooms domain.RoomRepository, participants domain.ParticipantRepository, users domain.UserRepository, messages domain.ChatMessageRepository, hub ws.HubBroadcaster) *roomCommandService {
	return &roomCommandService{
		rooms:        rooms,
		users:        users,
		participants: participants,
		messages:     messages,
		hub:          hub,
	}
}

func (s *roomCommandService) CreateRoom(id uuid.UUID, name string, userId uuid.UUID) error {
	room := domain.NewRoom(id, name)
	err := s.rooms.Store(room)

	if err != nil {
		return err
	}

	err = s.JoinRoom(room.Id, userId)

	if err != nil {
		return err
	}

	return nil
}

func (s *roomCommandService) JoinRoom(roomID uuid.UUID, userId uuid.UUID) error {
	err := s.participants.Store(domain.NewParticipant(roomID, userId))

	if err != nil {
		return err
	}

	user, err := s.users.FindByID(userId)
	if err != nil {
		return err
	}

	err = s.SendMessage(fmt.Sprintf(" %s joined", user.Name), "system", roomID, user.Id)

	if err != nil {
		return err
	}

	return nil
}

func (s *roomCommandService) LeaveRoom(roomID uuid.UUID, userId uuid.UUID) error {
	err := s.participants.DeleteByRoomIDAndUserID(roomID, userId)

	if err != nil {
		return err
	}

	user, err := s.users.FindByID(userId)

	if err != nil {
		return err
	}

	err = s.SendMessage(fmt.Sprintf("%s left", user.Name), "system", roomID, user.Id)

	if err != nil {
		return err
	}

	return nil
}

func (s *roomCommandService) DeleteRoom(id uuid.UUID) error {
	message := struct {
		RoomId uuid.UUID `json:"room_id"`
	}{
		RoomId: id,
	}

	err := s.notifyAllParticipants(id, "room_deleted", message)

	if err != nil {
		return err
	}

	err = s.rooms.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *roomCommandService) notifyAllParticipants(roomID uuid.UUID, messageType string, message interface{}) error {
	participants, err := s.participants.FindAllByRoomID(roomID)

	if err != nil {
		return err
	}

	for _, participant := range participants {
		s.hub.BroadcastNotification(messageType, message, participant.UserId)
	}

	return nil
}

func (s *roomCommandService) SendMessage(messageText string, messageType string, roomId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewChatMessage(messageText, messageType, roomId, userId)

	err := s.messages.Store(message)

	if err != nil {
		return err
	}

	user, err := s.users.FindByID(message.UserId)

	if err != nil {
		return err
	}

	fullMessage := MessageFull{
		User:        user,
		ChatMessage: message,
	}

	err = s.notifyAllParticipants(roomId, "message", fullMessage)

	if err != nil {
		return err
	}

	return nil
}
