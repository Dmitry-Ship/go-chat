package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToGroupConversationPersistence(t *testing.T) {
	name, _ := domain.NewConversationName("cool room")
	userName, _ := domain.NewUserName("test")
	userPassword, _ := domain.NewUserPassword("test", func(p []byte) ([]byte, error) { return p, nil })
	creator := domain.NewUser(uuid.New(), userName, userPassword)

	conversation, _ := domain.NewGroupConversation(uuid.New(), name, *creator)

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

	groupConversation := GroupConversation{
		ID:             uuid.New(),
		ConversationID: conversationId,
		Name:           "cool room",
		Avatar:         "avatar",
		OwnerID:        userID,
	}

	conversation := Conversation{
		ID:   conversationId,
		Type: 0,
	}

	participant := Participant{
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

func TestToDirectConversationDomain(t *testing.T) {
	conversationId := uuid.New()
	userID := uuid.New()

	conversation := &Conversation{
		ID:   conversationId,
		Type: 1,
	}

	participants := []*Participant{
		{
			ID:             uuid.New(),
			ConversationID: conversationId,
			UserID:         userID,
		},
		{
			ID:             uuid.New(),
			ConversationID: conversationId,
			UserID:         uuid.New(),
		},
	}

	domain := toDirectConversationDomain(conversation, participants)

	assert.Equal(t, domain.ID, conversationId)
	assert.Equal(t, domain.Type, domain.Type)
	assert.Equal(t, len(domain.Participants), 2)
}
