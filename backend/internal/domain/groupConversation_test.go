package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createTestGroupConversation() (*GroupConversation, *User, *Participant) {
	conversationID := uuid.New()
	creatorId := uuid.New()
	userName := "test"
	userPassword, _ := HashPassword("test")
	creatorUser := NewUser(creatorId, userName, userPassword)
	conversation, _ := NewGroupConversation(conversationID, "test", *creatorUser)

	return conversation, creatorUser, &conversation.Owner
}

func createTestUser() *User {
	userPassword, _ := HashPassword("test")
	userName := "test"
	userID := uuid.New()
	user := NewUser(userID, userName, userPassword)

	return user
}

func TestNewGroupConversation(t *testing.T) {
	conversationID := uuid.New()
	creatorId := uuid.New()
	userName := "test"
	userPassword, _ := HashPassword("test")
	creator := NewUser(creatorId, userName, userPassword)
	name := "test"

	conversation, err := NewGroupConversation(conversationID, name, *creator)

	assert.Equal(t, conversation.Conversation.ID, conversationID)
	assert.Equal(t, name, conversation.Name)
	assert.Equal(t, string(name[0]), conversation.Avatar)
	assert.Equal(t, conversation.Type, ConversationTypeGroup)
	assert.Equal(t, conversation.Conversation.ID, conversation.Owner.ConversationID)
	assert.Equal(t, creatorId, conversation.Owner.UserID)
	assert.NotNil(t, conversation.Owner.ID)
	assert.Equal(t, conversation.IsActive, true)
	assert.Equal(t, conversation.GetEvents()[len(conversation.GetEvents())-1], newGroupConversationCreatedEvent(conversation.Conversation.ID, creatorId))
	assert.Nil(t, err)
}

func TestValidateConversationName(t *testing.T) {
	name := "test"
	err := ValidateConversationName(name)

	assert.Nil(t, err)
	assert.Equal(t, "test", name)
}

func TestNewConversationNameErrors(t *testing.T) {
	type testCase struct {
		name        string
		expectedErr error
	}

	longName := ""

	for i := 0; i < 101; i++ {
		longName += "a"
	}

	testCases := []testCase{
		{
			name:        "",
			expectedErr: errors.New("name is empty"),
		}, {
			name:        longName,
			expectedErr: errors.New("name is too long"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateConversationName(tc.name)

			assert.Equal(t, err, tc.expectedErr)
		})
	}
}

func TestRename(t *testing.T) {
	conversation, _, creatorParticipant := createTestGroupConversation()

	newName := "new name"

	err := conversation.Rename(newName, creatorParticipant)

	assert.Nil(t, err)
	assert.Equal(t, "new name", conversation.Name)
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

	assert.Equal(t, ErrorUserNotInConversation, err)
}

func TestSendTextMessageNotActive(t *testing.T) {
	conversation, _, creatorParticipant := createTestGroupConversation()
	messageID := uuid.New()
	_ = conversation.Delete(creatorParticipant)

	_, err := conversation.SendTextMessage(messageID, "new message", creatorParticipant)

	assert.Equal(t, ErrorConversationNotActive, err)
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
	newName := "new name"
	user := createTestUser()
	participant, _ := conversation.Join(*user)

	err := conversation.Rename(newName, participant)

	assert.NotNil(t, err)
	assert.Equal(t, ErrorUserNotOwner, err)
	assert.NotEqual(t, newName, conversation.Name)
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
	user := createTestUser()
	participant, _ := conversation.Join(*user)

	err := conversation.Delete(participant)

	assert.NotNil(t, err)
	assert.Equal(t, ErrorUserNotOwner, err)
	assert.Equal(t, true, conversation.IsActive)
	assert.Equal(t, conversation.GetEvents()[len(conversation.GetEvents())-1], newGroupConversationCreatedEvent(conversation.Conversation.ID, creator.ID))
}

func TestDeleteNotActive(t *testing.T) {
	conversation, _, creatorParticipant := createTestGroupConversation()
	_ = conversation.Delete(creatorParticipant)

	err := conversation.Delete(creatorParticipant)

	assert.Equal(t, ErrorConversationNotActive, err)
}

