package domain

import (
	"github.com/google/uuid"
)

type ConversationCommandRepository interface {
	Store(conversation *Conversation) error
	RenameConversation(conversationId uuid.UUID, name string) error
	Delete(id uuid.UUID) error
}

type UserCommandRepository interface {
	Store(user *User) error
	FindByUsername(username string) (*User, error)
	StoreRefreshToken(userID uuid.UUID, refreshToken string) error
	GetRefreshTokenByUserID(userID uuid.UUID) (string, error)
	DeleteRefreshToken(userID uuid.UUID) error
}

type MessageCommandRepository interface {
	Store(message *Message) error
}

type ParticipantCommandRepository interface {
	Store(participant *Participant) error
	DeleteByConversationIDAndUserID(conversationId uuid.UUID, userId uuid.UUID) error
}
