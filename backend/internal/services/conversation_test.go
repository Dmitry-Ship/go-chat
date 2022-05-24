package services

import (
	"GitHub/go-chat/backend/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type publicConversationRepositoryMock struct {
	publicConversationOwnerID uuid.UUID
	methodsCalled             map[string]int
}

func (m *publicConversationRepositoryMock) Store(conversation *domain.PublicConversation) error {
	m.methodsCalled["StorePublicConversation"]++
	return nil
}

func (m *publicConversationRepositoryMock) Update(conversation *domain.PublicConversation) error {
	m.methodsCalled["UpdatePublicConversation"]++
	return nil
}

func (m *publicConversationRepositoryMock) GetByID(id uuid.UUID) (*domain.PublicConversation, error) {
	m.methodsCalled["GetPublicConversation"]++

	conversation, err := domain.NewPublicConversation(id, "cool room", m.publicConversationOwnerID)

	return conversation, err
}

type privateConversationRepositoryMock struct {
	methodsCalled map[string]int
}

func (m *privateConversationRepositoryMock) Store(conversation *domain.PrivateConversation) error {
	m.methodsCalled["StorePrivateConversation"]++
	return nil
}

func (m *privateConversationRepositoryMock) Update(conversation *domain.PrivateConversation) error {
	m.methodsCalled["UpdatePrivateConversation"]++
	return nil
}

func (m *privateConversationRepositoryMock) GetByID(id uuid.UUID) (*domain.PrivateConversation, error) {
	m.methodsCalled["GetPrivateConversation"]++

	conversation, err := domain.NewPrivateConversation(id, uuid.New(), uuid.New())

	return conversation, err
}

func (m *privateConversationRepositoryMock) GetID(fromUserID uuid.UUID, toUserID uuid.UUID) (uuid.UUID, error) {
	m.methodsCalled["GetPrivateConversationID"]++

	return uuid.New(), nil
}

type messagesRepositoryMock struct {
	methodsCalled map[string]int
}

func (m *messagesRepositoryMock) StoreTextMessage(message *domain.TextMessage) error {
	m.methodsCalled["StoreTextMessage"]++
	return nil
}

func (m *messagesRepositoryMock) StoreJoinedConversationMessage(message *domain.JoinedConversationMessage) error {
	m.methodsCalled["StoreJoinedConversationMessage"]++
	return nil
}

func (m *messagesRepositoryMock) StoreInvitedConversationMessage(message *domain.JoinedConversationMessage) error {
	m.methodsCalled["StoreInvitedConversationMessage"]++
	return nil
}

func (m *messagesRepositoryMock) StoreRenamedConversationMessage(message *domain.ConversationRenamedMessage) error {
	m.methodsCalled["StoreRenamedConversationMessage"]++
	return nil
}

func (m *messagesRepositoryMock) StoreLeftConversationMessage(message *domain.LeftConversationMessage) error {
	m.methodsCalled["StoreLeftConversationMessage"]++
	return nil
}

type participantsRepositoryMock struct {
	methodsCalled map[string]int
}

func (m *participantsRepositoryMock) Store(participant *domain.Participant) error {
	m.methodsCalled["Store"]++
	return nil
}

func (m *participantsRepositoryMock) Update(participant *domain.Participant) error {
	m.methodsCalled["Update"]++
	return nil
}

func (m *participantsRepositoryMock) GetByConversationIDAndUserID(conversationID uuid.UUID, userID uuid.UUID) (*domain.Participant, error) {
	m.methodsCalled["GetByConversationIDAndUserID"]++
	return domain.NewJoinedParticipant(conversationID, userID), nil
}

func TestCreatePublicConversation(t *testing.T) {
	publicConversationRepository := &publicConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	privateConversationRepository := &privateConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(publicConversationRepository, privateConversationRepository, participantsRepository, messagesRepository)

	err := conversationService.CreatePublicConversation(uuid.New(), "test", uuid.New())

	assert.Nil(t, err)
	assert.Equal(t, 1, publicConversationRepository.methodsCalled["StorePublicConversation"])
}

