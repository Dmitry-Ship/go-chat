package services

import (
	"GitHub/go-chat/backend/pkg/readModel"
	ws "GitHub/go-chat/backend/pkg/websocket"

	"github.com/google/uuid"
)

type ConversationWSResolver interface {
	NotifyAboutMessage(conversationId uuid.UUID, messageID uuid.UUID, userId uuid.UUID) error
	NotifyAboutConversationDeletion(id uuid.UUID) error
	NotifyAboutConversationRenamed(conversationId uuid.UUID, newName string) error
}

type conversationWSResolver struct {
	participants readModel.ParticipantQueryRepository
	messages     readModel.MessageQueryRepository
	hub          ws.Hub
}

func NewConversationWSResolver(
	participants readModel.ParticipantQueryRepository,
	messages readModel.MessageQueryRepository,
	hub ws.Hub,
) *conversationWSResolver {
	return &conversationWSResolver{
		participants: participants,
		messages:     messages,
		hub:          hub,
	}
}

func (s *conversationWSResolver) notifyParticipants(conversationID uuid.UUID, notification ws.OutgoingNotification) error {
	ids, err := s.participants.GetUserIdsByConversationID(conversationID)

	if err != nil {
		return err
	}

	s.hub.BroadcastToClients(notification, ids)

	return nil
}

func (s *conversationWSResolver) NotifyAboutMessage(conversationId uuid.UUID, messageID uuid.UUID, userId uuid.UUID) error {
	messageDTO, err := s.messages.GetMessageByID(messageID, userId)

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

func (s *conversationWSResolver) NotifyAboutConversationDeletion(id uuid.UUID) error {
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

	return nil
}

func (s *conversationWSResolver) NotifyAboutConversationRenamed(conversationId uuid.UUID, newName string) error {
	notification := ws.OutgoingNotification{
		Type: "conversation_renamed",
		Payload: struct {
			ConversationId uuid.UUID `json:"conversation_id"`
			NewName        string    `json:"new_name"`
		}{
			ConversationId: conversationId,
			NewName:        newName,
		},
	}

	err := s.notifyParticipants(conversationId, notification)

	if err != nil {
		return err
	}

	return nil
}
