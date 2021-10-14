package inmemory

import (
	"GitHub/go-chat/backend/domain"
	"errors"

	"github.com/google/uuid"
)

type chatMessageRepository struct {
	chatMessages map[uuid.UUID]*domain.ChatMessage
}

func NewChatMessageRepository() *chatMessageRepository {
	return &chatMessageRepository{
		chatMessages: make(map[uuid.UUID]*domain.ChatMessage),
	}
}

func (r *chatMessageRepository) FindByID(id uuid.UUID) (*domain.ChatMessage, error) {
	chatMessage, ok := r.chatMessages[id]
	if !ok {
		return nil, errors.New("message not found")
	}
	return chatMessage, nil
}

func (r *chatMessageRepository) FindAll() ([]*domain.ChatMessage, error) {
	chatMessages := make([]*domain.ChatMessage, 0, len(r.chatMessages))
	for _, chatMessage := range r.chatMessages {
		chatMessages = append(chatMessages, chatMessage)
	}
	return chatMessages, nil
}

func (r *chatMessageRepository) Create(chatMessage *domain.ChatMessage) (*domain.ChatMessage, error) {
	r.chatMessages[chatMessage.Id] = chatMessage
	return chatMessage, nil
}

func (r *chatMessageRepository) Update(chatMessage *domain.ChatMessage) error {
	_, ok := r.chatMessages[chatMessage.Id]
	if !ok {
		return errors.New("message not found")
	}
	r.chatMessages[chatMessage.Id] = chatMessage
	return nil
}

func (r *chatMessageRepository) Delete(id uuid.UUID) error {
	_, ok := r.chatMessages[id]
	if !ok {
		return errors.New("message not found")
	}
	delete(r.chatMessages, id)
	return nil
}

func (r *chatMessageRepository) FindByRoomID(roomID uuid.UUID) ([]*domain.ChatMessage, error) {
	chatMessages := make([]*domain.ChatMessage, 0, len(r.chatMessages))
	for _, message := range r.chatMessages {

		if message.RoomId == roomID {

			chatMessages = append(chatMessages, message)
		}
	}

	return chatMessages, nil
}
