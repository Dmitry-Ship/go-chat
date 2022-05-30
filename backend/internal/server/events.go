package server

import (
	"GitHub/go-chat/backend/internal/domain"
	"log"
)

func (h *Server) logHandlerError(err error) {
	log.Panicln("Error occurred event", err)
}

func (h *Server) listenForEvents() {
	for {
		select {
		case event := <-h.subscriber.Subscribe(domain.DomainEventChannel):
			switch e := event.Data.(type) {
			case *domain.GroupConversationRenamed:
				go h.sendRenamedConversationMessage(e)
				go h.sendUpdatedConversationNotification(e.ConversationID)
			case *domain.GroupConversationDeleted:
				go h.sendGroupConversationDeletedNotification(e)
			case *domain.GroupConversationLeft:
				go h.sendGroupConversationLeftMessage(e)
				go h.sendUpdatedConversationNotification(e.ConversationID)
			case *domain.GroupConversationJoined:
				go h.sendGroupConversationJoinedMessage(e)
			case *domain.GroupConversationInvited:
				go h.sendGroupConversationInvitedMessage(e)
				go h.sendUpdatedConversationNotification(e.ConversationID)
			case *domain.GroupConversationCreated:
				continue
			case *domain.MessageSent:
				go h.sendMessageNotification(e)
			case *domain.DirectConversationCreated:
				continue
			}

		case <-h.ctx.Done():
			return
		}
	}
}
