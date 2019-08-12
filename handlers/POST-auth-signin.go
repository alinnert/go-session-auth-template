package handlers

import (
	"auth-server/models"
	"auth-server/services/validator"
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
		// #region Read request body
		input := &signinHandlerRequestBody{}
		err := json.NewDecoder(r.Body).Decode(input)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while parsing request body.")
			return
		}
		// #endregion Read request body

		// #region Validate request body
		err = validator.Validate(
			validator.Check(
				input.Email, "email",
				validator.StringIsNotEmpty(),
				validator.StringIsEmail(),
			),
			validator.Check(
				input.Password, "password",
				validator.StringIsNotEmpty(),
				validator.StringMinLength(6),
			),
		)
		if handled := handleValidationErrors(w, err); handled {
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
		sessionManager := values.GetSessionManager()

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
