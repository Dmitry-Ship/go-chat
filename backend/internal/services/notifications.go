package services

import (
	"context"
	"log"
	"sync"
	"time"

	"GitHub/go-chat/backend/internal/domain"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type BatchBroadcastMessage struct {
	Payload ws.BatchOutgoingNotification `json:"notification"`
	UserID  uuid.UUID                    `json:"user_id"`
}

type SubscriptionEvent struct {
	Action string    `json:"action"`
	UserID uuid.UUID `json:"user_id"`
}

type BroadcastMessage struct {
	Payload   ws.OutgoingNotification `json:"notification"`
	UserID    uuid.UUID               `json:"user_id"`
	MessageID string                  `json:"message_id"`
	ServerID  string                  `json:"server_id"`
}

type batchKey struct {
	channelID uuid.UUID
	userID    uuid.UUID
}

type batchNotification struct {
	channelID    uuid.UUID
	notification ws.OutgoingNotification
}

type notificationService struct {
	ctx           context.Context
	activeClients ws.ActiveClients
	batcher       NotificationBatcher
	broadcaster   RedisBroadcaster
	deduplicator  MessageDeduplicator
	serverID      string
	wg            sync.WaitGroup
	cancel        context.CancelFunc
}

type NotificationService interface {
	Broadcast(ctx context.Context, channelID uuid.UUID, notification ws.OutgoingNotification) error
	RegisterClient(ctx context.Context, conn *websocket.Conn, userID uuid.UUID, handleNotification func(userID uuid.UUID, message []byte)) uuid.UUID
	Run()
	InvalidateMembership(ctx context.Context, userID uuid.UUID) error
	Shutdown()
}

type NotificationServiceOption func(*notificationService)

func WithActiveClients(ac ws.ActiveClients) NotificationServiceOption {
	return func(ns *notificationService) {
		ns.activeClients = ac
	}
}

func WithBatchTimeout(d time.Duration) NotificationServiceOption {
	return func(ns *notificationService) {
		ns.batcher = NewNotificationBatcher(10, d).(*notificationBatcher)
	}
}

func WithMaxBatchSize(n int) NotificationServiceOption {
	return func(ns *notificationService) {
		ns.batcher = NewNotificationBatcher(n, 10*time.Millisecond).(*notificationBatcher)
	}
}

func WithServerID(serverID string) NotificationServiceOption {
	return func(ns *notificationService) {
		ns.serverID = serverID
	}
}

func WithDeduplicationCapacity(capacity int) NotificationServiceOption {
	return func(ns *notificationService) {
		ns.deduplicator = NewMessageDeduplicator(capacity, 5*time.Minute).(*messageDeduplicator)
	}
}

func NewNotificationService(
	ctx context.Context,
	redisClient interface{},
	participants domain.ParticipantRepository,
	serverID string,
	opts ...NotificationServiceOption,
) NotificationService {
	broadcastRedisClient, _ := redisClient.(*redis.Client)

	nsCtx, cancel := context.WithCancel(ctx)

	ns := &notificationService{
		ctx:          nsCtx,
		serverID:     serverID,
		broadcaster:  NewRedisBroadcaster(broadcastRedisClient),
		batcher:      NewNotificationBatcher(10, 10*time.Millisecond).(*notificationBatcher),
		deduplicator: NewMessageDeduplicator(100000, 5*time.Minute).(*messageDeduplicator),
		cancel:       cancel,
	}

	for _, opt := range opts {
		opt(ns)
	}

	if ns.activeClients == nil {
		ns.activeClients = ws.NewActiveClients(participants)
	}

	return ns
}

func (s *notificationService) Broadcast(ctx context.Context, channelID uuid.UUID, notification ws.OutgoingNotification) error {
	return s.batcher.Add(channelID, notification.UserID, notification)
}

func (s *notificationService) RegisterClient(ctx context.Context, conn *websocket.Conn, userID uuid.UUID, handleNotification func(userID uuid.UUID, message []byte)) uuid.UUID {
	newClient := ws.NewClient(conn, s.activeClients.RemoveClient, handleNotification, userID)

	clientID := s.activeClients.AddClient(newClient)

	if err := s.activeClients.InvalidateMembership(ctx, userID); err != nil {
		log.Printf("Error invalidating membership: %v", err)
	}

	go newClient.WritePump()
	newClient.ReadPump()

	return clientID
}

func (s *notificationService) Run() {
	s.wg.Add(3)
	go s.deduplicator.StartCleanup(s.ctx)
	go s.runBatcher()
	go s.runSubscriber()
}

func (s *notificationService) runBatcher() {
	defer s.wg.Done()
	s.batcher.Start(s.ctx, s.flushBatch)
}

func (s *notificationService) runSubscriber() {
	defer s.wg.Done()

	chatMessages := s.broadcaster.SubscribeToChat(s.ctx)
	invalidationEvents := s.broadcaster.SubscribeToInvalidation(s.ctx)

	for {
		select {
		case bMessage := <-chatMessages:
			if s.deduplicator.AlreadySent(bMessage.MessageID) {
				continue
			}

			_ = s.batcher.Add(uuid.Nil, bMessage.UserID, bMessage.Payload)

		case event := <-invalidationEvents:
			if event.Action == "invalidate" {
				if err := s.activeClients.InvalidateMembership(s.ctx, event.UserID); err != nil {
					log.Printf("Error invalidating membership: %v", err)
				}
			}

		case <-s.ctx.Done():
			return
		}
	}
}

func (s *notificationService) flushBatch(key batchKey, notifications []ws.OutgoingNotification) {
	if len(notifications) == 0 {
		return
	}

	if len(notifications) == 1 {
		s.flushSingleNotification(key, notifications[0])
		return
	}

	s.flushBatchNotifications(key, notifications)
}

func (s *notificationService) flushSingleNotification(key batchKey, notification ws.OutgoingNotification) {
	if key.channelID != uuid.Nil {
		clients := s.activeClients.GetClientsByChannel(key.channelID)
		for _, client := range clients {
			client.SendNotification(notification)
		}
	}

	clients := s.activeClients.GetClientsByUser(notification.UserID)
	for _, client := range clients {
		client.SendNotification(notification)
	}

	messageID := uuid.New().String()
	s.deduplicator.MarkSent(messageID)

	if err := s.broadcaster.PublishNotification(s.ctx, notification, s.serverID); err != nil {
		log.Printf("Error publishing notification: %v", err)
	}
}

func (s *notificationService) flushBatchNotifications(key batchKey, notifications []ws.OutgoingNotification) {
	batch := ws.BatchOutgoingNotification{
		UserID: notifications[0].UserID,
		Events: make([]ws.NotificationEvent, len(notifications)),
	}

	for i, n := range notifications {
		batch.Events[i] = ws.NotificationEvent{
			Type:    n.Type,
			Payload: n.Payload,
		}
	}

	if key.channelID != uuid.Nil {
		clients := s.activeClients.GetClientsByChannel(key.channelID)
		for _, client := range clients {
			client.SendNotification(ws.OutgoingNotification{
				Type:    "batch",
				UserID:  batch.UserID,
				Payload: batch,
			})
		}
	}

	clients := s.activeClients.GetClientsByUser(batch.UserID)
	for _, client := range clients {
		client.SendNotification(ws.OutgoingNotification{
			Type:    "batch",
			UserID:  batch.UserID,
			Payload: batch,
		})
	}

	for _, n := range notifications {
		messageID := uuid.New().String()
		s.deduplicator.MarkSent(messageID)
		if err := s.broadcaster.PublishNotification(s.ctx, n, s.serverID); err != nil {
			log.Printf("Error publishing notification: %v", err)
		}
	}
}

func (s *notificationService) InvalidateMembership(ctx context.Context, userID uuid.UUID) error {
	if err := s.activeClients.InvalidateMembership(ctx, userID); err != nil {
		return err
	}

	if err := s.broadcaster.PublishInvalidate(ctx, userID); err != nil {
		log.Printf("Error publishing invalidate: %v", err)
	}

	return nil
}

func (s *notificationService) Shutdown() {
	s.cancel()
	s.wg.Wait()
	s.broadcaster.Close()
}
