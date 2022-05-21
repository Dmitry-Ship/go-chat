package postgres

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID        uuid.UUID `gorm:"type:uuid"`
	Type      uint8
	CreatedAt time.Time
	IsActive  bool
}

type PublicConversation struct {
	ID             uuid.UUID `gorm:"type:uuid"`
	Name           string
	Avatar         string
	ConversationID uuid.UUID `gorm:"type:uuid"`
	OwnerID        uuid.UUID `gorm:"type:uuid"`
}

type PrivateConversation struct {
	ID             uuid.UUID `gorm:"type:uuid"`
	FromUserID     uuid.UUID `gorm:"type:uuid"`
	ToUserID       uuid.UUID `gorm:"type:uuid"`
	ConversationID uuid.UUID `gorm:"type:uuid"`
}

type Message struct {
	ID             uuid.UUID `gorm:"type:uuid"`
	ConversationID uuid.UUID `gorm:"type:uuid"`
	UserID         uuid.UUID `gorm:"type:uuid"`
	CreatedAt      time.Time
	Type           uint8
}

type TextMessage struct {
	ID        uuid.UUID `gorm:"type:uuid"`
	MessageID uuid.UUID `gorm:"type:uuid"`
	Text      string
}

type ConversationRenamedMessage struct {
	ID        uuid.UUID `gorm:"type:uuid"`
	MessageID uuid.UUID `gorm:"type:uuid"`
	NewName   string
}

type Participant struct {
	ID             uuid.UUID `gorm:"type:uuid"`
	ConversationID uuid.UUID `gorm:"type:uuid"`
	UserID         uuid.UUID `gorm:"type:uuid"`
	CreatedAt      time.Time
	IsActive       bool
}

type User struct {
	ID           uuid.UUID `gorm:"type:uuid"`
	Avatar       string
	Name         string
	Password     string
	CreatedAt    time.Time
	RefreshToken string `gorm:"column:refresh_token"`
}

type UserNotificationTopic struct {
	ID     uuid.UUID `gorm:"type:uuid"`
	UserID uuid.UUID `gorm:"type:uuid"`
	Topic  string
}
