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
	Send(message ws.OutgoingNotification) error
	RegisterClient(conn *websocket.Conn, userID uuid.UUID, handleNotification func(userID uuid.UUID, message []byte)) uuid.UUID
	Run()
	SubscribeUserToChannel(userID uuid.UUID, channelID uuid.UUID) error
	UnsubscribeUserFromChannel(userID uuid.UUID, channelID uuid.UUID) error
	Notify(event domain.DomainEvent) error
}

type notificationService struct {
	ctx              context.Context
	activeClients    ws.ActiveClients
	redisClient      *redis.Client
	participants     domain.ParticipantRepository
	subscriptionSync ws.SubscriptionSync
	queries          readModel.QueriesRepository
}

func NewNotificationService(
	ctx context.Context,
	redisClient *redis.Client,
	participants domain.ParticipantRepository,
	subscriptionSync ws.SubscriptionSync,
	queries readModel.QueriesRepository,
) *notificationService {
	return &notificationService{
		ctx:              ctx,
		activeClients:    ws.NewActiveClients(),
		redisClient:      redisClient,
		participants:     participants,
		subscriptionSync: subscriptionSync,
		queries:          queries,
	}
}

func NewNotificationServiceWithClients(
	ctx context.Context,
	redisClient *redis.Client,
	participants domain.ParticipantRepository,
	subscriptionSync ws.SubscriptionSync,
	queries readModel.QueriesRepository,
	activeClients ws.ActiveClients,
) *notificationService {
	return &notificationService{
		ctx:              ctx,
		activeClients:    activeClients,
		redisClient:      redisClient,
		participants:     participants,
		subscriptionSync: subscriptionSync,
		queries:          queries,
	}
}

func (s *notificationService) Send(message ws.OutgoingNotification) error {
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

func (s *notificationService) RegisterClient(conn *websocket.Conn, userID uuid.UUID, handleNotification func(userID uuid.UUID, message []byte)) uuid.UUID {
	newClient := ws.NewClient(conn, s.activeClients.RemoveClient, handleNotification, userID)

	clientID := s.activeClients.AddClient(newClient)

	conversationIDs, err := s.participants.GetConversationIDsByUserID(userID)
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

		if event.Action == "subscribe" {
			s.activeClients.SubscribeChannel(client, event.ChannelID)
		} else if event.Action == "unsubscribe" {
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

func (s *notificationService) SubscribeUserToChannel(userID uuid.UUID, channelID uuid.UUID) error {
	s.activeClients.SubscribeUserToChannel(userID, channelID)

	clients := s.activeClients.GetClientsByUser(userID)
	for _, client := range clients {
		if err := s.subscriptionSync.PublishSubscribe(client.Id, channelID, userID); err != nil {
			log.Printf("Error publishing subscription: %v", err)
		}
	}

	return nil
}

func (s *notificationService) UnsubscribeUserFromChannel(userID uuid.UUID, channelID uuid.UUID) error {
	s.activeClients.UnsubscribeUserFromChannel(userID, channelID)

	clients := s.activeClients.GetClientsByUser(userID)
	for _, client := range clients {
		if err := s.subscriptionSync.PublishUnsubscribe(client.Id, channelID, userID); err != nil {
			log.Printf("Error publishing unsubscription: %v", err)
		}
	}

	return nil
}

func (s *notificationService) Notify(event domain.DomainEvent) error {
	recipients := s.getRecipients(event)

	for _, recipient := range recipients {
		message, err := s.buildMessage(recipient.UserID, event)
		if err != nil {
			log.Printf("Error building message for user %s: %v", recipient.UserID, err)
			continue
		}
		s.sendDirectly(recipient, message)
	}

	return nil
}

func (s *notificationService) getRecipients(event domain.DomainEvent) []*ws.Client {
	var clients []*ws.Client

	switch e := event.(type) {
	case
		domain.GroupConversationRenamed,
		domain.GroupConversationLeft,
		domain.GroupConversationJoined,
		domain.GroupConversationInvited,
		domain.MessageSent,
		domain.GroupConversationDeleted,
		domain.DirectConversationCreated:
		if e, ok := e.(domain.ConversationEvent); ok {
			clients = s.activeClients.GetClientsByChannel(e.GetConversationID())
		}
	}

	return clients
}

func (s *notificationService) buildMessage(userID uuid.UUID, event domain.DomainEvent) (ws.OutgoingNotification, error) {
	var buildMessage func(userID uuid.UUID) (ws.OutgoingNotification, error)

	switch e := event.(type) {
	case domain.GroupConversationRenamed, domain.GroupConversationLeft, domain.GroupConversationJoined, domain.GroupConversationInvited:
		if e, ok := e.(domain.ConversationEvent); ok {
			buildMessage = s.getConversationUpdatedBuilder(e.GetConversationID())
		}
	case domain.GroupConversationDeleted:
		buildMessage = s.getConversationDeletedBuilder(e.GetConversationID())
	case domain.MessageSent:
		buildMessage = s.getMessageSentBuilder(e.MessageID)
	}

	return buildMessage(userID)
}

func (s *notificationService) getConversationUpdatedBuilder(conversationID uuid.UUID) func(userID uuid.UUID) (ws.OutgoingNotification, error) {
	return func(userID uuid.UUID) (ws.OutgoingNotification, error) {
		conversation, err := s.queries.GetConversation(conversationID, userID)

		if err != nil {
			return ws.OutgoingNotification{}, fmt.Errorf("get conversation error: %w", err)
		}

		return ws.OutgoingNotification{
			Type:    "conversation_updated",
			UserID:  userID,
			Payload: conversation,
		}, nil
	}
}

func (s *notificationService) getConversationDeletedBuilder(conversationID uuid.UUID) func(userID uuid.UUID) (ws.OutgoingNotification, error) {
	return func(userID uuid.UUID) (ws.OutgoingNotification, error) {
		return ws.OutgoingNotification{
			Type:   "conversation_deleted",
			UserID: userID,
			Payload: struct {
				ConversationId uuid.UUID `json:"conversation_id"`
			}{
				ConversationId: conversationID,
			},
		}, nil
	}
}

func (s *notificationService) getMessageSentBuilder(messageID uuid.UUID) func(userID uuid.UUID) (ws.OutgoingNotification, error) {
	return func(userID uuid.UUID) (ws.OutgoingNotification, error) {
		messageDTO, err := s.queries.GetNotificationMessage(messageID, userID)

		if err != nil {
			return ws.OutgoingNotification{}, err
		}

		return ws.OutgoingNotification{
			Type:    "message",
			UserID:  userID,
			Payload: messageDTO,
		}, nil
	}
}

func (s *notificationService) sendDirectly(client *ws.Client, message ws.OutgoingNotification) {
	s.activeClients.SendToUserClients(message.UserID, message)
}
