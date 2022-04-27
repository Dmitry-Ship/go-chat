package postgres

import (
	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type notificationTopicRepository struct {
	db *gorm.DB
}

func NewNotificationTopicRepository(db *gorm.DB) *notificationTopicRepository {
	return &notificationTopicRepository{
		db: db,
	}
}

func (r *notificationTopicRepository) Store(notificationTopic *domain.NotificationTopic) error {
	persistence := UserNotificationTopic{
		Topic:  notificationTopic.Name,
		UserID: notificationTopic.UserID,
		ID:     notificationTopic.ID,
	}

	err := r.db.Create(persistence).Error

	return err

}

func (r *notificationTopicRepository) DeleteByUserIDAndTopic(userID uuid.UUID, topic string) error {
	persistence := UserNotificationTopic{}

	err := r.db.Where("user_id = ?", userID).Where("topic = ?", topic).Delete(persistence).Error

	return err
}

func (r *notificationTopicRepository) DeleteByTopic(topic string) error {
	persistence := UserNotificationTopic{}

	err := r.db.Where("topic = ?", topic).Delete(persistence).Error

	return err
}

func (r *notificationTopicRepository) GetUserIDsByTopic(topic string) ([]uuid.UUID, error) {
	var ids []uuid.UUID

	err := r.db.Model(&UserNotificationTopic{}).Where("topic = ?", topic).Select("user_id").Find(&ids).Error

	return ids, err
}
