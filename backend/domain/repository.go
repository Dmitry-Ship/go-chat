package domain

import "github.com/google/uuid"

type ConversationRepository interface {
	Store(conversation *Conversation) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*Conversation, error)
	FindAll() ([]*Conversation, error)
}

type UserRepository interface {
	Store(user *User) error
	FindByID(id uuid.UUID) (*User, error)
	FindByUsername(username string) (*User, error)
	StoreRefreshToken(userID uuid.UUID, refreshToken string) error
	GetRefreshTokenByUserID(userID uuid.UUID) (string, error)
	DeleteRefreshToken(userID uuid.UUID) error
	FindAll() ([]*User, error)
}

type ChatMessageRepository interface {
	Store(message *Message) error
	FindAllByConversationID(conversationId uuid.UUID) ([]*Message, error)
}

type ParticipantRepository interface {
	Store(participant *Participant) error
	DeleteByConversationIDAndUserID(conversationId uuid.UUID, userId uuid.UUID) error
	FindAllByConversationID(conversationId uuid.UUID) ([]*Participant, error)
	FindByConversationIDAndUserID(conversationId uuid.UUID, userId uuid.UUID) (*Participant, error)
}
