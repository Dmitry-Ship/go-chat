package domain

import (
	"github.com/google/uuid"
)

type ConversationRepository interface {
	Store(conversation *Conversation) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*ConversationDTO, error)
	FindAll() ([]*ConversationDTO, error)
}

type UserRepository interface {
	Store(user *User) error
	FindByID(id uuid.UUID) (*UserDTO, error)
	FindByUsername(username string) (*User, error)
	StoreRefreshToken(userID uuid.UUID, refreshToken string) error
	GetRefreshTokenByUserID(userID uuid.UUID) (string, error)
	DeleteRefreshToken(userID uuid.UUID) error
	FindAll() ([]*UserDTO, error)
}

type MessageRepository interface {
	Store(message *Message) error
	FindAllByConversationID(conversationId uuid.UUID) ([]*MessageDTO, error)
}

type ParticipantRepository interface {
	Store(participant *Participant) error
	DeleteByConversationIDAndUserID(conversationId uuid.UUID, userId uuid.UUID) error
	FindAllByConversationID(conversationId uuid.UUID) ([]*Participant, error)
	FindByConversationIDAndUserID(conversationId uuid.UUID, userId uuid.UUID) (*Participant, error)
}
