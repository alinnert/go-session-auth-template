package handlers

import (
	"auth-server/models"
	"auth-server/values"
	"encoding/json"
	"net/http"

	"github.com/dgraph-io/badger"
	"golang.org/x/crypto/bcrypt"
)

type signinHandlerRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SigninHandler POST /auth/signin
func SigninHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &signinHandlerRequestBody{}
		err := json.NewDecoder(r.Body).Decode(input)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while parsing request body.")
		}

		sessionManager := values.GetSessionManager()

		matchingUser, err := models.GetUserByEmail(
			r.Context().Value(values.DBContext).(*badger.DB),
			input.Email,
		)
		if err != nil {
			WriteErrorResponse(w, http.StatusUnauthorized, err,
				"Could not signin with provided credentials.")
			return
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(matchingUser.Password),
			[]byte(input.Password))
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

		sessionManager.Put(r.Context(), "user", input.Email)

		WriteResponse(w, nil)
	}
}
