package domain

import (
	"errors"

	"github.com/google/uuid"
)

type MessageRepository interface {
	StoreTextMessage(message *TextMessage) error
	StoreLeftConversationMessage(message *Message) error
	StoreJoinedConversationMessage(message *Message) error
	StoreInvitedConversationMessage(message *Message) error
	StoreRenamedConversationMessage(message *ConversationRenamedMessage) error
}

type BaseMessage interface {
	GetBaseData() *Message
}

const (
	MessageTypeText                = "text"
	MessageTypeRenamedConversation = "renamed_conversation"
	MessageTypeLeftConversation    = "left_conversation"
	MessageTypeJoinedConversation  = "joined_conversation"
	MessageTypeInvitedConversation = "invited_conversation"
)

type Message struct {
	aggregate
	ID             uuid.UUID
	ConversationID uuid.UUID
	UserID         uuid.UUID
	Type           string
}

func newMessage(conversationID uuid.UUID, userID uuid.UUID, messageType string) *Message {
	message := Message{
		ID:             uuid.New(),
		ConversationID: conversationID,
		Type:           messageType,
		UserID:         userID,
	}

	message.AddEvent(NewMessageSent(conversationID, message.ID, userID))

	return &message
}

func (m *Message) GetBaseData() *Message {
	return m
}

type TextMessage struct {
	Message
	ID   uuid.UUID
	Text string
}

func newTextMessage(conversationID uuid.UUID, userID uuid.UUID, text string) (*TextMessage, error) {
	if text == "" {
		return nil, errors.New("text is empty")
	}

	if len(text) > 1000 {
		return nil, errors.New("text is too long")
	}

	return &TextMessage{
		Message: *newMessage(conversationID, userID, MessageTypeText),
		ID:      uuid.New(),
		Text:    text,
	}, nil
}

type ConversationRenamedMessage struct {
	Message
	ID      uuid.UUID
	NewName string
}

func newConversationRenamedMessage(conversationID uuid.UUID, userID uuid.UUID, newName string) *ConversationRenamedMessage {
	return &ConversationRenamedMessage{
		Message: *newMessage(conversationID, userID, MessageTypeRenamedConversation),
		ID:      uuid.New(),
		NewName: newName,
	}

}

func newLeftConversationMessage(conversationID uuid.UUID, userID uuid.UUID) *Message {
	return newMessage(conversationID, userID, MessageTypeLeftConversation)
}

func newJoinedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) *Message {
	return newMessage(conversationID, userID, MessageTypeJoinedConversation)
}

func newInvitedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) *Message {
	return newMessage(conversationID, userID, MessageTypeInvitedConversation)
}
