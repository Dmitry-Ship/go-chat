package domain

type aggregate struct {
	events []DomainEvent
}

func (a *aggregate) GetEvents() []DomainEvent {
	return a.events
}

func (a *aggregate) AddEvent(event DomainEvent) {
	a.events = append(a.events, event)
}

type Aggregate interface {
	AddEvent(event DomainEvent)
	GetEvents() []DomainEvent
}
