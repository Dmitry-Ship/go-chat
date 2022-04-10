package domain

import (
	"github.com/google/uuid"
)

type ConversationCommandRepository interface {
	Store(conversation *Conversation) error
	RenameConversation(conversationId uuid.UUID, name string) error
	Delete(id uuid.UUID) error
}

type ConversationQueryRepository interface {
	GetConversationByID(id uuid.UUID, userId uuid.UUID) (*ConversationDTOFull, error)
	FindAll() ([]*ConversationDTO, error)
}

type UserCommandRepository interface {
	Store(user *User) error
	FindByUsername(username string) (*User, error)
	StoreRefreshToken(userID uuid.UUID, refreshToken string) error
	GetUserByID(id uuid.UUID) (*UserDTO, error)
	GetRefreshTokenByUserID(userID uuid.UUID) (string, error)
	DeleteRefreshToken(userID uuid.UUID) error
}

type UserQueryRepository interface {
	FindAll() ([]*UserDTO, error)
}

type MessageCommandRepository interface {
	Store(message *Message) error
}

type MessageQueryRepository interface {
	FindAllByConversationID(conversationId uuid.UUID, requestUserId uuid.UUID) ([]*MessageDTO, error)
	GetMessageByID(messageID uuid.UUID, requestUserID uuid.UUID) (*MessageDTO, error)
}

type ParticipantCommandRepository interface {
	Store(participant *Participant) error
	DeleteByConversationIDAndUserID(conversationId uuid.UUID, userId uuid.UUID) error
}

type ParticipantQueryRepository interface {
	GetUserIdsByConversationID(conversationID uuid.UUID) ([]uuid.UUID, error)
}
