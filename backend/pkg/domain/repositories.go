package domain

import (
	"github.com/google/uuid"
)

type ConversationCommandRepository interface {
	StorePublicConversation(conversation *PublicConversation) error
	StorePrivateConversation(conversation *PrivateConversation) error
	UpdatePublicConversation(conversation *PublicConversation) error
	GetPublicConversation(id uuid.UUID) (*PublicConversation, error)
	GetPrivateConversationID(firstUserId uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error)
	Delete(id uuid.UUID) error
}

type UserCommandRepository interface {
	Store(user *User) error
	Update(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	FindByUsername(username string) (*User, error)
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
