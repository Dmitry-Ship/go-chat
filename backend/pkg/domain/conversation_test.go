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
	assert.Equal(t, name, conversation.Data.Name)
	assert.Equal(t, string(name[0]), conversation.Data.Avatar)
	assert.Equal(t, conversation.Type, "public")
}

func TestNewPrivateConversation(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationId := uuid.New()
	conversation := NewPrivateConversation(conversationId, to, from)
	assert.Equal(t, conversation.ID, conversationId)
	assert.Equal(t, to, conversation.Data.ToUserId)
	assert.Equal(t, from, conversation.Data.FromUserId)
	assert.Equal(t, conversation.Type, "private")
}
