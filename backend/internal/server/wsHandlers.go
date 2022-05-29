package server

import (
	"encoding/json"
	"log"

	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
)

func (s *Server) handleNotification(notification *ws.IncomingNotification) {
	switch notification.Type {
	case "group_message":
		s.handleReceiveWSGroupChatMessage(notification.Data, notification.UserID)
	case "direct_message":
		s.handleReceiveWSDirectChatMessage(notification.Data, notification.UserID)
	default:
		log.Println("Unknown notification type:", notification.Type)
	}
}

func (s *Server) handleReceiveWSGroupChatMessage(data json.RawMessage, userID uuid.UUID) {
	request := struct {
		Content        string    `json:"content"`
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		log.Println(err)
		return
	}

	err := s.conversationCommands.SendGroupTextMessage(request.ConversationId, userID, request.Content)

	if err != nil {
		log.Println(err)
		return
	}
}

func (s *Server) handleReceiveWSDirectChatMessage(data json.RawMessage, userID uuid.UUID) {
	request := struct {
		Content        string    `json:"content"`
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		log.Println(err)
		return
	}

	err := s.conversationCommands.SendDirectTextMessage(request.ConversationId, userID, request.Content)

	if err != nil {
		log.Println(err)
		return
	}
}
