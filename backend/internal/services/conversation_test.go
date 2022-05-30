package services

import (
	"GitHub/go-chat/backend/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type groupConversationRepositoryMock struct {
	groupConversationOwnerID uuid.UUID
	methodsCalled            map[string]int
}

func (m *groupConversationRepositoryMock) Store(conversation *domain.GroupConversation) error {
	m.methodsCalled["Store"]++
	return nil
}

func (m *groupConversationRepositoryMock) Update(conversation *domain.GroupConversation) error {
	m.methodsCalled["Update"]++
	return nil
}

func (m *groupConversationRepositoryMock) GetByID(id uuid.UUID) (*domain.GroupConversation, error) {
	m.methodsCalled["GetByID"]++

	name, _ := domain.NewConversationName("cool room")

	conversation, err := domain.NewGroupConversation(id, name, m.groupConversationOwnerID)

	return conversation, err
}

type directConversationRepositoryMock struct {
	methodsCalled map[string]int
}

func (m *directConversationRepositoryMock) Store(conversation *domain.DirectConversation) error {
	m.methodsCalled["Store"]++
	return nil
}

func (m *directConversationRepositoryMock) Update(conversation *domain.DirectConversation) error {
	m.methodsCalled["Update"]++
	return nil
}

func (m *directConversationRepositoryMock) GetByID(id uuid.UUID) (*domain.DirectConversation, error) {
	m.methodsCalled["GetByID"]++

	conversation, err := domain.NewDirectConversation(id, uuid.New(), uuid.New())

	return conversation, err
}

func (m *directConversationRepositoryMock) GetID(fromUserID uuid.UUID, toUserID uuid.UUID) (uuid.UUID, error) {
	m.methodsCalled["GetID"]++

	return uuid.New(), nil
}

type messagesRepositoryMock struct {
	methodsCalled map[string]int
}

func (m *messagesRepositoryMock) Store(message *domain.Message) error {
	m.methodsCalled["Store"]++
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
	return domain.NewParticipant(uuid.New(), conversationID, userID), nil
}

func (m *participantsRepositoryMock) GetIDsByConversationID(conversationID uuid.UUID) ([]uuid.UUID, error) {
	m.methodsCalled["GetIDsByConversationID"]++

	return []uuid.UUID{uuid.New()}, nil
}

func TestCreateGroupConversation(t *testing.T) {
	groupConversationRepository := &groupConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	directConversationRepository := &directConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(groupConversationRepository, directConversationRepository, participantsRepository, messagesRepository)

	err := conversationService.CreateGroupConversation(uuid.New(), "test", uuid.New())

	assert.Nil(t, err)
	assert.Equal(t, 1, groupConversationRepository.methodsCalled["Store"])
}

func TestRenameGroupConversation(t *testing.T) {
	groupConversationRepository := &groupConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	directConversationRepository := &directConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(groupConversationRepository, directConversationRepository, participantsRepository, messagesRepository)

	err := conversationService.RenameGroupConversation(uuid.New(), groupConversationRepository.groupConversationOwnerID, "test")

	assert.Nil(t, err)
	assert.Equal(t, 1, groupConversationRepository.methodsCalled["Update"])
}

func TestDeleteGroupConversation(t *testing.T) {
	groupConversationRepository := &groupConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	directConversationRepository := &directConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(groupConversationRepository, directConversationRepository, participantsRepository, messagesRepository)

	err := conversationService.DeleteGroupConversation(uuid.New(), groupConversationRepository.groupConversationOwnerID)

	assert.Nil(t, err)
	assert.Equal(t, 1, groupConversationRepository.methodsCalled["Update"])
}

func TestSendTextMessage(t *testing.T) {
	groupConversationRepository := &groupConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	directConversationRepository := &directConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(groupConversationRepository, directConversationRepository, participantsRepository, messagesRepository)

	conversationID := uuid.New()
	userID := uuid.New()

	err := conversationService.SendGroupTextMessage(conversationID, userID, "test")

	assert.Nil(t, err)
	assert.Equal(t, 1, messagesRepository.methodsCalled["Store"])
	assert.Equal(t, 1, participantsRepository.methodsCalled["GetByConversationIDAndUserID"])
}

func TestSendJoinedConversationMessage(t *testing.T) {
	groupConversationRepository := &groupConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	directConversationRepository := &directConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(groupConversationRepository, directConversationRepository, participantsRepository, messagesRepository)
	conversationID := uuid.New()
	userID := uuid.New()

	err := conversationService.SendJoinedConversationMessage(conversationID, userID)

	assert.Nil(t, err)
	assert.Equal(t, 1, messagesRepository.methodsCalled["Store"])
}

func TestSendInvitedConversationMessage(t *testing.T) {
	groupConversationRepository := &groupConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	directConversationRepository := &directConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(groupConversationRepository, directConversationRepository, participantsRepository, messagesRepository)
	conversationID := uuid.New()
	userID := uuid.New()

	err := conversationService.SendInvitedConversationMessage(conversationID, userID)

	assert.Nil(t, err)
	assert.Equal(t, 1, messagesRepository.methodsCalled["Store"])
}

func TestSendRenamedConversationMessage(t *testing.T) {
	groupConversationRepository := &groupConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	directConversationRepository := &directConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(groupConversationRepository, directConversationRepository, participantsRepository, messagesRepository)
	conversationID := uuid.New()
	newName := "new name"
	userID := uuid.New()

	err := conversationService.SendRenamedConversationMessage(conversationID, userID, newName)

	assert.Nil(t, err)
	assert.Equal(t, 1, messagesRepository.methodsCalled["Store"])
}

func TestSendLeftConversationMessage(t *testing.T) {
	groupConversationRepository := &groupConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	directConversationRepository := &directConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(groupConversationRepository, directConversationRepository, participantsRepository, messagesRepository)
	conversationID := uuid.New()
	userID := uuid.New()

	err := conversationService.SendLeftConversationMessage(conversationID, userID)

	assert.Nil(t, err)
	assert.Equal(t, 1, messagesRepository.methodsCalled["Store"])
}

func TestJoinGroupConversation(t *testing.T) {
	groupConversationRepository := &groupConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	directConversationRepository := &directConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(groupConversationRepository, directConversationRepository, participantsRepository, messagesRepository)

	err := conversationService.JoinGroupConversation(uuid.New(), uuid.New())

	assert.Nil(t, err)
	assert.Equal(t, 1, participantsRepository.methodsCalled["Store"])
}

func TestLeaveGroupConversation(t *testing.T) {
	groupConversationRepository := &groupConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	directConversationRepository := &directConversationRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	participantsRepository := &participantsRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	conversationService := NewConversationService(groupConversationRepository, directConversationRepository, participantsRepository, messagesRepository)
	userID := uuid.New()

	err := conversationService.LeaveGroupConversation(uuid.New(), userID)

	assert.Nil(t, err)
	assert.Equal(t, 1, participantsRepository.methodsCalled["Update"])
}
