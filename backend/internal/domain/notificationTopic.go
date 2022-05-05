package domain

import (
	"github.com/google/uuid"
)

type NotificationTopicRepository interface {
	Store(notificationTopic *NotificationTopic) error
	DeleteByUserIDAndTopic(userId uuid.UUID, topic string) error
	DeleteAllByTopic(topic string) error
	GetUserIDsByTopic(topic string) ([]uuid.UUID, error)
}

type NotificationTopic struct {
	ID     uuid.UUID
	Name   string
	UserID uuid.UUID
}

func NewNotificationTopic(name string, userID uuid.UUID) *NotificationTopic {
	return &NotificationTopic{
		ID:     uuid.New(),
		Name:   name,
		UserID: userID,
	}
}
