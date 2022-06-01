package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createTestGroupConversation() (*GroupConversation, *User, *Participant) {
	name, _ := NewConversationName("test")
	conversationID := uuid.New()
	creatorId := uuid.New()
	userName, _ := NewUserName("test")
	userPassword, _ := NewUserPassword("test", func(p []byte) ([]byte, error) { return p, nil })
	creatorUser := NewUser(creatorId, userName, userPassword)
	conversation, _ := NewGroupConversation(conversationID, name, creatorUser)

	return conversation, creatorUser, &conversation.Owner
}

func TestNewGroupConversation(t *testing.T) {
	conversationID := uuid.New()
	creatorId := uuid.New()
	userName, _ := NewUserName("test")
	userPassword, _ := NewUserPassword("test", func(p []byte) ([]byte, error) { return p, nil })
	creator := NewUser(creatorId, userName, userPassword)
	name, _ := NewConversationName("test")

	conversation, err := NewGroupConversation(conversationID, name, creator)

	assert.Equal(t, conversation.Conversation.ID, conversationID)
	assert.Equal(t, name, &conversation.Name)
	assert.Equal(t, string(name.String()[0]), conversation.Avatar)
	assert.Equal(t, conversation.Type, ConversationTypeGroup)
	assert.Equal(t, conversation.Conversation.ID, conversation.Owner.ConversationID)
	assert.Equal(t, creatorId, conversation.Owner.UserID)
	assert.NotNil(t, conversation.Owner.ID)
	assert.Equal(t, conversation.IsActive, true)
	assert.Equal(t, conversation.GetEvents()[len(conversation.GetEvents())-1], newGroupConversationCreatedEvent(conversation.Conversation.ID, creatorId))
	assert.Nil(t, err)
}

func TestNewConversationName(t *testing.T) {
	name, err := NewConversationName("test")

	assert.Nil(t, err)
	assert.Equal(t, "test", name.String())
}

func TestNewConversationNameEmptyName(t *testing.T) {
	_, err := NewConversationName("")

	assert.Equal(t, "name is empty", err.Error())
}

func TestNewConversationNameLongName(t *testing.T) {
	name := ""

	for i := 0; i < 101; i++ {
		name += "a"
	}

	_, err := NewConversationName(name)

	assert.Equal(t, "name is too long", err.Error())
}

func TestRename(t *testing.T) {
	conversation, _, creatorParticipant := createTestGroupConversation()

	newName, _ := NewConversationName("new name")

	err := conversation.Rename(newName, creatorParticipant)

	assert.Nil(t, err)
	assert.Equal(t, "new name", conversation.Name.String())
	assert.Equal(t, conversation.GetEvents()[len(conversation.GetEvents())-1], newGroupConversationRenamedEvent(conversation.Conversation.ID, creatorParticipant.UserID, "new name"))
}

func TestSendTextMessage(t *testing.T) {
	conversation, _, creatorParticipant := createTestGroupConversation()
	messageID := uuid.New()

	message, err := conversation.SendTextMessage(messageID, "new message", creatorParticipant)

	assert.Nil(t, err)
	assert.Equal(t, "new message", message.Content.String())
	assert.Equal(t, conversation.Conversation.ID, message.ConversationID)
	assert.Equal(t, MessageTypeText, message.Type)
}

func TestSendTextMessageUserNotParticipant(t *testing.T) {
	conversation, _, _ := createTestGroupConversation()
	messageID := uuid.New()
	participant := NewParticipant(uuid.New(), uuid.New(), uuid.New())

	_, err := conversation.SendTextMessage(messageID, "new message", participant)

	assert.Equal(t, "user is not in conversation", err.Error())
}

func TestSendTextMessageNotActive(t *testing.T) {
	conversation, _, creatorParticipant := createTestGroupConversation()
	messageID := uuid.New()
	_ = conversation.Delete(creatorParticipant)

	_, err := conversation.SendTextMessage(messageID, "new message", creatorParticipant)

	assert.Equal(t, "conversation is not active", err.Error())
}