func TestJoin(t *testing.T) {
	conversation, _, _ := createTestGroupConversation()
	user := createTestUser()

	participant, err := conversation.Join(*user)

	assert.Nil(t, err)
	assert.Equal(t, conversation.Conversation.ID, participant.ConversationID)
	assert.Equal(t, user.ID, participant.UserID)
	assert.NotNil(t, participant.ID)
	assert.Equal(t, participant.IsActive, true)
	assert.Equal(t, participant.GetEvents()[len(participant.GetEvents())-1], newGroupConversationJoinedEvent(conversation.Conversation.ID, user.ID))
}

func TestJoinNotActive(t *testing.T) {
	conversation, creator, creatorParticipant := createTestGroupConversation()

	_ = conversation.Delete(creatorParticipant)

	_, err := conversation.Join(*creator)

	assert.Equal(t, ErrorConversationNotActive, err)
}

func TestKick(t *testing.T) {
	conversation, _, creatorParticipant := createTestGroupConversation()
	user := createTestUser()
	participant, _ := conversation.Join(*user)

	participant, err := conversation.Kick(creatorParticipant, participant)

	assert.Nil(t, err)
	assert.Equal(t, participant.IsActive, false)
	assert.Equal(t, participant.GetEvents()[len(participant.GetEvents())-1], newGroupConversationLeftEvent(conversation.Conversation.ID, user.ID))
}

func TestKickErrors(t *testing.T) {
	conversation, _, owner := createTestGroupConversation()
	conversation2, _, _ := createTestGroupConversation()
	user := createTestUser()
	participant, _ := conversation.Join(*user)
	user2 := createTestUser()
	participant2, _ := conversation.Join(*user2)
	participantFromAnotherConversation, _ := conversation2.Join(*user)

	type testCase struct {
		name        string
		kicker      Participant
		target      Participant
		expectedErr error
	}

	testCases := []testCase{
		{
			name:        "not owner",
			kicker:      *participant,
			target:      *participant2,
			expectedErr: ErrorUserNotOwner,
		}, {
			name:        "not in conversation",
			kicker:      *owner,
			target:      *participantFromAnotherConversation,
			expectedErr: ErrorUserNotInConversation,
		}, {
			name:        "kick oneself",
			kicker:      *owner,
			target:      *owner,
			expectedErr: ErrorCannotKickOneself,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kicked, err := conversation.Kick(&tc.kicker, &tc.target)

			assert.Nil(t, kicked)
			assert.Equal(t, err, tc.expectedErr)
		})
	}
}

func TestInvite(t *testing.T) {
	conversation, _, creatorParticipant := createTestGroupConversation()
	user := createTestUser()

	participant, err := conversation.Invite(creatorParticipant, user)

	assert.Nil(t, err)
	assert.Equal(t, conversation.Conversation.ID, participant.ConversationID)
	assert.Equal(t, user.ID, participant.UserID)
	assert.NotNil(t, participant.ID)
	assert.Equal(t, participant.IsActive, true)
	assert.Equal(t, participant.GetEvents()[len(participant.GetEvents())-1], newGroupConversationInvitedEvent(conversation.Conversation.ID, creatorParticipant.UserID, user.ID))
}

func TestInviteNotActive(t *testing.T) {
	conversation, _, creatorParticipant := createTestGroupConversation()
	user := createTestUser()
	_ = conversation.Delete(creatorParticipant)

	_, err := conversation.Invite(creatorParticipant, user)

	assert.Equal(t, ErrorConversationNotActive, err)
}

func TestInviteSelf(t *testing.T) {
	conversation, _, _ := createTestGroupConversation()
	user := createTestUser()
	participant, _ := conversation.Join(*user)

	_, err := conversation.Invite(participant, user)

	assert.Equal(t, err, ErrorCannotInviteOneself)
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

	assert.Equal(t, ErrorConversationNotActive, err)
}

func TestLeaveNotMember(t *testing.T) {
	conversation, _, _ := createTestGroupConversation()

	participant := NewParticipant(uuid.New(), uuid.New(), uuid.New())

	_, err := conversation.Leave(participant)

	assert.Equal(t, ErrorUserNotInConversation, err)
}

func TestLeaveAlreadyLeft(t *testing.T) {
	conversation, _, _ := createTestGroupConversation()
	_, _ = conversation.Leave(&conversation.Owner)

	_, err := conversation.Leave(&conversation.Owner)

	assert.Equal(t, ErrorUserNotInConversation, err)
}
