package services

import (
	"GitHub/go-chat/backend/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type conversationRepositoryMock struct {
	publicConversationOwnerID uuid.UUID
	methodsCalled             map[string]int
}

func (m *conversationRepositoryMock) StorePublicConversation(conversation *domain.PublicConversation) error {
	m.methodsCalled["StorePublicConversation"]++
	return nil
}

func (m *conversationRepositoryMock) StorePrivateConversation(conversation *domain.PrivateConversation) error {
	m.methodsCalled["StorePrivateConversation"]++
	return nil
}

func (m *conversationRepositoryMock) UpdatePublicConversation(conversation *domain.PublicConversation) error {
	m.methodsCalled["UpdatePublicConversation"]++
	return nil
}

func (m *conversationRepositoryMock) GetPublicConversation(id uuid.UUID) (*domain.PublicConversation, error) {
	m.methodsCalled["GetPublicConversation"]++
	return domain.NewPublicConversation(id, "cool room", m.publicConversationOwnerID), nil
}

func (m *conversationRepositoryMock) GetPrivateConversationID(firstUserId uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error) {
	m.methodsCalled["GetPrivateConversationID"]++
	return uuid.Nil, nil
}

type participantRepositoryMock struct {
	methodsCalled map[string]int
}

func (m *participantRepositoryMock) Store(participant *domain.Participant) error {
	m.methodsCalled["Store"]++
	return nil
}

func (m *participantRepositoryMock) Update(participant *domain.Participant) error {
	m.methodsCalled["Update"]++
	return nil
}

func (m *participantRepositoryMock) GetByConversationIDAndUserID(conversationID uuid.UUID, userID uuid.UUID) (*domain.Participant, error) {
	m.methodsCalled["GetByConversationIDAndUserID"]++
	return domain.NewJoinedParticipant(conversationID, userID), nil
}

func TestCreatePublicConversation(t *testing.T) {
	conversationRepository := &conversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantRepository := &participantRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(conversationRepository, participantRepository)

	err := conversationService.CreatePublicConversation(uuid.New(), "test", uuid.New())

	assert.Nil(t, err)
	assert.Equal(t, 1, conversationRepository.methodsCalled["StorePublicConversation"])
}

func TestJoinPublicConversation(t *testing.T) {
	conversationRepository := &conversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantRepository := &participantRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(conversationRepository, participantRepository)

	err := conversationService.JoinPublicConversation(uuid.New(), uuid.New())

	assert.Nil(t, err)
	assert.Equal(t, 1, participantRepository.methodsCalled["Store"])
}

func TestRenamePublicConversation(t *testing.T) {
	conversationRepository := &conversationRepositoryMock{
		publicConversationOwnerID: uuid.New(),
		methodsCalled:             make(map[string]int),
	}
	participantRepository := &participantRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(conversationRepository, participantRepository)

	err := conversationService.RenamePublicConversation(uuid.New(), conversationRepository.publicConversationOwnerID, "test")

	assert.Nil(t, err)
	assert.Equal(t, 1, conversationRepository.methodsCalled["UpdatePublicConversation"])
}

func TestLeavePublicConversation(t *testing.T) {
	conversationRepository := &conversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantRepository := &participantRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(conversationRepository, participantRepository)
	userID := uuid.New()

	err := conversationService.LeavePublicConversation(uuid.New(), userID)

	assert.Nil(t, err)
	assert.Equal(t, 1, participantRepository.methodsCalled["Update"])
}

func TestDeletePublicConversation(t *testing.T) {
	conversationRepository := &conversationRepositoryMock{
		publicConversationOwnerID: uuid.New(),
		methodsCalled:             make(map[string]int),
	}
	participantRepository := &participantRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(conversationRepository, participantRepository)

	err := conversationService.DeletePublicConversation(uuid.New(), conversationRepository.publicConversationOwnerID)

	assert.Nil(t, err)
	assert.Equal(t, 1, conversationRepository.methodsCalled["UpdatePublicConversation"])
}
