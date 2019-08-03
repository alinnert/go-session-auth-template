package handlers

import (
	"net/http"
)

// GetSecret GET /secret
func GetSecret() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(w, "It's a secret to everybody.")
	}
}
