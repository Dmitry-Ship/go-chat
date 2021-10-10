package domain

import (
	"math/rand"
	"time"
)

type ChatMessage struct {
	Id        int32  `json:"id"`
	RoomId    int32  `json:"room_id"`
	Content   string `json:"content"`
	CreatedAt int32  `json:"created_at"`
	Type      string `json:"type"`
	UserId    int32  `json:"user_id"`
}

func NewChatMessage(content string, messageType string, roomId int32, userId int32) ChatMessage {
	return ChatMessage{
		Id:        int32(rand.Int31()),
		RoomId:    roomId,
		Content:   content,
		CreatedAt: int32(time.Now().Unix()),
		Type:      messageType,
		UserId:    userId,
	}
}
