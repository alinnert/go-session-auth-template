package handlers

import (
	"encoding/json"
	"net/http"
)

// WriteResponse writes the response if it's not an error.
func WriteResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	responseMap := map[string]interface{}{"status": "ok"}

	if data != nil {
		responseMap["data"] = data
	}

	json.NewEncoder(w).Encode(responseMap)
}

// WriteErrorResponse writes an error response and sets a status.
func WriteErrorResponse(
	w http.ResponseWriter, status int, err error, message string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	responseMap := map[string]string{"status": "error", "message": message}

	if err != nil {
		responseMap["details"] = err.Error()
	}

	json.NewEncoder(w).Encode(responseMap)
}
