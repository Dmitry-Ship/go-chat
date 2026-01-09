package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
)

type MessageRepository interface {
	Store(ctx context.Context, message *Message) error
	StoreSystemMessage(ctx context.Context, message *Message) (bool, error)
}

type MessageType struct {
	slug string
}

func (r MessageType) String() string {
	return r.slug
}

var (
	MessageTypeText                = MessageType{"text"}
	MessageTypeRenamedConversation = MessageType{"renamed_conversation"}
	MessageTypeLeftConversation    = MessageType{"left_conversation"}
	MessageTypeJoinedConversation  = MessageType{"joined_conversation"}
	MessageTypeInvitedConversation = MessageType{"invited_conversation"}
)

type messageContent interface {
	String() string
}

type Message struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	UserID         uuid.UUID
	Content        messageContent
	Type           MessageType
}

func newMessage(messageID uuid.UUID, conversationID uuid.UUID, userID uuid.UUID, messageType MessageType, content messageContent) *Message {
	message := Message{
		ID:             messageID,
		ConversationID: conversationID,
		UserID:         userID,
		Type:           messageType,
		Content:        content,
	}

	return &message
}

type textMessageContent struct {
	text string
}

var sanitizer = bluemonday.UGCPolicy()

func newTextMessageContent(text string) (textMessageContent, error) {
	if text == "" {
		return textMessageContent{}, errors.New("text is empty")
	}

	if len(text) > 1000 {
		return textMessageContent{}, errors.New("text is too long")
	}

	sanitizedText := sanitizer.Sanitize(text)

	return textMessageContent{
		text: sanitizedText,
	}, nil
}

func (m textMessageContent) String() string {
	return m.text
}

func newTextMessage(messageID uuid.UUID, conversationID uuid.UUID, userID uuid.UUID, text textMessageContent) *Message {
	return newMessage(messageID, conversationID, userID, MessageTypeText, text)
}

type renamedMessageContent struct {
	newName string
}

func newRenamedMessageContent(newName string) renamedMessageContent {
	return renamedMessageContent{
		newName: newName,
	}
}

func (m renamedMessageContent) String() string {
	return m.newName
}

func newConversationRenamedMessage(messageID uuid.UUID, conversationID uuid.UUID, userID uuid.UUID, content renamedMessageContent) *Message {
	return newMessage(messageID, conversationID, userID, MessageTypeRenamedConversation, content)
}

type emptyMessageContent struct {
}

func newEmptyMessageContent() emptyMessageContent {
	return emptyMessageContent{}
}

func (m emptyMessageContent) String() string {
	return ""
}

func newLeftConversationMessage(messageID uuid.UUID, conversationID uuid.UUID, userID uuid.UUID) *Message {
	return newMessage(messageID, conversationID, userID, MessageTypeLeftConversation, newEmptyMessageContent())
}

func newJoinedConversationMessage(messageID uuid.UUID, conversationID uuid.UUID, userID uuid.UUID) *Message {
	return newMessage(messageID, conversationID, userID, MessageTypeJoinedConversation, newEmptyMessageContent())
}

func newInvitedConversationMessage(messageID uuid.UUID, conversationID uuid.UUID, userID uuid.UUID) *Message {
	return newMessage(messageID, conversationID, userID, MessageTypeInvitedConversation, newEmptyMessageContent())
}

// NewSystemMessage creates a system message (joined, left, invited, renamed) without requiring full entity validation.
// Validation is performed at the database level for optimal performance.
func NewSystemMessage(messageID uuid.UUID, conversationID uuid.UUID, userID uuid.UUID, messageType MessageType, content string) *Message {
	var msgContent messageContent
	if messageType == MessageTypeRenamedConversation {
		msgContent = newRenamedMessageContent(content)
	} else {
		msgContent = newEmptyMessageContent()
	}
	return newMessage(messageID, conversationID, userID, messageType, msgContent)
}
