package infra

import (
	"sync"
)

type Event struct {
	Topic string
	Data  interface{}
}

type EventBus struct {
	mu                  sync.RWMutex
	topicSubscribersMap map[string][]chan Event
	isClosed            bool
}

func NewEventBus() *EventBus {
	return &EventBus{
		mu:                  sync.RWMutex{},
		topicSubscribersMap: make(map[string][]chan Event),
	}
}

func (eb *EventBus) Subscribe(topic string) <-chan Event {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	subscriptionChannel := make(chan Event, 1)

	eb.topicSubscribersMap[topic] = append(eb.topicSubscribersMap[topic], subscriptionChannel)

	return subscriptionChannel
}

func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if eb.isClosed {
		return
	}

	for _, subscriptionChannel := range eb.topicSubscribersMap[topic] {
		subscriptionChannel <- Event{Topic: topic, Data: data}
	}
}

func (eb *EventBus) Close() {
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
