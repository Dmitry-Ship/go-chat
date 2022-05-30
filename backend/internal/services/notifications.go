package services

import (
	"GitHub/go-chat/backend/internal/domain"
	pubsub "GitHub/go-chat/backend/internal/infra/redis"
	ws "GitHub/go-chat/backend/internal/websocket"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type BroadcastMessage struct {
	Payload ws.OutgoingNotification `json:"notification"`
	UserID  uuid.UUID               `json:"user_id"`
}

type NotificationService interface {
	SubscribeToTopic(topic string, userID uuid.UUID) error
	UnsubscribeFromTopic(topic string, userID uuid.UUID) error
	SendToTopic(topic string, buildMessage func(userID uuid.UUID) (*ws.OutgoingNotification, error)) error
	RegisterClient(conn *websocket.Conn, userID uuid.UUID, handleNotification func(notification *ws.IncomingNotification))
	Run()
}

type notificationService struct {
	ctx                context.Context
	notificationTopics domain.NotificationTopicRepository
	activeClients      ws.ActiveClients
	redisClient        *redis.Client
}

func NewNotificationService(
	ctx context.Context,
	notificationTopics domain.NotificationTopicRepository,
	activeClients ws.ActiveClients,
	redisClient *redis.Client,
) *notificationService {
	return &notificationService{
		ctx:                ctx,
		notificationTopics: notificationTopics,
		activeClients:      activeClients,
		redisClient:        redisClient,
	}
}

func (s *notificationService) SubscribeToTopic(topic string, userID uuid.UUID) error {
	notificationTopicID := uuid.New()

	notificationTopic, err := domain.NewNotificationTopic(notificationTopicID, topic, userID)

	if err != nil {
		return err
	}

	return s.notificationTopics.Store(notificationTopic)
}

func (s *notificationService) UnsubscribeFromTopic(topic string, userID uuid.UUID) error {
	return s.notificationTopics.DeleteByUserIDAndTopic(userID, topic)
}

func (s *notificationService) sendWorker(
	userID uuid.UUID,
	buildMessage func(userID uuid.UUID) (*ws.OutgoingNotification, error),
	wg *sync.WaitGroup,
	sem chan struct{},
	errorChan chan error,
) {
	defer func() {
		wg.Done()
		<-sem
	}()

	sem <- struct{}{}

	notification, err := buildMessage(userID)

	if err != nil {
		errorChan <- err
		return
	}

	message := BroadcastMessage{
		Payload: *notification,
		UserID:  userID,
	}

	json, err := json.Marshal(message)

	if err != nil {
		errorChan <- err
		return
	}

	err = s.redisClient.Publish(s.ctx, pubsub.ChatChannel, []byte(json)).Err()

	if err != nil {
		errorChan <- err
	}
}

func (s *notificationService) SendToTopic(topic string, buildMessage func(userID uuid.UUID) (*ws.OutgoingNotification, error)) error {
	ids, err := s.notificationTopics.GetUserIDsByTopic(topic)

	if err != nil {
		return err
	}

	sem := make(chan struct{}, 100)
	errorChan := make(chan error, len(ids))
	var wg sync.WaitGroup
	wg.Add(len(ids))
	for _, id := range ids {
		go s.sendWorker(id, buildMessage, &wg, sem, errorChan)
	}
	wg.Wait()
	close(errorChan)

	for er := range errorChan {
		fmt.Println(er)
	}

	return nil
}

func (s *notificationService) RegisterClient(conn *websocket.Conn, userID uuid.UUID, handleNotification func(notification *ws.IncomingNotification)) {
	newClient := ws.NewClient(conn, s.activeClients.RemoveClient, handleNotification, userID)
	newClient.Listen()

	s.activeClients.AddClient(newClient)
}

func (s *notificationService) Run() {
	redisPubsub := s.redisClient.Subscribe(s.ctx, pubsub.ChatChannel)
	chatChannel := redisPubsub.Channel()
	defer redisPubsub.Close()

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
