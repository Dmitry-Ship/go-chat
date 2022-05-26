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
	assert.Equal(t, conversation.GetEvents()[len(conversation.events)-1], newGroupConversationCreatedEvent(conversationId, creatorId))
	assert.Nil(t, err)
}

func TestNewGroupConversationEmptyName(t *testing.T) {
	name := ""
	conversationId := uuid.New()
	creatorId := uuid.New()

	_, err := NewGroupConversation(conversationId, name, creatorId)

	assert.Equal(t, "name is empty", err.Error())
}

func TestRename(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)

	err := conversation.Rename("new name", creatorId)

	assert.Nil(t, err)
	assert.Equal(t, "new name", conversation.Data.Name)
	assert.Equal(t, conversation.GetEvents()[len(conversation.events)-1], newGroupConversationRenamedEvent(conversationId, creatorId, "new name"))
}

func TestSendTextMessage(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)
	participant := NewParticipant(conversationId, creatorId)

	message, err := conversation.SendTextMessage("new message", participant)

	assert.Nil(t, err)
	assert.Equal(t, "new message", message.Data.Text)
	assert.Equal(t, conversationId, message.ConversationID)
	assert.Equal(t, "text", message.Type)
}

func TestSendTextMessageUserNotParticipant(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)
	participant := NewParticipant(uuid.New(), uuid.New())

	_, err := conversation.SendTextMessage("new message", participant)

	assert.Equal(t, "user is not participant", err.Error())
}

func TestSendTextMessageNotActive(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)
	participant := NewParticipant(conversationId, creatorId)
	_ = conversation.Delete(creatorId)

	_, err := conversation.SendTextMessage("new message", participant)

	assert.Equal(t, "conversation is not active", err.Error())
}

func TestSendJoinedConversationMessage(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)

	message, err := conversation.SendJoinedConversationMessage(conversationId, creatorId)

	assert.Nil(t, err)
	assert.Equal(t, conversationId, message.ConversationID)
	assert.Equal(t, "joined_conversation", message.Type)
}

func TestSendInvitedConversationMessage(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)

	message, err := conversation.SendInvitedConversationMessage(conversationId, creatorId)

	assert.Nil(t, err)
	assert.Equal(t, conversationId, message.ConversationID)
	assert.Equal(t, "invited_conversation", message.Type)
}

func TestSendRenamedConversationMessage(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)

	message, err := conversation.SendRenamedConversationMessage(conversationId, creatorId, "new name")

	assert.Nil(t, err)
	assert.Equal(t, conversationId, message.ConversationID)
	assert.Equal(t, "renamed_conversation", message.Type)
}

func TestSendLeftConversationMessage(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)

	message, err := conversation.SendLeftConversationMessage(conversationId, creatorId)

	assert.Nil(t, err)
	assert.Equal(t, conversationId, message.ConversationID)
	assert.Equal(t, "left_conversation", message.Type)
}

func TestRenameNotOwner(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)

	err := conversation.Rename("new name", uuid.New())

	assert.NotNil(t, err)
	assert.Equal(t, "user is not owner", err.Error())
	assert.Equal(t, name, conversation.Data.Name)
	assert.Equal(t, conversation.GetEvents()[len(conversation.events)-1], newGroupConversationCreatedEvent(conversationId, creatorId))
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
	assert.Equal(t, conversation.GetEvents()[len(conversation.events)-1], newDirectConversationCreatedEvent(conversationId, to, from))
	assert.Nil(t, err)
}

func TestNewDirectConversationWithOneself(t *testing.T) {
	to := uuid.New()
	conversationId := uuid.New()

	_, err := NewDirectConversation(conversationId, to, to)

	assert.Equal(t, err.Error(), "cannot chat with yourself")
}

func TestGetFromUser(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationId := uuid.New()

	conversation, _ := NewDirectConversation(conversationId, to, from)

	assert.Equal(t, from, conversation.GetFromUser().UserID)
}

func TestGetToUser(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationId := uuid.New()

	conversation, _ := NewDirectConversation(conversationId, to, from)

	assert.Equal(t, to, conversation.GetToUser().UserID)
}

func TestDelete(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)

	err := conversation.Delete(creatorId)

	assert.Nil(t, err)
	assert.Equal(t, false, conversation.IsActive)
	assert.Equal(t, conversation.GetEvents()[len(conversation.events)-1], newGroupConversationDeletedEvent(conversation.ID))
}

func TestDeleteNotOwner(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)

	err := conversation.Delete(uuid.New())

	assert.NotNil(t, err)
	assert.Equal(t, "user is not owner", err.Error())
	assert.Equal(t, true, conversation.IsActive)
	assert.Equal(t, conversation.GetEvents()[len(conversation.events)-1], newGroupConversationCreatedEvent(conversationId, creatorId))
}

func TestJoin(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)
	userID := uuid.New()

	participant, err := conversation.Join(userID)

	assert.Nil(t, err)
	assert.Equal(t, conversationId, participant.ConversationID)
	assert.Equal(t, userID, participant.UserID)
	assert.NotNil(t, participant.ID)
	assert.NotNil(t, participant.CreatedAt)
	assert.Equal(t, participant.IsActive, true)
	assert.Equal(t, participant.GetEvents()[len(participant.events)-1], newGroupConversationJoinedEvent(conversationId, userID))
}

func TestJoinNotActive(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)
	userID := uuid.New()

	_ = conversation.Delete(creatorId)

	_, err := conversation.Join(userID)

	assert.Equal(t, err.Error(), "conversation is not active")
}

func TestInvite(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)
	userID := uuid.New()
	inviteeId := uuid.New()

	participant, err := conversation.Invite(userID, inviteeId)

	assert.Nil(t, err)
	assert.Equal(t, conversationId, participant.ConversationID)
	assert.Equal(t, inviteeId, participant.UserID)
	assert.NotNil(t, participant.ID)
	assert.NotNil(t, participant.CreatedAt)
	assert.Equal(t, participant.IsActive, true)
	assert.Equal(t, participant.GetEvents()[len(participant.events)-1], newGroupConversationInvitedEvent(conversationId, userID, inviteeId))
}

func TestInviteNotActive(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)
	userID := uuid.New()
	inviteeId := uuid.New()

	_ = conversation.Delete(creatorId)

	_, err := conversation.Invite(userID, inviteeId)

	assert.Equal(t, err.Error(), "conversation is not active")
}

func TestInviteOwner(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)

	_, err := conversation.Invite(uuid.New(), creatorId)

	assert.Equal(t, err.Error(), "user is owner")
}

func TestInviteSelf(t *testing.T) {
	name := "test"
	conversationId := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationId, name, creatorId)
	inviteeId := uuid.New()

	_, err := conversation.Invite(inviteeId, inviteeId)

	assert.Equal(t, err.Error(), "cannot invite yourself")
}
