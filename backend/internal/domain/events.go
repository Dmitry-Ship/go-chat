package domain

import "github.com/google/uuid"

var MessageSentName = "message_sent"
var PublicConversationDeletedName = "public_conversation_deleted"
var PublicConversationRenamedName = "public_conversation_renamed"
var PublicConversationLeftName = "left_public_conversation"
var PrivateConversationCreatedName = "private_conversation_created"
var PublicConversationJoinedName = "public_conversation_joined"

type DomainEvent interface {
	GetName() string
}

type domainEvent struct {
	name string
}

func (e *domainEvent) GetName() string {
	return e.name
}

type MessageSent struct {
	domainEvent
	ConversationID uuid.UUID
	UserID         uuid.UUID
	MessageID      uuid.UUID
}

func NewMessageSent(conversationID uuid.UUID, messageID uuid.UUID, userID uuid.UUID) *MessageSent {
	return &MessageSent{
		domainEvent: domainEvent{
			name: MessageSentName,
		},
		ConversationID: conversationID,
		UserID:         userID,
		MessageID:      messageID,
	}
}

type PublicConversationDeleted struct {
	domainEvent
	ConversationID uuid.UUID
}

func NewPublicConversationDeleted(conversationID uuid.UUID) *PublicConversationDeleted {
	return &PublicConversationDeleted{
		domainEvent: domainEvent{
			name: PublicConversationDeletedName,
		},
		ConversationID: conversationID,
	}
}

type PublicConversationRenamed struct {
	domainEvent
	ConversationID uuid.UUID
	UserID         uuid.UUID
	NewName        string
}

func NewPublicConversationRenamed(conversationID uuid.UUID, userID uuid.UUID, newName string) *PublicConversationRenamed {
	return &PublicConversationRenamed{
		domainEvent: domainEvent{
			name: PublicConversationRenamedName,
		},
		ConversationID: conversationID,
		UserID:         userID,
		NewName:        newName,
	}
}

type PublicConversationLeft struct {
	domainEvent
	ConversationID uuid.UUID
	UserID         uuid.UUID
}

func NewPublicConversationLeft(conversationID uuid.UUID, userID uuid.UUID) *PublicConversationLeft {
	return &PublicConversationLeft{
		domainEvent: domainEvent{
			name: PublicConversationLeftName,
		},
		ConversationID: conversationID,
		UserID:         userID,
	}
}

type PublicConversationJoined struct {
	domainEvent
	ConversationID uuid.UUID
	UserID         uuid.UUID
}

func NewPublicConversationJoined(conversationID uuid.UUID, userID uuid.UUID) *PublicConversationJoined {
	return &PublicConversationJoined{
		domainEvent: domainEvent{
			name: PublicConversationJoinedName,
		},
		ConversationID: conversationID,
		UserID:         userID,
	}
}

type PrivateConversationCreated struct {
	domainEvent
	ConversationID uuid.UUID
	ToUserID       uuid.UUID
	FromUserID     uuid.UUID
}

func NewPrivateConversationCreated(conversationID uuid.UUID, toUserID uuid.UUID, fromUserID uuid.UUID) *PrivateConversationCreated {
	return &PrivateConversationCreated{
		domainEvent: domainEvent{
			name: PrivateConversationCreatedName,
		},
		ConversationID: conversationID,
		ToUserID:       toUserID,
		FromUserID:     fromUserID,
	}
}
