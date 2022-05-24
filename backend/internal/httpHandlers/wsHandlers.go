package httpHandlers

import (
	"encoding/json"
	"log"

	"GitHub/go-chat/backend/internal/app"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
)

type WSHandlers interface {
	HandleNotification(notification *ws.IncomingNotification)
}

type wsHandlers struct {
	commands *app.Commands
}

func NewWSHandlers(commands *app.Commands) *wsHandlers {
	return &wsHandlers{
		commands: commands,
	}
}

func (s *wsHandlers) HandleNotification(notification *ws.IncomingNotification) {
	switch notification.Type {
	case "public_message":
		s.handleReceiveWSPublicChatMessage(notification.Data, notification.UserID)
	case "private_message":
		s.handleReceiveWSPrivateChatMessage(notification.Data, notification.UserID)
	default:
		log.Println("Unknown notification type:", notification.Type)
	}
}

func (s *wsHandlers) handleReceiveWSPublicChatMessage(data json.RawMessage, userID uuid.UUID) {
	request := struct {
		Content        string    `json:"content"`
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		log.Println(err)
		return
	}

	err := s.commands.ConversationService.SendPublicTextMessage(request.Content, request.ConversationId, userID)

	if err != nil {
		log.Println(err)
		return
	}
}

func (s *wsHandlers) handleReceiveWSPrivateChatMessage(data json.RawMessage, userID uuid.UUID) {
	request := struct {
		Content        string    `json:"content"`
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		log.Println(err)
		return
	}

	err := s.commands.ConversationService.SendPrivateTextMessage(request.Content, request.ConversationId, userID)

	if err != nil {
		log.Println(err)
		return
	}
}
