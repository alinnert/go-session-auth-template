package handlers

import (
	"auth-server/globals"
	"auth-server/models"
	"net/http"

	"github.com/dgraph-io/badger"
	"github.com/go-chi/chi"
)

// GetUser GET /user
// This is a route for debugging purposes. It fetches the data of one user.
// That way you can check if everything works on the database side.
func GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// #region Validate query params
		email := chi.URLParam(r, "email")
		// #endregion Get and validate query params

		// #region Get user
		db := r.Context().Value(globals.DBContext).(*badger.DB)
		user, err := models.GetUserByEmail(db, email)
		if err != nil {
			WriteErrorResponse(w, http.StatusUnauthorized, err,
				"Error while retrieving user.")
			return
		}
		// #endregion Get user

		if user == nil {
			WriteErrorResponse(w, http.StatusNotFound, nil, "User not found.")
			return
		}

		WriteResponse(w, user)
	}
}
