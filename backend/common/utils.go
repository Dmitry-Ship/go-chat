package common

import (
	"encoding/json"
	"net/http"
)

func SendJSONresponse(response interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "https://agitated-lalande-771c50.netlify.app")

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
