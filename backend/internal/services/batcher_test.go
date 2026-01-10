package services

import (
	"context"
	"sync"
	"testing"
	"time"

	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewNotificationBatcher(t *testing.T) {
	batcher := NewNotificationBatcher(10, time.Second)

	assert.NotNil(t, batcher)
	assert.IsType(t, &notificationBatcher{}, batcher)
}

func TestNotificationBatcher_Add(t *testing.T) {
	batcher := NewNotificationBatcher(10, time.Second)

	channelID := uuid.New()
	userID := uuid.New()
	notification := ws.OutgoingNotification{
		Type:    "test",
		UserID:  userID,
		Payload: "test payload",
	}

	err := batcher.Add(channelID, userID, notification)

	assert.NoError(t, err)
}

func TestNotificationBatcher_AddMultiple(t *testing.T) {
	batcher := NewNotificationBatcher(10, time.Second)

	channelID := uuid.New()
	userID := uuid.New()

	for i := 0; i < 5; i++ {
		notification := ws.OutgoingNotification{
			Type:    "test",
			UserID:  userID,
			Payload: i,
		}
		err := batcher.Add(channelID, userID, notification)
		assert.NoError(t, err)
	}
}

func TestNotificationBatcher_Start_Stop(t *testing.T) {
	batcher := NewNotificationBatcher(10, 50*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())

	var flushedKeys []batchKey
	var mu sync.Mutex
	var flushCount int

	flushFunc := func(key batchKey, notifications []ws.OutgoingNotification) {
		mu.Lock()
		flushedKeys = append(flushedKeys, key)
		flushCount++
		mu.Unlock()
	}

	go batcher.Start(ctx, flushFunc)

	channelID := uuid.New()
	userID := uuid.New()

	for i := 0; i < 3; i++ {
		notification := ws.OutgoingNotification{
			Type:    "test",
			UserID:  userID,
			Payload: i,
		}
		_ = batcher.Add(channelID, userID, notification)
	}

	time.Sleep(150 * time.Millisecond)
	cancel()
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	assert.GreaterOrEqual(t, flushCount, 1)
	mu.Unlock()
}

func TestNotificationBatcher_MaxBatchSize(t *testing.T) {
	batcher := NewNotificationBatcher(3, 10*time.Second)

	ctx, cancel := context.WithCancel(context.Background())

	var flushedKeys []batchKey
	var mu sync.Mutex
	var flushCount int

	flushFunc := func(key batchKey, notifications []ws.OutgoingNotification) {
		mu.Lock()
		flushedKeys = append(flushedKeys, key)
		flushCount++
		mu.Unlock()
	}

	go batcher.Start(ctx, flushFunc)

	channelID := uuid.New()
	userID := uuid.New()

	for i := 0; i < 5; i++ {
		notification := ws.OutgoingNotification{
			Type:    "test",
			UserID:  userID,
			Payload: i,
		}
		_ = batcher.Add(channelID, userID, notification)
	}

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	assert.GreaterOrEqual(t, flushCount, 1)
	mu.Unlock()

	cancel()
}

func TestNotificationBatcher_ContextDone(t *testing.T) {
	batcher := NewNotificationBatcher(10, 10*time.Second)

	ctx, cancel := context.WithCancel(context.Background())

	var flushedKeys []batchKey
	var mu sync.Mutex
	var flushCount int

	flushFunc := func(key batchKey, notifications []ws.OutgoingNotification) {
		mu.Lock()
		flushedKeys = append(flushedKeys, key)
		flushCount++
		mu.Unlock()
	}

	go batcher.Start(ctx, flushFunc)

	channelID := uuid.New()
	userID := uuid.New()

	for i := 0; i < 3; i++ {
		notification := ws.OutgoingNotification{
			Type:    "test",
			UserID:  userID,
			Payload: i,
		}
		_ = batcher.Add(channelID, userID, notification)
	}

	time.Sleep(50 * time.Millisecond)
	cancel()
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	assert.GreaterOrEqual(t, flushCount, 1)
	mu.Unlock()
}
