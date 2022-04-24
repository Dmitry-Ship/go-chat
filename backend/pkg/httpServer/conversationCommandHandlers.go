package httpServer

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (s *HTTPServer) handleCreatePrivateConversationIfNotExists(w http.ResponseWriter, r *http.Request) {
	request := struct {
		ToUserId uuid.UUID `json:"to_user_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value("userId").(uuid.UUID)

	conversationId, err := s.app.Commands.ConversationService.CreatePrivateConversationIfNotExists(userID, request.ToUserId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
	}{
		ConversationId: conversationId,
	}

	json.NewEncoder(w).Encode(response)
}

func (s *HTTPServer) handleCreatePublicConversation(w http.ResponseWriter, r *http.Request) {
	request := struct {
		ConversationName string    `json:"conversation_name"`
		ConversationId   uuid.UUID `json:"conversation_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value("userId").(uuid.UUID)

	err = s.app.Commands.ConversationService.CreatePublicConversation(request.ConversationId, request.ConversationName, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("OK")
}

func (s *HTTPServer) handleDeleteConversation(w http.ResponseWriter, r *http.Request) {
	request := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.app.Commands.ConversationService.DeleteConversation(request.ConversationId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("OK")
}

func (s *HTTPServer) handleJoinPublicConversation(w http.ResponseWriter, r *http.Request) {
	request := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value("userId").(uuid.UUID)

	err = s.app.Commands.ConversationService.JoinPublicConversation(request.ConversationId, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("OK")
}

func (s *HTTPServer) handleLeavePublicConversation(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userId").(uuid.UUID)
	request := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.app.Commands.ConversationService.LeavePublicConversation(request.ConversationId, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("OK")
}

func (s *HTTPServer) handleRenamePublicConversation(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userId").(uuid.UUID)
	request := struct {
		ConversationId   uuid.UUID `json:"conversation_id"`
		ConversationName string    `json:"new_name"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.app.Commands.ConversationService.RenamePublicConversation(request.ConversationId, userID, request.ConversationName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("OK")
}
