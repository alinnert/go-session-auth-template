package handlers

import (
	"net/http"
)

// GetIndex GET /
func GetIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(w, "Demo Auth Server is running.")
	}
}
