package inmemory

import (
	"GitHub/go-chat/backend/domain"
	"errors"
)

type chatMessageRepository struct {
	chatMessages map[int32]*domain.ChatMessage
}

func NewChatMessageRepository() *chatMessageRepository {
	return &chatMessageRepository{
		chatMessages: make(map[int32]*domain.ChatMessage),
	}
}

func (r *chatMessageRepository) FindByID(id int32) (*domain.ChatMessage, error) {
	chatMessage, ok := r.chatMessages[id]
	if !ok {
		return nil, errors.New("not found")
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
		return errors.New("not found")
	}
	r.chatMessages[chatMessage.Id] = chatMessage
	return nil
}

func (r *chatMessageRepository) Delete(id int32) error {
	_, ok := r.chatMessages[id]
	if !ok {
		return errors.New("not found")
	}
	delete(r.chatMessages, id)
	return nil
}

func (r *chatMessageRepository) FindByRoomID(roomID int32) ([]*domain.ChatMessage, error) {
	chatMessages := make([]*domain.ChatMessage, 0, len(r.chatMessages))
	for _, message := range r.chatMessages {

		if message.RoomId == roomID {

			chatMessages = append(chatMessages, message)
		}
	}

	return chatMessages, nil
}
