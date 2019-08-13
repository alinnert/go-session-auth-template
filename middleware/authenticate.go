package middleware

import (
	"auth-server/globals"
	"auth-server/handlers"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

// Authenticate Adds authentication
func Authenticate() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionManager := r.Context().Value(globals.SessionContext).(*scs.SessionManager)
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
