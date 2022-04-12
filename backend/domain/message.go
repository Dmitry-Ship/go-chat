package domain

import (
	"time"

	"github.com/google/uuid"
)

type MessageAggregate struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	UserID         uuid.UUID
	CreatedAt      time.Time
	Type           string
}

func (m *MessageAggregate) GetBaseMessage() *MessageAggregate {
	return m
}

func newMessage(conversationId uuid.UUID, userID uuid.UUID, messageType string) *MessageAggregate {
	return &MessageAggregate{
		ID:             uuid.New(),
		ConversationID: conversationId,
		CreatedAt:      time.Now(),
		Type:           messageType,
		UserID:         userID,
	}
}

type TextMessage struct {
	ID        uuid.UUID
	MessageID uuid.UUID
	Text      string
}

type TextMessageAggregate struct {
	MessageAggregate
	Text TextMessage
}

func (tm *TextMessageAggregate) GetText() TextMessage {
	return tm.Text
}

func NewTextMessage(conversationId uuid.UUID, userID uuid.UUID, text string) *TextMessageAggregate {
	baseMessage := newMessage(conversationId, userID, "text")

	return &TextMessageAggregate{
		MessageAggregate: *baseMessage,
		Text: TextMessage{
			ID:        uuid.New(),
			MessageID: baseMessage.ID,
			Text:      text,
		},
	}
}

type ConversationRenamedMessage struct {
	ID        uuid.UUID
	MessageID uuid.UUID
	NewName   string
}

type ConversationRenamedMessageAggregate struct {
	MessageAggregate
	ConversationRenamedMessage ConversationRenamedMessage
}

func (crm *ConversationRenamedMessageAggregate) GetConversationRenamedMessage() ConversationRenamedMessage {
	return crm.ConversationRenamedMessage
}

func NewConversationRenamedMessage(conversationId uuid.UUID, userID uuid.UUID, newName string) *ConversationRenamedMessageAggregate {

	baseMessage := newMessage(conversationId, userID, "conversation_renamed")
	return &ConversationRenamedMessageAggregate{
		MessageAggregate: *baseMessage,
		ConversationRenamedMessage: ConversationRenamedMessage{
			ID:        uuid.New(),
			MessageID: baseMessage.ID,
			NewName:   newName,
		},
	}
}

func NewLeftConversationMessage(conversationId uuid.UUID, userID uuid.UUID) *MessageAggregate {
	return newMessage(conversationId, userID, "left_conversation")
}

func NewJoinedConversationMessage(conversationId uuid.UUID, userID uuid.UUID) *MessageAggregate {
	return newMessage(conversationId, userID, "joined_conversation")
}
