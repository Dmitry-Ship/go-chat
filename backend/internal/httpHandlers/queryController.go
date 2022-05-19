package httpHandlers

import (
	"GitHub/go-chat/backend/internal/readModel"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type queryController struct {
	queries readModel.QueriesRepository
}

func NewQueryController(queries readModel.QueriesRepository) *queryController {
	return &queryController{
		queries: queries,
	}
}

func (s *queryController) handleGetContacts(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userId not found in context", http.StatusInternalServerError)
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

func (s *queryController) handleGetPotentialInvitees(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	conversationIdQuery := query.Get("conversation_id")
	conversationId, err := uuid.Parse(conversationIdQuery)

	if err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	paginationInfo, ok := r.Context().Value(paginationKey).(pagination)

	if !ok {
		http.Error(w, "pagination info not found in context", http.StatusInternalServerError)
		return
	}

	contacts, err := s.queries.GetPotentialInvitees(conversationId, paginationInfo)

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

func (s *queryController) handleGetConversations(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userId not found in context", http.StatusInternalServerError)
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

func (s *queryController) handleGetConversationsMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	conversationIdQuery := query.Get("conversation_id")
	conversationId, err := uuid.Parse(conversationIdQuery)

	if err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userId not found in context", http.StatusInternalServerError)
		return
	}

	paginationInfo, ok := r.Context().Value(paginationKey).(pagination)

	if !ok {
		http.Error(w, "pagination info not found in context", http.StatusInternalServerError)
		return
	}

	messages, err := s.queries.GetConversationMessages(conversationId, userID, paginationInfo)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}

	err = json.NewEncoder(w).Encode(messages)

	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *queryController) handleGetConversation(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userId not found in context", http.StatusInternalServerError)
		return
	}

	query := r.URL.Query()

	conversationIdQuery := query.Get("conversation_id")
	conversationId, err := uuid.Parse(conversationIdQuery)

	if err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	conversation, err := s.queries.GetConversation(conversationId, userID)

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

func (s *queryController) handleGetUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userId not found in context", http.StatusInternalServerError)
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
