package inmemory

import (
	"GitHub/go-chat/backend/domain"

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

func (r *chatMessageRepository) Store(chatMessage *domain.ChatMessage) error {
	r.chatMessages[chatMessage.Id] = chatMessage
	return nil
}

func (r *chatMessageRepository) FindAllByRoomID(roomID uuid.UUID) ([]*domain.ChatMessage, error) {
	chatMessages := make([]*domain.ChatMessage, 0, len(r.chatMessages))
	for _, message := range r.chatMessages {

		if message.RoomId == roomID {

			chatMessages = append(chatMessages, message)
		}
	}

	return chatMessages, nil
}
