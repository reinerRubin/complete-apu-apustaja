package handler

import (
	"encoding/json"
	"net/http"
)

// TODO make me clever
func WriteError(w http.ResponseWriter, r *ErrorResponse) {
	w.WriteHeader(400)

	payload, err := json.Marshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}
