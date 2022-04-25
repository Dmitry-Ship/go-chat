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

func TestRename(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	conversation := NewPublicConversation(conversationId, name)

	conversation.Rename("new name")

	assert.Equal(t, "new name", conversation.Data.Name)
}

func TestNewPrivateConversation(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationId := uuid.New()
	conversation := NewPrivateConversation(conversationId, to, from)
	assert.Equal(t, conversation.ID, conversationId)
	assert.Equal(t, to, conversation.Data.ToUser.UserID)
	assert.Equal(t, from, conversation.Data.FromUser.UserID)
	assert.Equal(t, conversationId, conversation.Data.FromUser.ConversationID)
	assert.Equal(t, conversationId, conversation.Data.ToUser.ConversationID)
	assert.NotNil(t, conversation.Data.ToUser.CreatedAt)
	assert.NotNil(t, conversation.Data.FromUser.CreatedAt)
	assert.NotNil(t, conversation.Data.FromUser.ID)
	assert.NotNil(t, conversation.Data.ToUser.ID)
	assert.Equal(t, conversation.Type, "private")
}

func TestGetFromUser(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationId := uuid.New()
	conversation := NewPrivateConversation(conversationId, to, from)
	assert.Equal(t, from, conversation.GetFromUser().UserID)
}

func TestGetToUser(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationId := uuid.New()
	conversation := NewPrivateConversation(conversationId, to, from)
	assert.Equal(t, to, conversation.GetToUser().UserID)
}
