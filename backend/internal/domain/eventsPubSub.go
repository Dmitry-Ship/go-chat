package domain

import (
	"sync"
)

type EventsSubscriber interface {
	Subscribe(topic string) <-chan DomainEvent
}

type pubsub struct {
	mu                  sync.RWMutex
	topicSubscribersMap map[string][]chan DomainEvent
	isClosed            bool
}

func NewPubsub() *pubsub {
	return &pubsub{
		mu:                  sync.RWMutex{},
		topicSubscribersMap: make(map[string][]chan DomainEvent),
	}
}

func (ps *pubsub) Subscribe(topic string) <-chan DomainEvent {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	subscriptionChannel := make(chan DomainEvent, 10)

	ps.topicSubscribersMap[topic] = append(ps.topicSubscribersMap[topic], subscriptionChannel)

	return subscriptionChannel
}

func (ps *pubsub) Publish(event DomainEvent) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.isClosed {
		return
	}

	for _, subscriptionChannel := range ps.topicSubscribersMap[event.GetName()] {
		go func(subscriptionChannel chan DomainEvent) {
			subscriptionChannel <- event
		}(subscriptionChannel)
	}
}

func (ps *pubsub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.isClosed {
		ps.isClosed = true
		for _, topicSubscribersMap := range ps.topicSubscribersMap {
			for _, ch := range topicSubscribersMap {
				close(ch)
			}
		}
	}
}
