package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	UserID         uuid.UUID
	CreatedAt      time.Time
	Type           string
}

func newMessage(conversationId uuid.UUID, userID uuid.UUID, messageType string) *Message {
	return &Message{
		ID:             uuid.New(),
		ConversationID: conversationId,
		CreatedAt:      time.Now(),
		Type:           messageType,
		UserID:         userID,
	}
}

func (m *Message) GetBaseData() *Message {
	return m
}

type TextMessageData struct {
	ID   uuid.UUID
	Text string
}

type TextMessage struct {
	Message
	Data TextMessageData
}

func NewTextMessage(conversationId uuid.UUID, userID uuid.UUID, text string) *TextMessage {
	baseMessage := newMessage(conversationId, userID, "text")

	return &TextMessage{
		Message: *baseMessage,
		Data: TextMessageData{
			ID:   uuid.New(),
			Text: text,
		},
	}
}

func (tm *TextMessage) GetTextMessageData() TextMessageData {
	return tm.Data
}

type conversationRenamedMessageData struct {
	ID      uuid.UUID
	NewName string
}

type ConversationRenamedMessage struct {
	Message
	Data conversationRenamedMessageData
}

func NewConversationRenamedMessage(conversationId uuid.UUID, userID uuid.UUID, newName string) *ConversationRenamedMessage {
	baseMessage := newMessage(conversationId, userID, "renamed_conversation")
	return &ConversationRenamedMessage{
		Message: *baseMessage,
		Data: conversationRenamedMessageData{
			ID:      uuid.New(),
			NewName: newName,
		},
	}

}

func (crm *ConversationRenamedMessage) GetConversationRenamedMessage() *conversationRenamedMessageData {
	return &crm.Data
}

type LeftConversationMessage = Message

func NewLeftConversationMessage(conversationId uuid.UUID, userID uuid.UUID) *LeftConversationMessage {
	return newMessage(conversationId, userID, "left_conversation")
}

type JoinedConversationMessage = Message

func NewJoinedConversationMessage(conversationId uuid.UUID, userID uuid.UUID) *JoinedConversationMessage {
	return newMessage(conversationId, userID, "joined_conversation")
}
