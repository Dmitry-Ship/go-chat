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
			case *domain.GroupConversationRenamed:
				go h.sendRenamedConversationMessage(e)
			case *domain.GroupConversationLeft:
				go h.sendGroupConversationLeftMessage(e)
			case *domain.GroupConversationJoined:
				go h.sendGroupConversationJoinedMessage(e)
			case *domain.GroupConversationInvited:
				go h.sendGroupConversationInvitedMessage(e)
			}

		case <-h.ctx.Done():
			return
		}
	}
}

func (h *messageEventHandlers) sendRenamedConversationMessage(e *domain.GroupConversationRenamed) {
	err := h.commands.ConversationService.SendRenamedConversationMessage(e.ConversationID, e.UserID, e.NewName)

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *messageEventHandlers) sendGroupConversationLeftMessage(e *domain.GroupConversationLeft) {
	err := h.commands.ConversationService.SendLeftConversationMessage(e.ConversationID, e.UserID)

	if err != nil {
		h.logHandlerError(err)
		return
	}
}

func (h *messageEventHandlers) sendGroupConversationJoinedMessage(e *domain.GroupConversationJoined) {
	err := h.commands.ConversationService.SendJoinedConversationMessage(e.ConversationID, e.UserID)

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *messageEventHandlers) sendGroupConversationInvitedMessage(e *domain.GroupConversationInvited) {
	err := h.commands.ConversationService.SendInvitedConversationMessage(e.ConversationID, e.UserID)

	if err != nil {
		h.logHandlerError(err)
	}
}
