package domain

import (
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid"`
	RoomId    uuid.UUID `json:"room_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt int64     `json:"created_at"`
}

func NewParticipant(roomId uuid.UUID, userId uuid.UUID) *Participant {
	return &Participant{
		ID:        uuid.New(),
		RoomId:    roomId,
		UserID:    userId,
		CreatedAt: int64(time.Now().Unix()),
	}
}
