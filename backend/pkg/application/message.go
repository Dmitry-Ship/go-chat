package application

import (
	"GitHub/go-chat/backend/domain"

	"github.com/google/uuid"
)

type MessageService interface {
	GetMessages(roomId uuid.UUID) ([]*domain.ChatMessage, error)
	SendMessage(messageText string, messageType string, roomId uuid.UUID, userId uuid.UUID) (*MessageFull, error)
	GetMessagesFull(roomId uuid.UUID) ([]*MessageFull, error)
	MakeMessageFull(message *domain.ChatMessage) (*MessageFull, error)
}

type messageService struct {
	messages  domain.ChatMessageRepository
	users     domain.UserRepository
	Broadcast chan *MessageFull
}

type MessageFull struct {
	User *domain.User `json:"user"`
	*domain.ChatMessage
}

func NewMessageService(messages domain.ChatMessageRepository, users domain.UserRepository, broadcast chan *MessageFull) *messageService {
	return &messageService{
		messages:  messages,
		users:     users,
		Broadcast: broadcast,
	}
}

func (s *messageService) GetMessages(roomId uuid.UUID) ([]*domain.ChatMessage, error) {
	return s.messages.FindByRoomID(roomId)
}

func (s *messageService) GetMessagesFull(roomId uuid.UUID) ([]*MessageFull, error) {
	messages, err := s.GetMessages(roomId)

	if err != nil {
		return nil, err
	}

	var messagesFull []*MessageFull

	for _, message := range messages {
		messageFull, err := s.MakeMessageFull(message)

		if err != nil {
			return nil, err
		}

		messagesFull = append(messagesFull, messageFull)
	}

	return messagesFull, nil

}

func (s *messageService) MakeMessageFull(message *domain.ChatMessage) (*MessageFull, error) {
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

func (s *messageService) SendMessage(messageText string, messageType string, roomId uuid.UUID, userId uuid.UUID) (*MessageFull, error) {
	message := domain.NewChatMessage(messageText, messageType, roomId, userId)

	newMessage, err := s.messages.Create(message)

	if err != nil {
		return nil, err
	}

	fullMessage, err := s.MakeMessageFull(newMessage)

	if err != nil {
		return nil, err
	}

	s.Broadcast <- fullMessage

	return fullMessage, nil
}
