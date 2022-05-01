package httpHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (s *QueryController) handleGetContacts(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

	if !ok {
		http.Error(w, "userId not found in context", http.StatusInternalServerError)
		return
	}

	contacts, err := s.queries.GetContacts(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(contacts)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
