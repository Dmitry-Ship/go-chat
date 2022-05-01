package httpHandlers

import (
	"GitHub/go-chat/backend/internal/readModel"

	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (s *QueryController) handleGetConversations(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userId not found in context", http.StatusInternalServerError)
		return
	}

	conversations, err := s.queries.GetUserConversations(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(conversations)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *QueryController) handleGetConversationsMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	conversationIdQuery := query.Get("conversation_id")
	conversationId, err := uuid.Parse(conversationIdQuery)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userId not found in context", http.StatusInternalServerError)
		return
	}

	messages, err := s.queries.GetConversationMessages(conversationId, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Messages []*readModel.MessageDTO `json:"messages"`
	}{
		Messages: messages,
	}

	err = json.NewEncoder(w).Encode(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *QueryController) handleGetConversation(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userId not found in context", http.StatusInternalServerError)
		return
	}

	query := r.URL.Query()

	conversationIdQuery := query.Get("conversation_id")
	conversationId, err := uuid.Parse(conversationIdQuery)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conversation, err := s.queries.GetConversation(conversationId, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(conversation)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
