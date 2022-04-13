package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewPublicConversation(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	conversation := NewPublicConversation(conversationId, name)
	assert.Equal(t, conversation.ID, conversationId)
	assert.Equal(t, name, conversation.Name)
	assert.Equal(t, conversation.Type, "public")
}
