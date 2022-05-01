package httpHandlers

import (
	"GitHub/go-chat/backend/internal/app"
	ws "GitHub/go-chat/backend/internal/infra/websocket"
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type wsHandlers struct {
	incomingNotificationsChan chan *ws.IncomingNotification
	commands                  *app.Commands
}

func NewWSHandlers(commands *app.Commands, incomingNotificationsChan chan *ws.IncomingNotification) *wsHandlers {
	return &wsHandlers{
		commands:                  commands,
		incomingNotificationsChan: incomingNotificationsChan,
	}
}

func (s *wsHandlers) Run() {
	for notification := range s.incomingNotificationsChan {
		switch notification.Type {
		case "message":
			s.handleReceiveWSChatMessage(notification.Data, notification.UserID)
		default:
			log.Println("Unknown notification type:", notification.Type)
		}
	}
}

func (s *wsHandlers) handleReceiveWSChatMessage(data json.RawMessage, userID uuid.UUID) {
	request := struct {
		Content        string    `json:"content"`
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		log.Println(err)
		return
	}

	err := s.commands.MessagingService.SendTextMessage(request.Content, request.ConversationId, userID)

	if err != nil {
		log.Println(err)
		return
	}
}
