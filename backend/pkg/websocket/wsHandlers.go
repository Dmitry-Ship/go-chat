package ws

import (
	"encoding/json"
	"log"
)

type WSHandlers interface {
	SendMessage(notification IncomingNotification)
}

type WSHandler func(message IncomingNotification, data json.RawMessage)

type wsHandlers struct {
	messageChannel  chan IncomingNotification
	messageHandlers map[string]WSHandler
}

func NewWSHandlers() *wsHandlers {
	return &wsHandlers{
		messageChannel:  make(chan IncomingNotification, 100),
		messageHandlers: make(map[string]WSHandler),
	}
}

func (s *wsHandlers) Run() {
	for {
		select {
		case message := <-s.messageChannel:
			var data json.RawMessage

			notification := struct {
				Type string      `json:"type"`
				Data interface{} `json:"data"`
			}{
				Data: &data,
			}

			if err := json.Unmarshal(message.Data, &notification); err != nil {
				log.Println(err)
				continue
			}

			if handler, ok := s.messageHandlers[notification.Type]; ok {
				go handler(message, data)
			}
		}

	}
}

func (s *wsHandlers) SendMessage(notification IncomingNotification) {
	s.messageChannel <- notification
}

func (s *wsHandlers) SetWSHandler(notificationType string, handler WSHandler) {
	s.messageHandlers[notificationType] = handler
}
