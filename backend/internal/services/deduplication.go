package services

import (
	"container/list"
	"context"
	"sync"
	"time"
)

type lruEntry struct {
	messageID string
	timestamp time.Time
}

type MessageDeduplicator interface {
	AlreadySent(messageID string) bool
	MarkSent(messageID string)
	StartCleanup(ctx context.Context)
}

type messageDeduplicator struct {
	capacity  int
	ttl       time.Duration
	entries   map[string]*list.Element
	evictList *list.List
	mu        sync.RWMutex
}

func NewMessageDeduplicator(capacity int, ttl time.Duration) MessageDeduplicator {
	return &messageDeduplicator{
		capacity:  capacity,
		ttl:       ttl,
		entries:   make(map[string]*list.Element),
		evictList: list.New(),
	}
}

func (d *messageDeduplicator) AlreadySent(messageID string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	elem, exists := d.entries[messageID]
	if !exists {
		return false
	}

	entry := elem.Value.(*lruEntry)
	return time.Since(entry.timestamp) <= d.ttl
}

func (d *messageDeduplicator) MarkSent(messageID string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if elem, exists := d.entries[messageID]; exists {
		entry := elem.Value.(*lruEntry)
		entry.timestamp = time.Now()
		d.evictList.MoveToFront(elem)
		return
	}

	if d.evictList.Len() >= d.capacity {
		d.evictOldest()
	}

	entry := &lruEntry{
		messageID: messageID,
		timestamp: time.Now(),
	}
	elem := d.evictList.PushFront(entry)
	d.entries[messageID] = elem
}

func (d *messageDeduplicator) evictOldest() {
	elem := d.evictList.Back()
	if elem != nil {
		d.evictList.Remove(elem)
		entry := elem.Value.(*lruEntry)
		delete(d.entries, entry.messageID)
	}
}

func (d *messageDeduplicator) StartCleanup(ctx context.Context) {
	ticker := time.NewTicker(d.ttl)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d.cleanupStale()
		case <-ctx.Done():
			return
		}
	}
}

func (d *messageDeduplicator) cleanupStale() {
	d.mu.Lock()
	defer d.mu.Unlock()

	cutoff := time.Now().Add(-d.ttl)
	for elem := d.evictList.Back(); elem != nil; {
		entry := elem.Value.(*lruEntry)
		if entry.timestamp.After(cutoff) {
			break
		}

		prev := elem.Prev()
		d.evictList.Remove(elem)
		delete(d.entries, entry.messageID)
		elem = prev
	}
}
