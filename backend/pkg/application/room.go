package application

import (
	"GitHub/go-chat/backend/domain"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type RoomService interface {
	CreateRoom(id uuid.UUID, name string, userId uuid.UUID) error
	HasJoined(roomId uuid.UUID, userId uuid.UUID) bool
	JoinRoom(roomId uuid.UUID, userId uuid.UUID) error
	LeaveRoom(roomId uuid.UUID, userId uuid.UUID) error
	DeleteRoom(id uuid.UUID) error
	SendMessage(messageText string, messageType string, roomId uuid.UUID, userId uuid.UUID) error
	NotifyAllParticipants(roomId uuid.UUID, messageType string, message interface{}) error
	GetRoom(id uuid.UUID) (*domain.Room, error)
	GetRooms() ([]*domain.Room, error)
	GetRoomMessages(roomId uuid.UUID) ([]*MessageFull, error)
}

type MessageFull struct {
	User *domain.User `json:"user"`
	*domain.ChatMessage
}

type roomService struct {
	rooms        domain.RoomRepository
	participants domain.ParticipantRepository
	users        domain.UserRepository
	messages     domain.ChatMessageRepository
	hub          Hub
}

func NewRoomService(rooms domain.RoomRepository, participants domain.ParticipantRepository, users domain.UserRepository, messages domain.ChatMessageRepository, hub Hub) *roomService {
	return &roomService{
		rooms:        rooms,
		users:        users,
		participants: participants,
		messages:     messages,
		hub:          hub,
	}
}

func (s *roomService) CreateRoom(id uuid.UUID, name string, userId uuid.UUID) error {
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

func (s *roomService) GetRoom(id uuid.UUID) (*domain.Room, error) {
	return s.rooms.FindByID(id)
}

func (s *roomService) GetRooms() ([]*domain.Room, error) {
	return s.rooms.FindAll()
}

func (s *roomService) JoinRoom(roomID uuid.UUID, userId uuid.UUID) error {
	if s.HasJoined(roomID, userId) {
		return errors.New("user already joined")
	}

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

func (s *roomService) LeaveRoom(roomID uuid.UUID, userId uuid.UUID) error {
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

func (s *roomService) DeleteRoom(id uuid.UUID) error {
	message := struct {
		RoomId uuid.UUID `json:"room_id"`
	}{
		RoomId: id,
	}

	err := s.NotifyAllParticipants(id, "room_deleted", message)

	if err != nil {
		return err
	}

	err = s.rooms.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *roomService) HasJoined(roomID uuid.UUID, userId uuid.UUID) bool {
	_, err := s.participants.FindByRoomIDAndUserID(roomID, userId)

	return err == nil
}

func (s *roomService) NotifyAllParticipants(roomID uuid.UUID, messageType string, message interface{}) error {
	participants, err := s.participants.FindAllByRoomID(roomID)

	if err != nil {
		return err
	}

	for _, participant := range participants {
		s.hub.BroadcastNotification(messageType, message, participant.UserId)
	}

	return nil
}

func (s *roomService) SendMessage(messageText string, messageType string, roomId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewChatMessage(messageText, messageType, roomId, userId)

	err := s.messages.Store(message)

	if err != nil {
		return err
	}

	fullMessage, err := s.makeMessageFull(message)

	if err != nil {
		return err
	}

	err = s.NotifyAllParticipants(roomId, "message", fullMessage)

	if err != nil {
		return err
	}

	return nil
}

func (s *roomService) GetRoomMessages(roomId uuid.UUID) ([]*MessageFull, error) {
	messages, err := s.messages.FindAllByRoomID(roomId)

	if err != nil {
		return nil, err
	}

	var messagesFull []*MessageFull

	for _, message := range messages {
		messageFull, err := s.makeMessageFull(message)

		if err != nil {
			return nil, err
		}

		messagesFull = append(messagesFull, messageFull)
	}

	return messagesFull, nil

}

func (s *roomService) makeMessageFull(message *domain.ChatMessage) (*MessageFull, error) {
	user, err := s.users.FindByID(message.UserId)

	if err != nil {
		return nil, err
	}

	m := MessageFull{
		User:        user,
		ChatMessage: message,
	}

	return &m, nil

}
