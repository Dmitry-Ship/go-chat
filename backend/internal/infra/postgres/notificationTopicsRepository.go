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

func (r *notificationTopicRepository) GetAllNotificationTopics(userID uuid.UUID) ([]string, error) {
	var topics []string

	err := r.db.Model(&UserNotificationTopic{}).Where("user_id = ?", userID).Select("topic").Find(&topics).Error

	return topics, err
}
