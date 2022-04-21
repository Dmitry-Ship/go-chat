package httpServer

import (
	"encoding/json"
	"net/http"
)

func (s *HTTPServer) handleGetContacts(w http.ResponseWriter, r *http.Request) {
	contacts, err := s.app.Queries.UsersRepository.FindAllUsers()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(contacts)
}
