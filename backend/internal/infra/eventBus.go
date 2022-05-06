package infra

import (
	"sync"
)

type event struct {
	Topic string
	Data  interface{}
}

type EventPublisher interface {
	Publish(topic string, data interface{})
}

type EventsSubscriber interface {
	Subscribe(topic string) <-chan event
}

type eventBus struct {
	mu                  sync.RWMutex
	topicSubscribersMap map[string][]chan event
	isClosed            bool
}

func NewEventBus() *eventBus {
	return &eventBus{
		mu:                  sync.RWMutex{},
		topicSubscribersMap: make(map[string][]chan event),
	}
}

func (eb *eventBus) Subscribe(topic string) <-chan event {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	subscriptionChannel := make(chan event, 100)

	eb.topicSubscribersMap[topic] = append(eb.topicSubscribersMap[topic], subscriptionChannel)

	return subscriptionChannel
}

func (eb *eventBus) Publish(topic string, data interface{}) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if eb.isClosed {
		return
	}

	for _, subscriptionChannel := range eb.topicSubscribersMap[topic] {
		go func(subscriptionChannel chan<- event) {
			subscriptionChannel <- event{Topic: topic, Data: data}
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
