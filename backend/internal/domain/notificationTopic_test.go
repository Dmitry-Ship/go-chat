package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewNotificationTopic(t *testing.T) {
	name := "test"
	userID := uuid.New()

	notificationTopic := NewNotificationTopic(name, userID)

	assert.Equal(t, name, notificationTopic.Name)
	assert.Equal(t, userID, notificationTopic.UserID)
	assert.NotNil(t, notificationTopic.ID)
}
