package domain

import (
	"github.com/google/uuid"
)

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
