package readModel

import (
	"github.com/google/uuid"
)

type UserQueryRepository interface {
	GetContacts(userID uuid.UUID) ([]*ContactDTO, error)
	GetUserByID(userID uuid.UUID) (*UserDTO, error)
}

type ConversationQueryRepository interface {
	GetConversation(id uuid.UUID, userId uuid.UUID) (*ConversationFullDTO, error)
	GetUserConversations(userId uuid.UUID) ([]*ConversationDTO, error)
}

type MessageQueryRepository interface {
	GetConversationMessages(conversationId uuid.UUID, requestUserId uuid.UUID) ([]*MessageDTO, error)
	GetNotificationMessage(messageID uuid.UUID, requestUserID uuid.UUID) (*MessageDTO, error)
}
