package httpHandlers

import (
	"GitHub/go-chat/backend/pkg/readModel"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func HandleGetUser(userQueryRepository readModel.UserQueryRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userId").(uuid.UUID)
		user, err := userQueryRepository.GetUserByID(userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(user)

	}
}
