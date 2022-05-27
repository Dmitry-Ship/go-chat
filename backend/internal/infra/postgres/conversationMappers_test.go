package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToGroupConversationPersistence(t *testing.T) {
	conversation, _ := domain.NewGroupConversation(uuid.New(), "cool room", uuid.New())

	persistence := toGroupConversationPersistence(conversation)

	assert.Equal(t, persistence.ID, conversation.Data.ID)
	assert.Equal(t, persistence.ConversationID, conversation.GetBaseData().ID)
	assert.Equal(t, persistence.Name, conversation.Data.Name)
	assert.Equal(t, persistence.Avatar, conversation.Data.Avatar)
	assert.Equal(t, persistence.OwnerID, conversation.Data.Owner.UserID)
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

	assert.Equal(t, domain.Data.ID, groupConversation.ID)
	assert.Equal(t, domain.GetBaseData().ID, groupConversation.ConversationID)
	assert.Equal(t, domain.Data.Name, groupConversation.Name)
	assert.Equal(t, domain.Data.Avatar, groupConversation.Avatar)
	assert.Equal(t, domain.Data.Owner.UserID, groupConversation.OwnerID)
}

func TestToDirectConversationPersistence(t *testing.T) {
	conversation, _ := domain.NewDirectConversation(uuid.New(), uuid.New(), uuid.New())

	persistence := toDirectConversationPersistence(conversation)

	assert.Equal(t, persistence.ID, conversation.Data.ID)
	assert.Equal(t, persistence.ConversationID, conversation.GetBaseData().ID)
	assert.Equal(t, persistence.FromUserID, conversation.Data.FromUser.UserID)
	assert.Equal(t, persistence.ToUserID, conversation.Data.ToUser.UserID)
}
