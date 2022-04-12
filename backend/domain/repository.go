package domain

import (
	"github.com/google/uuid"
)

type BaseMessage interface {
	GetBaseMessage() *MessageAggregate
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
	StoreTextMessage(message *TextMessageAggregate) error
	StoreLeftConversation(message *MessageAggregate) error
	StoreJoinedConversation(message *MessageAggregate) error
	StoreRenamedConversation(message *ConversationRenamedMessageAggregate) error
}

type ParticipantCommandRepository interface {
	Store(participant *ParticipantAggregate) error
	DeleteByConversationIDAndUserID(conversationId uuid.UUID, userId uuid.UUID) error
}
