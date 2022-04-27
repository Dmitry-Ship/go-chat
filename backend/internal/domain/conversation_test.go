package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewPublicConversation(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()

	conversation := NewPublicConversation(conversationId, name, creatorId)

	assert.Equal(t, conversation.ID, conversationId)
	assert.Equal(t, name, conversation.Data.Name)
	assert.Equal(t, string(name[0]), conversation.Data.Avatar)
	assert.Equal(t, conversation.Type, "public")
	assert.Equal(t, conversationId, conversation.Data.Owner.ConversationID)
	assert.Equal(t, creatorId, conversation.Data.Owner.UserID)
	assert.NotNil(t, conversation.Data.Owner.CreatedAt)
	assert.NotNil(t, conversation.Data.Owner.ID)
	assert.Equal(t, conversation.IsActive, true)
	assert.Equal(t, conversation.events[len(conversation.events)-1], NewPublicConversationCreated(conversationId, creatorId))
}

func TestRenameSuccess(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation := NewPublicConversation(conversationId, name, creatorId)

	err := conversation.Rename("new name", creatorId)

	assert.Nil(t, err)
	assert.Equal(t, "new name", conversation.Data.Name)
	assert.Equal(t, conversation.events[len(conversation.events)-1], NewPublicConversationRenamed(conversationId, creatorId, "new name"))
}

func TestRenameFailure(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation := NewPublicConversation(conversationId, name, creatorId)

	err := conversation.Rename("new name", uuid.New())

	assert.NotNil(t, err)
	assert.Equal(t, "user is not owner", err.Error())
	assert.Equal(t, name, conversation.Data.Name)
	assert.Equal(t, conversation.events[len(conversation.events)-1], NewPublicConversationCreated(conversationId, creatorId))
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
	assert.Equal(t, true, conversation.IsActive)
	assert.Equal(t, conversation.events[len(conversation.events)-1], NewPrivateConversationCreated(conversationId, to, from))
}

func TestPrivateConversation_GetFromUser(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationId := uuid.New()
	conversation := NewPrivateConversation(conversationId, to, from)
	assert.Equal(t, from, conversation.GetFromUser().UserID)
}

func TestPrivateConversation_GetToUser(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationId := uuid.New()
	conversation := NewPrivateConversation(conversationId, to, from)
	assert.Equal(t, to, conversation.GetToUser().UserID)
}

func TestPublicConversation_DeleteSuccess(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()

	conversation := NewPublicConversation(conversationId, name, creatorId)

	err := conversation.Delete(creatorId)

	assert.Nil(t, err)
	assert.Equal(t, false, conversation.IsActive)
	assert.Equal(t, conversation.events[len(conversation.events)-1], NewPublicConversationDeleted(conversation.ID))
}

func TestPublicConversation_DeleteFailure(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()

	conversation := NewPublicConversation(conversationId, name, creatorId)

	err := conversation.Delete(uuid.New())

	assert.NotNil(t, err)
	assert.Equal(t, "user is not owner", err.Error())
	assert.Equal(t, true, conversation.IsActive)
	assert.Equal(t, conversation.events[len(conversation.events)-1], NewPublicConversationCreated(conversationId, creatorId))
}
