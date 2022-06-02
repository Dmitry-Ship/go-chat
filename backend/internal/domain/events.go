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

type conversationEvent struct {
	conversationID uuid.UUID
}

func (e *conversationEvent) GetConversationID() uuid.UUID {
	return e.conversationID
}

type ConversationEvent interface {
	DomainEvent
	GetConversationID() uuid.UUID
}

type MessageSent struct {
	domainEvent
	conversationEvent
	UserID    uuid.UUID
	MessageID uuid.UUID
}

func NewMessageSent(conversationID uuid.UUID, messageID uuid.UUID, userID uuid.UUID) *MessageSent {
	return &MessageSent{
		domainEvent: domainEvent{
			name: MessageSentName,
		},
		conversationEvent: conversationEvent{
			conversationID: conversationID,
		},
		UserID:    userID,
		MessageID: messageID,
	}
}

type GroupConversationDeleted struct {
	domainEvent
	conversationEvent
}

func newGroupConversationDeletedEvent(conversationID uuid.UUID) *GroupConversationDeleted {
	return &GroupConversationDeleted{
		domainEvent: domainEvent{
			name: GroupConversationDeletedName,
		},
		conversationEvent: conversationEvent{
			conversationID: conversationID,
		},
	}
}

type GroupConversationRenamed struct {
	domainEvent
	conversationEvent
	UserID  uuid.UUID
	NewName string
}

func newGroupConversationRenamedEvent(conversationID uuid.UUID, userID uuid.UUID, newName string) *GroupConversationRenamed {
	return &GroupConversationRenamed{
		domainEvent: domainEvent{
			name: GroupConversationRenamedName,
		},
		conversationEvent: conversationEvent{
			conversationID: conversationID,
		},
		UserID:  userID,
		NewName: newName,
	}
}

type GroupConversationCreated struct {
	domainEvent
	conversationEvent
	OwnerID uuid.UUID
}

func newGroupConversationCreatedEvent(conversationID uuid.UUID, ownerId uuid.UUID) *GroupConversationCreated {
	return &GroupConversationCreated{
		domainEvent: domainEvent{
			name: GroupConversationCreatedName,
		},
		conversationEvent: conversationEvent{
			conversationID: conversationID,
		},
		OwnerID: ownerId,
	}
}

type GroupConversationLeft struct {
	domainEvent
	conversationEvent
	UserID uuid.UUID
}

func newGroupConversationLeftEvent(conversationID uuid.UUID, userID uuid.UUID) *GroupConversationLeft {
	return &GroupConversationLeft{
		domainEvent: domainEvent{
			name: GroupConversationLeftName,
		},
		conversationEvent: conversationEvent{
			conversationID: conversationID,
		},
		UserID: userID,
	}
}

type GroupConversationJoined struct {
	domainEvent
	conversationEvent
	UserID uuid.UUID
}

func newGroupConversationJoinedEvent(conversationID uuid.UUID, userID uuid.UUID) *GroupConversationJoined {
	return &GroupConversationJoined{
		domainEvent: domainEvent{
			name: GroupConversationJoinedName,
		},
		conversationEvent: conversationEvent{
			conversationID: conversationID,
		},
		UserID: userID,
	}
}

type GroupConversationInvited struct {
	domainEvent
	conversationEvent
	UserID    uuid.UUID
	InvitedBy uuid.UUID
}

func newGroupConversationInvitedEvent(conversationID uuid.UUID, userID uuid.UUID, invitee uuid.UUID) *GroupConversationInvited {
	return &GroupConversationInvited{
		domainEvent: domainEvent{
			name: GroupConversationInvitedName,
		},
		conversationEvent: conversationEvent{
			conversationID: conversationID,
		},
		UserID:    invitee,
		InvitedBy: userID,
	}
}

type DirectConversationCreated struct {
	domainEvent
	conversationEvent
	UserIDs []uuid.UUID
}

func newDirectConversationCreatedEvent(conversationID uuid.UUID, userIDs []uuid.UUID) *DirectConversationCreated {
	return &DirectConversationCreated{
		domainEvent: domainEvent{
			name: DirectConversationCreatedName,
		},
		conversationEvent: conversationEvent{
			conversationID: conversationID,
		},
		UserIDs: userIDs,
	}
}
