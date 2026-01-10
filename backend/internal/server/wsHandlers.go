package server

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

func (s *Server) handleNotification(userID uuid.UUID, message []byte) {
	var data json.RawMessage

	notification := struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		Data: &data,
	}

	if err := json.Unmarshal(message, &notification); err != nil {
		log.Println(err)
		return
	}

	log.Println("Received WebSocket notification type:", notification.Type)
}
