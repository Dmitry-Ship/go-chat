package readModel

import (
	"github.com/google/uuid"
)

type UserQueryRepository interface {
	FindAll() ([]*UserDTO, error)
	GetUserByID(id uuid.UUID) (*UserDTO, error)
}

type ConversationQueryRepository interface {
	GetConversationByID(id uuid.UUID, userId uuid.UUID) (*ConversationDTOFull, error)
	FindAll() ([]*ConversationDTO, error)
}

type MessageQueryRepository interface {
	FindAllByConversationID(conversationId uuid.UUID, requestUserId uuid.UUID) ([]*MessageDTO, error)
	GetMessageByID(messageID uuid.UUID, requestUserID uuid.UUID) (*MessageDTO, error)
}

type ParticipantQueryRepository interface {
	GetUserIdsByConversationID(conversationID uuid.UUID) ([]uuid.UUID, error)
}
