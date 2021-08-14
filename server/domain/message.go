package domain

import "time"

type Message struct {
	Id        int64  `json:"id"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	Type      string `json:"type"`
	Sender    string `json:"sender"`
}

func NewMessage(content string, messageType string, sender string) Message {
	return Message{
		Id:        int64(time.Now().UnixNano()),
		Content:   content,
		CreatedAt: int64(time.Now().Unix()),
		Type:      messageType,
		Sender:    sender,
	}
}
