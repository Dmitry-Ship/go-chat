package services

import (
	"context"
	"log"

	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type NotificationService interface {
	Broadcast(ctx context.Context, conversationID uuid.UUID, notification ws.OutgoingNotification) error
	RegisterClient(ctx context.Context, conn *websocket.Conn, userID uuid.UUID) uuid.UUID
	Run()
	InvalidateMembership(ctx context.Context, userID uuid.UUID) error
	Shutdown()
}

type broadcastMessage struct {
	notification   ws.OutgoingNotification
	conversationID uuid.UUID
}

type BroadcastMessage struct {
	Payload        ws.OutgoingNotification `json:"notification"`
	UserID         uuid.UUID               `json:"user_id"`
	MessageID      string                  `json:"message_id"`
	ServerID       string                  `json:"server_id"`
	ConversationID uuid.UUID               `json:"conversation_id"`
}

type SubscriptionEvent struct {
	Action string    `json:"action"`
	UserID uuid.UUID `json:"user_id"`
}

type notificationService struct {
	ctx           context.Context
	cancel        context.CancelFunc
	serverID      string
	activeClients ws.ActiveClients

	broadcast chan broadcastMessage

	registerClient chan *ws.Client
	removeClient   chan *ws.Client
}

func NewNotificationService(
	ctx context.Context,
	serverID string,
	activeClients ws.ActiveClients,
) NotificationService {
	nmCtx, cancel := context.WithCancel(ctx)

	nm := &notificationService{
		ctx:            nmCtx,
		cancel:         cancel,
		serverID:       serverID,
		activeClients:  activeClients,
		broadcast:      make(chan broadcastMessage, 1000),
		registerClient: make(chan *ws.Client, 100),
		removeClient:   make(chan *ws.Client, 100),
	}

	return nm
}

func (ns *notificationService) Broadcast(ctx context.Context, conversationID uuid.UUID, notification ws.OutgoingNotification) error {
	ns.broadcast <- broadcastMessage{
		conversationID: conversationID,
		notification:   notification,
	}
	return nil
}

func (ns *notificationService) RegisterClient(ctx context.Context, conn *websocket.Conn, userID uuid.UUID) uuid.UUID {
	client := ws.NewClient(conn, ns.removeClient, userID)
	ns.registerClient <- client
	return client.Id
}

func (ns *notificationService) InvalidateMembership(ctx context.Context, userID uuid.UUID) error {
	return ns.activeClients.InvalidateMembership(ctx, userID)
}

func (ns *notificationService) Run() {
	for {
		select {
		case msg := <-ns.broadcast:
			clients := ns.activeClients.GetClientsByChannel(msg.conversationID)
			for _, client := range clients {
				if err := client.SendNotification(msg.notification); err != nil {
					log.Printf("Error sending notification to client %s: %v", client.Id, err)
					ns.removeClient <- client
				}
			}

		case client := <-ns.registerClient:
			ns.activeClients.AddClient(client)

			go client.WritePump()
			client.ReadPump()

		case client := <-ns.removeClient:
			ns.activeClients.RemoveClient(client)

		case <-ns.ctx.Done():
			return
		}
	}
}

func (ns *notificationService) Shutdown() {
	ns.cancel()
}
