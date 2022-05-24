package httpHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (s *commandController) handleCreateDirectConversationIfNotExists(w http.ResponseWriter, r *http.Request) {
	request := struct {
		ToUserId uuid.UUID `json:"to_user_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userId not found in context", http.StatusInternalServerError)
		return
	}

	conversationId, err := s.commands.ConversationService.StartDirectConversation(userID, request.ToUserId)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	response := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
	}{
		ConversationId: conversationId,
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *commandController) handleCreateGroupConversation(w http.ResponseWriter, r *http.Request) {
	request := struct {
		ConversationName string    `json:"conversation_name"`
		ConversationId   uuid.UUID `json:"conversation_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userId not found in context", http.StatusInternalServerError)
		return
	}

	err = s.commands.ConversationService.CreateGroupConversation(request.ConversationId, request.ConversationName, userID)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	err = json.NewEncoder(w).Encode("OK")

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *commandController) handleDeleteConversation(w http.ResponseWriter, r *http.Request) {
	request := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userId not found in context", http.StatusInternalServerError)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	err = s.commands.ConversationService.DeleteGroupConversation(request.ConversationId, userID)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	err = json.NewEncoder(w).Encode("OK")

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *commandController) handleJoinGroupConversation(w http.ResponseWriter, r *http.Request) {
	request := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "userId not found in context", http.StatusInternalServerError)
		return
	}

	err = s.commands.ConversationService.JoinGroupConversation(request.ConversationId, userID)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	err = json.NewEncoder(w).Encode("OK")

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *commandController) handleLeaveGroupConversation(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userId not found in context", http.StatusInternalServerError)
		return
	}

	request := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	err = s.commands.ConversationService.LeaveGroupConversation(request.ConversationId, userID)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	err = json.NewEncoder(w).Encode("OK")

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *commandController) handleRenameGroupConversation(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userId not found in context", http.StatusInternalServerError)
		return
	}

	request := struct {
		ConversationId   uuid.UUID `json:"conversation_id"`
		ConversationName string    `json:"new_name"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	err = s.commands.ConversationService.RenameGroupConversation(request.ConversationId, userID, request.ConversationName)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	err = json.NewEncoder(w).Encode("OK")

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *commandController) handleInviteToGroupConversation(w http.ResponseWriter, r *http.Request) {
	request := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
		InviteeId      uuid.UUID `json:"user_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	err = s.commands.ConversationService.InviteToGroupConversation(request.ConversationId, request.InviteeId)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	err = json.NewEncoder(w).Encode("OK")

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}
