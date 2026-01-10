package services

import (
	"context"
	"sync"
	"time"

	"GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
)

type pendingBatch struct {
	notifications []ws.OutgoingNotification
	timer         *time.Timer
	mu            sync.Mutex
}

type FlushFunc func(key batchKey, notifications []ws.OutgoingNotification)

type NotificationBatcher interface {
	Add(channelID, userID uuid.UUID, notification ws.OutgoingNotification) error
	Start(ctx context.Context, flushFunc FlushFunc)
}

type notificationBatcher struct {
	batchChannel chan batchNotification
	maxBatchSize int
	batchTimeout time.Duration
	batches      map[batchKey]*pendingBatch
	mu           sync.RWMutex
}

func NewNotificationBatcher(maxBatchSize int, batchTimeout time.Duration) NotificationBatcher {
	return &notificationBatcher{
		batchChannel: make(chan batchNotification, 1000),
		maxBatchSize: maxBatchSize,
		batchTimeout: batchTimeout,
		batches:      make(map[batchKey]*pendingBatch),
	}
}

func (b *notificationBatcher) Add(channelID, userID uuid.UUID, notification ws.OutgoingNotification) error {
	b.batchChannel <- batchNotification{
		channelID:    channelID,
		notification: notification,
	}
	return nil
}

func (b *notificationBatcher) Start(ctx context.Context, flushFunc FlushFunc) {
	defer func() {
		b.flushAll(flushFunc)
	}()

	for {
		select {
		case bn, ok := <-b.batchChannel:
			if !ok {
				return
			}
			b.addToBatch(bn, flushFunc)

		case <-ctx.Done():
			b.flushAll(flushFunc)
			return
		}
	}
}

func (b *notificationBatcher) addToBatch(bn batchNotification, flushFunc FlushFunc) {
	key := batchKey{channelID: bn.channelID, userID: bn.notification.UserID}

	b.mu.Lock()
	pb, exists := b.batches[key]
	if !exists {
		pb = &pendingBatch{
			notifications: []ws.OutgoingNotification{bn.notification},
		}
		b.batches[key] = pb
		b.mu.Unlock()

		pb.mu.Lock()
		defer pb.mu.Unlock()

		pb.timer = time.AfterFunc(b.batchTimeout, func() {
			b.flushBatch(key, pb, flushFunc)
		})
		return
	}

	b.mu.Unlock()

	pb.mu.Lock()
	defer pb.mu.Unlock()

	pb.notifications = append(pb.notifications, bn.notification)
	if len(pb.notifications) >= b.maxBatchSize {
		pb.timer.Stop()
		pb.notifications = nil
		pb.timer = nil
		flushFunc(key, []ws.OutgoingNotification{bn.notification})
	}
}

func (b *notificationBatcher) flushBatch(key batchKey, pb *pendingBatch, flushFunc FlushFunc) {
	pb.mu.Lock()
	notifications := pb.notifications
	pb.notifications = nil
	pb.timer = nil
	pb.mu.Unlock()

	if len(notifications) == 0 {
		return
	}

	flushFunc(key, notifications)

	b.mu.Lock()
	delete(b.batches, key)
	b.mu.Unlock()
}

func (b *notificationBatcher) flushAll(flushFunc FlushFunc) {
	b.mu.Lock()
	batches := b.batches
	b.batches = make(map[batchKey]*pendingBatch)
	b.mu.Unlock()

	for key, pb := range batches {
		pb.mu.Lock()
		notifications := pb.notifications
		timer := pb.timer
		pb.notifications = nil
		pb.timer = nil
		pb.mu.Unlock()

		if timer != nil {
			timer.Stop()
		}

		if len(notifications) > 0 {
			flushFunc(key, notifications)
		}
	}
}