func TestSendJoinedConversationMessage(t *testing.T) {
	conversation, creator, _ := createTestGroupConversation()
	messageID := uuid.New()

	message, err := conversation.SendJoinedConversationMessage(messageID, creator)

	assert.Nil(t, err)
	assert.Equal(t, conversation.Conversation.ID, message.ConversationID)
	assert.Equal(t, MessageTypeJoinedConversation, message.Type)
}

func TestSendInvitedConversationMessage(t *testing.T) {
	conversation, creator, _ := createTestGroupConversation()
	messageID := uuid.New()

	message, err := conversation.SendInvitedConversationMessage(messageID, creator)

	assert.Nil(t, err)
	assert.Equal(t, conversation.Conversation.ID, message.ConversationID)
	assert.Equal(t, MessageTypeInvitedConversation, message.Type)
}

func TestSendRenamedConversationMessage(t *testing.T) {
	conversation, _, creatorParticipant := createTestGroupConversation()
	messageID := uuid.New()

	message, err := conversation.SendRenamedConversationMessage(messageID, creatorParticipant, "new name")

	assert.Nil(t, err)
	assert.Equal(t, conversation.Conversation.ID, message.ConversationID)
	assert.Equal(t, MessageTypeRenamedConversation, message.Type)
}

func TestSendLeftConversationMessage(t *testing.T) {
	conversation, _, creatorParticipant := createTestGroupConversation()
	messageID := uuid.New()

	message, err := conversation.SendLeftConversationMessage(messageID, creatorParticipant)

	assert.Nil(t, err)
	assert.Equal(t, conversation.Conversation.ID, message.ConversationID)
	assert.Equal(t, MessageTypeLeftConversation, message.Type)
}

func TestRenameNotOwner(t *testing.T) {
	conversation, creator, _ := createTestGroupConversation()
	newName, _ := NewConversationName("new name")
	userPassword, _ := NewUserPassword("test", func(p []byte) ([]byte, error) { return p, nil })
	userName, _ := NewUserName("test")
	userID := uuid.New()
	user := NewUser(userID, userName, userPassword)
	participant, _ := conversation.Join(user)

	err := conversation.Rename(newName, participant)

	assert.NotNil(t, err)
	assert.Equal(t, "user is not owner", err.Error())
	assert.NotEqual(t, newName.String(), conversation.Name.String())
	assert.Equal(t, conversation.GetEvents()[len(conversation.GetEvents())-1], newGroupConversationCreatedEvent(conversation.Conversation.ID, creator.ID))
}

func TestDelete(t *testing.T) {
	conversation, _, participantCreator := createTestGroupConversation()

	err := conversation.Delete(participantCreator)

	assert.Nil(t, err)
	assert.Equal(t, false, conversation.IsActive)
	assert.Equal(t, conversation.GetEvents()[len(conversation.GetEvents())-1], newGroupConversationDeletedEvent(conversation.Conversation.ID))
}

func TestDeleteNotOwner(t *testing.T) {
	conversation, creator, _ := createTestGroupConversation()
	userPassword, _ := NewUserPassword("test", func(p []byte) ([]byte, error) { return p, nil })
	userName, _ := NewUserName("test")
	userID := uuid.New()
	user := NewUser(userID, userName, userPassword)
	participant, _ := conversation.Join(user)

	err := conversation.Delete(participant)

	assert.NotNil(t, err)
	assert.Equal(t, "user is not owner", err.Error())
	assert.Equal(t, true, conversation.IsActive)
	assert.Equal(t, conversation.GetEvents()[len(conversation.GetEvents())-1], newGroupConversationCreatedEvent(conversation.Conversation.ID, creator.ID))
}

func TestDeleteNotActive(t *testing.T) {
	conversation, _, creatorParticipant := createTestGroupConversation()
	_ = conversation.Delete(creatorParticipant)

	err := conversation.Delete(creatorParticipant)

	assert.Equal(t, "conversation is not active", err.Error())
}

