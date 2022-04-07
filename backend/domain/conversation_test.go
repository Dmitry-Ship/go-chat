package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewConversation(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	conversation := NewConversation(conversationId, name, false)
	assert.Equal(t, conversation.ID, conversationId)
	assert.Equal(t, name, conversation.Name)
	assert.False(t, conversation.IsPrivate)
}
