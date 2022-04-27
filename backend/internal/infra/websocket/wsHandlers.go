package ws

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type WSHandlers interface {
	HandleNotification(notificationType string, message json.RawMessage, userID uuid.UUID)
}

type WSHandler func(data json.RawMessage, userID uuid.UUID)

type wsHandlers struct {
	messageHandlers map[string]WSHandler
}

func NewWSHandlers() *wsHandlers {
	return &wsHandlers{
		messageHandlers: make(map[string]WSHandler),
	}
}

func (s *wsHandlers) HandleNotification(notificationType string, data json.RawMessage, userID uuid.UUID) {
	if handler, ok := s.messageHandlers[notificationType]; ok {
		handler(data, userID)
		return
	}

	log.Printf("No handler for notification type %s", notificationType)
}

func (s *wsHandlers) SetWSHandler(notificationType string, handler WSHandler) {
	s.messageHandlers[notificationType] = handler
}
