package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"GitHub/go-chat/backend/internal/domain"
	pubsub "GitHub/go-chat/backend/internal/infra/redis"
	"GitHub/go-chat/backend/internal/readModel"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type BroadcastMessage struct {
	Payload ws.OutgoingNotification `json:"notification"`
	UserID  uuid.UUID               `json:"user_id"`
}

type NotificationService interface {
	Send(ctx context.Context, message ws.OutgoingNotification) error
	RegisterClient(ctx context.Context, conn *websocket.Conn, userID uuid.UUID, handleNotification func(userID uuid.UUID, message []byte)) uuid.UUID
	Run()
	SubscribeUserToChannel(ctx context.Context, userID uuid.UUID, channelID uuid.UUID) error
	UnsubscribeUserFromChannel(ctx context.Context, userID uuid.UUID, channelID uuid.UUID) error
	NotifyConversationUpdated(ctx context.Context, conversation readModel.ConversationFullDTO) error
	NotifyConversationDeleted(ctx context.Context, conversationID uuid.UUID) error
	NotifyMessageSent(ctx context.Context, conversationID uuid.UUID, message readModel.MessageDTO) error
}

type notificationService struct {
	ctx              context.Context
	activeClients    ws.ActiveClients
	redisClient      *redis.Client
	participants     domain.ParticipantRepository
	subscriptionSync ws.SubscriptionSync
}

func NewNotificationService(
	ctx context.Context,
	redisClient *redis.Client,
	participants domain.ParticipantRepository,
	subscriptionSync ws.SubscriptionSync,
) *notificationService {
	return &notificationService{
		ctx:              ctx,
		activeClients:    ws.NewActiveClients(),
		redisClient:      redisClient,
		participants:     participants,
		subscriptionSync: subscriptionSync,
	}
}

func NewNotificationServiceWithClients(
	ctx context.Context,
	redisClient *redis.Client,
	participants domain.ParticipantRepository,
	subscriptionSync ws.SubscriptionSync,
	activeClients ws.ActiveClients,
) *notificationService {
	return &notificationService{
		ctx:              ctx,
		activeClients:    activeClients,
		redisClient:      redisClient,
		participants:     participants,
		subscriptionSync: subscriptionSync,
	}
}

func (s *notificationService) Send(ctx context.Context, message ws.OutgoingNotification) error {
	bMessage := BroadcastMessage{
		Payload: message,
		UserID:  message.UserID,
	}

	json, err := json.Marshal(bMessage)

	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	if err := s.redisClient.Publish(s.ctx, pubsub.ChatChannel, []byte(json)).Err(); err != nil {
		return fmt.Errorf("redis publish error: %w", err)
	}

	return nil
}

func (s *notificationService) RegisterClient(ctx context.Context, conn *websocket.Conn, userID uuid.UUID, handleNotification func(userID uuid.UUID, message []byte)) uuid.UUID {
	newClient := ws.NewClient(conn, s.activeClients.RemoveClient, handleNotification, userID)

	clientID := s.activeClients.AddClient(newClient)

	conversationIDs, err := s.participants.GetConversationIDsByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error getting user conversations: %v", err)
	} else {
		for _, conversationID := range conversationIDs {
			s.activeClients.SubscribeChannel(newClient, conversationID)
			if err := s.subscriptionSync.PublishSubscribe(clientID, conversationID, userID); err != nil {
				log.Printf("Error publishing subscription: %v", err)
			}
		}
	}

	go newClient.WritePump()
	newClient.ReadPump()

	return clientID
}

func (s *notificationService) Run() {
	redisPubsub := s.redisClient.Subscribe(s.ctx, pubsub.ChatChannel)
	chatChannel := redisPubsub.Channel()
	defer func() {
		_ = redisPubsub.Close()
	}()

	go s.subscriptionSync.Run(s.ctx, func(event ws.SubscriptionEvent) {
		client := s.activeClients.GetClient(event.ClientID)
		if client == nil {
			log.Printf("Client not found for ID: %s", event.ClientID)
			return
		}

		switch event.Action {
		case "subscribe":
			s.activeClients.SubscribeChannel(client, event.ChannelID)
		case "unsubscribe":
			s.activeClients.UnsubscribeChannel(client, event.ChannelID)
		}
	})

	for {
		select {
		case message := <-chatChannel:
			if message.Payload == "ping" {
				s.redisClient.Publish(s.ctx, pubsub.ChatChannel, "pong")
				continue
			}

			var bMessage BroadcastMessage

			if err := json.Unmarshal([]byte(message.Payload), &bMessage); err != nil {
				log.Println(err)
				continue
			}

			s.activeClients.SendToUserClients(bMessage.UserID, bMessage.Payload)

		case <-s.ctx.Done():
			return
		}
	}
}

func (s *notificationService) SubscribeUserToChannel(ctx context.Context, userID uuid.UUID, channelID uuid.UUID) error {
	s.activeClients.SubscribeUserToChannel(userID, channelID)

	clients := s.activeClients.GetClientsByUser(userID)
	for _, client := range clients {
		if err := s.subscriptionSync.PublishSubscribe(client.Id, channelID, userID); err != nil {
			log.Printf("Error publishing subscription: %v", err)
		}
	}

	return nil
}

func (s *notificationService) UnsubscribeUserFromChannel(ctx context.Context, userID uuid.UUID, channelID uuid.UUID) error {
	s.activeClients.UnsubscribeUserFromChannel(userID, channelID)

	clients := s.activeClients.GetClientsByUser(userID)
	for _, client := range clients {
		if err := s.subscriptionSync.PublishUnsubscribe(client.Id, channelID, userID); err != nil {
			log.Printf("Error publishing unsubscription: %v", err)
		}
	}

	return nil
}

func (s *notificationService) NotifyConversationUpdated(ctx context.Context, conversation readModel.ConversationFullDTO) error {
	clients := s.activeClients.GetClientsByChannel(conversation.ID)

	for _, client := range clients {
		message := ws.OutgoingNotification{
			Type:    "conversation_updated",
			UserID:  client.UserID,
			Payload: conversation,
		}
		s.sendDirectly(client, message)
	}

	return nil
}

func (s *notificationService) NotifyConversationDeleted(ctx context.Context, conversationID uuid.UUID) error {
	clients := s.activeClients.GetClientsByChannel(conversationID)

	for _, client := range clients {
		message := ws.OutgoingNotification{
			Type:   "conversation_deleted",
			UserID: client.UserID,
			Payload: struct {
				ConversationId uuid.UUID `json:"conversation_id"`
			}{
				ConversationId: conversationID,
			},
		}
		s.sendDirectly(client, message)
	}

	return nil
}

func (s *notificationService) NotifyMessageSent(ctx context.Context, conversationID uuid.UUID, message readModel.MessageDTO) error {
	clients := s.activeClients.GetClientsByChannel(conversationID)

	for _, client := range clients {
		notification := ws.OutgoingNotification{
			Type:    "message",
			UserID:  client.UserID,
			Payload: message,
		}
		s.sendDirectly(client, notification)
	}

	return nil
}

func (s *notificationService) sendDirectly(client *ws.Client, message ws.OutgoingNotification) {
	s.activeClients.SendToUserClients(message.UserID, message)
}
