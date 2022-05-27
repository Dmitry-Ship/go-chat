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

type TextMessageData struct {
	ID   uuid.UUID
	Text string
}

type TextMessage struct {
	Message
	Data TextMessageData
}

func newTextMessage(conversationID uuid.UUID, userID uuid.UUID, text string) (*TextMessage, error) {
	if text == "" {
		return nil, errors.New("text is empty")
	}

	if len(text) > 1000 {
		return nil, errors.New("text is too long")
	}

	baseMessage := newMessage(conversationID, userID, MessageTypeText)

	return &TextMessage{
		Message: *baseMessage,
		Data: TextMessageData{
			ID:   uuid.New(),
			Text: text,
		},
	}, nil
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

func newConversationRenamedMessage(conversationID uuid.UUID, userID uuid.UUID, newName string) *ConversationRenamedMessage {
	baseMessage := newMessage(conversationID, userID, MessageTypeRenamedConversation)
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

func newLeftConversationMessage(conversationID uuid.UUID, userID uuid.UUID) *LeftConversationMessage {
	return newMessage(conversationID, userID, MessageTypeLeftConversation)
}

type JoinedConversationMessage = Message

func newJoinedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) *JoinedConversationMessage {
	return newMessage(conversationID, userID, MessageTypeJoinedConversation)
}

type InvitedConversationMessage = Message

func newInvitedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) *InvitedConversationMessage {
	return newMessage(conversationID, userID, MessageTypeInvitedConversation)
}
