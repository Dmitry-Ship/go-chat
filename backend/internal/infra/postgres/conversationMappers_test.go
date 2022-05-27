package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToGroupConversationPersistence(t *testing.T) {
	name, _ := domain.NewConversationName("cool room")
	conversation, _ := domain.NewGroupConversation(uuid.New(), name, uuid.New())

	persistence := toGroupConversationPersistence(conversation)

	assert.Equal(t, persistence.ID, conversation.ID)
	assert.Equal(t, persistence.ConversationID, conversation.GetBaseData().ID)
	assert.Equal(t, persistence.Name, conversation.Name.String())
	assert.Equal(t, persistence.Avatar, conversation.Avatar)
	assert.Equal(t, persistence.OwnerID, conversation.Owner.UserID)
}

func TestToGroupConversationDomain(t *testing.T) {
	conversationId := uuid.New()
	userID := uuid.New()

	groupConversation := &GroupConversation{
		ID:             uuid.New(),
		ConversationID: conversationId,
		Name:           "cool room",
		Avatar:         "avatar",
		OwnerID:        userID,
	}

	conversation := &Conversation{
		ID:   conversationId,
		Type: 0,
	}

	participant := &Participant{
		ID:             uuid.New(),
		ConversationID: conversationId,
		UserID:         userID,
	}

	domain := toGroupConversationDomain(conversation, groupConversation, participant)

	assert.Equal(t, domain.ID, groupConversation.ID)
	assert.Equal(t, domain.GetBaseData().ID, groupConversation.ConversationID)
	assert.Equal(t, domain.Name.String(), groupConversation.Name)
	assert.Equal(t, domain.Avatar, groupConversation.Avatar)
	assert.Equal(t, domain.Owner.UserID, groupConversation.OwnerID)
}

func TestToDirectConversationPersistence(t *testing.T) {
	conversation, _ := domain.NewDirectConversation(uuid.New(), uuid.New(), uuid.New())

	persistence := toDirectConversationPersistence(conversation)

	assert.Equal(t, persistence.ID, conversation.ID)
	assert.Equal(t, persistence.ConversationID, conversation.GetBaseData().ID)
	assert.Equal(t, persistence.FromUserID, conversation.FromUser.UserID)
	assert.Equal(t, persistence.ToUserID, conversation.ToUser.UserID)
}
