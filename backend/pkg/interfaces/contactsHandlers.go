package interfaces

import (
	"GitHub/go-chat/backend/pkg/application"

	"encoding/json"
	"net/http"
)

func HandleGetContacts(contactsService application.ContactsQueryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contacts, err := contactsService.GetContacts()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(contacts)
	}
}
