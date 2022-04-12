package application

import (
	"GitHub/go-chat/backend/pkg/readModel"
	ws "GitHub/go-chat/backend/pkg/websocket"
	"fmt"

	"github.com/google/uuid"
)

type ConversationWSResolver interface {
	DispatchMessage(messageID uuid.UUID, conversationId uuid.UUID, userId uuid.UUID)
	DispatchConversationDeleted(conversationId uuid.UUID)
	DispatchConversationRenamed(conversationId uuid.UUID, newName string)
}

type messageChannelItem struct {
	senderUserID   uuid.UUID
	messageID      uuid.UUID
	conversationID uuid.UUID
}

type conversationRenamedItem struct {
	conversationID uuid.UUID
	newName        string
}

type conversationWSResolver struct {
	participants               readModel.ParticipantQueryRepository
	messages                   readModel.MessageQueryRepository
	hub                        ws.Hub
	messageChannel             chan messageChannelItem
	conversationDeletedChannel chan uuid.UUID
	conversationRenamedChannel chan conversationRenamedItem
}

func NewConversationWSResolver(
	participants readModel.ParticipantQueryRepository,
	messages readModel.MessageQueryRepository,
	hub ws.Hub,
) *conversationWSResolver {
	return &conversationWSResolver{
		participants:               participants,
		messages:                   messages,
		hub:                        hub,
		messageChannel:             make(chan messageChannelItem, 1000),
		conversationDeletedChannel: make(chan uuid.UUID, 1000),
		conversationRenamedChannel: make(chan conversationRenamedItem, 1000),
	}
}

func (s *conversationWSResolver) Run() {
	for {
		select {
		case message := <-s.messageChannel:
			err := s.notifyAboutMessage(message.messageID, message.conversationID, message.senderUserID)

			if err != nil {
				fmt.Println(err)
			}
		case id := <-s.conversationDeletedChannel:
			err := s.notifyAboutConversationDeletion(id)

			if err != nil {
				fmt.Println(err)
			}
		case message := <-s.conversationRenamedChannel:
			err := s.notifyAboutConversationRenamed(message)

			if err != nil {
				fmt.Println(err)
			}

		}
	}

}

func (s *conversationWSResolver) DispatchMessage(messageID uuid.UUID, conversationId uuid.UUID, userId uuid.UUID) {
	s.messageChannel <- messageChannelItem{
		senderUserID:   userId,
		messageID:      messageID,
		conversationID: conversationId,
	}
}

func (s *conversationWSResolver) DispatchConversationDeleted(conversationId uuid.UUID) {
	s.conversationDeletedChannel <- conversationId
}

func (s *conversationWSResolver) DispatchConversationRenamed(conversationId uuid.UUID, newName string) {
	s.conversationRenamedChannel <- conversationRenamedItem{
		newName:        newName,
		conversationID: conversationId,
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

func (s *conversationWSResolver) notifyAboutMessage(messageID uuid.UUID, conversationId uuid.UUID, userId uuid.UUID) error {
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

func (s *conversationWSResolver) notifyAboutConversationDeletion(id uuid.UUID) error {
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

func (s *conversationWSResolver) notifyAboutConversationRenamed(item conversationRenamedItem) error {
	notification := ws.OutgoingNotification{
		Type: "conversation_renamed",
		Payload: struct {
			ConversationId uuid.UUID `json:"conversation_id"`
			NewName        string    `json:"new_name"`
		}{
			ConversationId: item.conversationID,
			NewName:        item.newName,
		},
	}

	err := s.notifyParticipants(item.conversationID, notification)

	if err != nil {
		return err
	}

	return nil
}
