package domainEventsHandlers

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/readModel"
	"context"
	"log"
)

type eventHandlers struct {
	ctx        context.Context
	subscriber infra.EventsSubscriber
	commands   *app.Commands
	queries    readModel.QueriesRepository
}

func NewEventHandlers(ctx context.Context, subscriber infra.EventsSubscriber, commands *app.Commands, queries readModel.QueriesRepository) *eventHandlers {
	return &eventHandlers{
		ctx:        ctx,
		subscriber: subscriber,
		commands:   commands,
		queries:    queries,
	}
}

func (h *eventHandlers) logHandlerError(err error) {
	log.Panicln("Error occurred event", err)
}

func (h *eventHandlers) ListenForEvents() {
	for {
		select {
		case event := <-h.subscriber.Subscribe(domain.DomainEventChannel):
			switch e := event.Data.(type) {
			case *domain.GroupConversationRenamed:
				go h.sendRenamedConversationMessage(e)
				go h.sendUpdatedConversationNotification(e.ConversationID)
			case *domain.GroupConversationDeleted:
				go h.sendGroupConversationDeletedNotification(e)
				go h.deleteConversationTopic(e)
			case *domain.GroupConversationLeft:
				go h.sendGroupConversationLeftMessage(e)
				go h.unsubscribeFromConversation(e)
				go h.sendUpdatedConversationNotification(e.ConversationID)
			case *domain.GroupConversationJoined:
				go h.sendGroupConversationJoinedMessage(e)
			case *domain.GroupConversationInvited:
				go h.sendGroupConversationInvitedMessage(e)
				go h.subscribeToConversationNotifications(e.ConversationID, e.UserID)
				go h.sendUpdatedConversationNotification(e.ConversationID)
			case *domain.GroupConversationCreated:
				go h.subscribeToConversationNotifications(e.ConversationID, e.OwnerID)
			case *domain.MessageSent:
				go h.sendMessageNotification(e)
			case *domain.DirectConversationCreated:
				go h.subscribeToConversationNotifications(e.ConversationID, e.FromUserID)
				go h.subscribeToConversationNotifications(e.ConversationID, e.ToUserID)
			}

		case <-h.ctx.Done():
			return
		}
	}
}
