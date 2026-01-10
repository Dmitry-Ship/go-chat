package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMessageDeduplicator(t *testing.T) {
	dedup := NewMessageDeduplicator(100, 5*time.Minute)

	assert.NotNil(t, dedup)
	assert.IsType(t, &messageDeduplicator{}, dedup)
}

func TestMessageDeduplicator_AlreadySent(t *testing.T) {
	dedup := NewMessageDeduplicator(100, 5*time.Minute)

	result := dedup.AlreadySent("test-message-1")

	assert.False(t, result)
}

func TestMessageDeduplicator_MarkSent(t *testing.T) {
	dedup := NewMessageDeduplicator(100, 5*time.Minute)

	dedup.MarkSent("test-message-1")

	result := dedup.AlreadySent("test-message-1")
	assert.True(t, result)
}

func TestMessageDeduplicator_MarkSentMultiple(t *testing.T) {
	dedup := NewMessageDeduplicator(100, 5*time.Minute)

	dedup.MarkSent("test-message-1")
	dedup.MarkSent("test-message-2")
	dedup.MarkSent("test-message-3")

	assert.True(t, dedup.AlreadySent("test-message-1"))
	assert.True(t, dedup.AlreadySent("test-message-2"))
	assert.True(t, dedup.AlreadySent("test-message-3"))
	assert.False(t, dedup.AlreadySent("test-message-4"))
}

func TestMessageDeduplicator_MarkSentUpdatesTimestamp(t *testing.T) {
	dedup := NewMessageDeduplicator(100, 5*time.Minute)

	dedup.MarkSent("test-message-1")
	time.Sleep(10 * time.Millisecond)
	dedup.MarkSent("test-message-1")

	result := dedup.AlreadySent("test-message-1")
	assert.True(t, result)
}

func TestMessageDeduplicator_LRU_Eviction(t *testing.T) {
	dedup := NewMessageDeduplicator(3, 5*time.Minute)

	for i := 0; i < 5; i++ {
		dedup.MarkSent("test-message-" + string(rune('1'+i)))
	}

	assert.True(t, dedup.AlreadySent("test-message-3"))
	assert.True(t, dedup.AlreadySent("test-message-4"))
	assert.True(t, dedup.AlreadySent("test-message-5"))
	assert.False(t, dedup.AlreadySent("test-message-1"))
	assert.False(t, dedup.AlreadySent("test-message-2"))
}

func TestMessageDeduplicator_StartCleanup(t *testing.T) {
	dedup := NewMessageDeduplicator(100, 50*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())

	go dedup.StartCleanup(ctx)

	dedup.MarkSent("test-message-1")
	time.Sleep(20 * time.Millisecond)
	assert.True(t, dedup.AlreadySent("test-message-1"))

	time.Sleep(100 * time.Millisecond)

	cancel()
}

func TestMessageDeduplicator_ContextDone(t *testing.T) {
	dedup := NewMessageDeduplicator(100, 50*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())

	go dedup.StartCleanup(ctx)

	dedup.MarkSent("test-message-1")

	cancel()

	time.Sleep(20 * time.Millisecond)

	result := dedup.AlreadySent("test-message-1")
	assert.True(t, result)
}

func TestMessageDeduplicator_AlreadySent_EmptyMessageID(t *testing.T) {
	dedup := NewMessageDeduplicator(100, 5*time.Minute)

	result := dedup.AlreadySent("")

	assert.False(t, result)
}

func TestMessageDeduplicator_AlreadySent_NonExistent(t *testing.T) {
	dedup := NewMessageDeduplicator(100, 5*time.Minute)

	dedup.MarkSent("test-message-1")
	dedup.MarkSent("test-message-2")

	result := dedup.AlreadySent("test-message-3")

	assert.False(t, result)
}
