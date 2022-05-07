package services

import (
	"GitHub/go-chat/backend/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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

func TestJoinPublicConversation(t *testing.T) {
	participantRepository := &participantRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewPublicConversationParticipationService(participantRepository)

	err := conversationService.Join(uuid.New(), uuid.New())

	assert.Nil(t, err)
	assert.Equal(t, 1, participantRepository.methodsCalled["Store"])
}

func TestLeavePublicConversation(t *testing.T) {
	participantRepository := &participantRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewPublicConversationParticipationService(participantRepository)
	userID := uuid.New()

	err := conversationService.Leave(uuid.New(), userID)

	assert.Nil(t, err)
	assert.Equal(t, 1, participantRepository.methodsCalled["Update"])
}
