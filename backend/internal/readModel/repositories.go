package readModel

import (
	"github.com/google/uuid"
)

type userQueryRepository interface {
	GetContacts(userID uuid.UUID) ([]*ContactDTO, error)
	GetPotentialInvitees(conversationId uuid.UUID) ([]*ContactDTO, error)
	GetUserByID(userID uuid.UUID) (*UserDTO, error)
}

type conversationQueryRepository interface {
	GetConversation(id uuid.UUID, userId uuid.UUID) (*ConversationFullDTO, error)
	GetUserConversations(userId uuid.UUID) ([]*ConversationDTO, error)
}

type messageQueryRepository interface {
	GetConversationMessages(conversationId uuid.UUID, requestUserId uuid.UUID) ([]*MessageDTO, error)
	GetNotificationMessage(messageID uuid.UUID, requestUserID uuid.UUID) (*MessageDTO, error)
}

type QueriesRepository interface {
	messageQueryRepository
	conversationQueryRepository
	userQueryRepository
}
