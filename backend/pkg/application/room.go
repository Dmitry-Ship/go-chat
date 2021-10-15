package application

import (
	"GitHub/go-chat/backend/domain"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type RoomService interface {
	CreateRoom(name string, userId uuid.UUID) (*domain.Room, error)
	GetRoom(id uuid.UUID) (*domain.Room, error)
	HasJoined(userId uuid.UUID, roomId uuid.UUID) bool
	GetRooms() ([]*domain.Room, error)
	JoinRoom(userId uuid.UUID, roomId uuid.UUID) (*domain.Participant, error)
	LeaveRoom(userId uuid.UUID, roomId uuid.UUID) error
	DeleteRoom(id uuid.UUID) error
	SendMessage(messageText string, messageType string, roomId uuid.UUID, userId uuid.UUID) (*MessageFull, error)
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

func (s *roomService) CreateRoom(name string, userId uuid.UUID) (*domain.Room, error) {
	room := domain.NewRoom(name)
	newRoom, err := s.rooms.Create(room)

	if err != nil {
		return nil, err
	}

	_, err = s.JoinRoom(userId, room.Id)

	if err != nil {
		return nil, err
	}

	return newRoom, nil
}

func (s *roomService) GetRoom(id uuid.UUID) (*domain.Room, error) {
	return s.rooms.FindByID(id)
}

func (s *roomService) GetRooms() ([]*domain.Room, error) {
	return s.rooms.FindAll()
}

func (s *roomService) JoinRoom(userId uuid.UUID, roomID uuid.UUID) (*domain.Participant, error) {
	if s.HasJoined(userId, roomID) {
		return nil, errors.New("user already joined")
	}

	newParticipant, err := s.participants.Create(domain.NewParticipant(roomID, userId))

	if err != nil {
		return nil, err
	}

	user, err := s.users.FindByID(userId)
	if err != nil {
		return nil, err
	}

	s.SendMessage(fmt.Sprintf(" %s joined", user.Name), "system", roomID, user.Id)

	return newParticipant, nil
}

func (s *roomService) LeaveRoom(userId uuid.UUID, roomID uuid.UUID) error {
	err := s.participants.DeleteByRoomIDAndUserID(roomID, userId)

	if err != nil {
		return err
	}

	user, err := s.users.FindByID(userId)

	if err != nil {
		return err
	}

	s.SendMessage(fmt.Sprintf("%s left", user.Name), "system", roomID, user.Id)
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

	err = s.participants.DeleteAllByRoomID(id)

	if err != nil {
		return err
	}

	err = s.rooms.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *roomService) HasJoined(userId uuid.UUID, roomID uuid.UUID) bool {
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

func (s *roomService) SendMessage(messageText string, messageType string, roomId uuid.UUID, userId uuid.UUID) (*MessageFull, error) {
	message := domain.NewChatMessage(messageText, messageType, roomId, userId)

	newMessage, err := s.messages.Create(message)

	if err != nil {
		return nil, err
	}

	fullMessage, err := s.makeMessageFull(newMessage)

	if err != nil {
		return nil, err
	}

	err = s.NotifyAllParticipants(roomId, "message", fullMessage)

	if err != nil {
		return nil, err
	}

	return fullMessage, nil
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
