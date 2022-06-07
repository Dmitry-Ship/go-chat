package readModel

import (
	"github.com/google/uuid"
)

type PaginationInfo interface {
	GetPage() int
	GetPageSize() int
}

type userQueryRepository interface {
	GetContacts(userID uuid.UUID, paginationInfo PaginationInfo) ([]ContactDTO, error)
	GetPotentialInvitees(conversationID uuid.UUID, paginationInfo PaginationInfo) ([]ContactDTO, error)
	GetParticipants(conversationID uuid.UUID, userID uuid.UUID, paginationInfo PaginationInfo) ([]ContactDTO, error)
	GetUserByID(userID uuid.UUID) (UserDTO, error)
}

type conversationQueryRepository interface {
	GetConversation(id uuid.UUID, userID uuid.UUID) (ConversationFullDTO, error)
	GetUserConversations(userID uuid.UUID, paginationInfo PaginationInfo) ([]ConversationDTO, error)
}

type messageQueryRepository interface {
	GetConversationMessages(conversationID uuid.UUID, requestUserID uuid.UUID, paginationInfo PaginationInfo) ([]MessageDTO, error)
	GetNotificationMessage(messageID uuid.UUID, requestUserID uuid.UUID) (MessageDTO, error)
}

type QueriesRepository interface {
	messageQueryRepository
	conversationQueryRepository
	userQueryRepository
}
