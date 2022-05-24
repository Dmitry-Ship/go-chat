package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewGroupConversation(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()

	conversation, err := NewGroupConversation(conversationId, name, creatorId)

	assert.Equal(t, conversation.ID, conversationId)
	assert.Equal(t, name, conversation.Data.Name)
	assert.Equal(t, string(name[0]), conversation.Data.Avatar)
	assert.Equal(t, conversation.Type, "group")
	assert.Equal(t, conversationId, conversation.Data.Owner.ConversationID)
	assert.Equal(t, creatorId, conversation.Data.Owner.UserID)
	assert.NotNil(t, conversation.Data.Owner.CreatedAt)
	assert.NotNil(t, conversation.Data.Owner.ID)
	assert.Equal(t, conversation.IsActive, true)
	assert.Equal(t, conversation.events[len(conversation.events)-1], NewGroupConversationCreated(conversationId, creatorId))
	assert.Nil(t, err)
}

func TestNewGroupConversationEmptyName(t *testing.T) {
	name := ""
	conversationId := uuid.New()
	creatorId := uuid.New()

	_, err := NewGroupConversation(conversationId, name, creatorId)

	assert.Equal(t, "name is empty", err.Error())
}

func TestRenameSuccess(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)

	err := conversation.Rename("new name", creatorId)

	assert.Nil(t, err)
	assert.Equal(t, "new name", conversation.Data.Name)
	assert.Equal(t, conversation.events[len(conversation.events)-1], NewGroupConversationRenamed(conversationId, creatorId, "new name"))
}

func TestRenameFailure(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)

	err := conversation.Rename("new name", uuid.New())

	assert.NotNil(t, err)
	assert.Equal(t, "user is not owner", err.Error())
	assert.Equal(t, name, conversation.Data.Name)
	assert.Equal(t, conversation.events[len(conversation.events)-1], NewGroupConversationCreated(conversationId, creatorId))
}

func TestNewDirectConversation(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationId := uuid.New()

	conversation, err := NewDirectConversation(conversationId, to, from)

	assert.Equal(t, conversation.ID, conversationId)
	assert.Equal(t, to, conversation.Data.ToUser.UserID)
	assert.Equal(t, from, conversation.Data.FromUser.UserID)
	assert.Equal(t, conversationId, conversation.Data.FromUser.ConversationID)
	assert.Equal(t, conversationId, conversation.Data.ToUser.ConversationID)
	assert.NotNil(t, conversation.Data.ToUser.CreatedAt)
	assert.NotNil(t, conversation.Data.FromUser.CreatedAt)
	assert.NotNil(t, conversation.Data.FromUser.ID)
	assert.NotNil(t, conversation.Data.ToUser.ID)
	assert.Equal(t, conversation.Type, "direct")
	assert.Equal(t, true, conversation.IsActive)
	assert.Equal(t, conversation.events[len(conversation.events)-1], NewDirectConversationCreated(conversationId, to, from))
	assert.Nil(t, err)
}

func TestNewDirectConversationWithOneself(t *testing.T) {
	to := uuid.New()
	conversationId := uuid.New()

	_, err := NewDirectConversation(conversationId, to, to)

	assert.Equal(t, err.Error(), "cannot chat with yourself")
}

func TestDirectConversation_GetFromUser(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationId := uuid.New()

	conversation, _ := NewDirectConversation(conversationId, to, from)

	assert.Equal(t, from, conversation.GetFromUser().UserID)
}

func TestDirectConversation_GetToUser(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationId := uuid.New()

	conversation, _ := NewDirectConversation(conversationId, to, from)

	assert.Equal(t, to, conversation.GetToUser().UserID)
}

func TestGroupConversation_DeleteSuccess(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)

	err := conversation.Delete(creatorId)

	assert.Nil(t, err)
	assert.Equal(t, false, conversation.IsActive)
	assert.Equal(t, conversation.events[len(conversation.events)-1], NewGroupConversationDeleted(conversation.ID))
}

func TestGroupConversation_DeleteFailure(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)

	err := conversation.Delete(uuid.New())

	assert.NotNil(t, err)
	assert.Equal(t, "user is not owner", err.Error())
	assert.Equal(t, true, conversation.IsActive)
	assert.Equal(t, conversation.events[len(conversation.events)-1], NewGroupConversationCreated(conversationId, creatorId))
}

func TestGroupConversation_JoinSuccess(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)
	userId := uuid.New()

	participant, err := conversation.Join(userId)

	assert.Nil(t, err)
	assert.Equal(t, conversationId, participant.ConversationID)
	assert.Equal(t, userId, participant.UserID)
	assert.NotNil(t, participant.ID)
	assert.NotNil(t, participant.CreatedAt)
	assert.Equal(t, participant.IsActive, true)
	assert.Equal(t, participant.events[len(participant.events)-1], NewGroupConversationJoined(conversationId, userId))
}

func TestGroupConversation_JoinFailure(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)
	userId := uuid.New()

	conversation.Delete(creatorId)

	_, err := conversation.Join(userId)

	assert.Equal(t, err.Error(), "conversation is not active")
}

func TestGroupConversation_InviteSuccess(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)
	userId := uuid.New()

	participant, err := conversation.Invite(userId)

	assert.Nil(t, err)
	assert.Equal(t, conversationId, participant.ConversationID)
	assert.Equal(t, userId, participant.UserID)
	assert.NotNil(t, participant.ID)
	assert.NotNil(t, participant.CreatedAt)
	assert.Equal(t, participant.IsActive, true)
	assert.Equal(t, participant.events[len(participant.events)-1], NewGroupConversationInvited(conversationId, userId))
}

func TestGroupConversation_InviteFailure(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)
	userId := uuid.New()
	conversation.Delete(creatorId)

	_, err := conversation.Invite(userId)

	assert.Equal(t, err.Error(), "conversation is not active")
}
