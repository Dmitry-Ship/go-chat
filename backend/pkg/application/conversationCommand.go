package application

import (
	"GitHub/go-chat/backend/domain"
	ws "GitHub/go-chat/backend/pkg/websocket"

	"fmt"

	"github.com/google/uuid"
)

type ConversationCommandService interface {
	CreatePublicConversation(id uuid.UUID, name string, userId uuid.UUID) error
	JoinPublicConversation(conversationId uuid.UUID, userId uuid.UUID) error
	LeavePublicConversation(conversationId uuid.UUID, userId uuid.UUID) error
	DeleteConversation(id uuid.UUID) error
	SendUserMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error
	SendSystemMessage(messageText string, conversationId uuid.UUID) error
}

type conversationCommandService struct {
	conversations domain.ConversationCommandRepository
	participants  domain.ParticipantCommandRepository
	users         domain.UserCommandRepository
	messages      domain.MessageCommandRepository
	hub           ws.Hub
}

func NewConversationCommandService(
	conversations domain.ConversationCommandRepository,
	participants domain.ParticipantCommandRepository,
	users domain.UserCommandRepository,
	messages domain.MessageCommandRepository,
	hub ws.Hub,
) *conversationCommandService {
	return &conversationCommandService{
		conversations: conversations,
		users:         users,
		participants:  participants,
		messages:      messages,
		hub:           hub,
	}
}

func (s *conversationCommandService) CreatePublicConversation(id uuid.UUID, name string, userId uuid.UUID) error {
	conversation := domain.NewConversation(id, name, false)
	err := s.conversations.Store(conversation)

	if err != nil {
		return err
	}

	err = s.JoinPublicConversation(conversation.ID, userId)

	if err != nil {
		return err
	}

	return nil
}

func (s *conversationCommandService) JoinPublicConversation(conversationID uuid.UUID, userId uuid.UUID) error {
	err := s.participants.Store(domain.NewParticipant(conversationID, userId))

	if err != nil {
		return err
	}

	user, err := s.users.FindByID(userId)

	if err != nil {
		return err
	}

	err = s.SendSystemMessage(fmt.Sprintf("%s joined", user.Name), conversationID)

	if err != nil {
		return err
	}

	return nil
}

func (s *conversationCommandService) LeavePublicConversation(conversationID uuid.UUID, userId uuid.UUID) error {
	err := s.participants.DeleteByConversationIDAndUserID(conversationID, userId)

	if err != nil {
		return err
	}

	user, err := s.users.FindByID(userId)

	if err != nil {
		return err
	}

	err = s.SendSystemMessage(fmt.Sprintf("%s left", user.Name), conversationID)

	if err != nil {
		return err
	}

	return nil
}

func (s *conversationCommandService) DeleteConversation(id uuid.UUID) error {
	notification := ws.OutgoingNotification{
		Type: "conversation_deleted",
		Payload: struct {
			ConversationId uuid.UUID `json:"conversation_id"`
		}{
			ConversationId: id,
		},
	}

	err := s.notifyParticipants(id, notification)

	if err != nil {
		return err
	}

	err = s.conversations.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *conversationCommandService) notifyParticipants(conversationID uuid.UUID, notification ws.OutgoingNotification) error {
	ids, err := s.participants.GetUserIdsByConversationID(conversationID)

	if err != nil {
		return err
	}

	s.hub.BroadcastToClients(notification, ids)

	return nil
}

func (s *conversationCommandService) SendUserMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewUserMessage(messageText, conversationId, userId)

	err := s.messages.Store(message)

	if err != nil {
		return err
	}

	messageDTO, err := s.messages.FindByID(message.ID, userId)

	if err != nil {
		return err
	}

	notification := ws.OutgoingNotification{
		Type:    "message",
		Payload: messageDTO,
	}

	err = s.notifyParticipants(conversationId, notification)

	if err != nil {
		return err
	}

	return nil
}

func (s *conversationCommandService) SendSystemMessage(messageText string, conversationId uuid.UUID) error {
	message := domain.NewSystemMessage(messageText, conversationId)

	err := s.messages.Store(message)

	if err != nil {
		return err
	}

	messageDTO, err := s.messages.FindByID(message.ID, uuid.Nil)

	if err != nil {
		return err
	}

	notification := ws.OutgoingNotification{
		Type:    "message",
		Payload: messageDTO,
	}

	err = s.notifyParticipants(conversationId, notification)

	if err != nil {
		return err
	}

	return nil
}
