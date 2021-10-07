package common

import (
	"encoding/json"
	"net/http"
)

func SendJSONresponse(response interface{}, w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "https://agitated-lalande-771c50.netlify.app")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
