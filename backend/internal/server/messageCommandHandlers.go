package server

import (
	"encoding/json"
	"net/http"
	"time"

	"GitHub/go-chat/backend/internal/domain"

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

	message, err := domain.NewMessage(request.ConversationId, userID, domain.MessageTypeUser, request.Content)
	if err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	_, err = s.message.Send(r.Context(), message, userID)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	response := struct {
		MessageId string    `json:"message_id"`
		CreatedAt time.Time `json:"created_at"`
	}{
		MessageId: message.ID.String(),
		CreatedAt: time.Now(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}
