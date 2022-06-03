package server

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"
	ws "GitHub/go-chat/backend/internal/websocket"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func genericWorker[T domain.DomainEvent](eventChan <-chan infra.Event, handler func(T) error) {
	for event := range eventChan {
		e, ok := event.Data.(T)

		if !ok {
			fmt.Println("invalid event type: ", event.Topic)
			continue
		}

		err := handler(e)

		if err != nil {
			log.Println("Error occurred while handling event: ", e.GetName(), err)
		}
	}
}

func spawnWorkers[T domain.DomainEvent](numberOfWorkers int, topic string, handler func(T) error, subscriber infra.EventsSubscriber) {
	eventChan := subscriber.Subscribe(topic)
	for i := 0; i < numberOfWorkers; i++ {
		go genericWorker(eventChan, handler)
	}
}

func (h *Server) listenForEvents() {
	spawnWorkers(5, domain.DomainEventTopic, h.sendWSNotification, h.subscriber)
	spawnWorkers(2, domain.DomainEventTopic, h.sendMessage, h.subscriber)
}

func (h *Server) sendMessage(event domain.DomainEvent) error {
	switch e := event.(type) {
	case *domain.GroupConversationRenamed:
		return h.conversationCommands.SendRenamedConversationMessage(e.GetConversationID(), e.UserID, e.NewName)
	case *domain.GroupConversationLeft:
		return h.conversationCommands.SendLeftConversationMessage(e.GetConversationID(), e.UserID)
	case *domain.GroupConversationJoined:
		return h.conversationCommands.SendJoinedConversationMessage(e.GetConversationID(), e.UserID)
	case *domain.GroupConversationInvited:
		return h.conversationCommands.SendInvitedConversationMessage(e.GetConversationID(), e.UserID)
	}

	return nil
}

func (h *Server) sendWSNotification(event domain.DomainEvent) error {
	var receiversIDs []uuid.UUID
	var err error
	var buildMessage func(userID uuid.UUID) (*ws.OutgoingNotification, error)

	// find notification receivers
	switch e := event.(type) {
	case
		*domain.GroupConversationRenamed,
		*domain.GroupConversationLeft,
		*domain.GroupConversationJoined,
		*domain.GroupConversationInvited,
		*domain.MessageSent,
		*domain.GroupConversationDeleted:
		if e, ok := e.(domain.ConversationEvent); ok {
			receiversIDs, err = h.notificationResolver.GetConversationRecipients(e.GetConversationID())

			if err != nil {
				return err
			}
		}
	}

	// build notification
	switch e := event.(type) {
	case *domain.GroupConversationRenamed, *domain.GroupConversationLeft, *domain.GroupConversationJoined, *domain.GroupConversationInvited:
		if e, ok := e.(domain.ConversationEvent); ok {
			buildMessage = h.notificationBuilder.GetConversationUpdatedBuilder(e.GetConversationID())
		}
	case *domain.GroupConversationDeleted:
		buildMessage = h.notificationBuilder.GetConversationDeletedBuilder(e.GetConversationID())
	case *domain.MessageSent:
		buildMessage = h.notificationBuilder.GetMessageSentBuilder(e.MessageID)
	}

	return h.notificationCommands.Broadcast(receiversIDs, buildMessage)
}
