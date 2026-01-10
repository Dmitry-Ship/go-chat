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
	RenameConversationAndReturn(conversationID uuid.UUID, name string) error
}

type messageQueryRepository interface {
	GetConversationMessages(conversationID uuid.UUID, requestUserID uuid.UUID, paginationInfo PaginationInfo) ([]MessageDTO, error)
	GetNotificationMessage(messageID uuid.UUID, requestUserID uuid.UUID) (MessageDTO, error)
	StoreMessageAndReturnWithUser(id uuid.UUID, conversationID uuid.UUID, userID uuid.UUID, content string, messageType int32) (MessageDTO, error)
}

type authorizationQueryRepository interface {
	IsMember(conversationID uuid.UUID, userID uuid.UUID) (bool, error)
	IsMemberOwner(conversationID uuid.UUID, userID uuid.UUID) (bool, error)
	InviteToConversationAtomic(conversationID uuid.UUID, inviteeID uuid.UUID, participantID uuid.UUID) (uuid.UUID, error)
	KickParticipantAtomic(conversationID uuid.UUID, targetID uuid.UUID) (int64, error)
	LeaveConversationAtomic(conversationID uuid.UUID, userID uuid.UUID) (int64, error)
}

type QueriesRepository interface {
	messageQueryRepository
	conversationQueryRepository
	userQueryRepository
	authorizationQueryRepository
}
