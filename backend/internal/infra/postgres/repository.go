package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"

	"gorm.io/gorm"
)

type repository struct {
	db             *gorm.DB
	eventPublisher infra.EventPublisher
}

func newRepository(db *gorm.DB, eventPublisher infra.EventPublisher) *repository {
	return &repository{
		db:             db,
		eventPublisher: eventPublisher,
	}
}

func (r *repository) dispatchEvents(aggregate domain.Aggregate) {
	for _, event := range aggregate.GetEvents() {
		r.eventPublisher.Publish(domain.DomainEventChannel, event)
	}
}
