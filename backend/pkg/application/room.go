package application

import (
	"GitHub/go-chat/backend/domain"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type roomData struct {
	Room   domain.Room `json:"room"`
	Joined bool        `json:"joined"`
}

type MessageFull struct {
	User *domain.User `json:"user"`
	*domain.ChatMessage
}

type RoomService interface {
	CreateRoom(id uuid.UUID, name string, userId uuid.UUID) error
	JoinRoom(roomId uuid.UUID, userId uuid.UUID) error
	LeaveRoom(roomId uuid.UUID, userId uuid.UUID) error
	DeleteRoom(id uuid.UUID) error
	SendMessage(messageText string, messageType string, roomId uuid.UUID, userId uuid.UUID) error
	GetRoom(roomId uuid.UUID, userId uuid.UUID) (*roomData, error)
	GetRooms() ([]*domain.Room, error)
	GetRoomMessages(roomId uuid.UUID) ([]*MessageFull, error)
}

type roomService struct {
	rooms        domain.RoomRepository
	participants domain.ParticipantRepository
	users        domain.UserRepository
	messages     domain.ChatMessageRepository
	hub          HubBroadcaster
}

func NewRoomService(rooms domain.RoomRepository, participants domain.ParticipantRepository, users domain.UserRepository, messages domain.ChatMessageRepository, hub HubBroadcaster) *roomService {
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

func (s *roomService) GetRoom(roomId uuid.UUID, userId uuid.UUID) (*roomData, error) {

	room, err := s.rooms.FindByID(roomId)

	if err != nil {
		return nil, err
	}

	data := roomData{
		Room:   *room,
		Joined: s.hasJoined(roomId, userId),
	}

	return &data, nil
}

func (s *roomService) GetRooms() ([]*domain.Room, error) {
	return s.rooms.FindAll()
}

func (s *roomService) JoinRoom(roomID uuid.UUID, userId uuid.UUID) error {
	if s.hasJoined(roomID, userId) {
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

func (s *roomService) hasJoined(roomID uuid.UUID, userId uuid.UUID) bool {
	_, err := s.participants.FindByRoomIDAndUserID(roomID, userId)

	return err == nil
}

func (s *roomService) notifyAllParticipants(roomID uuid.UUID, messageType string, message interface{}) error {
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

	err = s.notifyAllParticipants(roomId, "message", fullMessage)

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
