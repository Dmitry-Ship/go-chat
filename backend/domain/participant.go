package domain

import (
	"math/rand"
	"time"
)

type Participant struct {
	Id        int32 `json:"id"`
	RoomId    int32 `json:"room_id"`
	UserId    int32 `json:"user_id"`
	CreatedAt int64 `json:"created_at"`
}

func NewParticipant(roomId int32, userId int32) *Participant {
	return &Participant{
		Id:        rand.Int31(),
		RoomId:    roomId,
		UserId:    userId,
		CreatedAt: int64(time.Now().Unix()),
	}
}
