package domain

import "github.com/google/uuid"

const (
	MessageSentName               = "message_sent"
	DirectConversationCreatedName = "direct_conversation_created"
	GroupConversationCreatedName  = "group_conversation_created"
	GroupConversationRenamedName  = "group_conversation_renamed"
	GroupConversationDeletedName  = "group_conversation_deleted"
	GroupConversationLeftName     = "group_conversation_left"
	GroupConversationJoinedName   = "group_conversation_joined"
	GroupConversationInvitedName  = "group_conversation_invited"
)

const DomainEventChannel = "domain_event"

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

type GroupConversationDeleted struct {
	domainEvent
	ConversationID uuid.UUID
}

func NewGroupConversationDeleted(conversationID uuid.UUID) *GroupConversationDeleted {
	return &GroupConversationDeleted{
		domainEvent: domainEvent{
			name: GroupConversationDeletedName,
		},
		ConversationID: conversationID,
	}
}

type GroupConversationRenamed struct {
	domainEvent
	ConversationID uuid.UUID
	UserID         uuid.UUID
	NewName        string
}

func NewGroupConversationRenamed(conversationID uuid.UUID, userID uuid.UUID, newName string) *GroupConversationRenamed {
	return &GroupConversationRenamed{
		domainEvent: domainEvent{
			name: GroupConversationRenamedName,
		},
		ConversationID: conversationID,
		UserID:         userID,
		NewName:        newName,
	}
}

type GroupConversationCreated struct {
	domainEvent
	ConversationID uuid.UUID
	OwnerID        uuid.UUID
}

func NewGroupConversationCreated(conversationID uuid.UUID, ownerId uuid.UUID) *GroupConversationCreated {
	return &GroupConversationCreated{
		domainEvent: domainEvent{
			name: GroupConversationCreatedName,
		},
		ConversationID: conversationID,
		OwnerID:        ownerId,
	}
}

type GroupConversationLeft struct {
	domainEvent
	ConversationID uuid.UUID
	UserID         uuid.UUID
}

func NewGroupConversationLeft(conversationID uuid.UUID, userID uuid.UUID) *GroupConversationLeft {
	return &GroupConversationLeft{
		domainEvent: domainEvent{
			name: GroupConversationLeftName,
		},
		ConversationID: conversationID,
		UserID:         userID,
	}
}

type GroupConversationJoined struct {
	domainEvent
	ConversationID uuid.UUID
	UserID         uuid.UUID
}

func NewGroupConversationJoined(conversationID uuid.UUID, userID uuid.UUID) *GroupConversationJoined {
	return &GroupConversationJoined{
		domainEvent: domainEvent{
			name: GroupConversationJoinedName,
		},
		ConversationID: conversationID,
		UserID:         userID,
	}
}

type GroupConversationInvited struct {
	domainEvent
	ConversationID uuid.UUID
	UserID         uuid.UUID
	InvitedBy      uuid.UUID
}

func NewGroupConversationInvited(conversationID uuid.UUID, userID uuid.UUID, invitee uuid.UUID) *GroupConversationInvited {
	return &GroupConversationInvited{
		domainEvent: domainEvent{
			name: GroupConversationInvitedName,
		},
		ConversationID: conversationID,
		UserID:         invitee,
		InvitedBy:      userID,
	}
}

type DirectConversationCreated struct {
	domainEvent
	ConversationID uuid.UUID
	ToUserID       uuid.UUID
	FromUserID     uuid.UUID
}

func NewDirectConversationCreated(conversationID uuid.UUID, toUserID uuid.UUID, fromUserID uuid.UUID) *DirectConversationCreated {
	return &DirectConversationCreated{
		domainEvent: domainEvent{
			name: DirectConversationCreatedName,
		},
		ConversationID: conversationID,
		ToUserID:       toUserID,
		FromUserID:     fromUserID,
	}
}