func TestJoin(t *testing.T) {
	conversation, _, _ := createTestGroupConversation()
	userPassword, _ := NewUserPassword("test", func(p []byte) ([]byte, error) { return p, nil })
	userName, _ := NewUserName("test")
	userID := uuid.New()
	user := NewUser(userID, userName, userPassword)

	participant, err := conversation.Join(user)

	assert.Nil(t, err)
	assert.Equal(t, conversation.Conversation.ID, participant.ConversationID)
	assert.Equal(t, userID, participant.UserID)
	assert.NotNil(t, participant.ID)
	assert.Equal(t, participant.IsActive, true)
	assert.Equal(t, participant.GetEvents()[len(participant.GetEvents())-1], newGroupConversationJoinedEvent(conversation.Conversation.ID, userID))
}

func TestJoinNotActive(t *testing.T) {
	conversation, creator, creatorParticipant := createTestGroupConversation()

	_ = conversation.Delete(creatorParticipant)

	_, err := conversation.Join(creator)

	assert.Equal(t, err.Error(), "conversation is not active")
}

func TestInvite(t *testing.T) {
	conversation, _, creatorParticipant := createTestGroupConversation()
	userPassword, _ := NewUserPassword("test", func(p []byte) ([]byte, error) { return p, nil })
	userName, _ := NewUserName("test")
	inviteeId := uuid.New()
	user := NewUser(inviteeId, userName, userPassword)

	participant, err := conversation.Invite(creatorParticipant, user)

	assert.Nil(t, err)
	assert.Equal(t, conversation.Conversation.ID, participant.ConversationID)
	assert.Equal(t, inviteeId, participant.UserID)
	assert.NotNil(t, participant.ID)
	assert.Equal(t, participant.IsActive, true)
	assert.Equal(t, participant.GetEvents()[len(participant.GetEvents())-1], newGroupConversationInvitedEvent(conversation.Conversation.ID, creatorParticipant.UserID, inviteeId))
}

func TestInviteNotActive(t *testing.T) {
	conversation, _, creatorParticipant := createTestGroupConversation()
	userName, _ := NewUserName("test")
	userPassword, _ := NewUserPassword("test", func(p []byte) ([]byte, error) { return p, nil })
	inviteeId := uuid.New()
	user := NewUser(inviteeId, userName, userPassword)
	_ = conversation.Delete(creatorParticipant)

	_, err := conversation.Invite(creatorParticipant, user)

	assert.Equal(t, err.Error(), "conversation is not active")
}

func TestInviteOwner(t *testing.T) {
	conversation, creator, creatorParticipant := createTestGroupConversation()

	_, err := conversation.Invite(creatorParticipant, creator)

	assert.Equal(t, err.Error(), "user is owner")
}

func TestInviteSelf(t *testing.T) {
	conversation, _, _ := createTestGroupConversation()
	userPassword, _ := NewUserPassword("test", func(p []byte) ([]byte, error) { return p, nil })
	userName, _ := NewUserName("test")
	userID := uuid.New()
	user := NewUser(userID, userName, userPassword)
	participant, _ := conversation.Join(user)

	_, err := conversation.Invite(participant, user)

	assert.Equal(t, err.Error(), "cannot invite yourself")
}

func TestLeave(t *testing.T) {
	conversation, _, _ := createTestGroupConversation()

	participant, err := conversation.Leave(&conversation.Owner)

	assert.Nil(t, err)
	assert.Equal(t, participant.IsActive, false)
}

func TestLeaveNotActive(t *testing.T) {
	conversation, _, creatorParticipant := createTestGroupConversation()
	_ = conversation.Delete(creatorParticipant)

	_, err := conversation.Leave(&conversation.Owner)

	assert.Equal(t, err.Error(), "conversation is not active")
}

func TestLeaveNotMember(t *testing.T) {
	conversation, _, _ := createTestGroupConversation()

	participant := NewParticipant(uuid.New(), uuid.New(), uuid.New())

	_, err := conversation.Leave(participant)

	assert.Equal(t, err.Error(), "user is not in conversation")
}

func TestLeaveAlreadyLeft(t *testing.T) {
	conversation, _, _ := createTestGroupConversation()
	_, _ = conversation.Leave(&conversation.Owner)

	_, err := conversation.Leave(&conversation.Owner)

	assert.Equal(t, err.Error(), "user is not in conversation")
}
