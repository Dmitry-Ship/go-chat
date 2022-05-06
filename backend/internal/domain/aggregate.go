package domain

type aggregate struct {
	events []DomainEvent
}

type DomainEventPublisher interface {
	Publish(topic string, data interface{})
}

func (a *aggregate) Dispatch(publisher DomainEventPublisher) {
	for _, event := range a.events {
		publisher.Publish(DomainEventChannel, event)
	}
}

func (a *aggregate) AddEvent(event DomainEvent) {
	a.events = append(a.events, event)
}

type Aggregate interface {
	AddEvent(event *DomainEvent)
	Dispatch(publisher DomainEventPublisher)
}
