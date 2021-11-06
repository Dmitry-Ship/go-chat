package interfaces

import (
	"GitHub/go-chat/backend/pkg/application"
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type wsMessageHandler struct {
	roomCommandService application.RoomCommandService
	MessageChannel     chan json.RawMessage
}

type incomingNotification struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewWSMessageHandler(
	roomCommandService application.RoomCommandService,
) *wsMessageHandler {
	return &wsMessageHandler{
		roomCommandService: roomCommandService,
		MessageChannel:     make(chan json.RawMessage, 100),
	}
}

func (h *wsMessageHandler) Run() {
	for message := range h.MessageChannel {
		var data json.RawMessage

		notification := incomingNotification{
			Data: &data,
		}

		if err := json.Unmarshal(message, &notification); err != nil {
			log.Println(err)
			continue
		}

		switch notification.Type {
		case "message":
			request := struct {
				Content string    `json:"content"`
				RoomId  uuid.UUID `json:"room_id"`
				UserId  uuid.UUID `json:"user_id"`
			}{}

			if err := json.Unmarshal([]byte(data), &request); err != nil {
				log.Println(err)
				continue
			}

			err := h.roomCommandService.SendMessage(request.Content, "user", request.RoomId, request.UserId)

			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}
