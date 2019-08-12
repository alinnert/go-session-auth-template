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
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// SigninHandler POST /auth/signin
func SigninHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// #region Validate request body
		input := &signinHandlerRequestBody{}
		err := json.NewDecoder(r.Body).Decode(input)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while parsing request body.")
			return
		}

		err = values.ValidateRequest.Struct(input)
		if err != nil {
			WriteErrorListResponse(w, http.StatusBadRequest, err,
				"Request body is not valid.")
			return
		}
		// #endregion Validate request body

		// #region Fetch user
		db := r.Context().Value(values.DBContext).(*badger.DB)

		matchingUser, err := models.GetUserByEmail(db, input.Email)
		if err != nil || matchingUser == nil {
			WriteErrorResponse(w, http.StatusUnauthorized, err,
				"Could not signin with provided credentials.")
			return
		}
		// #endregion Fetch user

		// #region Check password
		actualPassword := []byte(matchingUser.Password)
		providedPassword := []byte(input.Password)

		err = bcrypt.CompareHashAndPassword(actualPassword, providedPassword)
		if err != nil {
			WriteErrorResponse(w, http.StatusUnauthorized, err,
				"Could not signin with provided credentials.")
			return
		}
		// #endregion Check password

		// #region Set session cookie
		sessionManager := values.SessionManager

		err = sessionManager.RenewToken(r.Context())
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while renewing session token.")
			return
		}

		sessionManager.Put(r.Context(), "user", input.Email)
		// #endregion Set session cookie

		WriteResponse(w, nil)
	}
}
