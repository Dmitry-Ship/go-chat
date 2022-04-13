package httpHandlers

import (
	"GitHub/go-chat/backend/pkg/readModel"
	"encoding/json"
	"net/http"
)

func HandleGetContacts(userQueryRepository readModel.UserQueryRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contacts, err := userQueryRepository.FindAllUsers()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(contacts)
	}
}
