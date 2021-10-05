package domain

import "time"

type Message struct {
	Id        int64  `json:"id"`
	RoomId    int64  `json:"room_id"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	Type      string `json:"type"`
	User      *User  `json:"user"`
}

func NewMessage(content string, messageType string, roomId int64, user *User) Message {
	return Message{
		Id:        int64(time.Now().UnixNano()),
		RoomId:    roomId,
		Content:   content,
		CreatedAt: int64(time.Now().Unix()),
		Type:      messageType,
		User:      user,
	}
}
