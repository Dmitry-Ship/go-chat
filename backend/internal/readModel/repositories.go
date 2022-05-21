package readModel

import (
	"github.com/google/uuid"
)

type PaginationInfo interface {
	GetPage() int
	GetPageSize() int
}

type userQueryRepository interface {
	GetContacts(userID uuid.UUID, paginationInfo PaginationInfo) ([]*ContactDTO, error)
	GetPotentialInvitees(conversationId uuid.UUID, paginationInfo PaginationInfo) ([]*ContactDTO, error)
	GetParticipants(conversationId uuid.UUID, userID uuid.UUID, paginationInfo PaginationInfo) ([]*ContactDTO, error)
	GetUserByID(userID uuid.UUID) (*UserDTO, error)
}

type conversationQueryRepository interface {
	GetConversation(id uuid.UUID, userId uuid.UUID) (*ConversationFullDTO, error)
	GetUserConversations(userId uuid.UUID, paginationInfo PaginationInfo) ([]*ConversationDTO, error)
}

type messageQueryRepository interface {
	GetConversationMessages(conversationId uuid.UUID, requestUserId uuid.UUID, paginationInfo PaginationInfo) ([]*MessageDTO, error)
	GetNotificationMessage(messageID uuid.UUID, requestUserID uuid.UUID) (*MessageDTO, error)
}

type QueriesRepository interface {
	messageQueryRepository
	conversationQueryRepository
	userQueryRepository
}
