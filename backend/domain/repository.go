package domain

import (
	"github.com/google/uuid"
)

type BaseMessage interface {
	GetBaseData() *Message
}

type ConversationCommandRepository interface {
	Store(conversation *ConversationAggregate) error
	RenameConversation(conversationId uuid.UUID, name string) error
	Delete(id uuid.UUID) error
}

type UserCommandRepository interface {
	Store(user *UserAggregate) error
	FindByUsername(username string) (*UserAggregate, error)
	StoreRefreshToken(userID uuid.UUID, refreshToken string) error
	GetRefreshTokenByUserID(userID uuid.UUID) (string, error)
	DeleteRefreshToken(userID uuid.UUID) error
}

type MessageCommandRepository interface {
	StoreTextMessage(message *TextMessage) error
	StoreLeftConversation(message *Message) error
	StoreJoinedConversation(message *Message) error
	StoreRenamedConversation(message *ConversationRenamedMessage) error
}

type ParticipantCommandRepository interface {
	Store(participant *ParticipantAggregate) error
	DeleteByConversationIDAndUserID(conversationId uuid.UUID, userId uuid.UUID) error
}
