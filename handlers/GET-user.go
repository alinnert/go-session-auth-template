package handlers

import (
	"auth-server/globals"
	"auth-server/models"
	"net/http"

	"github.com/dgraph-io/badger"
)

// GetUser GET /user
// This is a route for debugging purposes. It fetches the data of one user.
// That way you can check if everything works on the database side.
func GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// #region Validate query params
		emails, ok := r.URL.Query()["email"]
		if !ok || len(emails[0]) < 1 {
			WriteErrorResponse(w, http.StatusBadRequest, nil,
				"Parameter 'email' is missing.")
			return
		}
		// #endregion Get and validate query params

		// #region Get user
		user, err := models.GetUserByEmail(
			r.Context().Value(globals.DBContext).(*badger.DB),
			emails[0],
		)
		if err != nil {
			WriteErrorResponse(w, http.StatusUnauthorized, err,
				"Error while retrieving user.")
			return
		}
		if user == nil {
			WriteErrorResponse(w, http.StatusNotFound, nil, "User not found.")
			return
		}
		// #endregion Get user

		WriteResponse(w, user)
	}
}
