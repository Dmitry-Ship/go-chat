package httpHandlers

import (
	"GitHub/go-chat/backend/pkg/readModel"

	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func HandleGetConversations(conversationQueryRepository readModel.ConversationQueryRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conversations, err := conversationQueryRepository.FindAllConversations()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(conversations)
	}
}

func HandleGetConversationsMessages(messageQueryRepository readModel.MessageQueryRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		conversationIdQuery := query.Get("conversation_id")
		conversationId, err := uuid.Parse(conversationIdQuery)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID, _ := r.Context().Value("userId").(uuid.UUID)

		messages, err := messageQueryRepository.FindAllByConversationID(conversationId, userID)

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
}

func HandleGetConversation(conversationQueryRepository readModel.ConversationQueryRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userId").(uuid.UUID)
		query := r.URL.Query()

		conversationIdQuery := query.Get("conversation_id")
		conversationId, err := uuid.Parse(conversationIdQuery)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		conversation, err := conversationQueryRepository.GetConversationByID(conversationId, userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(conversation)
	}
}
