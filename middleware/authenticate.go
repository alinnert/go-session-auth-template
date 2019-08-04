package middleware

import (
	"auth-server/handlers"
	"auth-server/values"
	"net/http"
)

// Authenticate Adds authentication
func Authenticate() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionManager := values.GetSessionManager()
			user := sessionManager.GetString(r.Context(), "user")

			if user == "" {
				handlers.WriteErrorResponse(w, http.StatusForbidden, nil,
					"You're not authenticated.")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
