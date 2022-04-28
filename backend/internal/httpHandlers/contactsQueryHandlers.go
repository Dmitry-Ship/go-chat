package httpHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (s *QueryController) handleGetContacts(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userId").(uuid.UUID)
	contacts, err := s.queries.Users.GetContacts(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(contacts)
}
