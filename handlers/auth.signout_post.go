package handlers

import (
	"auth-server/values"
	"net/http"
)

// SignoutHandler POST /auth/signout
func SignoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionManager := values.GetSession()
		err := sessionManager.Destroy(r.Context())
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while destroying the session.")
			return
		}

		WriteResponse(w, nil)
	}
}
