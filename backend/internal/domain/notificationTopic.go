package domain

import (
	"errors"

	"github.com/google/uuid"
)

type NotificationTopicRepository interface {
	Store(notificationTopic *NotificationTopic) error
	DeleteByUserIDAndTopic(userID uuid.UUID, topic string) error
	DeleteAllByTopic(topic string) error
	GetUserIDsByTopic(topic string) ([]uuid.UUID, error)
}

type NotificationTopic struct {
	ID     uuid.UUID
	Name   string
	UserID uuid.UUID
}

func NewNotificationTopic(topic string, userID uuid.UUID) (*NotificationTopic, error) {
	if topic == "" {
		return nil, errors.New("topic is empty")
	}

	if len(topic) > 100 {
		return nil, errors.New("topic is too long")
	}

	return &NotificationTopic{
		ID:     uuid.New(),
		Name:   topic,
		UserID: userID,
	}, nil
}
