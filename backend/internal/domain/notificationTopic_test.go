package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewNotificationTopic(t *testing.T) {
	name := "test"
	userID := uuid.New()
	notificationTopicID := uuid.New()

	notificationTopic, err := NewNotificationTopic(notificationTopicID, name, userID)

	assert.Equal(t, name, notificationTopic.Name)
	assert.Equal(t, userID, notificationTopic.UserID)
	assert.Equal(t, notificationTopicID, notificationTopic.ID)
	assert.Nil(t, err)
}

func TestNewNotificationTopicEmptyTopic(t *testing.T) {
	name := ""
	userID := uuid.New()
	notificationTopicID := uuid.New()

	_, err := NewNotificationTopic(notificationTopicID, name, userID)

	assert.Equal(t, "topic is empty", err.Error())
}

func TestNewNotificationTopicLongTopic(t *testing.T) {
	name := ""

	for i := 0; i < 101; i++ {
		name += "a"
	}
	userID := uuid.New()
	notificationTopicID := uuid.New()

	_, err := NewNotificationTopic(notificationTopicID, name, userID)

	assert.Equal(t, "topic is too long", err.Error())
}
