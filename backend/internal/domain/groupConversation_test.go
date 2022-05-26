package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewGroupConversation(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()

	conversation, err := NewGroupConversation(conversationID, name, creatorId)

	assert.Equal(t, conversation.ID, conversationID)
	assert.Equal(t, name, conversation.Data.Name)
	assert.Equal(t, string(name[0]), conversation.Data.Avatar)
	assert.Equal(t, conversation.Type, "group")
	assert.Equal(t, conversationID, conversation.Data.Owner.ConversationID)
	assert.Equal(t, creatorId, conversation.Data.Owner.UserID)
	assert.NotNil(t, conversation.Data.Owner.CreatedAt)
	assert.NotNil(t, conversation.Data.Owner.ID)
	assert.Equal(t, conversation.IsActive, true)
	assert.Equal(t, conversation.GetEvents()[len(conversation.GetEvents())-1], newGroupConversationCreatedEvent(conversationID, creatorId))
	assert.Nil(t, err)
}

func TestNewGroupConversationEmptyName(t *testing.T) {
	name := ""
	conversationID := uuid.New()
	creatorId := uuid.New()

	_, err := NewGroupConversation(conversationID, name, creatorId)

	assert.Equal(t, "name is empty", err.Error())
}

func TestRename(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)

	err := conversation.Rename("new name", creatorId)

	assert.Nil(t, err)
	assert.Equal(t, "new name", conversation.Data.Name)
	assert.Equal(t, conversation.GetEvents()[len(conversation.GetEvents())-1], newGroupConversationRenamedEvent(conversationID, creatorId, "new name"))
}

func TestSendTextMessage(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)
	participant := NewParticipant(conversationID, creatorId)

	message, err := conversation.SendTextMessage("new message", participant)

	assert.Nil(t, err)
	assert.Equal(t, "new message", message.Data.Text)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, "text", message.Type)
}

func TestSendTextMessageUserNotParticipant(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)
	participant := NewParticipant(uuid.New(), uuid.New())

	_, err := conversation.SendTextMessage("new message", participant)

	assert.Equal(t, "user is not participant", err.Error())
}

func TestSendTextMessageNotActive(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)
	participant := NewParticipant(conversationID, creatorId)
	_ = conversation.Delete(creatorId)

	_, err := conversation.SendTextMessage("new message", participant)

	assert.Equal(t, "conversation is not active", err.Error())
}

func TestSendJoinedConversationMessage(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)

	message, err := conversation.SendJoinedConversationMessage(conversationID, creatorId)

	assert.Nil(t, err)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, "joined_conversation", message.Type)
}

func TestSendInvitedConversationMessage(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)

	message, err := conversation.SendInvitedConversationMessage(conversationID, creatorId)

	assert.Nil(t, err)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, "invited_conversation", message.Type)
}

func TestSendRenamedConversationMessage(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)

	message, err := conversation.SendRenamedConversationMessage(conversationID, creatorId, "new name")

	assert.Nil(t, err)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, "renamed_conversation", message.Type)
}

func TestSendLeftConversationMessage(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)

	message, err := conversation.SendLeftConversationMessage(conversationID, creatorId)

	assert.Nil(t, err)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, "left_conversation", message.Type)
}

func TestRenameNotOwner(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)

	err := conversation.Rename("new name", uuid.New())

	assert.NotNil(t, err)
	assert.Equal(t, "user is not owner", err.Error())
	assert.Equal(t, name, conversation.Data.Name)
	assert.Equal(t, conversation.GetEvents()[len(conversation.GetEvents())-1], newGroupConversationCreatedEvent(conversationID, creatorId))
}

func TestDelete(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)

	err := conversation.Delete(creatorId)

	assert.Nil(t, err)
	assert.Equal(t, false, conversation.IsActive)
	assert.Equal(t, conversation.GetEvents()[len(conversation.GetEvents())-1], newGroupConversationDeletedEvent(conversation.ID))
}

func TestDeleteNotOwner(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)

	err := conversation.Delete(uuid.New())

	assert.NotNil(t, err)
	assert.Equal(t, "user is not owner", err.Error())
	assert.Equal(t, true, conversation.IsActive)
	assert.Equal(t, conversation.GetEvents()[len(conversation.GetEvents())-1], newGroupConversationCreatedEvent(conversationID, creatorId))
}

func TestDeleteNotActive(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)
	_ = conversation.Delete(creatorId)

	err := conversation.Delete(creatorId)

	assert.Equal(t, "conversation is not active", err.Error())
}

func TestJoin(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)
	userID := uuid.New()

	participant, err := conversation.Join(userID)

	assert.Nil(t, err)
	assert.Equal(t, conversationID, participant.ConversationID)
	assert.Equal(t, userID, participant.UserID)
	assert.NotNil(t, participant.ID)
	assert.NotNil(t, participant.CreatedAt)
	assert.Equal(t, participant.IsActive, true)
	assert.Equal(t, participant.GetEvents()[len(participant.GetEvents())-1], newGroupConversationJoinedEvent(conversationID, userID))
}

func TestJoinNotActive(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)
	userID := uuid.New()

	_ = conversation.Delete(creatorId)

	_, err := conversation.Join(userID)

	assert.Equal(t, err.Error(), "conversation is not active")
}

func TestInvite(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)
	userID := uuid.New()
	inviteeId := uuid.New()

	participant, err := conversation.Invite(userID, inviteeId)

	assert.Nil(t, err)
	assert.Equal(t, conversationID, participant.ConversationID)
	assert.Equal(t, inviteeId, participant.UserID)
	assert.NotNil(t, participant.ID)
	assert.NotNil(t, participant.CreatedAt)
	assert.Equal(t, participant.IsActive, true)
	assert.Equal(t, participant.GetEvents()[len(participant.GetEvents())-1], newGroupConversationInvitedEvent(conversationID, userID, inviteeId))
}

func TestInviteNotActive(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)
	userID := uuid.New()
	inviteeId := uuid.New()

	_ = conversation.Delete(creatorId)

	_, err := conversation.Invite(userID, inviteeId)

	assert.Equal(t, err.Error(), "conversation is not active")
}

func TestInviteOwner(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)

	_, err := conversation.Invite(uuid.New(), creatorId)

	assert.Equal(t, err.Error(), "user is owner")
}

func TestInviteSelf(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)
	inviteeId := uuid.New()

	_, err := conversation.Invite(inviteeId, inviteeId)

	assert.Equal(t, err.Error(), "cannot invite yourself")
}

func TestLeave(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)

	participant, err := conversation.Leave(&conversation.Data.Owner)

	assert.Nil(t, err)
	assert.Equal(t, participant.IsActive, false)
}

func TestLeaveNotActive(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)
	_ = conversation.Delete(creatorId)

	_, err := conversation.Leave(&conversation.Data.Owner)

	assert.Equal(t, err.Error(), "conversation is not active")
}

func TestLeaveNotMember(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)

	participant := NewParticipant(uuid.New(), uuid.New())

	_, err := conversation.Leave(participant)

	assert.Equal(t, err.Error(), "participant is not in conversation")
}

func TestLeaveAlreadyLeft(t *testing.T) {
	name := "test"
	conversationID := uuid.New()
	creatorId := uuid.New()
	conversation, _ := NewGroupConversation(conversationID, name, creatorId)
	_, _ = conversation.Leave(&conversation.Data.Owner)

	_, err := conversation.Leave(&conversation.Data.Owner)

	assert.Equal(t, err.Error(), "participant already left")
}
