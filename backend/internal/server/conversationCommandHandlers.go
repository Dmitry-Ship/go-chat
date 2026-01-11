package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) handleStartDirectConversation(w http.ResponseWriter, r *http.Request) {
	request := struct {
		ToUserID uuid.UUID `json:"to_user_id"`
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

	conversationID, err := s.directConversation.StartDirectConversation(r.Context(), userID, request.ToUserID)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	response := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
	}{
		ConversationId: conversationID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) handleCreateGroupConversation(w http.ResponseWriter, r *http.Request) {
	request := struct {
		ConversationName string    `json:"conversation_name"`
		ConversationId   uuid.UUID `json:"conversation_id"`
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

	err := s.groupConversation.CreateGroupConversation(r.Context(), request.ConversationId, request.ConversationName, userID)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode("OK"); err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) handleDeleteConversation(w http.ResponseWriter, r *http.Request) {
	request := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	err := s.groupConversation.DeleteGroupConversation(r.Context(), request.ConversationId, userID)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode("OK"); err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) handleJoin(w http.ResponseWriter, r *http.Request) {
	request := struct {
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

	err := s.membership.Join(r.Context(), request.ConversationId, userID)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode("OK"); err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) handleLeave(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
		return
	}

	request := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	err := s.membership.Leave(r.Context(), request.ConversationId, userID)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode("OK"); err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) handleRename(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
		return
	}

	request := struct {
		ConversationId   uuid.UUID `json:"conversation_id"`
		ConversationName string    `json:"new_name"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	err := s.groupConversation.Rename(r.Context(), request.ConversationId, userID, request.ConversationName)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode("OK"); err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) handleInvite(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
		return
	}

	request := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
		InviteeId      uuid.UUID `json:"user_id"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	err := s.membership.Invite(r.Context(), request.ConversationId, userID, request.InviteeId)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode("OK"); err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) handleKick(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
		return
	}

	request := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
		TargetId       uuid.UUID `json:"user_id"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	err := s.membership.Kick(r.Context(), request.ConversationId, userID, request.TargetId)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode("OK"); err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}
