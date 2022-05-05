package services

import (
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/gorilla/websocket"

	"github.com/google/uuid"
)

type ClientRegister interface {
	RegisterClient(conn *websocket.Conn, userID uuid.UUID, handleNotification func(notification *ws.IncomingNotification))
}

type clientRegister struct {
	activeClients ws.ActiveClients
}

func NewClientRegister(activeClients ws.ActiveClients) *clientRegister {
	return &clientRegister{
		activeClients: activeClients,
	}
}

func (s *clientRegister) RegisterClient(conn *websocket.Conn, userID uuid.UUID, handleNotification func(notification *ws.IncomingNotification)) {
	newClient := ws.NewClient(conn, s.activeClients.RemoveClient, handleNotification, userID)
	newClient.Listen()

	s.activeClients.AddClient(newClient)
}
