package domain

import (
	"github.com/google/uuid"
)

type ConversationRepository interface {
	StorePublicConversation(conversation *PublicConversation) error
	StorePrivateConversation(conversation *PrivateConversation) error
	UpdatePublicConversation(conversation *PublicConversation) error
	GetPublicConversation(id uuid.UUID) (*PublicConversation, error)
	GetPrivateConversationID(firstUserId uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error)
}

type UserRepository interface {
	Store(user *User) error
	Update(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	FindByUsername(username string) (*User, error)
}

type MessageRepository interface {
	StoreTextMessage(message *TextMessage) error
	StoreLeftConversationMessage(message *Message) error
	StoreJoinedConversationMessage(message *Message) error
	StoreRenamedConversationMessage(message *ConversationRenamedMessage) error
}

type ParticipantRepository interface {
	Store(participant *Participant) error
	GetByConversationIDAndUserID(conversationID uuid.UUID, userID uuid.UUID) (*Participant, error)
	Update(participant *Participant) error
}

type NotificationTopicRepository interface {
	Store(notificationTopic *NotificationTopic) error
	DeleteByUserIDAndTopic(userId uuid.UUID, topic string) error
	DeleteAllByTopic(topic string) error
	GetUserIDsByTopic(topic string) ([]uuid.UUID, error)
}
