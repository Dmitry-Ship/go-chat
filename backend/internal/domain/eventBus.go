package domain

import (
	"sync"
)

type EventsSubscriber interface {
	Subscribe(topic string) <-chan DomainEvent
}

type eventBus struct {
	mu                  sync.RWMutex
	topicSubscribersMap map[string][]chan DomainEvent
	isClosed            bool
}

func NewEventBus() *eventBus {
	return &eventBus{
		mu:                  sync.RWMutex{},
		topicSubscribersMap: make(map[string][]chan DomainEvent),
	}
}

func (eb *eventBus) Subscribe(topic string) <-chan DomainEvent {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	subscriptionChannel := make(chan DomainEvent, 100)

	eb.topicSubscribersMap[topic] = append(eb.topicSubscribersMap[topic], subscriptionChannel)

	return subscriptionChannel
}

func (eb *eventBus) Publish(topic string, event DomainEvent) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if eb.isClosed {
		return
	}

	for _, subscriptionChannel := range eb.topicSubscribersMap[topic] {
		go func(subscriptionChannel chan<- DomainEvent) {
			subscriptionChannel <- event
		}(subscriptionChannel)
	}
}

func (eb *eventBus) Close() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if !eb.isClosed {
		eb.isClosed = true
		for _, topicGroup := range eb.topicSubscribersMap {
			for _, subscriberChannel := range topicGroup {
				close(subscriberChannel)
			}
		}
	}
}
