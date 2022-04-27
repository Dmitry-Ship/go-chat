package domain

type aggregate struct {
	events []DomainEvent
}

type EventPublisher interface {
	Publish(event DomainEvent)
}

func (a *aggregate) Raise(publisher EventPublisher) {
	for _, event := range a.events {
		publisher.Publish(event)
	}
}

func (a *aggregate) AddEvent(event DomainEvent) {
	a.events = append(a.events, event)
}

type Aggregate interface {
	AddEvent(event *DomainEvent)
	Raise(publisher EventPublisher)
}
