package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type notificationTopicRepository struct {
	repository
}

func NewNotificationTopicRepository(db *gorm.DB, eventPublisher infra.EventPublisher) *notificationTopicRepository {
	return &notificationTopicRepository{
		repository: *newRepository(db, eventPublisher),
	}
}

func (r *notificationTopicRepository) Store(notificationTopic *domain.NotificationTopic) error {
	persistence := &UserNotificationTopic{
		Topic:  notificationTopic.Name,
		UserID: notificationTopic.UserID,
		ID:     notificationTopic.ID,
	}

	return r.store(notificationTopic, persistence)
}

func (r *notificationTopicRepository) DeleteByUserIDAndTopic(userID uuid.UUID, topic string) error {
	persistence := &UserNotificationTopic{}

	err := r.db.Where("user_id = ?", userID).Where("topic = ?", topic).Delete(persistence).Error

	return err
}

func (r *notificationTopicRepository) DeleteAllByTopic(topic string) error {
	persistence := &UserNotificationTopic{}

	err := r.db.Where("topic = ?", topic).Delete(persistence).Error

	return err
}

func (r *notificationTopicRepository) GetUserIDsByTopic(topic string) ([]uuid.UUID, error) {
	var ids []uuid.UUID

	err := r.db.Model(&UserNotificationTopic{}).Where("topic = ?", topic).Select("user_id").Find(&ids).Error

	if err != nil {
		return nil, err
	}

	return ids, nil
}
