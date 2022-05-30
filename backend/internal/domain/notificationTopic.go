package domain

import (
	"errors"

	"github.com/google/uuid"
)

type NotificationTopicRepository interface {
	Store(notificationTopic *NotificationTopic) error
	DeleteByUserIDAndTopic(userID uuid.UUID, topic string) error
	GetUserIDsByTopic(topic string) ([]uuid.UUID, error)
}

type NotificationTopic struct {
	aggregate
	ID     uuid.UUID
	Name   string
	UserID uuid.UUID
}

func NewNotificationTopic(notificationTopicID uuid.UUID, topic string, userID uuid.UUID) (*NotificationTopic, error) {
	if topic == "" {
		return nil, errors.New("topic is empty")
	}

	if len(topic) > 100 {
		return nil, errors.New("topic is too long")
	}

	return &NotificationTopic{
		ID:     notificationTopicID,
		Name:   topic,
		UserID: userID,
	}, nil
}
