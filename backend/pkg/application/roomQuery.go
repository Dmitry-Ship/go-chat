package application

import (
	"GitHub/go-chat/backend/domain"

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

type RoomQueryService interface {
	GetRoom(roomId uuid.UUID, userId uuid.UUID) (roomData, error)
	GetRooms() ([]*domain.Room, error)
	GetRoomMessages(roomId uuid.UUID) ([]MessageFull, error)
}

type roomQueryService struct {
	rooms        domain.RoomRepository
	participants domain.ParticipantRepository
	users        domain.UserRepository
	messages     domain.ChatMessageRepository
}

func NewRoomQueryService(rooms domain.RoomRepository, participants domain.ParticipantRepository, users domain.UserRepository, messages domain.ChatMessageRepository) *roomQueryService {
	return &roomQueryService{
		rooms:        rooms,
		users:        users,
		participants: participants,
		messages:     messages,
	}
}

func (s *roomQueryService) GetRoom(roomId uuid.UUID, userId uuid.UUID) (roomData, error) {

	room, err := s.rooms.FindByID(roomId)

	if err != nil {
		return roomData{}, err
	}

	data := roomData{
		Room:   *room,
		Joined: s.hasJoined(roomId, userId),
	}

	return data, nil
}

func (s *roomQueryService) GetRooms() ([]*domain.Room, error) {
	return s.rooms.FindAll()
}

func (s *roomQueryService) hasJoined(roomID uuid.UUID, userId uuid.UUID) bool {
	_, err := s.participants.FindByRoomIDAndUserID(roomID, userId)

	return err == nil
}

func (s *roomQueryService) GetRoomMessages(roomId uuid.UUID) ([]MessageFull, error) {
	messages, err := s.messages.FindAllByRoomID(roomId)

	if err != nil {
		return nil, err
	}

	var messagesFull []MessageFull

	for _, message := range messages {
		messageFull, err := s.makeMessageFull(message)

		if err != nil {
			return nil, err
		}

		messagesFull = append(messagesFull, messageFull)
	}

	return messagesFull, nil
}

func (s *roomQueryService) makeMessageFull(message *domain.ChatMessage) (MessageFull, error) {
	user, err := s.users.FindByID(message.UserId)

	if err != nil {
		return MessageFull{}, err
	}

	m := MessageFull{
		User:        user,
		ChatMessage: message,
	}

	return m, nil
}
