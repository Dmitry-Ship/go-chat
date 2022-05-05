package hub

import (
	"encoding/json"
	"log"

	"GitHub/go-chat/backend/internal/app"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/gorilla/websocket"

	"github.com/google/uuid"
)

type ClientRegister interface {
	RegisterClient(conn *websocket.Conn, userID uuid.UUID)
}

type clientRegister struct {
	commands      *app.Commands
	activeClients ws.ActiveClients
}

func NewClientRegister(commands *app.Commands, activeClients ws.ActiveClients) *clientRegister {
	return &clientRegister{
		commands:      commands,
		activeClients: activeClients,
	}
}

func (s *clientRegister) RegisterClient(conn *websocket.Conn, userID uuid.UUID) {
	newClient := ws.NewClient(conn, s.activeClients.RemoveClient, s.handleNotification, userID)
	newClient.Listen()

	s.activeClients.AddClient(newClient)
}

func (s *clientRegister) handleNotification(notification *ws.IncomingNotification) {
	switch notification.Type {
	case "message":
		s.handleReceiveWSChatMessage(notification.Data, notification.UserID)
	default:
		log.Println("Unknown notification type:", notification.Type)
	}
}

func (s *clientRegister) handleReceiveWSChatMessage(data json.RawMessage, userID uuid.UUID) {
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
