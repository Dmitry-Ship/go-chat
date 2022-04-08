package domain

import (
	"github.com/google/uuid"
)

type ConversationCommandRepository interface {
	Store(conversation *Conversation) error
	Delete(id uuid.UUID) error
}

type ConversationQueryRepository interface {
	FindByID(id uuid.UUID) (*ConversationDTO, error)
	FindAll() ([]*ConversationDTO, error)
}

type UserCommandRepository interface {
	Store(user *User) error
	FindByUsername(username string) (*User, error)
	StoreRefreshToken(userID uuid.UUID, refreshToken string) error
	FindByID(id uuid.UUID) (*UserDTO, error)
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
	FindAllByConversationID(conversationId uuid.UUID) ([]*MessageDTO, error)
}

type ParticipantCommandRepository interface {
	Store(participant *Participant) error
	DeleteByConversationIDAndUserID(conversationId uuid.UUID, userId uuid.UUID) error
	FindAllByConversationID(conversationId uuid.UUID) ([]*Participant, error)
	FindByConversationIDAndUserID(conversationId uuid.UUID, userId uuid.UUID) (*Participant, error)
}
