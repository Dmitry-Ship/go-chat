package domain

import (
	"github.com/google/uuid"
)

type ConversationCommandRepository interface {
	StorePublicConversation(conversation *PublicConversation) error
	RenamePublicConversation(conversationId uuid.UUID, name string) error
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
	StoreTextMessage(message *TextMessage) error
	StoreLeftConversationMessage(message *Message) error
	StoreJoinedConversationMessage(message *Message) error
	StoreRenamedConversationMessage(message *ConversationRenamedMessage) error
}

type ParticipantCommandRepository interface {
	Store(participant *Participant) error
	DeleteByConversationIDAndUserID(conversationId uuid.UUID, userId uuid.UUID) error
}

type NotificationTopicCommandRepository interface {
	Store(notificationTopic *NotificationTopic) error
	DeleteByUserIDAndTopic(userId uuid.UUID, topic string) error
	GetAllNotificationTopics(userID uuid.UUID) ([]string, error)
}
