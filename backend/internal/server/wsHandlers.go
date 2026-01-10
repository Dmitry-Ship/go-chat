package server

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

func (s *Server) handleNotification(userID uuid.UUID, message []byte) {
	var data json.RawMessage

	notification := struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		Data: &data,
	}

	if err := json.Unmarshal(message, &notification); err != nil {
		log.Println(err)
		return
	}

	switch notification.Type {
	case "group_message":
		s.handleReceiveWSGroupChatMessage(data, userID)
	case "direct_message":
		s.handleReceiveWSDirectChatMessage(data, userID)
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

	err := s.conversationCommands.SendTextMessage(context.Background(), request.ConversationId, userID, request.Content)

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

	err := s.conversationCommands.SendTextMessage(context.Background(), request.ConversationId, userID, request.Content)

	if err != nil {
		log.Println(err)
		return
	}
}
