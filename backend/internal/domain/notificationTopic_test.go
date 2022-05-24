package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewNotificationTopic(t *testing.T) {
	name := "test"
	userID := uuid.New()

	notificationTopic, err := NewNotificationTopic(name, userID)

	assert.Equal(t, name, notificationTopic.Name)
	assert.Equal(t, userID, notificationTopic.UserID)
	assert.NotNil(t, notificationTopic.ID)
	assert.Nil(t, err)
}

func TestNewNotificationTopicEmptyTopic(t *testing.T) {
	name := ""
	userID := uuid.New()

	_, err := NewNotificationTopic(name, userID)

	assert.Equal(t, "topic is empty", err.Error())
}
