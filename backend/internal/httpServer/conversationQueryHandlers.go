package httpServer

import (
	"GitHub/go-chat/backend/internal/readModel"

	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (s *QueryController) handleGetConversations(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userId").(uuid.UUID)
	conversations, err := s.queries.Conversations.GetUserConversations(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(conversations)
}

func (s *QueryController) handleGetConversationsMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	conversationIdQuery := query.Get("conversation_id")
	conversationId, err := uuid.Parse(conversationIdQuery)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value("userId").(uuid.UUID)

	messages, err := s.queries.Messages.GetConversationMessages(conversationId, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Messages []*readModel.MessageDTO `json:"messages"`
	}{
		Messages: messages,
	}

	json.NewEncoder(w).Encode(data)
}

func (s *QueryController) handleGetConversation(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userId").(uuid.UUID)
	query := r.URL.Query()

	conversationIdQuery := query.Get("conversation_id")
	conversationId, err := uuid.Parse(conversationIdQuery)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conversation, err := s.queries.Conversations.GetConversation(conversationId, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(conversation)
}
