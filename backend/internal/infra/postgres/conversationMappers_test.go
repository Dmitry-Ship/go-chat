package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToPublicConversationPersistence(t *testing.T) {
	conversation, _ := domain.NewPublicConversation(uuid.New(), "cool room", uuid.New())

	persistence := toPublicConversationPersistence(conversation)

	assert.Equal(t, persistence.ID, conversation.Data.ID)
	assert.Equal(t, persistence.ConversationID, conversation.GetBaseData().ID)
	assert.Equal(t, persistence.Name, conversation.Data.Name)
	assert.Equal(t, persistence.Avatar, conversation.Data.Avatar)
	assert.Equal(t, persistence.OwnerID, conversation.Data.Owner.UserID)
}

func TestToPublicConversationDomain(t *testing.T) {
	conversationId := uuid.New()
	userId := uuid.New()

	publicConversation := &PublicConversation{
		ID:             uuid.New(),
		ConversationID: conversationId,
		Name:           "cool room",
		Avatar:         "avatar",
		OwnerID:        userId,
	}

	conversation := &Conversation{
		ID:        conversationId,
		CreatedAt: time.Now(),
		Type:      0,
	}

	participant := &Participant{
		ID:             uuid.New(),
		ConversationID: conversationId,
		UserID:         userId,
	}

	domain := toPublicConversationDomain(conversation, publicConversation, participant)

	assert.Equal(t, domain.Data.ID, publicConversation.ID)
	assert.Equal(t, domain.GetBaseData().ID, publicConversation.ConversationID)
	assert.Equal(t, domain.Data.Name, publicConversation.Name)
	assert.Equal(t, domain.Data.Avatar, publicConversation.Avatar)
	assert.Equal(t, domain.Data.Owner.UserID, publicConversation.OwnerID)
}

func TestToPrivateConversationPersistence(t *testing.T) {
	conversation, _ := domain.NewPrivateConversation(uuid.New(), uuid.New(), uuid.New())

	persistence := toPrivateConversationPersistence(conversation)

	assert.Equal(t, persistence.ID, conversation.Data.ID)
	assert.Equal(t, persistence.ConversationID, conversation.GetBaseData().ID)
	assert.Equal(t, persistence.FromUserID, conversation.Data.FromUser.UserID)
	assert.Equal(t, persistence.ToUserID, conversation.Data.ToUser.UserID)
}
