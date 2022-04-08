package application

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/mappers"
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
	conversations domain.ConversationRepository
	participants  domain.ParticipantRepository
	users         domain.UserRepository
	messages      domain.ChatMessageRepository
	hub           ws.Hub
}

func NewConversationCommandService(conversations domain.ConversationRepository, participants domain.ParticipantRepository, users domain.UserRepository, messages domain.ChatMessageRepository, hub ws.Hub) *conversationCommandService {
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

	err := s.notifyAllParticipants(id, notification)

	if err != nil {
		return err
	}

	err = s.conversations.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *conversationCommandService) notifyAllParticipants(conversationID uuid.UUID, notification ws.OutgoingNotification) error {
	participants, err := s.participants.FindAllByConversationID(conversationID)

	if err != nil {
		return err
	}

	for _, participant := range participants {
		if notification.Type == "message" {
			message := notification.Payload.(MessageFullDTO)

			if message.Type != "system" {
				message.IsInbound = participant.UserID != message.User.ID
				notification.Payload = message
			}

		}

		s.hub.BroadcastToClients(notification, participant.UserID)
	}

	return nil
}

func (s *conversationCommandService) SendUserMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewUserMessage(messageText, conversationId, userId)

	err := s.messages.Store(message)

	if err != nil {
		return err
	}

	user, err := s.users.FindByID(userId)

	if err != nil {
		return err
	}

	notification := ws.OutgoingNotification{
		Type: "message",
		Payload: MessageFullDTO{
			User:       mappers.ToUserDTO(user),
			MessageDTO: mappers.ToMessageDTO(message),
		},
	}

	err = s.notifyAllParticipants(conversationId, notification)

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

	notification := ws.OutgoingNotification{
		Type: "message",
		Payload: MessageFullDTO{
			MessageDTO: mappers.ToMessageDTO(message),
		},
	}

	err = s.notifyAllParticipants(conversationId, notification)

	if err != nil {
		return err
	}

	return nil
}
