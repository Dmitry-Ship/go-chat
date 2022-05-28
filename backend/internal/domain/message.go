package domain

import (
	"errors"

	"github.com/google/uuid"
)

type MessageRepository interface {
	Store(message *Message) error
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
	aggregate
	ID             uuid.UUID
	ConversationID uuid.UUID
	UserID         uuid.UUID
	Content        messageContent
	Type           MessageType
}

func newMessage(conversationID uuid.UUID, userID uuid.UUID, messageType MessageType, content messageContent) *Message {
	message := Message{
		ID:             uuid.New(),
		ConversationID: conversationID,
		UserID:         userID,
		Type:           messageType,
		Content:        content,
	}

	message.AddEvent(NewMessageSent(conversationID, message.ID, userID))

	return &message
}

type textMessageContent struct {
	text string
}

func newTextMessageContent(text string) (*textMessageContent, error) {
	if text == "" {
		return nil, errors.New("text is empty")
	}

	if len(text) > 1000 {
		return nil, errors.New("text is too long")
	}

	return &textMessageContent{
		text: text,
	}, nil
}

func (m *textMessageContent) String() string {
	return m.text
}

func newTextMessage(conversationID uuid.UUID, userID uuid.UUID, text *textMessageContent) *Message {
	return newMessage(conversationID, userID, MessageTypeText, text)
}

type renamedMessageContent struct {
	newName string
}

func newRenamedMessageContent(newName string) *renamedMessageContent {
	return &renamedMessageContent{
		newName: newName,
	}
}

func (m *renamedMessageContent) String() string {
	return m.newName
}

func newConversationRenamedMessage(conversationID uuid.UUID, userID uuid.UUID, content *renamedMessageContent) *Message {
	return newMessage(conversationID, userID, MessageTypeRenamedConversation, content)
}

type emptyMessageContent struct {
	text string
}

func newEmptyMessageContent() *emptyMessageContent {
	return &emptyMessageContent{
		text: "newName",
	}
}

func (m *emptyMessageContent) String() string {
	return m.text
}

func newLeftConversationMessage(conversationID uuid.UUID, userID uuid.UUID) *Message {
	return newMessage(conversationID, userID, MessageTypeLeftConversation, newEmptyMessageContent())
}

func newJoinedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) *Message {
	return newMessage(conversationID, userID, MessageTypeJoinedConversation, newEmptyMessageContent())
}

func newInvitedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) *Message {
	return newMessage(conversationID, userID, MessageTypeInvitedConversation, newEmptyMessageContent())
}
