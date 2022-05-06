package domain

type aggregate struct {
	events []DomainEvent
}

type EventPublisher interface {
	Publish(topic string, event DomainEvent)
}

func (a *aggregate) Dispatch(publisher EventPublisher) {
	for _, event := range a.events {
		publisher.Publish(DomainEventChannel, event)
	}
}

func (a *aggregate) AddEvent(event DomainEvent) {
	a.events = append(a.events, event)
}

type Aggregate interface {
	AddEvent(event *DomainEvent)
	Dispatch(publisher EventPublisher)
}
