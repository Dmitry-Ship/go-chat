package domain

import (
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	Id        uuid.UUID `json:"id"`
	RoomId    uuid.UUID `json:"room_id"`
	UserId    uuid.UUID `json:"user_id"`
	CreatedAt int64     `json:"created_at"`
}

func NewParticipant(roomId uuid.UUID, userId uuid.UUID) *Participant {
	return &Participant{
		Id:        uuid.New(),
		RoomId:    roomId,
		UserId:    userId,
		CreatedAt: int64(time.Now().Unix()),
	}
}
