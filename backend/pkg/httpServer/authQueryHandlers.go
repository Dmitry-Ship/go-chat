package httpServer

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (s *HTTPServer) handleGetUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userId").(uuid.UUID)
	user, err := s.app.Queries.UsersRepository.GetUserByID(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)

}
