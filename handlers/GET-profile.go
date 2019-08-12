package handlers

import (
	"auth-server/values"
	"net/http"
)

// GetProfile GET /profile
func GetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionManager := values.SessionManager
		currentUser := sessionManager.GetString(r.Context(), "user")

		if currentUser == "" {
			WriteResponse(w, "You're not logged in.")
			return
		}

		WriteResponse(w, "You're logged in as "+currentUser)
	}
}
