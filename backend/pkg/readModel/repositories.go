package readModel

import (
	"github.com/google/uuid"
)

type UserQueryRepository interface {
	FindContacts(userID uuid.UUID) ([]*UserDTO, error)
	GetUserByID(userID uuid.UUID) (*UserDTO, error)
}

type ConversationQueryRepository interface {
	GetConversationByID(id uuid.UUID, userId uuid.UUID) (*ConversationFullDTO, error)
	FindMyConversations(userId uuid.UUID) ([]*ConversationDTO, error)
}

type MessageQueryRepository interface {
	FindAllByConversationID(conversationId uuid.UUID, requestUserId uuid.UUID) ([]*MessageDTO, error)
	GetMessageByID(messageID uuid.UUID, requestUserID uuid.UUID) (*MessageDTO, error)
}

type ParticipantQueryRepository interface {
	GetUserIdsByConversationID(conversationID uuid.UUID) ([]uuid.UUID, error)
}
