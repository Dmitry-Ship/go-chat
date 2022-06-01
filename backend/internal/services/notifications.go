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

type buildFunc func(userID uuid.UUID) (*ws.OutgoingNotification, error)

type NotificationService interface {
	SendToConversation(conversationId uuid.UUID, buildMessage buildFunc) error
	RegisterClient(conn *websocket.Conn, userID uuid.UUID, handleNotification func(userID uuid.UUID, message []byte))
	Run()
}

type notificationService struct {
	ctx           context.Context
	participants  domain.ParticipantRepository
	activeClients ws.ActiveClients
	redisClient   *redis.Client
}

func NewNotificationService(
	ctx context.Context,
	participants domain.ParticipantRepository,
	redisClient *redis.Client,
) *notificationService {
	return &notificationService{
		ctx:           ctx,
		participants:  participants,
		activeClients: ws.NewActiveClients(),
		redisClient:   redisClient,
	}
}

func (s *notificationService) sendWorker(
	userID uuid.UUID,
	buildMessage buildFunc,
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

func (s *notificationService) broadcast(ids []uuid.UUID, buildMessage buildFunc) error {
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

func (s *notificationService) SendToConversation(conversationId uuid.UUID, buildMessage buildFunc) error {
	ids, err := s.participants.GetIDsByConversationID(conversationId)

	if err != nil {
		return err
	}

	return s.broadcast(ids, buildMessage)
}

func (s *notificationService) RegisterClient(conn *websocket.Conn, userID uuid.UUID, handleNotification func(userID uuid.UUID, message []byte)) {
	newClient := ws.NewClient(conn, s.activeClients.RemoveClient, handleNotification, userID)

	s.activeClients.AddClient(newClient)

	go newClient.WritePump()
	newClient.ReadPump()
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
