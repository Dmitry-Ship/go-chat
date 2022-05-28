package postgres

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Conversation struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid"`
	Type     uint8
	IsActive bool
}

type GroupConversation struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid"`
	Name           string
	Avatar         string
	ConversationID uuid.UUID `gorm:"type:uuid"`
	OwnerID        uuid.UUID `gorm:"type:uuid"`
}

type DirectConversation struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid"`
	FromUserID     uuid.UUID `gorm:"type:uuid"`
	ToUserID       uuid.UUID `gorm:"type:uuid"`
	ConversationID uuid.UUID `gorm:"type:uuid"`
}

type Message struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid"`
	ConversationID uuid.UUID `gorm:"type:uuid"`
	UserID         uuid.UUID `gorm:"type:uuid"`
	Content        string
	Type           uint8
}

type Participant struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid"`
	ConversationID uuid.UUID `gorm:"type:uuid"`
	UserID         uuid.UUID `gorm:"type:uuid"`
	IsActive       bool
}

type User struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid"`
	Avatar       string
	Name         string
	Password     string
	RefreshToken string `gorm:"column:refresh_token"`
}

type UserNotificationTopic struct {
	gorm.Model
	ID     uuid.UUID `gorm:"type:uuid"`
	UserID uuid.UUID `gorm:"type:uuid"`
	Topic  string
}
