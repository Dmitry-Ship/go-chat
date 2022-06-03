package domain

import "github.com/google/uuid"

const (
	MessageSentEventName               = "message_sent"
	DirectConversationCreatedEventName = "direct_conversation_created"
	GroupConversationCreatedEventName  = "group_conversation_created"
	GroupConversationRenamedEventName  = "group_conversation_renamed"
	GroupConversationDeletedEventName  = "group_conversation_deleted"
	GroupConversationLeftEventName     = "group_conversation_left"
	GroupConversationJoinedEventName   = "group_conversation_joined"
	GroupConversationInvitedEventName  = "group_conversation_invited"
)

const DomainEventTopic = "domain_event"

type DomainEvent interface {
	GetName() string
	GetTopic() string
}

type domainEvent struct {
	name string
}

func (e *domainEvent) GetName() string {
	return e.name
}

func (e *domainEvent) GetTopic() string {
	return DomainEventTopic
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
			name: MessageSentEventName,
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
			name: GroupConversationDeletedEventName,
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
			name: GroupConversationRenamedEventName,
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
			name: GroupConversationCreatedEventName,
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
			name: GroupConversationLeftEventName,
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
			name: GroupConversationJoinedEventName,
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
			name: GroupConversationInvitedEventName,
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
			name: DirectConversationCreatedEventName,
		},
		conversationEvent: conversationEvent{
			conversationID: conversationID,
		},
		UserIDs: userIDs,
	}
}
