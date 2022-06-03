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

func spawnWorkers[T domain.DomainEvent](numberOfWorkers int, topic string, handler func(T) error, subscriber infra.EventsSubscriber) {
	eventChan := subscriber.Subscribe(topic)
	for i := 0; i < numberOfWorkers; i++ {
		go genericWorker(eventChan, handler)
	}
}

func (h *Server) listenForEvents() {
	spawnWorkers(10, domain.MessageSentEventName, h.sendWSNotification, h.subscriber)
	spawnWorkers(1, domain.GroupConversationRenamedEventName, h.sendMessage, h.subscriber)
	spawnWorkers(1, domain.GroupConversationRenamedEventName, h.sendWSNotification, h.subscriber)
	spawnWorkers(1, domain.GroupConversationDeletedEventName, h.sendWSNotification, h.subscriber)
	spawnWorkers(1, domain.GroupConversationLeftEventName, h.sendMessage, h.subscriber)
	spawnWorkers(1, domain.GroupConversationLeftEventName, h.sendWSNotification, h.subscriber)
	spawnWorkers(1, domain.GroupConversationJoinedEventName, h.sendMessage, h.subscriber)
	spawnWorkers(1, domain.GroupConversationInvitedEventName, h.sendMessage, h.subscriber)
	spawnWorkers(1, domain.GroupConversationInvitedEventName, h.sendWSNotification, h.subscriber)
}
