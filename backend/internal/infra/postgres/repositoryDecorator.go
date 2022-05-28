package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"
)

type repoInterface interface {
	Store(domain.Aggregate) error
	Update(domain.Aggregate) error
}

type repositoryDecorator struct {
	repo           repoInterface
	eventPublisher infra.EventPublisher
}

func newRepositoryDecorator(repo repoInterface, eventPublisher infra.EventPublisher) *repositoryDecorator {
	return &repositoryDecorator{
		repo:           repo,
		eventPublisher: eventPublisher,
	}
}

func (r *repositoryDecorator) dispatchEvents(aggregate domain.Aggregate) {
	for _, event := range aggregate.GetEvents() {
		r.eventPublisher.Publish(domain.DomainEventChannel, event)
	}
}

func (r *repositoryDecorator) Store(aggregate domain.Aggregate) error {

	err := r.repo.Store(aggregate)

	if err != nil {
		return err
	}

	r.dispatchEvents(aggregate)

	return nil
}