func TestRenamePublicConversation(t *testing.T) {
	publicConversationRepository := &publicConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	privateConversationRepository := &privateConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(publicConversationRepository, privateConversationRepository, participantsRepository, messagesRepository)

	err := conversationService.RenamePublicConversation(uuid.New(), publicConversationRepository.publicConversationOwnerID, "test")

	assert.Nil(t, err)
	assert.Equal(t, 1, publicConversationRepository.methodsCalled["UpdatePublicConversation"])
}

func TestDeletePublicConversation(t *testing.T) {
	publicConversationRepository := &publicConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	privateConversationRepository := &privateConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(publicConversationRepository, privateConversationRepository, participantsRepository, messagesRepository)

	err := conversationService.DeletePublicConversation(uuid.New(), publicConversationRepository.publicConversationOwnerID)

	assert.Nil(t, err)
	assert.Equal(t, 1, publicConversationRepository.methodsCalled["UpdatePublicConversation"])
}

func TestSendTextMessage(t *testing.T) {
	publicConversationRepository := &publicConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	privateConversationRepository := &privateConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(publicConversationRepository, privateConversationRepository, participantsRepository, messagesRepository)

	conversationID := uuid.New()
	userID := uuid.New()

	err := conversationService.SendPublicTextMessage("test", conversationID, userID)

	assert.Nil(t, err)
	assert.Equal(t, 1, messagesRepository.methodsCalled["StoreTextMessage"])
	assert.Equal(t, 1, participantsRepository.methodsCalled["GetByConversationIDAndUserID"])
}

func TestSendJoinedConversationMessage(t *testing.T) {
	publicConversationRepository := &publicConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	privateConversationRepository := &privateConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(publicConversationRepository, privateConversationRepository, participantsRepository, messagesRepository)
	conversationID := uuid.New()
	userID := uuid.New()

	err := conversationService.SendJoinedConversationMessage(conversationID, userID)

	assert.Nil(t, err)
	assert.Equal(t, 1, messagesRepository.methodsCalled["StoreJoinedConversationMessage"])
}

func TestSendInvitedConversationMessage(t *testing.T) {
	publicConversationRepository := &publicConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	privateConversationRepository := &privateConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(publicConversationRepository, privateConversationRepository, participantsRepository, messagesRepository)
	conversationID := uuid.New()
	userID := uuid.New()

	err := conversationService.SendInvitedConversationMessage(conversationID, userID)

	assert.Nil(t, err)
	assert.Equal(t, 1, messagesRepository.methodsCalled["StoreInvitedConversationMessage"])
}

func TestSendRenamedConversationMessage(t *testing.T) {
	publicConversationRepository := &publicConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	privateConversationRepository := &privateConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(publicConversationRepository, privateConversationRepository, participantsRepository, messagesRepository)
	conversationID := uuid.New()
	newName := "new name"
	userID := uuid.New()

	err := conversationService.SendRenamedConversationMessage(conversationID, userID, newName)

	assert.Nil(t, err)
	assert.Equal(t, 1, messagesRepository.methodsCalled["StoreRenamedConversationMessage"])
}

func TestSendLeftConversationMessage(t *testing.T) {
	publicConversationRepository := &publicConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	privateConversationRepository := &privateConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(publicConversationRepository, privateConversationRepository, participantsRepository, messagesRepository)
	conversationID := uuid.New()
	userID := uuid.New()

	err := conversationService.SendLeftConversationMessage(conversationID, userID)

	assert.Nil(t, err)
	assert.Equal(t, 1, messagesRepository.methodsCalled["StoreLeftConversationMessage"])
}

func TestJoinPublicConversation(t *testing.T) {
	publicConversationRepository := &publicConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	privateConversationRepository := &privateConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(publicConversationRepository, privateConversationRepository, participantsRepository, messagesRepository)

	err := conversationService.JoinPublicConversation(uuid.New(), uuid.New())

	assert.Nil(t, err)
	assert.Equal(t, 1, participantsRepository.methodsCalled["Store"])
}

func TestLeavePublicConversation(t *testing.T) {
	publicConversationRepository := &publicConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	privateConversationRepository := &privateConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(publicConversationRepository, privateConversationRepository, participantsRepository, messagesRepository)
	userID := uuid.New()

	err := conversationService.LeavePublicConversation(uuid.New(), userID)

	assert.Nil(t, err)
	assert.Equal(t, 1, participantsRepository.methodsCalled["Update"])
}
