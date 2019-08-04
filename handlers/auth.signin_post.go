package handlers

import (
	"auth-server/models"
	"auth-server/values"
	"encoding/json"
	"net/http"

	"github.com/dgraph-io/badger"
	"golang.org/x/crypto/bcrypt"
)

// SigninHandler POST /auth/signin
func SigninHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		json.NewDecoder(r.Body).Decode(&user)
		sessionManager := values.GetSessionManager()

		matchingUser, err := models.GetUserByEmail(
			r.Context().Value(values.DBContext).(*badger.DB),
			user.Email,
		)
		if err != nil {
			WriteErrorResponse(w, http.StatusUnauthorized, err,
				"Could not signin with provided credentials.")
			return
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(matchingUser.Password),
			[]byte(user.Password))
		if err != nil {
			WriteErrorResponse(w, http.StatusUnauthorized, err,
				"Could not signin with provided credentials.")
			return
		}

		err = sessionManager.RenewToken(r.Context())
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while renewing session token.")
			return
		}

		sessionManager.Put(r.Context(), "user", user.Email)

		WriteResponse(w, nil)
	}
}
