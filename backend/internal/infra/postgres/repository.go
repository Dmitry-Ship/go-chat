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
		r.eventPublisher.Publish(event.GetTopic(), event)
	}
}

func (r *repository) beginTransaction(aggregate domain.Aggregate, callback func(tx *gorm.DB) error) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := callback(tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	r.dispatchEvents(aggregate)

	return nil
}

func (r *repository) store(aggregate domain.Aggregate, persistence interface{}) error {
	if err := r.db.Create(persistence).Error; err != nil {
		return err
	}

	r.dispatchEvents(aggregate)

	return nil
}

func (r *repository) update(aggregate domain.Aggregate, persistence interface{}) error {
	if err := r.db.Save(persistence).Error; err != nil {
		return err
	}

	r.dispatchEvents(aggregate)

	return nil
}
