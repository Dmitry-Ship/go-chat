package server

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"
	"fmt"
	"log"
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

func spawnWorkers[T domain.DomainEvent](numberOfWorkers int, topic string, handler func(T) error, subscriber *infra.EventBus) {
	eventChan := subscriber.Subscribe(topic)
	for i := 0; i < numberOfWorkers; i++ {
		go genericWorker(eventChan, handler)
	}
}

func (h *Server) listenForEvents() {
	spawnWorkers(1, domain.DomainEventTopic, h.sendWSNotification, h.subscriber)
	spawnWorkers(1, domain.DomainEventTopic, h.createMessage, h.subscriber)
	spawnWorkers(1, domain.DomainEventTopic, h.handleSubscriptionChanges, h.subscriber)
}

func (h *Server) createMessage(event domain.DomainEvent) error {
	switch e := event.(type) {
	case domain.GroupConversationRenamed:
		return h.conversationCommands.SendRenamedConversationMessage(e.GetConversationID(), e.UserID, e.NewName)
	case domain.GroupConversationLeft:
		return h.conversationCommands.SendLeftConversationMessage(e.GetConversationID(), e.UserID)
	case domain.GroupConversationJoined:
		return h.conversationCommands.SendJoinedConversationMessage(e.GetConversationID(), e.UserID)
	case domain.GroupConversationInvited:
		return h.conversationCommands.SendInvitedConversationMessage(e.GetConversationID(), e.UserID)
	}

	return nil
}

func (h *Server) sendWSNotification(event domain.DomainEvent) error {
	return h.notificationCommands.Notify(event)
}

func (h *Server) handleSubscriptionChanges(event domain.DomainEvent) error {
	switch e := event.(type) {
	case domain.GroupConversationJoined:
		return h.notificationCommands.SubscribeUserToChannel(e.UserID, e.GetConversationID())
	case domain.GroupConversationLeft:
		return h.notificationCommands.UnsubscribeUserFromChannel(e.UserID, e.GetConversationID())
	case domain.GroupConversationInvited:
		return h.notificationCommands.SubscribeUserToChannel(e.UserID, e.GetConversationID())
	case domain.DirectConversationCreated:
		for _, userID := range e.UserIDs {
			if err := h.notificationCommands.SubscribeUserToChannel(userID, e.GetConversationID()); err != nil {
				log.Printf("Error subscribing user %s to direct conversation: %v", userID, err)
			}
		}
	}

	return nil
}
