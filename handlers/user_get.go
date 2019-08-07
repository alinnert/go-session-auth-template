package handlers

import (
	"auth-server/models"
	"auth-server/values"
	"net/http"

	"github.com/dgraph-io/badger"
)

// GetUser GET /user
// This is a route for debugging purposes. It fetches the data of one user.
// That way you can check if everything works on the database side.
func GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		emails, ok := r.URL.Query()["email"]
		if !ok || len(emails[0]) < 1 {
			WriteErrorResponse(w, http.StatusBadRequest, nil,
				"Parameter 'email' is missing.")
			return
		}

		user, err := models.GetUserByEmail(
			r.Context().Value(values.DBContext).(*badger.DB),
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

		WriteResponse(w, user)
	}
}
