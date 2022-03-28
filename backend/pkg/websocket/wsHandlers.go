package ws

import (
	"encoding/json"
	"log"
)

type WSHandlers interface {
	HandleNotification(notification IncomingNotification)
}

type WSHandler func(message IncomingNotification, data json.RawMessage)

type wsHandlers struct {
	messageHandlers map[string]WSHandler
}

func NewWSHandlers() *wsHandlers {
	return &wsHandlers{
		messageHandlers: make(map[string]WSHandler),
	}
}

func (s *wsHandlers) HandleNotification(message IncomingNotification) {
	var data json.RawMessage

	notification := struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		Data: &data,
	}

	if err := json.Unmarshal(message.Data, &notification); err != nil {
		log.Println(err)
		return
	}

	if handler, ok := s.messageHandlers[notification.Type]; ok {
		handler(message, data)
	}
}

func (s *wsHandlers) SetWSHandler(notificationType string, handler WSHandler) {
	s.messageHandlers[notificationType] = handler
}
