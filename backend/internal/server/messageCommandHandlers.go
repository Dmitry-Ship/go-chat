package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (s *Server) handleSendMessage(w http.ResponseWriter, r *http.Request) {
	request := struct {
		Content        string    `json:"content"`
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
		return
	}

	err := s.message.SendTextMessage(r.Context(), request.ConversationId, userID, request.Content)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	response := struct {
		MessageId string    `json:"message_id"`
		CreatedAt time.Time `json:"created_at"`
	}{
		MessageId: uuid.New().String(),
		CreatedAt: time.Now(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}
