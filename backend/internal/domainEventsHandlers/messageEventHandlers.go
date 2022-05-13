package domainEventsHandlers

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"
	"context"
	"log"
)

type messageEventHandlers struct {
	ctx        context.Context
	subscriber infra.EventsSubscriber
	commands   *app.Commands
}

func NewMessageEventHandlers(ctx context.Context, subscriber infra.EventsSubscriber, commands *app.Commands) *messageEventHandlers {
	return &messageEventHandlers{
		ctx:        ctx,
		subscriber: subscriber,
		commands:   commands,
	}
}

func (h *messageEventHandlers) logHandlerError(err error) {
	log.Panicln("Error occurred event", err)
}

func (h *messageEventHandlers) ListenForEvents() {
	for {
		select {
		case event := <-h.subscriber.Subscribe(domain.DomainEventChannel):
			switch e := event.Data.(type) {
			case *domain.PublicConversationRenamed:
				go h.sendRenamedConversationMessage(e)
			case *domain.PublicConversationLeft:
				go h.sendPublicConversationLeftMessage(e)
			case *domain.PublicConversationJoined:
				go h.sendPublicConversationJoinedMessage(e)
			case *domain.PublicConversationInvited:
				go h.sendPublicConversationInvitedMessage(e)
			}

		case <-h.ctx.Done():
			return
		}
	}
}

func (h *messageEventHandlers) sendRenamedConversationMessage(e *domain.PublicConversationRenamed) {
	err := h.commands.MessagingService.SendRenamedConversationMessage(e.ConversationID, e.UserID, e.NewName)

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *messageEventHandlers) sendPublicConversationLeftMessage(e *domain.PublicConversationLeft) {
	err := h.commands.MessagingService.SendLeftConversationMessage(e.ConversationID, e.UserID)

	if err != nil {
		h.logHandlerError(err)
		return
	}
}

func (h *messageEventHandlers) sendPublicConversationJoinedMessage(e *domain.PublicConversationJoined) {
	err := h.commands.MessagingService.SendJoinedConversationMessage(e.ConversationID, e.UserID)

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *messageEventHandlers) sendPublicConversationInvitedMessage(e *domain.PublicConversationInvited) {
	err := h.commands.MessagingService.SendInvitedConversationMessage(e.ConversationID, e.UserID)

	if err != nil {
		h.logHandlerError(err)
	}
}
