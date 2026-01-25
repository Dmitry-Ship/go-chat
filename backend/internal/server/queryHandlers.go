package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"GitHub/go-chat/backend/internal/readModel"
	"github.com/google/uuid"
)

func (s *Server) handleGetContacts(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
		return
	}

	paginationInfo, ok := r.Context().Value(paginationKey).(pagination)

	if !ok {
		http.Error(w, "pagination info not found in context", http.StatusInternalServerError)
		return
	}

	contacts, err := s.queries.GetContacts(userID, paginationInfo)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	err = json.NewEncoder(w).Encode(contacts)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) handleGetPotentialInvitees(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	conversationIDQuery := query.Get("conversation_id")
	conversationID, err := uuid.Parse(conversationIDQuery)

	if err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	paginationInfo, ok := r.Context().Value(paginationKey).(pagination)

	if !ok {
		http.Error(w, "pagination info not found in context", http.StatusInternalServerError)
		return
	}

	contacts, err := s.queries.GetPotentialInvitees(conversationID, paginationInfo)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	err = json.NewEncoder(w).Encode(contacts)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) handleGetConversations(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
		return
	}

	paginationInfo, ok := r.Context().Value(paginationKey).(pagination)

	if !ok {
		http.Error(w, "pagination info not found in context", http.StatusInternalServerError)
		return
	}

	conversations, err := s.queries.GetUserConversations(userID, paginationInfo)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	err = json.NewEncoder(w).Encode(conversations)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) handleGetConversationsMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	conversationIDQuery := query.Get("conversation_id")
	conversationID, err := uuid.Parse(conversationIDQuery)

	if err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	cursor, err := parseMessageCursor(query.Get("cursor"))
	if err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	limit := parseMessageLimit(query)

	page, err := s.queries.GetConversationMessages(conversationID, cursor, limit)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	for i := range page.Messages {
		page.Messages[i].User = nil
	}

	err = json.NewEncoder(w).Encode(page)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) handleGetConversationUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	conversationIDQuery := query.Get("conversation_id")
	if _, err := uuid.Parse(conversationIDQuery); err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	idsParam := query.Get("ids")
	if idsParam == "" {
		response := readModel.ConversationUsersResponse{
			Users: map[uuid.UUID]readModel.UserDTO{},
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			returnError(w, http.StatusInternalServerError, err)
		}
		return
	}

	uniqueIDs := make([]uuid.UUID, 0)
	seenIDs := make(map[uuid.UUID]struct{})
	for _, idStr := range strings.Split(idsParam, ",") {
		if idStr == "" {
			continue
		}
		userID, err := uuid.Parse(idStr)
		if err != nil {
			returnError(w, http.StatusBadRequest, err)
			return
		}
		if _, exists := seenIDs[userID]; exists {
			continue
		}
		seenIDs[userID] = struct{}{}
		uniqueIDs = append(uniqueIDs, userID)
	}

	users, err := s.queries.GetUsersByIDs(uniqueIDs)
	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	usersMap := make(map[uuid.UUID]readModel.UserDTO, len(users))
	for _, user := range users {
		usersMap[user.ID] = user
	}

	response := readModel.ConversationUsersResponse{
		Users: usersMap,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) handleGetConversation(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
		return
	}

	query := r.URL.Query()

	conversationIDQuery := query.Get("conversation_id")
	conversationID, err := uuid.Parse(conversationIDQuery)

	if err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	conversation, err := s.queries.GetConversation(conversationID, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(conversation)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) handleGetParticipants(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	conversationIDQuery := query.Get("conversation_id")
	conversationID, err := uuid.Parse(conversationIDQuery)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
		return
	}

	paginationInfo, ok := r.Context().Value(paginationKey).(pagination)

	if !ok {
		http.Error(w, "pagination info not found in context", http.StatusInternalServerError)
		return
	}

	contacts, err := s.queries.GetParticipants(conversationID, userID, paginationInfo)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	err = json.NewEncoder(w).Encode(contacts)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
		return
	}

	user, err := s.queries.GetUserByID(userID)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	err = json.NewEncoder(w).Encode(user)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}
