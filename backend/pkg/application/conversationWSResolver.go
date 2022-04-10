package application

import (
	"GitHub/go-chat/backend/domain"
	ws "GitHub/go-chat/backend/pkg/websocket"
	"fmt"

	"github.com/google/uuid"
)

type ConversationWSResolver interface {
	DispatchUserMessage(messageID uuid.UUID, conversationId uuid.UUID, userId uuid.UUID)
	DispatchConversationDeleted(conversationId uuid.UUID)
	DispatchSystemMessage(messageID uuid.UUID, conversationId uuid.UUID)
}

type userMessageChannelItem struct {
	senderUserID   uuid.UUID
	messageID      uuid.UUID
	conversationID uuid.UUID
}

type systemMessageChannelItem struct {
	messageID      uuid.UUID
	conversationID uuid.UUID
}

type conversationWSResolver struct {
	participants               domain.ParticipantQueryRepository
	messages                   domain.MessageQueryRepository
	hub                        ws.Hub
	userMessageChannel         chan userMessageChannelItem
	systemMessageChannel       chan systemMessageChannelItem
	conversationDeletedChannel chan uuid.UUID
}

func NewConversationWSResolver(
	participants domain.ParticipantQueryRepository,
	messages domain.MessageQueryRepository,
	hub ws.Hub,
) *conversationWSResolver {
	return &conversationWSResolver{
		participants:               participants,
		messages:                   messages,
		hub:                        hub,
		userMessageChannel:         make(chan userMessageChannelItem, 1000),
		systemMessageChannel:       make(chan systemMessageChannelItem, 1000),
		conversationDeletedChannel: make(chan uuid.UUID, 1000),
	}
}

func (s *conversationWSResolver) Run() {
	for {
		select {
		case message := <-s.userMessageChannel:
			err := s.notifyAboutUserMessage(message.messageID, message.conversationID, message.senderUserID)

			if err != nil {
				fmt.Println(err)
			}
		case message := <-s.systemMessageChannel:
			err := s.notifyAboutSystemMessage(message.messageID, message.conversationID)

			if err != nil {
				fmt.Println(err)
			}
		case id := <-s.conversationDeletedChannel:
			err := s.notifyAboutConversationDeletion(id)

			if err != nil {
				fmt.Println(err)
			}
		}
	}

}

func (s *conversationWSResolver) DispatchUserMessage(messageID uuid.UUID, conversationId uuid.UUID, userId uuid.UUID) {
	s.userMessageChannel <- userMessageChannelItem{
		senderUserID:   userId,
		messageID:      messageID,
		conversationID: conversationId,
	}
}

func (s *conversationWSResolver) DispatchConversationDeleted(conversationId uuid.UUID) {
	s.conversationDeletedChannel <- conversationId
}

func (s *conversationWSResolver) DispatchSystemMessage(messageID uuid.UUID, conversationId uuid.UUID) {
	s.systemMessageChannel <- systemMessageChannelItem{
		messageID:      messageID,
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

func (s *conversationWSResolver) notifyAboutUserMessage(messageID uuid.UUID, conversationId uuid.UUID, userId uuid.UUID) error {
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

func (s *conversationWSResolver) notifyAboutSystemMessage(messageID uuid.UUID, conversationId uuid.UUID) error {
	messageDTO, err := s.messages.GetMessageByID(messageID, uuid.Nil)

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
