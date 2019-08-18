package handlers

import (
	"auth-server/globals"
	"auth-server/models"
	"net/http"
)

// GetUser GET /user
// This is a route for debugging purposes. It fetches the data of one user.
// That way you can check if everything works on the database side.
func GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// #region Get user
		user := r.Context().Value(globals.UserContext).(*models.User)
		// #endregion Get user

		if user == nil {
			WriteErrorResponse(w, http.StatusNotFound, nil, "User not found.")
			return
		}

		WriteResponse(w, user)
	}
}
