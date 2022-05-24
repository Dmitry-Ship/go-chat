package services

import (
	"GitHub/go-chat/backend/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type pubicConversationRepositoryMock struct {
	publicConversationOwnerID uuid.UUID
	methodsCalled             map[string]int
}

func (m *pubicConversationRepositoryMock) Store(conversation *domain.PublicConversation) error {
	m.methodsCalled["StorePublicConversation"]++
	return nil
}

func (m *pubicConversationRepositoryMock) Update(conversation *domain.PublicConversation) error {
	m.methodsCalled["UpdatePublicConversation"]++
	return nil
}

func (m *pubicConversationRepositoryMock) GetByID(id uuid.UUID) (*domain.PublicConversation, error) {
	m.methodsCalled["GetPublicConversation"]++

	conversation, err := domain.NewPublicConversation(id, "cool room", m.publicConversationOwnerID)

	return conversation, err
}

func TestCreatePublicConversation(t *testing.T) {
	conversationRepository := &pubicConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewPublicConversationEditingService(conversationRepository)

	err := conversationService.Create(uuid.New(), "test", uuid.New())

	assert.Nil(t, err)
	assert.Equal(t, 1, conversationRepository.methodsCalled["StorePublicConversation"])
}

func TestRenamePublicConversation(t *testing.T) {
	conversationRepository := &pubicConversationRepositoryMock{
		publicConversationOwnerID: uuid.New(),
		methodsCalled:             make(map[string]int),
	}
	conversationService := NewPublicConversationEditingService(conversationRepository)

	err := conversationService.Rename(uuid.New(), conversationRepository.publicConversationOwnerID, "test")

	assert.Nil(t, err)
	assert.Equal(t, 1, conversationRepository.methodsCalled["UpdatePublicConversation"])
}

func TestDeletePublicConversation(t *testing.T) {
	conversationRepository := &pubicConversationRepositoryMock{
		publicConversationOwnerID: uuid.New(),
		methodsCalled:             make(map[string]int),
	}

	conversationService := NewPublicConversationEditingService(conversationRepository)

	err := conversationService.Delete(uuid.New(), conversationRepository.publicConversationOwnerID)

	assert.Nil(t, err)
	assert.Equal(t, 1, conversationRepository.methodsCalled["UpdatePublicConversation"])
}
