package domain

import (
	"errors"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
)

type MessageType struct {
	slug string
}

func (r MessageType) String() string {
	return r.slug
}

var (
	MessageTypeUser   = MessageType{"user"}
	MessageTypeSystem = MessageType{"system"}
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

func NewMessage(conversationID uuid.UUID, userID uuid.UUID, messageType MessageType, content string) (*Message, error) {
	text, err := NewTextMessageContent(content)
	if err != nil {
		return nil, err
	}
	message := Message{
		ID:             uuid.New(),
		ConversationID: conversationID,
		UserID:         userID,
		Type:           messageType,
		Content:        text,
	}

	return &message, nil
}

type textMessageContent struct {
	text string
}

var sanitizer = bluemonday.UGCPolicy()

func NewTextMessageContent(text string) (textMessageContent, error) {
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
