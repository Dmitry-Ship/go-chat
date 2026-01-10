package ws

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	unregisterClient := func(c *Client) {}
	handleIncomingNotification := func(userID uuid.UUID, message []byte) {}
	userID := uuid.New()

	client := NewClient(nil, unregisterClient, handleIncomingNotification, userID)

	assert.NotNil(t, client)
	assert.NotEqual(t, uuid.Nil, client.Id)
	assert.Equal(t, userID, client.UserID)
	assert.NotNil(t, client.sendChannel)
	assert.Equal(t, SendChannelSize, cap(client.sendChannel))
}

func TestClient_SendNotification(t *testing.T) {
	unregisterClient := func(c *Client) {}
	handleIncomingNotification := func(userID uuid.UUID, message []byte) {}
	userID := uuid.New()

	client := NewClient(nil, unregisterClient, handleIncomingNotification, userID)

	notification := OutgoingNotification{
		Type:    "test",
		UserID:  userID,
		Payload: "test payload",
	}

	client.SendNotification(notification)

	assert.Len(t, client.sendChannel, 1)

	received := <-client.sendChannel
	assert.Equal(t, notification.Type, received.Type)
	assert.Equal(t, notification.Payload, received.Payload)
}

func TestClient_SendNotification_Multiple(t *testing.T) {
	unregisterClient := func(c *Client) {}
	handleIncomingNotification := func(userID uuid.UUID, message []byte) {}
	userID := uuid.New()

	client := NewClient(nil, unregisterClient, handleIncomingNotification, userID)

	for i := 0; i < 5; i++ {
		notification := OutgoingNotification{
			Type:    "test",
			UserID:  userID,
			Payload: i,
		}
		client.SendNotification(notification)
	}

	assert.Len(t, client.sendChannel, 5)
}

func TestClient_ConnectionOptions(t *testing.T) {
	unregisterClient := func(c *Client) {}
	handleIncomingNotification := func(userID uuid.UUID, message []byte) {}
	userID := uuid.New()

	client := NewClient(nil, unregisterClient, handleIncomingNotification, userID)

	assert.Equal(t, WriteWait, client.connectionOptions.writeWait)
	assert.Equal(t, PongWait, client.connectionOptions.pongWait)
	assert.Equal(t, PingPeriod, client.connectionOptions.pingPeriod)
	assert.Equal(t, int64(MaxMessageSize), client.connectionOptions.maxMessageSize)
}

func TestClient_ChannelIDsInitialized(t *testing.T) {
	unregisterClient := func(c *Client) {}
	handleIncomingNotification := func(userID uuid.UUID, message []byte) {}
	userID := uuid.New()

	client := NewClient(nil, unregisterClient, handleIncomingNotification, userID)

	assert.NotNil(t, client.channelIDs)
	assert.Equal(t, 0, len(client.channelIDs))
}
