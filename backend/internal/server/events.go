package server

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"
	ws "GitHub/go-chat/backend/internal/websocket"
	"context"
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
	spawnWorkers(1, domain.DomainEventTopic, h.sendWSNotification, h.subscriber)
	spawnWorkers(1, domain.DomainEventTopic, h.createMessage, h.subscriber)
}

func (h *Server) createMessage(event domain.DomainEvent) error {
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
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	receiversChan, err := h.notificationPipelineService.GetReceivers(ctx, event)

	if err != nil {
		return err
	}

	messageChans, buildErrorChans := infra.FanOut(100, func() (chan ws.OutgoingNotification, chan error) {
		return h.notificationPipelineService.BuildMessage(ctx, receiversChan, event)
	})

	sendError := h.notificationPipelineService.BroadcastMessage(ctx, infra.MergeChannels(ctx, messageChans...))

	for err := range infra.MergeChannels(ctx, buildErrorChans...) {
		log.Println("Error occurred while building message: ", err)
	}

	for err := range sendError {
		log.Println("Error occurred while sending message: ", err)
	}

	return nil
}
