package utils

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON sends a JSON response with the specified status code and data.
func RespondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Handle encoding error if needed
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}
